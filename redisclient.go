package main

import (
	"context"

	"github.com/go-redis/redis"
)

type redisClient struct {
	client *redis.Client
	ctx    context.Context
}

func (c redisClient) set(key string, value string) error {
	rdb := c.client
	rdbErr := rdb.Set(c.ctx, key, value, 0).Err()
	return rdbErr
}

func (c redisClient) get(key string) (string, error) {
	rdb := c.client
	val, err := rdb.Get(c.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
