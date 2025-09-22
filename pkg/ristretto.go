package pkg

import (
	"log"
	"time"

	"github.com/dgraph-io/ristretto"
)

// CacheService defines the public API for our cache.
type CacheService interface {
	Set(key, value any) bool
	SetWithTTL(key, value any, ttl time.Duration) bool
	Get(key any) (any, bool)
	Delete(key any)
	Close()
}

// RistrettoCache is an in-memory cache service.
type RistrettoCache struct {
	cache *ristretto.Cache
}

// NewRistrettoCache creates and configures a new Ristretto cache service.
func NewRistrettoCache() *RistrettoCache {
	// A good starting configuration.
	// NumCounters: 10x the number of items you expect to keep in the cache.
	// MaxCost: The maximum total cost (e.g., memory) the cache can hold.
	// BufferItems: Recommended to be 64 for most use cases.
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // 10M counters
		MaxCost:     1 << 30, // 1GB max cost
		BufferItems: 64,
	})
	if err != nil {
		log.Fatalf("failed to create ristretto cache: %v", err)
	}

	return &RistrettoCache{
		cache: cache,
	}
}

// Set stores a key-value item with a default cost of 1.
func (r *RistrettoCache) Set(key, value any) bool {
	return r.cache.Set(key, value, 1)
}

// SetWithTTL stores a key-value item with a TTL.
func (r *RistrettoCache) SetWithTTL(key, value any, ttl time.Duration) bool {
	return r.cache.SetWithTTL(key, value, 1, ttl)
}

// Get retrieves a value from the cache.
func (r *RistrettoCache) Get(key any) (any, bool) {
	return r.cache.Get(key)
}

// Delete removes a key-value item from the cache.
func (r *RistrettoCache) Delete(key any) {
	r.cache.Del(key)
}

// Close gracefully stops the cache's background goroutines.
func (r *RistrettoCache) Close() {
	r.cache.Close()
}
