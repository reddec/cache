package cache

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type item struct {
	value   interface{}
	expired bool
	owner   *Cache
	timer   *time.Ticker
}

//Cache of items which can be expired or cancelled (pulled)
type Cache struct {
	generator uint64
	lock      sync.RWMutex
	data      map[string]item
}

//NewCache - create new cache
func NewCache() *Cache {
	return &Cache{data: make(map[string]item)}
}

//Put value to cache with specified time-to-live and return generated unique
// (in terms of one cache) key.
func (cache *Cache) Put(value interface{}, ttl time.Duration) string {
	key := strconv.FormatUint(atomic.AddUint64(&cache.generator, 1), 10)
	it := item{
		value: value,
		owner: cache,
		timer: time.NewTicker(ttl)}
	cache.lock.Lock()
	cache.data[key] = it
	cache.lock.Unlock()
	go func() {
		defer it.timer.Stop()
		<-it.timer.C
		if !it.expired {
			it.owner.Pull(key)
		}
	}()
	return key
}

//Pull key from cache and remove it if it not already expired.
//Returns value and status
func (cache *Cache) Pull(key string) (interface{}, bool) {
	cache.lock.RLock()
	it, ok := cache.data[key]
	cache.lock.RUnlock()
	if !ok {
		return nil, false
	}
	cache.lock.Lock()
	it.expired = true
	it.timer.Stop()
	delete(cache.data, key)
	cache.lock.Unlock()
	return it.value, true
}
