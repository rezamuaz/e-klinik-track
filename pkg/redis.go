package pkg

import (
	"context"
	"e-klinik/config"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// --- 1. Interface and Struct Definitions ---

// RedisService defines a reusable interface for Redis operations.
// All methods that perform I/O (Set, Get, Delete) must accept a context.Context.
type RedisService interface {
	Set(ctx context.Context, key string, value any) bool
	SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) bool
	Get(ctx context.Context, key string, dest any) bool
	Delete(ctx context.Context, key string)
	Close() error
	Ping(ctx context.Context) bool
}

// RedisCache implements RedisService using Redis as backend.
type RedisCache struct {
	Client *redis.Client
	// REMOVED: ctx context.Context - The context should come from the request.
}

// --- 2. Constructor ---

// NewRedisCache initializes a new RedisCache service.
// NOTE: Ping still uses context.Background() as it's a startup check.
func NewRedisCache(cfg *config.Config) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})

	// Use context.Background() only for the initial connection health check.
	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("[Redis] ⚠️  Failed to connect to %s:%s — %v", cfg.Redis.Host, cfg.Redis.Port, err)
	} else {
		log.Printf("[Redis] ✅ Connected to %s:%s", cfg.Redis.Host, cfg.Redis.Port)
	}

	return &RedisCache{
		Client: client,
	}
}

// --- 3. Fixed Methods (Accepting Context) ---

// Set stores a key-value pair with no TTL.
func (r *RedisCache) Set(ctx context.Context, key string, value any) bool {
	return r.SetWithTTL(ctx, key, value, 0)
}

// SetWithTTL stores a key-value pair with an expiration duration.
func (r *RedisCache) SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) bool {
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("[Redis] ❌ Failed to marshal value for key %s: %v", key, err)
		return false
	}

	// Use the provided context (ctx)
	if err := r.Client.Set(ctx, key, data, ttl).Err(); err != nil {
		log.Printf("[Redis] ❌ Failed to set key %s: %v", key, err)
		return false
	}
	return true
}

// Get retrieves a value from Redis and unmarshals it into dest.
func (r *RedisCache) Get(ctx context.Context, key string, dest any) bool {
	// Use the provided context (ctx)
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Printf("[Redis] ⚠️  Get failed for key %s: %v", key, err)
		}
		return false
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		log.Printf("[Redis] ⚠️  Unmarshal failed for key %s: %v", key, err)
		return false
	}
	return true
}

func (r *RedisCache) GetRaw(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Printf("[Redis] ⚠️  GetRaw failed for key %s: %v", key, err)
		}
		// Kembalikan error, termasuk redis.Nil jika data tidak ditemukan
		return "", err
	}
	// Operasi berhasil
	return val, nil
}

// Delete removes a key from Redis.
func (r *RedisCache) Delete(ctx context.Context, key string) {
	// Use the provided context (ctx)
	if err := r.Client.Del(ctx, key).Err(); err != nil {
		log.Printf("[Redis] ⚠️  Failed to delete key %s: %v", key, err)
	}
}

// Ping checks Redis connectivity.
func (r *RedisCache) Ping(ctx context.Context) bool {
	// Use the provided context (ctx)
	if _, err := r.Client.Ping(ctx).Result(); err != nil {
		log.Printf("[Redis] ⚠️  Ping failed: %v", err)
		return false
	}
	return true
}

// Close gracefully closes the Redis client connection.
func (r *RedisCache) Close() error {
	return r.Client.Close()
}
func (r *RedisCache) GetClient() *redis.Client {
	return r.Client
}
