// @Contact:    huaxinrui
// @Time:       2019/9/25 下午4:03

package engine

import (
	"sync"
)

type EngineIF interface {
	TermQuery(query map[string]string, callback func(err error)) []string
	loopConnectUntilSuccess()
}

type EngineClient struct {
	typ   string // 存储类型
	mutex *sync.Mutex
}
