package services

import (
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// CacheService provides in-memory caching functionality
type CacheService struct {
	cache *gocache.Cache
	mutex sync.RWMutex
}

// NewCacheService creates a new cache service
func NewCacheService(defaultExpiration, cleanupInterval time.Duration) *CacheService {
	return &CacheService{
		cache: gocache.New(defaultExpiration, cleanupInterval),
	}
}

// Set stores a value in the cache
func (c *CacheService) Set(key string, value interface{}, duration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache.Set(key, value, duration)
}

// Get retrieves a value from the cache
func (c *CacheService) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.cache.Get(key)
}

// Delete removes a value from the cache
func (c *CacheService) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache.Delete(key)
}

// Flush clears all items from the cache
func (c *CacheService) Flush() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache.Flush()
}

// ItemCount returns the number of items in the cache
func (c *CacheService) ItemCount() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.cache.ItemCount()
}

// SetDefault stores a value with the default expiration
func (c *CacheService) SetDefault(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache.SetDefault(key, value)
}

// GetOrSet gets a value from cache or sets it if not present
func (c *CacheService) GetOrSet(key string, fetch func() (interface{}, error), duration time.Duration) (interface{}, error) {
	// First try to get from cache
	if value, found := c.Get(key); found {
		return value, nil
	}

	// Not in cache, fetch and store
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Double-check after acquiring write lock
	if value, found := c.cache.Get(key); found {
		return value, nil
	}

	value, err := fetch()
	if err != nil {
		return nil, err
	}

	c.cache.Set(key, value, duration)
	return value, nil
}

// Keys returns all cache keys
func (c *CacheService) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	items := c.cache.Items()
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	return keys
}

// CacheStats holds cache statistics
type CacheStats struct {
	ItemCount int `json:"itemCount"`
	HitRate   float64 `json:"hitRate"`
}

// GetStats returns cache statistics
func (c *CacheService) GetStats() CacheStats {
	return CacheStats{
		ItemCount: c.ItemCount(),
	}
}
