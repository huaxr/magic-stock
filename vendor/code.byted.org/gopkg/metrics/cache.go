package metrics

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type cachedTags struct {
	tb []byte
}

func (e *cachedTags) Bytes() []byte {
	if e == nil {
		return nil
	}
	return e.tb
}

type tcache struct {
	m sync.Map

	setn uint64
}

func (c *tcache) Get(key []byte) *cachedTags {
	t, ok := c.m.Load(ss(key))
	if ok {
		return t.(*cachedTags)
	}
	return nil
}

func (c *tcache) Set(key []byte, tt *cachedTags) {
	if atomic.AddUint64(&c.setn, 1)&0x3fff == 0 {
		// every 0x3fff times call, we clear the map for memory leak issue
		// there is no reason to have so many tags
		// FIXME: sync.Map don't have Len method and `setn` may not equal to the len in concurrency env
		samples := make([]interface{}, 0, 3)
		c.m.Range(func(key interface{}, value interface{}) bool {
			c.m.Delete(key)
			if len(samples) < cap(samples) {
				samples = append(samples, key)
			}
			return true
		}) // clear map
		fmt.Fprintln(os.Stderr, "WARN: gopkg/metrics: too many tags. samples:", samples)
	}
	c.m.Store(string(key), tt)

}

func (c *tcache) GetOrCreate(tags []T, extTagBytes []byte) *cachedTags {
	k := make([]byte, 0, 500)

	// XXX: we dont sort the tags to improve performance
	// for v2 api, the tags list should be stable all the time
	// for v1 api which use map to store tags, we sort it in Map2Tags
	k = appendTags(k, tags)
	if e := c.Get(k); e != nil {
		return e
	}
	b := make([]byte, 0, len(k)+1+len(extTagBytes))
	b = append(b, k...)
	if len(extTagBytes) > 0 {
		b = append(b, '|')
		b = append(b, extTagBytes...)
	}
	e := &cachedTags{b}
	c.Set(k, e)
	return e
}

func (c *tcache) MakeMetricEntry(mt metricsType, prefix string, name string, v float64, ts int64, tags []T) metricEntry {
	e := metricEntry{mt: mt, prefix: prefix, name: name, ts: ts, v: v}
	e.tt = c.GetOrCreate(tags, extTags)
	return e
}

type mcache struct {
	ch chan *mes
	s  *Sender

	interval int64
	block    int32

	countermm mcacheMM
	othermm   mcacheMM
}

func newMcache(s *Sender) *mcache {
	m := &mcache{
		ch:       make(chan *mes, 2*asyncGoroutines),
		s:        s,
		interval: int64(flushInterval),
	}
	go m.flushloop()
	for i := 0; i < asyncGoroutines; i++ {
		go m.sendLoop()
	}
	return m
}

func (m *mcache) Send(e metricEntry) error {
	var mm *mes
	if e.IsCounter() {
		mm = m.countermm.Add(e, maxPendingSize)
	} else {
		mm = m.othermm.Add(e, maxPendingSize)
	}
	if mm == nil {
		return nil
	}
	if atomic.LoadInt32(&m.block) != 0 {
		m.ch <- mm
		return nil
	}
	select {
	case m.ch <- mm:
		return nil
	default:
		return ErrEmitBufferFull
	}
}

func (m *mcache) Flush() {
	if mm := m.countermm.Reset(); mm != nil {
		m.s.SendCounter(*mm)
		mesPool.Put(mm)
	}
	if mm := m.othermm.Reset(); mm != nil {
		m.s.Send(*mm)
		mesPool.Put(mm)
	}
}

func (m *mcache) SetFlushInterval(d time.Duration) {
	atomic.StoreInt64(&m.interval, int64(d))
}

func (m *mcache) FlushInterval() time.Duration {
	return time.Duration(atomic.LoadInt64(&m.interval))
}

func (m *mcache) flushloop() {
	for {
		<-time.After(m.FlushInterval())
		m.ch <- m.countermm.Reset()
		m.ch <- m.othermm.Reset()
	}
}

func (m *mcache) sendLoop() {
	for {
		e := <-m.ch
		if e == nil {
			continue
		}
		mm := *e
		if len(mm) > 0 {
			if mm[0].IsCounter() {
				m.s.SendCounter(mm)
			} else {
				m.s.Send(mm)
			}
		}
		mesPool.Put(e)
	}
}

type mes []metricEntry

type mcacheMM struct {
	mu sync.Mutex
	mm *mes
}

func (s *mcacheMM) Add(m metricEntry, max int) (mm *mes) {
	s.mu.Lock()
	if s.mm == nil {
		s.mm = mesPoolGet()
	}
	*s.mm = append(*s.mm, m)
	if len(*s.mm) >= max {
		mm, s.mm = s.mm, mesPoolGet()
	}
	s.mu.Unlock()
	return
}

func (s *mcacheMM) Reset() (mm *mes) {
	s.mu.Lock()
	mm, s.mm = s.mm, mesPoolGet()
	s.mu.Unlock()
	return
}

func mesPoolGet() *mes {
	mm := mesPool.Get().(*mes)
	*mm = (*mm)[:0]
	return mm
}

var mesPool = sync.Pool{
	New: func() interface{} {
		mm := make(mes, maxPendingSize)
		return &mm
	},
}
