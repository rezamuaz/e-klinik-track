package pkg

import (
	"context"
	"e-klinik/config"
	"encoding/json"
	"fmt"

	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) (*redis.Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       0,
		// DialTimeout:  cfg.Redis.DialTimeout * time.Second,
		// ReadTimeout:  cfg.Redis.ReadTimeout * time.Second,
		// WriteTimeout: cfg.Redis.WriteTimeout * time.Second,
		// PoolSize:     cfg.Redis.PoolSize,
		// PoolTimeout:  cfg.Redis.PoolTimeout,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, WrapErrorf(err, ErrorCodeUnknown, "rdb.Ping")
	}

	return rdb, nil
}

func Set[T any](ctx context.Context, c *redis.Client, key string, value T, duration time.Duration) error {
	ct, cancel := context.WithTimeout(ctx, 30)
	defer cancel()
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ct, key, v, duration).Err()
}

func Get[T any](ctx context.Context, c *redis.Client, key string) (T, error) {
	ct, cancel := context.WithTimeout(ctx, 30)
	defer cancel()
	var dest T = *new(T)
	v, err := c.Get(ct, key).Result()
	if err != nil {
		return dest, err
	}
	err = json.Unmarshal([]byte(v), &dest)
	if err != nil {
		return dest, err
	}
	return dest, nil
}

func Del[T any](ctx context.Context, c *redis.Client, key string) error {
	ct, cancel := context.WithTimeout(ctx, 30)
	defer cancel()

	return c.Del(ct, key).Err()
}
