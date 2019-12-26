package databus_client

import (
	"code.byted.org/gopkg/metrics"
	"errors"
	"sync"
	"time"
)

const (
	DEFAULT_SOCKET_PATH         = "/opt/tmp/sock/databus_collector.seqpacket.sock"
	DEFAULT_STREAM_SOCKET_PATH  = "/opt/tmp/sock/databus_collector.stream.sock"
	PACKET_SIZE_LIMIT           = 212992 // 208KB
	DEFAULT_TIMEOUT             = 100 * time.Millisecond
	ERROR_TOLERATE              = 50
	DEFAULT_MAX_CONN_NUM        = 5
	READ_BUFFER_SIZE            = 1 << 10
	METRICS_PREFIX              = "inf.databus"
	METRICS_SUCC                = "collect.success"
	METRICS_FAIL                = "collect.fail"
	METRICS_RETRY               = "collect.retry"
	METRICS_CACHE_SUCCESS       = "collect.cache_success"
	METRICS_CACHE_FAIL          = "collect.cache_fail"
	METRICS_CACHE_TIME_EXPIRED  = "collect.cache_time_expired"
	RESP_SUCC_CODE              = 0
	RESP_UNKNOWN_CHANNEL_CODE   = -1
	RESP_BUFFER_FULL_CODE       = -2
	STREAM_HEADER_SIZE          = 64       // 1 byte for version, 4 bytes for body length, the rests are reserved
	MAX_STREAM_READ_BUFFER_SIZE = 10485760 // max message size 10MB
	DEFAULT_CACHE_SIZE          = 35 * 1024 * 1024
	DEFAULT_CACHE_MAX_TIME      = 0
	KVERSION                    = 1
)

var (
	ErrAlreadyClosed     = errors.New("client already closed")
	ErrUnknownChannel    = errors.New("databus agent said unknown channel")
	ErrAgentBufferFull   = errors.New("databus agent said buffer full, may be slow down send")
	ErrSeqpacketTooLarge = errors.New("packet too large, use NewStreamCollector instead")
	ErrStreamTooLarge    = errors.New("stream message too large, cannot send a packet larger than 10MB")
	ErrCacheFull         = errors.New("Cache Full")
	ErrCacheTimeExpired  = errors.New("Cache TimeExpired")
	metricsClient        = newClientMetrics()
)

func getNowTimeMs() int64 {
	return time.Now().UnixNano() / 1000 / 1000
}

type Condition struct {
	lock      sync.Locker
	ch        chan int
	interrupt chan int
}

func NewCondition(locker sync.Locker) *Condition {
	c := new(Condition)
	c.ch = make(chan int)
	c.interrupt = make(chan int)
	c.lock = locker
	return c
}

func (c *Condition) AwaitWithTimeOut(timeout_ms int64) bool {
	timeout := time.Duration(timeout_ms) * time.Millisecond
	c.lock.Unlock()
	defer c.lock.Lock()
	select {
	case c.ch <- 1:
		return true
	case <-time.After(timeout):
		return false
	case c.interrupt <- 1:
		return false
	}
}

func (c *Condition) Signal() {
	select {
	case _ = <-c.ch:
		return
	default:
		return
	}
}

func (c *Condition) SignalAll() {
	for {
		select {
		case _ = <-c.ch:
			break
		default:
			return
		}
	}
}

func (c *Condition) Interrupt() {
	for {
		select {
		case _ = <-c.interrupt:
			break
		default:
			return
		}
	}
}

// metrics
func newClientMetrics() *metrics.MetricsClientV2 {
	metricsClient := metrics.NewDefaultMetricsClientV2(METRICS_PREFIX, false)
	_ = metricsClient.DefineCounter(METRICS_SUCC)
	_ = metricsClient.DefineCounter(METRICS_FAIL)
	_ = metricsClient.DefineCounter(METRICS_RETRY)
	_ = metricsClient.DefineCounter(METRICS_CACHE_SUCCESS)
	_ = metricsClient.DefineCounter(METRICS_CACHE_FAIL)
	_ = metricsClient.DefineCounter(METRICS_CACHE_TIME_EXPIRED)
	return metricsClient
}

func recordSuccess(channel string, num int, is_cache bool) {
	if is_cache {
		_ = metricsClient.EmitCounter(METRICS_SUCC, num,
			metrics.T{"channel", channel},
			metrics.T{"from_cache", "True"},
		)
	} else {
		_ = metricsClient.EmitCounter(METRICS_SUCC, num,
			metrics.T{"channel", channel},
			metrics.T{"from_cache", "False"},
		)
	}
}

func recordFail(channel string, num int, needResp bool) {
	if needResp {
		_ = metricsClient.EmitCounter(METRICS_FAIL, num,
			metrics.T{"channel", channel},
			metrics.T{"need_resp", "True"},
		)
	} else {
		_ = metricsClient.EmitCounter(METRICS_FAIL, num,
			metrics.T{"channel", channel},
			metrics.T{"need_resp", "False"},
		)
	}
}

func recordRetry(channel string, num int, is_cache bool) {
	if is_cache {
		_ = metricsClient.EmitCounter(METRICS_RETRY, num,
			metrics.T{"channel", channel},
			metrics.T{"from_cache", "True"},
		)
	} else {
		_ = metricsClient.EmitCounter(METRICS_RETRY, num,
			metrics.T{"channel", channel},
			metrics.T{"from_cache", "False"},
		)
	}
}

func recordCacheSuccess(channel string, num int) {
	_ = metricsClient.EmitCounter(METRICS_CACHE_SUCCESS, num,
		metrics.T{"channel", channel},
	)
}

func recordCacheFail(channel string, num int) {
	_ = metricsClient.EmitCounter(METRICS_CACHE_FAIL, num,
		metrics.T{"channel", channel},
	)
}

func recordTimeExpired(channel string, num int) {
	_ = metricsClient.EmitCounter(METRICS_CACHE_TIME_EXPIRED, num,
		metrics.T{"channel", channel},
	)
}
