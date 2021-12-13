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

func (cache *ttlCacheSCache) set(key string, data string) {
	cache.mutex.Lock()
	item := &ttlCache{data: data}
	item.touch(cache.ttl)
	cache.items[key] = item
	cache.mutex.Unlock()
}

func (cache *ttlCacheSCache) get(key string) (data string, found bool) {
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

func (cache *ttlCacheSCache) count() int {
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

func (m *ttlCacheStruct) Set(key string, value string) {
	m.cache.set(key, value)
}

func (m *ttlCacheStruct) Get(key string) string {
	value, exists := m.cache.get(key)
	if exists != true {
		Panicerr("Key " + key + " not found in cache")
	}
	return value
}

func (m *ttlCacheStruct) Exists(key string) bool {
	_, exists := m.cache.get(key)
	return exists
}

func (m *ttlCacheStruct) Count() int {
	return m.cache.count()
}
