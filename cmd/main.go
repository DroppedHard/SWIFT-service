package main

import (
	"fmt"
	"log"

	"github.com/DroppedHard/SWIFT-service/cmd/api"
	"github.com/DroppedHard/SWIFT-service/config"
	"github.com/DroppedHard/SWIFT-service/db"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := db.NewRedisStorage(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", config.Envs.DBHost, config.Envs.DBPort),
		Password:     config.Envs.DBPassword,
		DB:           config.Envs.DBNum,
		PoolSize:     config.Envs.DBPoolSize,
		MinIdleConns: config.Envs.DBMinIdleConns,
	})

	db.TestClientConection(rdb)

	server := api.NewAPIServer(config.Envs.PublicHost, config.Envs.Port, rdb)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
