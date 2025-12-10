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
	rInteraction *redis.Client
)

func InitRedis() {
	rInteraction = redis.NewClient(&redis.Options{
		Addr:       config.Redis.Addr,
		Password:   config.Redis.Password,
		ClientName: "Interaction",
		DB:         2,
	})

	if err := redisotel.InstrumentTracing(rInteraction,
		redisotel.WithAttributes(attribute.String("peer.service", "redis-interaction")),
	); err != nil {
		logger.Fatalf("redis otel instrumentation error: %v", err)
	}

	if _, err := rInteraction.Ping(context.Background()).Result(); err != nil {
		logger.Fatalf("redis connect ping failed: %v", err)
	}
}
