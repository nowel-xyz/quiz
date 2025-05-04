package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var Ctx = context.Background()

func InitRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"), // e.g., "localhost:6379"
		Password: os.Getenv("REDIS_PASSWORD"), // "" if none
		DB:       0,
	})

	_, err := Redis.Ping(Ctx).Result()
	if err != nil {
		return err
	}

	

	log.Println("✅ Connected to Redis")
	return nil
}
