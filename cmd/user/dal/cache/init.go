package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/redis/go-redis/extra/redisotel/v9"
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

	// 配置 OTel，并给它加上标签
	if err := redisotel.InstrumentTracing(rUser,
		redisotel.WithAttributes(attribute.String("peer.service", "redis-user")),
	); err != nil {
		logger.Fatalf("redis otel instrumentation error: %v", err)
	}

	// Ping 一下，确保连接成功
	if _, err := rUser.Ping(context.Background()).Result(); err != nil {
		logger.Fatalf("redis connect ping failed: %v", err)
	}
}