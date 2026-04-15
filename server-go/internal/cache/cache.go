package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

type CacheConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func Init(cfg *CacheConfig) error {
	if cfg == nil || cfg.Host == "" {
		log.Println("[Cache] Redis not configured, running without cache")
		return nil
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       DB,
		PoolSize: 10,
	})

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Printf("[Cache] Failed to connect to Redis at %s: %v", addr, err)
		return err
	}

	log.Printf("[Cache] Connected to Redis at %s", addr)
	return nil
}

func IsEnabled() bool {
	return Client != nil
}

func Get(key string) (string, error) {
	if !IsEnabled() {
		return "", redis.Nil
	}
	return Client.Get(ctx, key).Result()
}

func Set(key string, value interface{}, expiration time.Duration) error {
	if !IsEnabled() {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("[Cache] marshal error for key %s: %w", key, err)
	}
	return Client.Set(ctx, key, data, expiration).Err()
}

func SetString(key string, value string, expiration time.Duration) error {
	if !IsEnabled() {
		return nil
	}
	return Client.Set(ctx, key, value, expiration).Err()
}

func Delete(keys ...string) error {
	if !IsEnabled() {
		return nil
	}
	return Client.Del(ctx, keys...).Err()
}

func DeleteByPattern(pattern string) (int64, error) {
	if !IsEnabled() {
		return 0, nil
	}
	iter := Client.Scan(ctx, 0, pattern, 100).Iterator()
	var count int64
	for iter.Next(ctx) {
		if err := Client.Del(ctx, iter.Val()).Err(); err != nil {
			log.Printf("[Cache] Error deleting key %s: %v", iter.Val(), err)
			continue
		}
		count++
	}
	if err := iter.Err(); err != nil {
		return count, err
	}
	return count, nil
}

func GetJSON(key string, dest interface{}) error {
	if !IsEnabled() {
		return redis.Nil
	}
	data, err := Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func SetJSON(key string, value interface{}, expiration time.Duration) error {
	return Set(key, value, expiration)
}

func InvalidateByPrefix(prefixes ...string) {
	for _, prefix := range prefixes {
		count, err := DeleteByPattern(prefix + "*")
		if err != nil {
			log.Printf("[Cache] Error invalidating prefix %s: %v", prefix, err)
			continue
		}
		if count > 0 {
			log.Printf("[Cache] Invalidated %d keys with prefix %s", count, prefix)
		}
	}
}

const (
	CacheTTLShort  = 5 * time.Minute
	CacheTTLMedium = 15 * time.Minute
	CacheTTLLong   = 1 * time.Hour

	PrefixCategory = "orange:category:"
	PrefixPlatform = "orange:platform:"
	PrefixProduct  = "orange:product:"
	PrefixTag      = "orange:tag:"
	PrefixBlogger  = "orange:blogger:"
	PrefixStat     = "orange:stat:"
)
