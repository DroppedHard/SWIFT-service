package main

import (
	"context"
	"fmt"
	"log"

	"github.com/DroppedHard/SWIFT-service/cmd/api"
	"github.com/DroppedHard/SWIFT-service/config"
	"github.com/DroppedHard/SWIFT-service/db"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := db.NewRedisStorage(&redis.Options{
		Addr:         config.Envs.DBAddress,
		Password:     config.Envs.DBPassword,
		DB:           config.Envs.DBNum,
		PoolSize:     config.Envs.DBPoolSize,
		MinIdleConns: config.Envs.DBMinIdleConns,
	})

	stats := rdb.PoolStats()
	fmt.Printf("Total connections: %d\n", stats.TotalConns)
	fmt.Printf("Idle connections: %d\n", stats.IdleConns)
	fmt.Printf("Active connections: %d\n", stats.StaleConns)
	fmt.Printf("Wait count: %d\n", stats.Timeouts)

	initStorage(rdb)

	server := api.NewAPIServer(config.Envs.Port, rdb)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(client *redis.Client) {
	ctx := context.Background()
	err := client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Get(ctx, "foo").Err()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Del(ctx, "foo").Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Succesfully connected!")
}
