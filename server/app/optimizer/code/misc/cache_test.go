package misc

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := New[string, int]()
	cache.Set("one", 1, 1*time.Second)
	cache.Set("two", 2, 2*time.Second)
	cache.Set("three", 1, 1*time.Second)

	value, found := cache.Get("two")
	assert.True(t, found)
	assert.Equal(t, value, 2)

	value, found = cache.Pop("two")
	assert.True(t, found)
	assert.Equal(t, value, 2)

	time.Sleep(3 * time.Second)

	_, found = cache.Get("one")
	assert.False(t, found)
	_, found = cache.Get("two")
	assert.False(t, found)
	_, found = cache.Get("three")
	assert.False(t, found)
}
