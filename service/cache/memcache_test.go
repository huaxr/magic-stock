// @Contact:    huaxinrui
// @Time:       2019/9/16 下午4:20

package cache

import (
	"testing"
	"time"
)

var cacheObject = NewMemoryCache(1e4, 1e4)

type cacheTest struct {
	key      string
	value    interface{}
	duration time.Duration
}

type cacheDataValue struct {
	id   int
	name string
}

var cacheTestValue = []cacheTest{
	cacheTest{key: "cache1", value: cacheDataValue{id: 1, name: "linux"}, duration: time.Minute * 2},
	cacheTest{key: "cache2", value: cacheDataValue{id: 2, name: "mac"}, duration: time.Minute * 2},
	cacheTest{key: "cache2", value: cacheDataValue{id: 3, name: "win"}, duration: time.Minute * 2},
}

func TestSet(t *testing.T) {
	for _, v := range cacheTestValue {
		err := cacheObject.Set(v.key, v.value, v.duration)
		if err != nil {
			t.Errorf("key:%s value:%v duration:%v", v.key, v.value, v.duration)
		}
	}
}

func TestGet(t *testing.T) {
	for _, v := range cacheTestValue {
		err := cacheObject.Set(v.key, v.value, v.duration)
		if err != nil {
			t.Errorf("key:%s value:%v duration:%v", v.key, v.value, v.duration)
		}
	}
	for _, v := range cacheTestValue {
		value, ok := cacheObject.Get(v.key)
		if !ok {
			t.Errorf("key:%s value:%v duration:%v", v.key, v.value, v.duration)
		} else {
			value := value.(cacheDataValue)

			t.Logf("id=%v name=%s", value.id, value.name)
		}

	}
}
