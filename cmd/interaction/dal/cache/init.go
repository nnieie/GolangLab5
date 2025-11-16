package cache

import (
	"github.com/redis/go-redis/v9"

	"github.com/nnieie/golanglab5/config"
)

var (
	rInteraction *redis.Client
)

func InitRedis() {
	rInteraction = redis.NewClient(&redis.Options{
		Addr:       config.Redis.Addr,
		Password:   config.Redis.Password,
		ClientName: "Interaction",
		DB:         2,
	})
}
