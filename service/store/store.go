package store

import "github.com/redis/go-redis/v9"

func NewStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: *client}
}
