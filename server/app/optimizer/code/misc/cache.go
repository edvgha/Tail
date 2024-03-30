package misc

import (
	"sync"
	"time"
)

type item[V any] struct {
	value          V
	expiry         time.Time
	expireCallBack func(V) bool
}

func (i item[V]) isExpired() bool {
	return time.Now().After(i.expiry)
}

type Cache[K comparable, V any] struct {
	items map[K]item[V]
	mu    sync.Mutex
}

func New[K comparable, V any]() *Cache[K, V] {
	c := &Cache[K, V]{
		items: make(map[K]item[V]),
	}

	go func() {
		for range time.Tick(1 * time.Second) {
			c.mu.Lock()

			for key, item := range c.items {
				if item.isExpired() {
					item.expireCallBack(item.value)
					delete(c.items, key)
				}
			}
			c.mu.Unlock()
		}
	}()

	return c
}

func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration, expireCallBack func(V) bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item[V]{
		value:          value,
		expiry:         time.Now().Add(ttl),
		expireCallBack: expireCallBack,
	}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return item.value, false
	}

	if item.isExpired() {
		delete(c.items, key)
		return item.value, false
	}

	return item.value, true
}

func (c *Cache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *Cache[K, V]) Pop(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return item.value, false
	}

	delete(c.items, key)

	if item.isExpired() {
		return item.value, false
	}

	return item.value, true
}
