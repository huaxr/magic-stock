// @Contact:    huaxinrui
// @Time:       2019/10/17 下午1:24

package core

import (
	"code.byted.org/byte_security/platform_api/core/engine"
	"code.byted.org/byte_security/platform_api/core/store"
	"code.byted.org/kv/goredis"
	"sync"
)

type backend struct {
	Store       store.StoreIF
	Engine      engine.EngineIF
	Redis       *goredis.Client
	expectation chan error
	lock        *sync.Mutex
}
