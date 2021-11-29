package cache

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func NewClient() *redis.Client {
	redisPort := os.Getenv("REDIS_PORT")
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%v", redisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
