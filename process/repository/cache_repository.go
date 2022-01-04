package process

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type ICacheRepository interface {
	Set(key string, val interface{}, expTime time.Duration) error
	Get(key string) (string, error)
	Delete(key string) (int64, error)
}

type CacheRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewCacheRepository(resource *redis.Client) ICacheRepository {
	cacheRepository := &CacheRepository{client: resource, ctx: context.Background()}
	return cacheRepository
}

func (t *CacheRepository) Set(key string, val interface{}, timeout time.Duration) error {
	err := t.client.Set(t.ctx, key, val, timeout).Err()
	if err != nil {
		return err
	}
	return nil
}

func (t *CacheRepository) Get(key string) (string, error) {
	res, err := t.client.Get(t.ctx, key).Result()
	log.Println("GET INTERFACE ->", t.client)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (t *CacheRepository) Delete(key string) (int64, error) {
	res, err := t.client.Del(t.ctx, key).Result()
	if err != nil {
		return -1, err
	}
	return res, nil
}
