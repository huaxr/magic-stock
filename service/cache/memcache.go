// @Contact:    huaxinrui
// @Time:       2019/9/16 下午4:19

package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type MemoryCache struct {
	CacheClient
	MaxEntries int
	bucket     *CacheContainer
	objectPool *sync.Pool
}

type CacheContainer struct {
	bucket  map[string]*list.Element
	lruList *list.List
}

type CacheEntry struct {
	Key        string
	Value      interface{}
	Expiration int64
}

func (this *MemoryCache) InitCache() {
	GlobalCache = NewMemoryCache(2000, 2000)
}

func NewMemoryCache(maxEntries int, capacity int) *MemoryCache {
	cache := new(MemoryCache)
	cache.rwMutex = new(sync.RWMutex)
	cache.MaxEntries = maxEntries
	cache.bucket = &CacheContainer{bucket: make(map[string]*list.Element, capacity), lruList: list.New()}
	cache.objectPool = &sync.Pool{
		New: func() interface{} {
			return &CacheEntry{}
		},
	}
	return cache
}

func (this *MemoryCache) pushExpiredCacheEntry(el *list.Element, forcibly bool) bool {
	item, ok := el.Value.(*CacheEntry)
	if time.Now().Unix() > item.Expiration || forcibly {
		this.bucket.lruList.Remove(el)
		if ok {
			delete(this.bucket.bucket, item.Key)
			this.objectPool.Put(&item)
		}
		return true
	}
	this.bucket.lruList.MoveToBack(el)
	return false
}

func (this *MemoryCache) Count() int {
	if this.bucket == nil {
		return 0
	}
	return this.bucket.lruList.Len()
}

func (this *MemoryCache) Contains(key string) bool {
	if this.bucket == nil {
		return false
	}
	this.rwMutex.RLock()
	if el, ok := this.bucket.bucket[key]; ok {
		this.rwMutex.RUnlock()
		this.rwMutex.Lock()
		defer this.rwMutex.Unlock()
		return this.pushExpiredCacheEntry(el, false)
	}
	this.rwMutex.RUnlock()
	return false
}

func (this *MemoryCache) Get(key string) (interface{}, bool) {
	if this.bucket == nil {
		return nil, false
	}
	this.rwMutex.Lock()
	defer this.rwMutex.Unlock()

	if el, ok := this.bucket.bucket[key]; ok {
		if ok = this.pushExpiredCacheEntry(el, false); ok {
			// 已经过期
			return nil, false
		}
		return el.Value.(*CacheEntry).Value, true
	}
	return nil, false
}

func (this *MemoryCache) GetCacheEntry(key string) (*CacheEntry, bool) {
	if this.bucket == nil {
		return &CacheEntry{}, false
	}
	this.rwMutex.Lock()
	defer this.rwMutex.Unlock()

	if el, ok := this.bucket.bucket[key]; ok {
		if ok = this.pushExpiredCacheEntry(el, false); ok {
			return &CacheEntry{}, false
		}
		return el.Value.(*CacheEntry), true
	}
	return &CacheEntry{}, false
}

func (this *MemoryCache) Set(key string, value interface{}, duration time.Duration) error {
	if this.bucket == nil {
		return errors.New("缓存容器没有初始化")
	}
	if len(key) >= 1024 {
		return errors.New("缓存键名长度必须小于1024字节")
	}

	this.rwMutex.Lock()
	defer this.rwMutex.Unlock()

	// 更新
	if el, ok := this.bucket.bucket[key]; ok {
		if item, ok := el.Value.(*CacheEntry); ok {
			item.Value = value
			item.Expiration = time.Now().Add(duration).Unix()
			return nil
		}
	}
	var cacheItem *CacheEntry
	var ok bool

	cacheElement := this.objectPool.Get()

	if cacheItem, ok = cacheElement.(*CacheEntry); ok == false {
		cacheItem = &CacheEntry{}
	}

	cacheItem.Key = key
	cacheItem.Value = value
	cacheItem.Expiration = time.Now().Add(duration).Unix()

	el := this.bucket.lruList.PushFront(cacheItem)
	this.bucket.bucket[key] = el
	if this.bucket.lruList.Len() > this.MaxEntries {
		temp := this.bucket.lruList.Back()
		this.pushExpiredCacheEntry(temp, true)
	}

	return nil
}

func (this *MemoryCache) Del(key string) bool {
	if this.bucket == nil {
		return false
	}
	this.rwMutex.RLock()
	if el, ok := this.bucket.bucket[key]; ok {
		this.rwMutex.RUnlock()
		this.rwMutex.Lock()
		defer this.rwMutex.Unlock()

		delete(this.bucket.bucket, key)
		this.bucket.lruList.Remove(el)
		if item, ok := el.Value.(*CacheEntry); ok {
			this.objectPool.Put(&item)
		}
	}
	return true
}

func (this *MemoryCache) Clear() {
	this.rwMutex.Lock()
	defer this.rwMutex.Unlock()

	this.bucket.lruList = list.New()
	this.bucket.bucket = make(map[string]*list.Element, this.MaxEntries)
}

func InitCache() {
	GlobalCache = &MemoryCache{}
	GlobalCache.InitCache()
}
