package golanglibs

import (
	"sync"
	"time"
)

type ttlCache struct {
	sync.RWMutex
	data    string
	expires *time.Time
}

func (item *ttlCache) touch(duration time.Duration) {
	item.Lock()
	expiration := time.Now().Add(duration)
	item.expires = &expiration
	item.Unlock()
}

func (item *ttlCache) expired() bool {
	var value bool
	item.RLock()
	if item.expires == nil {
		value = true
	} else {
		value = item.expires.Before(time.Now())
	}
	item.RUnlock()
	return value
}

type ttlCacheSCache struct {
	mutex sync.RWMutex
	ttl   time.Duration
	items map[string]*ttlCache
}

func (cache *ttlCacheSCache) Set(key string, data string) {
	cache.mutex.Lock()
	item := &ttlCache{data: data}
	item.touch(cache.ttl)
	cache.items[key] = item
	cache.mutex.Unlock()
}

func (cache *ttlCacheSCache) Get(key string) (data string, found bool) {
	cache.mutex.Lock()
	item, exists := cache.items[key]
	if !exists || item.expired() {
		data = ""
		found = false
	} else {
		item.touch(cache.ttl)
		data = item.data
		found = true
	}
	cache.mutex.Unlock()
	return
}

func (cache *ttlCacheSCache) Count() int {
	cache.mutex.RLock()
	count := len(cache.items)
	cache.mutex.RUnlock()
	return count
}

func (cache *ttlCacheSCache) cleanup() {
	cache.mutex.Lock()
	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
	cache.mutex.Unlock()
}

func (cache *ttlCacheSCache) startCleanupTimer() {
	ticker := time.Tick(3 * time.Second) // 3秒刷一次缓存
	go (func() {
		for {
			select {
			case <-ticker:
				cache.cleanup()
			}
		}
	})()
}

func NewCache(duration time.Duration) *ttlCacheSCache {
	cache := &ttlCacheSCache{
		ttl:   duration,
		items: map[string]*ttlCache{},
	}
	cache.startCleanupTimer()
	return cache
}

type ttlCacheStruct struct {
	cache *ttlCacheSCache
}

func getTTLCache(ttlsecond interface{}) *ttlCacheStruct {
	return &ttlCacheStruct{
		cache: NewCache(getTimeDuration(ttlsecond)),
	}
}

func (m *ttlCacheStruct) set(key string, value string) {
	m.cache.Set(key, value)
}

func (m *ttlCacheStruct) get(key string) string {
	value, exists := m.cache.Get(key)
	if exists != true {
		panicerr("Key " + key + " not found in cache")
	}
	return value
}

func (m *ttlCacheStruct) exists(key string) bool {
	_, exists := m.cache.Get(key)
	return exists
}

func (m *ttlCacheStruct) count() int {
	return m.cache.Count()
}
