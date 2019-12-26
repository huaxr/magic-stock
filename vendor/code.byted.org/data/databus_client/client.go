package databus_client

import (
	"container/list"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/golang/protobuf/proto"
)

// 1. 尽量复用DatabusCollector对象
// 2. 测试情况下BufferedDatabusCollector吞吐会比DatabusCollector吞吐低
// 3. 如果DatabusCollector对象太多会导致domain socket连接特别多

type conn struct {
	nc       *net.UnixConn
	config   *DatabusCollector
	errTimes int
	net      string
}

func (cn *conn) reconnect() error {
	con, err := net.DialUnix(cn.net, nil, &net.UnixAddr{cn.config.socketPath, cn.net})
	if err != nil {
		return err
	}
	cn.release()
	cn.nc = con
	return nil
}

// 检查错误类型，特定类型需要重置连接
func (cn *conn) checkError(err error) {
	if err == nil {
		return
	}

	resetConn := false

	if err == io.EOF {
		resetConn = true
	}

	if operr, ok := err.(*net.OpError); ok {
		causeErr := operr.Err
		if syscallerr, ok := causeErr.(*os.SyscallError); ok {
			if syscallerr.Err == syscall.EPIPE { // broken pipe
				resetConn = true
			}
		}
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() { // time out
		resetConn = true
	}

	if resetConn {
		cn.release()
		return
	}

	cn.errTimes += 1
}

func (cn *conn) Write(b []byte) error {
	// 发送失败次数达到阈值重建链接

	if cn.nc == nil || cn.errTimes > ERROR_TOLERATE {
		err := cn.reconnect()
		if err != nil {
			return err
		}
	}
	cn.nc.SetDeadline(time.Now().Add(cn.config.writeTimeout))
	if cn.net == "unix" {
		res := make([]byte, STREAM_HEADER_SIZE)
		res[0] = byte(1)
		binary.BigEndian.PutUint32(res[1:], uint32(len(b)))
		_, err := cn.nc.Write(res)
		if err != nil {
			cn.checkError(err)
			return err
		}
	}
	_, err := cn.nc.Write(b)
	cn.checkError(err)
	return err
}

func (cn *conn) Read(buf []byte) (n int, err error) {
	n, err = cn.nc.Read(buf)
	cn.checkError(err)
	return
}

func (cn *conn) release() {
	if cn.nc != nil {
		cn.nc.Close()
		cn.nc = nil
		cn.errTimes = 0
	}
}

type DatabusCollector struct {
	socketPath   string
	pool         chan *conn
	writeTimeout time.Duration
	net          string
	ctx          context.Context
	cache        *Cache
	consumerStop bool
	logTime      map[string]int64
	logTimeLock  sync.Mutex
}

func NewDefaultCollector() *DatabusCollector {
	return NewCollector(DEFAULT_SOCKET_PATH, DEFAULT_TIMEOUT, DEFAULT_MAX_CONN_NUM)
}

func NewCollectorWithTimeout(timeout time.Duration) *DatabusCollector {
	return NewCollector(DEFAULT_SOCKET_PATH, timeout, DEFAULT_MAX_CONN_NUM)
}

func NewStreamCollector() *DatabusCollector {
	return NewCollectorV1(DEFAULT_STREAM_SOCKET_PATH, DEFAULT_TIMEOUT, DEFAULT_MAX_CONN_NUM, "unix")
}

func NewStreamCollectorWithTimeout(timeout time.Duration) *DatabusCollector {
	return NewCollectorV1(DEFAULT_STREAM_SOCKET_PATH, timeout, DEFAULT_MAX_CONN_NUM, "unix")
}

func NewCollector(socketPath string, timeout time.Duration, maxConn int) *DatabusCollector {
	return NewCollectorV1(socketPath, timeout, maxConn, "unixpacket")
}

func NewCollectorV1(socketPath string, timeout time.Duration, maxConn int, net string) *DatabusCollector {
	cache := &Cache{cacheUsed: 0, cacheSize: DEFAULT_CACHE_SIZE,
		list: list.New(), cond: NewCondition(new(sync.Mutex)), maxCacheTime: DEFAULT_CACHE_MAX_TIME}

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	collector := &DatabusCollector{
		socketPath:   socketPath,
		writeTimeout: timeout,
		pool:         make(chan *conn, maxConn),
		net:          net,
		ctx:          context.Background(),
		cache:        cache,
		logTime:      make(map[string]int64),
	}
	go collector.Consumer()
	return collector
}

func (this *DatabusCollector) WithContext(ctx context.Context) *DatabusCollector {
	cache := &Cache{cacheUsed: 0, cacheSize: DEFAULT_CACHE_SIZE,
		list: list.New(), cond: NewCondition(new(sync.Mutex)), maxCacheTime: DEFAULT_CACHE_MAX_TIME}

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	cli := &DatabusCollector{
		socketPath:   this.socketPath,
		pool:         this.pool,
		writeTimeout: this.writeTimeout,
		net:          this.net,
		ctx:          ctx,
		consumerStop: false,
		cache:        cache,
		logTime:      make(map[string]int64),
	}
	return cli
}

// 从pool中取出一个连接
func (this *DatabusCollector) borrow() (*conn, error) {
	var cl *conn
	var ok bool
	select {
	case cl, ok = <-this.pool:
		if !ok {
			return nil, ErrAlreadyClosed
		}
	default:
		cl = this.newConn()
	}
	return cl, nil
}

// 归还连接
func (this *DatabusCollector) putBack(cl *conn) {
	select {
	case this.pool <- cl:
	default:
		cl.release()
	}
}

func (this *DatabusCollector) newConn() *conn {
	con := &conn{
		nc:       nil,
		config:   this,
		errTimes: 0,
		net:      this.net,
	}
	return con
}

func (this *DatabusCollector) translateChannel(channel string) *string {
	var stressTag string
	if v := this.ctx.Value("K_STRESS"); v != nil {
		stressTag = v.(string)
	}
	if len(stressTag) != 0 {
		newChannel := channel + "_stress"
		return &newChannel
	}
	return &channel
}

func (this *DatabusCollector) CollectArray(channel string, messages []*ApplicationMessage) error {
	err := this.CollectArrayWithTimeout(channel, messages, 0)
	return err
}

func (this *DatabusCollector) CollectArrayWithTimeout(channel string, messages []*ApplicationMessage, timeoutMs int64) error {
	var payload RequestPayload
	var buf []byte
	cacheErr := errors.New("-1")
	err := errors.New("-1")

	messageNum := len(messages)
	payload.Channel = this.translateChannel(channel)
	payload.Messages = messages
	buf, err = proto.Marshal(&payload)
	defer func() {
		if err == nil {
			recordSuccess(channel, messageNum, false)
		} else {
			if cacheErr != nil {
				recordFail(channel, messageNum, false)
			}
		}
	}()
	if err != nil {
		return err
	}
	if this.net == "unixpacket" {
		if len(buf) > PACKET_SIZE_LIMIT {
			err = ErrSeqpacketTooLarge
			return err
		}
	} else if len(buf) > MAX_STREAM_READ_BUFFER_SIZE {
		err = ErrStreamTooLarge
		return err
	}
	err = this.send(buf, channel, messageNum)
	if err == nil {
		return err
	}
	cacheErr = this.push(payload, int64(len(buf)), timeoutMs)
	if cacheErr == nil {
		return nil
	}
	return err
}

func (this *DatabusCollector) logInterval(interval_ms int64, formating string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	logType := fmt.Sprintf("%s,%d", file, line)
	this.logTimeLock.Lock()
	lastTime, ok := this.logTime[logType]
	if ok {
		if getNowTimeMs()-lastTime > interval_ms {
			this.logTime[logType] = getNowTimeMs()
			log.Printf("%s %s", logType, fmt.Sprintf(formating, args...))
		}
	} else {
		this.logTime[logType] = getNowTimeMs()
		log.Printf("%s %s", logType, fmt.Sprintf(formating, args...))
	}
	this.logTimeLock.Unlock()
}

func (this *DatabusCollector) send(buf []byte, channel string, messageNum int) error {
	var err error
	var con *conn
	for i := 0; i < cap(this.pool)+1; i++ {
		con, err = this.borrow()
		if err != nil {
			// only happen when pool is close, so not need to retry & put back
			return err
		}
		err = con.Write(buf)
		if err == nil {
			this.putBack(con)
			return nil
		}
		this.putBack(con)
		recordRetry(channel, messageNum, false)
		this.logInterval(10000, "channel:%s Retry: %d/%d for err: %s\n", channel, i,
			cap(this.pool), err.Error())
	}
	return err
}

func (this *DatabusCollector) CollectArrayWithResponse(channel string, messages []*ApplicationMessage) error {
	_, err := this.CollectArrayWithRespWithTimeout(channel, messages, 0)
	return err
}

func (this *DatabusCollector) CollectArrayWithResp(channel string, messages []*ApplicationMessage) (*ResponsePayload, error) {
	resp, err := this.CollectArrayWithRespWithTimeout(channel, messages, 0)
	return resp, err
}

func (this *DatabusCollector) CollectArrayWithResponseWithTimeout(channel string, messages []*ApplicationMessage, timeoutMs int64) error {
	_, err := this.CollectArrayWithRespWithTimeout(channel, messages, timeoutMs)
	return err
}

func (this *DatabusCollector) CollectArrayWithRespWithTimeout(channel string, messages []*ApplicationMessage, timeoutMs int64) (*ResponsePayload, error) {
	var payload RequestPayload
	var buf []byte
	var err error
	payload.Channel = this.translateChannel(channel)
	payload.Messages = messages
	payload.NeedResp = proto.Int32(1)
	buf, err = proto.Marshal(&payload)
	defer func() {
		if err == nil {
			recordSuccess(channel, len(messages), false)
		} else {
			recordFail(channel, len(messages), true)
		}
	}()
	if err != nil {
		return nil, err
	}
	if this.net == "unixpacket" {
		if len(buf) > PACKET_SIZE_LIMIT {
			err = ErrSeqpacketTooLarge
			return nil, err
		}
	} else if len(buf) > MAX_STREAM_READ_BUFFER_SIZE {
		err = ErrStreamTooLarge
		return nil, err
	}
	startTime := getNowTimeMs()
	nowTime := startTime
	var resp *ResponsePayload
	for (nowTime - startTime) <= timeoutMs {
		resp, err = this.sendWithResp(buf, channel)
		if err == nil {
			return resp, err
		}
		if timeoutMs > 0 {
			time.Sleep(time.Duration(10) * time.Millisecond)
		} else {
			break
		}
		nowTime = getNowTimeMs()
	}
	return nil, err
}

func (this *DatabusCollector) sendWithResp(buf []byte, channel string) (*ResponsePayload, error) {
	var err error
	var con *conn
	for i := 0; i < cap(this.pool)+1; i++ {
		con, err = this.borrow()
		if err != nil {
			return nil, err
		}
		err = con.Write(buf)
		if err == nil {
			var resp *ResponsePayload
			resp, err = this.recvResp(con)
			this.putBack(con)
			if err == nil {
				err = this.checkAgentError(resp)
				return resp, err
			} else {
				return nil, err
			}
		}
		this.putBack(con)
		recordRetry(channel, 1, false)
		this.logInterval(10000, "channel:%s Retry: %d/%d for err: %s\n", channel, i,
			cap(this.pool), err.Error())
	}
	return nil, err
}

func (this *DatabusCollector) checkAgentError(resp *ResponsePayload) error {
	var returnError error
	switch resp.GetCode() {
	case RESP_SUCC_CODE:
		returnError = nil
		break
	case RESP_UNKNOWN_CHANNEL_CODE:
		returnError = ErrUnknownChannel
		break
	case RESP_BUFFER_FULL_CODE:
		returnError = ErrAgentBufferFull
		break
	default:
		returnError = errors.New("agent return unknown error " + string(resp.GetCode()))
		break
	}
	return returnError
}

func (this *DatabusCollector) Collect(channel string, value []byte, key []byte, codec int32) error {
	var message ApplicationMessage
	message.Key = key
	message.Value = value
	message.Codec = &codec
	return this.CollectArray(channel, []*ApplicationMessage{&message})
}

func (this *DatabusCollector) CollectWithTimeout(channel string, value []byte, key []byte, codec int32, timeout_ms int64) error {
	var message ApplicationMessage
	message.Key = key
	message.Value = value
	message.Codec = &codec
	return this.CollectArrayWithTimeout(channel, []*ApplicationMessage{&message}, timeout_ms)
}

func (this *DatabusCollector) CollectWithResp(channel string, value []byte, key []byte, codec int32) (*ResponsePayload, error) {
	var message ApplicationMessage
	message.Key = key
	message.Value = value
	message.Codec = &codec
	return this.CollectArrayWithResp(channel, []*ApplicationMessage{&message})
}

func (this *DatabusCollector) CollectWithRespWithTimeout(channel string, value []byte, key []byte, codec int32, timeout_ms int64) (*ResponsePayload, error) {
	var message ApplicationMessage
	message.Key = key
	message.Value = value
	message.Codec = &codec
	return this.CollectArrayWithRespWithTimeout(channel, []*ApplicationMessage{&message}, timeout_ms)
}

// 内部实现，不要调用
func (this *DatabusCollector) SendProto(b []byte) error {
	con, err := this.borrow()
	if err != nil {
		return err
	}
	err = con.Write(b)
	this.putBack(con)
	if err != nil {
		return err
	}
	return nil
}

// 内部实现，不要调用
func (this *DatabusCollector) SendProtoRecvResp(b []byte) (*ResponsePayload, error) {
	con, err := this.borrow()
	if err != nil {
		return nil, err
	}
	defer this.putBack(con)
	err = con.Write(b)
	if err != nil {
		return nil, err
	}
	resp, err := this.recvResp(con)
	return resp, err
}

func (this *DatabusCollector) recvResp(con *conn) (resp *ResponsePayload, err error) {
	var header []byte
	var body []byte
	var receivedSize int
	var toReceiveSize int
	var bodySize uint32

	if this.net == "unixpacket" {
		body = make([]byte, READ_BUFFER_SIZE)
		receivedSize, err = con.Read(body)
		if err != nil {
			return nil, err
		}
		resp, err = this.extractResp(body, receivedSize)
		return resp, err
	}
	// read of stream collector
	// read header
	toReceiveSize = STREAM_HEADER_SIZE
	buf := make([]byte, toReceiveSize)
	for toReceiveSize > 0 {
		receivedSize, err = con.Read(buf)
		if err == nil {
			toReceiveSize -= receivedSize
			header = append(header, buf[:receivedSize]...)
		} else {
			return nil, err
		}
	}
	if uint8(header[0]) != KVERSION {
		err = errors.New("received stream response has incorrect KVERSION. client's is " +
			strconv.FormatInt(KVERSION, 10) + ", agent's is " + strconv.FormatUint(uint64(header[0]), 10))
		return nil, err
	}
	bodySize = binary.BigEndian.Uint32(header[1:])
	if bodySize > MAX_STREAM_READ_BUFFER_SIZE {
		err = errors.New("received stream response's body is too large")
		return nil, err
	}
	// read body
	toReceiveSize = int(bodySize)
	for toReceiveSize > 0 {
		buf := make([]byte, toReceiveSize)
		receivedSize, err = con.Read(buf)
		if err == nil {
			toReceiveSize -= receivedSize
			body = append(body, buf[:receivedSize]...)
		} else {
			return nil, err
		}
	}
	resp, err = this.extractResp(body, int(bodySize))
	return resp, err
}

func (this *DatabusCollector) extractResp(b []byte, n int) (*ResponsePayload, error) {
	resp := &ResponsePayload{}
	err := proto.Unmarshal(b[0:n], resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (this *DatabusCollector) Close() error {
	//连接等待gc关闭, close queue 会触发panic
	log.Printf("databus collector stop")
	this.cache.cond.Interrupt()
	this.consumerStop = true
	this.flush()
	return nil
}

type CacheData struct {
	payload RequestPayload
	time    int64
	length  int64
}

func (cacheData *CacheData) getCacheDataSize() int64 {
	return cacheData.length + 8 + 8
}

type Cache struct {
	list         *list.List
	cacheSize    int64
	cacheUsed    int64
	maxCacheTime int64
	cond         *Condition
}

func (this *DatabusCollector) SetCacheSize(size int64) {
	this.cache.cacheSize = size
}

func (this *DatabusCollector) SetCacheTimeMs(timeMs int64) {
	this.cache.maxCacheTime = timeMs
}

func (this *DatabusCollector) push(payload RequestPayload, length int64, timeoutMs int64) error {
	cacheData := &CacheData{payload: payload, length: length, time: getNowTimeMs()}
	messageNum := len(payload.Messages)
	defer func() {
		this.cache.cond.lock.Unlock()
	}()

	this.cache.cond.lock.Lock()

	if this.cache.list.Len() > 0 && this.cache.maxCacheTime > 0 &&
		(getNowTimeMs()-this.cache.list.Front().Value.(*CacheData).time > this.cache.maxCacheTime) {
		recordTimeExpired(payload.GetChannel(), messageNum)
		this.logInterval(10000, "channel:%s push cache fail, max_cache_time:%d, now:%d, "+
			"oldest_time:%d\n", payload.GetChannel(), this.cache.maxCacheTime, getNowTimeMs(),
			this.cache.list.Front().Value.(*CacheData).time)
		return ErrCacheTimeExpired
	}

	if this.cache.cacheSize > (this.cache.cacheUsed + cacheData.getCacheDataSize()) {
		this.cache.cacheUsed += cacheData.getCacheDataSize()
		this.cache.list.PushBack(cacheData)
		recordCacheSuccess(payload.GetChannel(), messageNum)
		return nil
	}
	if timeoutMs <= 0 {
		recordCacheFail(payload.GetChannel(), messageNum)
		this.logInterval(10000, "channel:%s push cache fail, cache full. cache_size:%d, "+
			"cache_used:%d, message_length:%d", payload.GetChannel(), this.cache.cacheSize, this.cache.cacheUsed,
			cacheData.getCacheDataSize())
		return ErrCacheFull
	}

	startTime := getNowTimeMs()
	timeLeft := timeoutMs
	for timeLeft > 0 && (this.cache.cacheSize < (this.cache.cacheUsed + cacheData.getCacheDataSize())) {
		is_notify := this.cache.cond.AwaitWithTimeOut(timeLeft)
		if is_notify {
			timeLeft = timeoutMs - (getNowTimeMs() - startTime)
			continue
		} else {
			recordCacheFail(payload.GetChannel(), messageNum)
			this.logInterval(10000, "channel:%s push cache fail, because timeout. "+
				"start_time:%d, now:%d, timeout:%d, cache_size:%d, cache_used:%d", payload.GetChannel(), startTime,
				getNowTimeMs(), timeoutMs, this.cache.cacheSize, this.cache.cacheUsed)
			return ErrCacheFull
		}
	}
	if timeLeft <= 0 {
		recordCacheFail(payload.GetChannel(), messageNum)
		this.logInterval(10000, "channel:%s push cache fail, because timeout. start_time:%d,"+
			"now:%d, timeout:%d, cache_size:%d, cache_used:%d", payload.GetChannel(), startTime, getNowTimeMs(),
			timeoutMs, this.cache.cacheSize, this.cache.cacheUsed)
		return ErrCacheFull
	}
	this.cache.cacheUsed += cacheData.getCacheDataSize()
	this.cache.list.PushBack(cacheData)
	recordCacheSuccess(payload.GetChannel(), messageNum)
	return nil
}

func (this *DatabusCollector) flush() bool {
	defer func() {
		this.cache.cond.lock.Unlock()
	}()

	ok := true
	log.Printf("begin flush, len:%d", this.cache.list.Len())
	this.cache.cond.lock.Lock()
	for e := this.cache.list.Front(); e != nil; e = e.Next() {
		cacheData := e.Value.(*CacheData)
		if cacheData != nil {
			buf, err := proto.Marshal(&cacheData.payload)
			if err != nil {
				this.logInterval(1000, "flush Marshal err:%s", err.Error())
				continue
			}
			err = this.SendProto(buf)
			if err != nil {
				this.logInterval(1000, "channel:%s flush send err:%s",
					cacheData.payload.GetChannel(), err.Error())
				recordFail(cacheData.payload.GetChannel(), len(cacheData.payload.Messages), false)
				ok = false
			} else {
				recordSuccess(cacheData.payload.GetChannel(), len(cacheData.payload.Messages), true)
			}
		}
	}
	this.cache.cacheUsed = 0
	var next *list.Element
	for e := this.cache.list.Front(); e != nil; e = next {
		next = e.Next()
		this.cache.list.Remove(e)
	}
	log.Printf("flush ret:%t", ok)
	return ok
}

func (this *DatabusCollector) Consumer() {
	log.Println("consumer start")
	for !this.consumerStop {
		should_paused := false
		should_remove := false
		var cacheData *CacheData
		cacheData = nil
		this.cache.cond.lock.Lock()
		if this.cache.list.Len() > 0 && this.cache.list.Front() != nil {
			cacheData = this.cache.list.Front().Value.(*CacheData)
		}
		this.cache.cond.lock.Unlock()

		if cacheData == nil {
			should_paused = true
		} else {
			buf, err := proto.Marshal(&cacheData.payload)
			err = this.SendProto(buf)
			if err == nil {
				should_remove = true
			} else {
				should_paused = true
				this.logInterval(10000, "cache thread send fail, err:%s", err.Error())
			}
		}

		this.cache.cond.lock.Lock()
		if should_remove && this.cache.list.Len() > 0 {
			recordSuccess(cacheData.payload.GetChannel(), len(cacheData.payload.Messages), true)
			this.cache.list.Remove(this.cache.list.Front())
			this.cache.cacheUsed = this.cache.cacheUsed - cacheData.getCacheDataSize()
			this.cache.cond.SignalAll()
		}
		this.cache.cond.lock.Unlock()

		if should_paused && cacheData != nil {
			recordRetry(cacheData.payload.GetChannel(), len(cacheData.payload.Messages), true)
			this.logInterval(100000, "cache thread will send retry")
		}

		if should_paused {
			time.Sleep(time.Duration(10) * time.Millisecond)
			continue
		}
	}
	log.Println("consumer exist")
}
