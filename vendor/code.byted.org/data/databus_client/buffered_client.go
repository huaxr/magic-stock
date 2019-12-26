package databus_client

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	DEFAULT_QUEUE_SIZE   = 1000
	DEFAULT_CONSUMER_NUM = 1
)

var ErrBufferFull = errors.New("buffer is full")

type BufferedDatabusCollector struct {
	queued       chan []byte
	socketPath   string
	writeTimeout time.Duration
	maxConn      int
	net          string
}

func NewDefaultBufferedDatabusCollector() *BufferedDatabusCollector {
	return NewBufferedDatabusCollector(DEFAULT_SOCKET_PATH, DEFAULT_TIMEOUT, DEFAULT_MAX_CONN_NUM, DEFAULT_QUEUE_SIZE, DEFAULT_CONSUMER_NUM)
}

func NewBufferedDatabusCollector(socketPath string, timeout time.Duration, maxConn int, queueSize int, consumerNum int) *BufferedDatabusCollector {
	return NewBufferedDatabusCollectorV1(socketPath, timeout, maxConn, queueSize, consumerNum, "unixpacket")
}

func NewStreamBufferedDatabusCollector() *BufferedDatabusCollector {
	return NewBufferedDatabusCollectorV1(DEFAULT_SOCKET_PATH, DEFAULT_TIMEOUT, DEFAULT_MAX_CONN_NUM, DEFAULT_QUEUE_SIZE, DEFAULT_CONSUMER_NUM, "unix")
}

func NewBufferedDatabusCollectorV1(socketPath string, timeout time.Duration, maxConn int, queueSize int, consumerNum int, net string) *BufferedDatabusCollector {
	bc := &BufferedDatabusCollector{
		queued:       make(chan []byte, queueSize),
		socketPath:   socketPath,
		writeTimeout: timeout,
		maxConn:      maxConn,
		net:          net,
	}
	for i := 0; i < consumerNum; i++ {
		go bc.consumer()
	}
	return bc
}

func (bc *BufferedDatabusCollector) consumer() {
	collector := NewCollectorV1(bc.socketPath, bc.writeTimeout, 2, bc.net)
	for {
		select {
		case msg, ok := <-bc.queued:
			if !ok {
				break
			}
			collector.SendProto(msg)
		}
	}
}

func (bc *BufferedDatabusCollector) CollectArray(channel string, messages []*ApplicationMessage) error {
	var payload RequestPayload
	payload.Channel = &channel
	payload.Messages = messages
	buf, err := proto.Marshal(&payload)
	messageNum := len(messages)
	defer func() {
		// 和之前语义保持一致，bufferclient只要放入队列就认为成功
		if err == nil {
			recordSuccess(channel, messageNum, false)
		} else {
			recordFail(channel, messageNum, false)
		}
	}()
	if err != nil {
		return err
	}
	if bc.net == "unixpacket" {
		if len(buf) > PACKET_SIZE_LIMIT {
			err = fmt.Errorf("CollectArray Err: Packet too large.")
			return err
		}
	} else if len(buf) > MAX_STREAM_READ_BUFFER_SIZE {
		err = fmt.Errorf("CollectArray Err: Stream message too large.")
		return err
	}
	select {
	case bc.queued <- buf:
		return nil
	default:
		err = ErrBufferFull
		return err
	}
}
func (bc *BufferedDatabusCollector) Collect(channel string, value []byte, key []byte, codec int32) error {
	var message ApplicationMessage
	message.Key = key
	message.Value = value
	message.Codec = &codec
	var payload RequestPayload
	payload.Channel = &channel
	payload.Messages = []*ApplicationMessage{&message}
	buf, err := proto.Marshal(&payload)
	messageNum := len(payload.Messages)
	defer func() {
		if err == nil {
			recordSuccess(channel, messageNum, false)
		} else {
			recordFail(channel, messageNum, false)
		}
	}()
	if err != nil {
		return err
	}

	select {
	case bc.queued <- buf:
		return nil
	default:
		err = ErrBufferFull
		return err
	}
}
func (bc *BufferedDatabusCollector) Close() error {
	//close(bc.queued)
	return nil
}
