package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type (
	CacheManager interface {
		Get(ctx context.Context, key string) *string
		Set(ctx context.Context, key string, val any) error
		Del(ctx context.Context, key string) error
	}

	cacheManager struct {
		redis         *redis.Client
		hitCounter    prometheus.Counter
		missedCounter prometheus.Counter
	}
)

func NewCacheManager(options *redis.Options) CacheManager {
	redisClient := redis.NewClient(options)

	hitCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_hit",
		Help: "Cache Hit",
	})
	missedCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "cache_missed",
		Help: "Cache Missed",
	})

	prometheus.MustRegister(hitCounter)
	prometheus.MustRegister(missedCounter)

	return cacheManager{
		redis:         redisClient,
		hitCounter:    hitCounter,
		missedCounter: missedCounter,
	}
}

func (c cacheManager) Get(ctx context.Context, key string) *string {
	val, err := c.redis.Get(ctx, key).Result()
	c.hitCounter.Inc()

	if err == redis.Nil {
		c.missedCounter.Inc()
		return nil
	}

	return &val
}

func (c cacheManager) Set(ctx context.Context, key string, val any) error {
	typeOf := reflect.TypeOf(val)
	kind := typeOf.Kind()
	if kind == reflect.Struct || (kind == reflect.Slice && typeOf.Elem().Kind() == reflect.Struct) {
		result, err := json.Marshal(val)
		if err != nil {
			fmt.Println("Error marshaling to JSON:", err)
		}
		val = string(result)
	}

	return c.redis.Set(ctx, key, val, 24*time.Hour).Err()
}

func (c cacheManager) Del(ctx context.Context, key string) error {
	return c.redis.Del(ctx, key).Err()
}

func GetValue[T any](ctx context.Context, cm CacheManager, key string) *T {
	val := cm.Get(ctx, key)
	if val == nil || *val == "" {
		return nil
	}

	var result T
	err := json.Unmarshal([]byte(*val), &result)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &result
}
