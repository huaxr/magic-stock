package metrics

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"unsafe"
)

const (
	_emit = "emit"
)

var debugMetrics = os.Getenv("DEBUG_GOPKG_METRICS") != ""

var unixdomainsock = ""

func init() {
	for _, path := range []string{"/opt/tmp/sock/metric.sock", "/tmp/metric.sock"} {
		for i := 0; i < 3; i++ {
			conn, err := net.Dial("unixgram", path)
			if err == nil {
				conn.Close()
				unixdomainsock = path
				return
			}
			if strings.Contains(err.Error(), "no such file") {
				break
			}
		}
	}
}

type metricEntry struct {
	mt     metricsType
	prefix string
	name   string
	v      float64
	ts     int64

	tt *cachedTags
}

func (m *metricEntry) IsCounter() bool {
	return m.mt == metricsTypeCounter || m.mt == metricsTypeRateCounter || m.mt == metricsTypeMeter
}

func (m *metricEntry) MarshalSize() int {
	// protocol: 6 fields: emit $type $prefix.name  $value $tag ""
	n := 0
	n += msgpackArrayHeaderSize
	n += msgpackStringSize(_emit)
	n += msgpackStringSize(m.mt.String())
	if len(m.prefix) > 0 {
		n += msgpackStringHeaderSize + (len(m.prefix) + 1 + len(m.name))
	} else {
		n += msgpackStringSize(m.name)
	}
	n += msgpackStringHeaderSize + floatStrSize(m.v) // int64 + "." + 5 prec float + str header
	n += msgpackStringHeaderSize + len(m.tt.Bytes())
	if m.ts > 0 {
		n += msgpackStringHeaderSize + int64StrSize(m.ts)
	} else {
		n += msgpackStringHeaderSize + 0
	}
	return n
}

func (m *metricEntry) AppendTo(p []byte) []byte {
	// protocol: 6 fields: emit $type $prefix.name  $value $tag ""
	p = msgpackAppendArrayHeader(p, 6)
	p = msgpackAppendString(p, _emit)
	p = msgpackAppendString(p, m.mt.String())
	if len(m.prefix) > 0 {
		p = msgpackAppendStringHeader(p, uint16(len(m.prefix)+1+len(m.name)))
		p = append(p, m.prefix...)
		p = append(p, '.')
		p = append(p, m.name...)
	} else {
		p = msgpackAppendString(p, m.name)
	}
	p = msgpackAppendStringHeader(p, uint16(floatStrSize(m.v)))
	p = appendFloat64(p, m.v)
	p = msgpackAppendStringHeader(p, uint16(len(m.tt.Bytes())))
	p = append(p, m.tt.Bytes()...)
	if m.ts > 0 {
		p = msgpackAppendStringHeader(p, uint16(int64StrSize(m.ts)))
		p = appendInt64(p, m.ts)
	} else {
		p = msgpackAppendString(p, "")
	}
	return p
}

func (m *metricEntry) MarshalTo(b []byte) {
	p := b[:0]
	p = m.AppendTo(p)
	if len(p) != len(b) {
		panic("buf size err")
	}
}

type metricsWriter struct {
	addr string

	mu    sync.RWMutex
	conns []net.Conn
}

func (w *metricsWriter) connect() (net.Conn, error) {
	if strings.HasPrefix(w.addr, "/") {
		return net.Dial("unixgram", w.addr)
	} else {
		return net.Dial("udp", w.addr)
	}
}

func (w *metricsWriter) getconn() (net.Conn, error) {
	w.mu.Lock()
	if len(w.conns) > 0 {
		conn := w.conns[len(w.conns)-1]
		w.conns = w.conns[:len(w.conns)-1]
		w.mu.Unlock()
		return conn, nil
	}
	w.mu.Unlock()
	return w.connect()
}

func (w *metricsWriter) putconn(conn net.Conn) {
	w.mu.Lock()
	if len(w.conns) < asyncGoroutines {
		w.conns = append(w.conns, conn)
	} else {
		conn.Close()
	}
	w.mu.Unlock()
}

func (w *metricsWriter) Write(b []byte) (int, error) {
	if w.addr == BlackholeAddr {
		if debugMetrics {
			fmt.Fprintf(os.Stderr, "gopkg/metrics: write %d bytes to blackhole\n", len(b))
		}
		return len(b), nil
	}
	conn, err := w.getconn()
	if err != nil {
		if debugMetrics {
			fmt.Fprintf(os.Stderr, "gopkg/metrics: conn err: %s\n", err)
		}
		return 0, err
	}
	n, err := conn.Write(b)
	if err != nil {
		if debugMetrics {
			fmt.Fprintf(os.Stderr, "gopkg/metrics: write err: %s\n", err)
		}
		conn.Close()
	} else {
		w.putconn(conn)
	}
	return n, err
}

type Sender struct {
	batch bool
	agg   bool

	w io.Writer
}

func NewSender(addr string) *Sender {
	s := &Sender{agg: true}
	if addr == DefaultMetricsServer && unixdomainsock != "" {
		addr = unixdomainsock
		s.batch = true
	}
	s.w = &metricsWriter{addr: addr}
	return s
}

type aggregatekey struct {
	prefix string
	name   string
	tt     uintptr // tags pointer
	mt     metricsType
}

type counterAggregator struct {
	keys []aggregatekey
	m    map[aggregatekey]*metricEntry
}

var counterAggregatorPool = sync.Pool{
	New: func() interface{} {
		return &counterAggregator{
			keys: make([]aggregatekey, 0, maxPendingSize),
			m:    make(map[aggregatekey]*metricEntry, maxPendingSize),
		}
	},
}

func (a *counterAggregator) Merge(ms []metricEntry) []metricEntry {
	for i := range ms {
		e := &ms[i]
		k := aggregatekey{
			prefix: e.prefix,
			name:   e.name,
			tt:     uintptr(unsafe.Pointer(e.tt)),
			mt:     e.mt,
		}
		v, ok := a.m[k]
		if ok {
			v.v += e.v
		} else {
			a.m[k] = e
			a.keys = append(a.keys, k)
		}
	}
	p := ms[:0]
	for _, k := range a.keys {
		p = append(p, *a.m[k])
		delete(a.m, k)
	}
	a.keys = a.keys[:0]
	return p
}

func (s *Sender) DisableCounterAggregator(v bool) {
	s.agg = !v
}

func (s *Sender) SendCounter(ms []metricEntry) {
	if len(ms) == 0 {
		return
	}
	if s.agg {
		a := counterAggregatorPool.Get().(*counterAggregator)
		s.Send(a.Merge(ms))
		counterAggregatorPool.Put(a)
	} else {
		s.Send(ms)
	}
}

func printMs(ms []metricEntry) {
	fmt.Fprintf(os.Stderr, "[DEBUG] gopkg/metrics: sending %d metrics to server:\n", len(ms))
	dup := make(map[aggregatekey]bool)
	for _, m := range ms {
		k := aggregatekey{prefix: m.prefix, name: m.name, mt: m.mt}
		if !dup[k] {
			fmt.Fprintf(os.Stderr, "[DEBUG] %s %s.%s\n", m.mt, m.prefix, m.name)
			dup[k] = true
		}
	}
}

func (s *Sender) Send(ms []metricEntry) {
	if debugMetrics {
		printMs(ms)
	}
	if !s.batch {
		p := wbufpool.Get().(*wbuf)
		defer wbufpool.Put(p)
		for _, m := range ms {
			s.w.Write(m.AppendTo(p.b[:0]))
		}
		return
	}
	// send bunch
	for len(ms) > 0 {
		ms = ms[s.sendbunch(ms):]
	}
}

type wbuf struct {
	b []byte

	mem [2 * maxBunchBytes]byte
}

var wbufpool = sync.Pool{
	New: func() interface{} {
		p := new(wbuf)
		p.b = p.mem[:0]
		return p
	},
}

func (s *Sender) sendbunch(ms []metricEntry) int {
	if len(ms) == 0 {
		return 0
	}

	// limit to send maxBunchBytes
	k := 0
	n := msgpackArrayHeaderSize
	for _, m := range ms {
		n += m.MarshalSize()
		if n >= maxBunchBytes {
			break
		}
		k++
	}

	if k == 0 {
		panic("metrics too large to send: " + ms[0].name) // single metrics > maxBunchBytes
	}

	ms = ms[:k]

	p := wbufpool.Get().(*wbuf)
	defer wbufpool.Put(p)
	p.b = p.b[:0]

	// marshal to p
	p.b = msgpackAppendArrayHeader(p.b, uint16(len(ms)))
	for _, m := range ms {
		p.b = m.AppendTo(p.b)
	}
	s.w.Write(p.b)
	return k
}
