package redis

import (
	"errors"

	"github.com/go-redis/redis"
	"github.com/phantom-atom/file-explorer/cache"
	"github.com/phantom-atom/file-explorer/config"
)

type redisCache struct {
	config func() *config.Config
	c      *redis.Client
}

//NewCache 创建redis缓存
func NewCache(configFunc func() *config.Config) (cache.Cache, error) {
	config := configFunc()
	cacheConf := &config.Cache

	if len(cacheConf.Locations) == 0 {
		return nil, errors.New("redis_cache: locations is empty")
	}

	options, err := redis.ParseURL(cacheConf.Locations[0])
	if err != nil {
		return nil, err
	}

	c := redis.NewClient(options)
	return &redisCache{
		config: configFunc,
		c:      c,
	}, nil
}

func (r *redisCache) Set(e ...*cache.Entity) error {
	for _, entity := range e {
		statusCmd := r.c.Set(entity.Key, entity.Value, entity.Expiration)
		if err := statusCmd.Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *redisCache) Get(keys ...string) ([]*cache.Entity, error) {
	entitys := make([]*cache.Entity, len(keys))

	for i, key := range keys {
		stringCmd := r.c.Get(key)
		byt, err := stringCmd.Bytes()
		if err != nil {
			if err == redis.Nil {
				return nil, cache.ErrNotFound
			}
			return nil, err
		}

		durationCmd := r.c.TTL(key)
		expiration, err := durationCmd.Result()
		if err != nil {
			return nil, err
		}
		entitys[i] = &cache.Entity{
			Key:        key,
			Value:      byt,
			Expiration: expiration,
		}
	}
	return entitys, nil
}

func (r *redisCache) Del(keys ...string) error {
	for _, key := range keys {
		if err := r.c.Del(key).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *redisCache) List() ([]*cache.Entity, error) {
	keys, err := r.c.Keys("*").Result()
	if err != nil {
		return nil, err
	}
	return r.Get(keys...)
}
