package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"good_api_test/models"
)

type Cache interface {
	Get(ctx context.Context, key string) (*models.Good, error)
	Set(ctx context.Context, key string, good *models.Good, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{client}
}

func (c *RedisCache) Get(ctx context.Context, key string) (*models.Good, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, redis.Nil
		}
		return nil, err
	}

	good := &models.Good{}
	err = json.Unmarshal([]byte(val), good)
	if err != nil {
		return nil, err
	}

	return good, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, good *models.Good, expiration time.Duration) error {
	val, err := json.Marshal(good)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, val, expiration).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
