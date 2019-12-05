// @Contact:    huaxinrui
// @Time:       2019/6/10 下午3:56

package core

import (
	"code.byted.org/byte_security/platform_api/core/engine"
	"code.byted.org/byte_security/platform_api/core/store"
	"code.byted.org/byte_security/platform_api/service/conf"
	"code.byted.org/kv/goredis"
)

func init() {
	Backend = new(backend)
	Backend.Store = store.InitStore(cc.Store, false)
	Backend.Engine, _ = engine.InitEngine(cc.ES, false, false)
	//Backend.Redis = o.initRedis()
}

func initRedis() *goredis.Client {
	opt := goredis.NewOption()
	client, _ := goredis.NewClientWithOption(cc.RedisCluster, opt)
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client
}

var cc = &conf.Config
var Backend *backend
