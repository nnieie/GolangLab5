package cache

import (
	"github.com/redis/go-redis/v9"

	"github.com/nnieie/golanglab5/config"
)

var (
	rChat *redis.Client
)

func InitRedis() {
	rChat = redis.NewClient(&redis.Options{
		Addr:       config.Redis.Addr,
		Password:   config.Redis.Password,
		ClientName: "Chat",
		DB:         1,
	})
}
