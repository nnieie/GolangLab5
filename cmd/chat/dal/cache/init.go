package cache

import (
	"github.com/redis/go-redis/v9"

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"go.opentelemetry.io/otel/attribute"
	"context"
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

	if err := redisotel.InstrumentTracing(rChat,
		redisotel.WithAttributes(attribute.String("peer.service", "redis-chat")),
	); err != nil {
		logger.Fatalf("redis otel instrumentation error: %v", err)
	}

	if _, err := rChat.Ping(context.Background()).Result(); err != nil {
		logger.Fatalf("redis connect ping failed: %v", err)
	}
}
