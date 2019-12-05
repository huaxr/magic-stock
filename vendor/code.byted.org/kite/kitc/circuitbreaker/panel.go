package circuit

import (
	"sync"
	"sync/atomic"
	"time"
)

// PanelStateChangeHandler .
type PanelStateChangeHandler func(key string, oldState, newState State, m Metricser)

// Panel manages a batch of circuitbreakers
// TODO(zhangyuanjia): remove unused breaker
type Panel struct {
	breakers       sync.Map
	defaultOptions Options
	changeHandler  PanelStateChangeHandler
	ticker         *time.Ticker
}

// NewPanel .
func NewPanel(changeHandler PanelStateChangeHandler,
	defaultOptions Options) (*Panel, error) {
	if defaultOptions.BucketTime <= 0 {
		defaultOptions.BucketTime = DEFAULT_BUCKET_TIME
	}

	if defaultOptions.BucketNums <= 0 {
		defaultOptions.BucketNums = DEFAULT_BUCKET_NUMS
	}

	if defaultOptions.CoolingTimeout <= 0 {
		defaultOptions.CoolingTimeout = DEFAULT_COOLING_TIMEOUT
	}

	if defaultOptions.DetectTimeout <= 0 {
		defaultOptions.DetectTimeout = DEFAULT_DETECT_TIMEOUT
	}
	_, err := NewBreaker(defaultOptions)
	if err != nil {
		return nil, err
	}
	p := &Panel{
		breakers:       sync.Map{},
		defaultOptions: defaultOptions,
		changeHandler:  changeHandler,
		ticker:         time.NewTicker(defaultOptions.BucketTime),
	}
	go p.tick()
	return p, nil
}

// GetBreaker .
func (p *Panel) GetBreaker(key string) *Breaker {
	cb, ok := p.breakers.Load(key)
	if ok {
		return cb.(*Breaker)
	}

	op := p.defaultOptions
	if p.changeHandler != nil {
		op.StateChangeHandler = func(oldState, newState State, m Metricser) {
			p.changeHandler(key, oldState, newState, m)
		}
	}
	ncb, _ := NewBreaker(op)
	cb, ok = p.breakers.LoadOrStore(key, ncb)
	return cb.(*Breaker)
}

// RemoveBreaker 不是并发安全的
func (p *Panel) RemoveBreaker(key string) {
	p.breakers.Delete(key)
}

// DumpBreakers .
func (p *Panel) DumpBreakers() map[string]*Breaker {
	breakers := make(map[string]*Breaker)
	p.breakers.Range(func(key, value interface{}) bool {
		breakers[key.(string)] = value.(*Breaker)
		return true
	})
	return breakers
}

// Succeed .
func (p *Panel) Succeed(key string) {
	p.GetBreaker(key).Succeed()
}

// Fail .
func (p *Panel) Fail(key string) {
	p.GetBreaker(key).Fail()
}

// FailWithTrip .
func (p *Panel) FailWithTrip(key string, trip TripFunc) {
	p.GetBreaker(key).FailWithTrip(trip)
}

// Timeout .
func (p *Panel) Timeout(key string) {
	p.GetBreaker(key).Timeout()
}

// TimeoutWithTrip .
func (p *Panel) TimeoutWithTrip(key string, trip TripFunc) {
	p.GetBreaker(key).TimeoutWithTrip(trip)
}

// IsAllowed .
func (p *Panel) IsAllowed(key string) bool {
	return p.GetBreaker(key).IsAllowed()
}

// tick .
func (p *Panel) tick() {
	for range p.ticker.C {
		p.breakers.Range(func(key, value interface{}) bool {
			if b, ok := value.(*Breaker); ok {
				if w, ok := b.Metricser.(*window); ok {
					w.Lock()
					// 这一段必须在前面，因为latest可能会覆盖oldest
					if w.inWindow == w.bucketNums {
						// the lastest covered the oldest(latest == oldest)
						oldBucket := &w.buckets[w.oldest]
						atomic.AddInt64(&w.allSuccess, -oldBucket.Successes())
						atomic.AddInt64(&w.allFailure, -oldBucket.Failures())
						atomic.AddInt64(&w.allTimeout, -oldBucket.Timeouts())
						w.oldest++
						if w.oldest >= w.bucketNums {
							w.oldest = 0
						}
					} else {
						w.inWindow++
					}

					w.latest++
					if w.latest >= w.bucketNums {
						w.latest = 0
					}
					(&w.buckets[w.latest]).Reset()
					w.Unlock()
				}
			}
			return true
		})
	}
}
