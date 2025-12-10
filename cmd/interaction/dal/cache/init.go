package cache

import (
	"context"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/pkg/logger"
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
