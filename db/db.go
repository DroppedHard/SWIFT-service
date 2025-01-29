package db

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisStorage(cfg *redis.Options) *redis.Client {
	db := redis.NewClient(cfg)
	return db
}

func TestClientConection(client *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	key, val := "foo", "bar"
	defer cancel()
	
	err := client.Set(ctx, key, val, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Get(ctx, key).Err()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Del(ctx, key).Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Succesfully connected!")
}
