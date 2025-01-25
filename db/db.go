package db

import (
	"github.com/redis/go-redis/v9"
)

func NewRedisStorage(cfg *redis.Options) *redis.Client {
	db := redis.NewClient(cfg)
	return db
}
