// @Contact:    huaxinrui
// @Time:       2019/9/19 下午8:21

package cache

import (
	"log"
	"sync"
	"time"
)

type CacheIF interface {
	InitCache()
	Set(key string, value interface{}, duration time.Duration) error
	Get(key string) (interface{}, bool)
	Del(key string) bool
	Count() int
	Clear()
	CacheDecorator(key string, min time.Duration, fc fc) interface{}
}

type fc func() interface{}

type CacheClient struct {
	rwMutex *sync.RWMutex
}

var GlobalCache CacheIF

func (this *MemoryCache) CacheDecorator(key string, min time.Duration, fc fc) interface{} {
	res, success := this.Get(key)
	if success {
		log.Println("cache found")
		return res
	}

	result := fc()
	err := this.Set(key, result, min*time.Minute)
	if err != nil {
		log.Println(err)
	}
	return result
}
