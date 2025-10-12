package cache

import (
	"github.com/redis/go-redis/v9"

	"github.com/nnieie/golanglab5/config"
)

var (
	rUser *redis.Client
)

func InitRedis() {
	rUser = redis.NewClient(&redis.Options{
		Addr:       config.Redis.Addr,
		Password:   config.Redis.Password,
		ClientName: "User",
		DB:         0,
	})
}
