package infrastructure

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",                   
		DB:       0,     
	})

	// 接続確認
	_, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis: %v", err)
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return client
}