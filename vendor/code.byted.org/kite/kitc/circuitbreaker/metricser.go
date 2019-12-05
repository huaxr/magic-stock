package circuit

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Metricser metricses errors, timeouts and successes
type Metricser interface {
	Fail()    // records a failure
	Succeed() // records a success
	Timeout() // records a timeout

	Failures() int64          // return the number of failures
	Successes() int64         // return the number of successes
	Timeouts() int64          // return the number of timeouts
	ConseErrors() int64       // return the consecutive errors recently
	ConseTime() time.Duration // return the consecutive error time
	ErrorRate() float64       // rate = (timeouts + failures) / (timeouts + failures + successes)
	Samples() int64           // (timeouts + failures + successes)
	Counts() (successes, failures, timeouts int64)

	Reset()
}

const (
	// bucket time is the time each bucket holds
	DEFAULT_BUCKET_TIME = time.Millisecond * 100

	// bucket nums is the number of buckets the metricser has;
	// the more buckets you have, the less counters you lose when
	// the oldest bucket expire;
	DEFAULT_BUCKET_NUMS = 100

	// default window size is (DEFAULT_BUCKET_TIME * DEFAULT_BUCKET_NUMS),
	// which is 10 seconds;
)

// bucket holds counts of failures and successes
type bucket struct {
	failure int64
	success int64
	timeout int64
}

// Reset resets the counts to 0 and refreshes the time stamp
func (b *bucket) Reset() {
	atomic.StoreInt64(&b.failure, 0)
	atomic.StoreInt64(&b.success, 0)
	atomic.StoreInt64(&b.timeout, 0)
}

func (b *bucket) Fail() {
	atomic.AddInt64(&b.failure, 1)
}

func (b *bucket) Succeed() {
	atomic.AddInt64(&b.success, 1)
}

func (b *bucket) Timeout() {
	atomic.AddInt64(&b.timeout, 1)
}

func (b *bucket) Failures() int64 {
	return atomic.LoadInt64(&b.failure)
}

func (b *bucket) Successes() int64 {
	return atomic.LoadInt64(&b.success)
}

func (b *bucket) Timeouts() int64 {
	return atomic.LoadInt64(&b.timeout)
}

// window maintains a ring of buckets and increments the failure and success
// counts of the current bucket.
type window struct {
	sync.RWMutex
	oldest  int32    // oldest bucket index
	latest  int32    // latest bucket index
	buckets []bucket // buckets this window holds

	bucketTime time.Duration // time each bucket holds
	bucketNums int32         // the numbe of buckets
	inWindow   int32         // the number of buckets in the window

	allSuccess int64
	allFailure int64
	allTimeout int64

	errStart int64
	conseErr int64
}

// NewWindow .
func NewWindow() Metricser {
	m, _ := NewWindowWithOptions(DEFAULT_BUCKET_TIME, DEFAULT_BUCKET_NUMS)
	return m
}

// NewWindowWithOptions creates a new window.
func NewWindowWithOptions(bucketTime time.Duration, bucketNums int32) (Metricser, error) {
	if bucketNums < 100 {
		return nil, fmt.Errorf("BucketNums can't be less than 100")
	}

	w := new(window)
	w.bucketNums = bucketNums
	w.bucketTime = bucketTime
	w.buckets = make([]bucket, w.bucketNums)

	w.Reset()
	return w, nil
}

// Fail records a failure in the current bucket.
func (w *window) Fail() {
	w.Lock()
	b := &w.buckets[atomic.LoadInt32(&w.latest)]
	atomic.AddInt64(&w.conseErr, 1)
	atomic.AddInt64(&w.allFailure, 1)
	if atomic.LoadInt64(&w.errStart) == 0 {
		atomic.StoreInt64(&w.errStart, time.Now().UnixNano())
	}
	w.Unlock()
	b.Fail()
}

// Success records a success in the current bucket.
func (w *window) Succeed() {
	w.RLock()
	b := &w.buckets[atomic.LoadInt32(&w.latest)]
	atomic.StoreInt64(&w.errStart, 0)
	atomic.StoreInt64(&w.conseErr, 0)
	atomic.AddInt64(&w.allSuccess, 1)
	w.RUnlock()
	b.Succeed()
}

// Timeout records a timeout in the current bucket
func (w *window) Timeout() {
	w.Lock()
	b := &w.buckets[atomic.LoadInt32(&w.latest)]
	atomic.AddInt64(&w.conseErr, 1)
	atomic.AddInt64(&w.allTimeout, 1)
	if atomic.LoadInt64(&w.errStart) == 0 {
		atomic.StoreInt64(&w.errStart, time.Now().UnixNano())
	}
	w.Unlock()
	b.Timeout()
}

func (w *window) Counts() (successes, failures, timeouts int64) {
	return atomic.LoadInt64(&w.allSuccess), atomic.LoadInt64(&w.allFailure), atomic.LoadInt64(&w.allTimeout)
}

// Successes returns the total number of successes recorded in all buckets.
func (w *window) Successes() int64 {
	return atomic.LoadInt64(&w.allSuccess)
}

// Failures returns the total number of failures recorded in all buckets.
func (w *window) Failures() int64 {
	return atomic.LoadInt64(&w.allFailure)
}

// Timeouts returns the total number of Timeout recorded in all buckets.
func (w *window) Timeouts() int64 {
	return atomic.LoadInt64(&w.allTimeout)
}

func (w *window) ConseErrors() int64 {
	return atomic.LoadInt64(&w.conseErr)
}

func (w *window) ConseTime() time.Duration {
	return time.Duration(time.Now().UnixNano() - atomic.LoadInt64(&w.errStart))
}

// ErrorRate returns the error rate calculated over all buckets, expressed as
// a floating point number (e.g. 0.9 for 90%)
func (w *window) ErrorRate() float64 {
	successes, failures, timeouts := w.Counts()

	if (successes + failures + timeouts) == 0 {
		return 0.0
	}

	return float64(failures+timeouts) / float64(successes+failures+timeouts)
}

func (w *window) Samples() int64 {
	successes, failures, timeouts := w.Counts()

	return successes + failures + timeouts
}

// Reset resets this window
func (w *window) Reset() {
	w.Lock()
	atomic.StoreInt32(&w.oldest, 0)
	atomic.StoreInt32(&w.latest, 0)
	atomic.StoreInt32(&w.inWindow, 1)
	atomic.StoreInt64(&w.conseErr, 0)
	atomic.StoreInt64(&w.allSuccess, 0)
	atomic.StoreInt64(&w.allFailure, 0)
	atomic.StoreInt64(&w.allTimeout, 0)
	(&w.buckets[w.latest]).Reset()
	w.Unlock() // don't use defer
}
