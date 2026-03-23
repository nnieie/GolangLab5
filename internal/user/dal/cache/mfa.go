package cache

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
)

func SetTOTPSecret(ctx context.Context, secret string, userID string) error {
	err := rUser.Set(ctx, constants.TOTPSecret+userID, secret, constants.TOTPSecretExpTime).Err()
	if err != nil {
		logger.Errorf("redis set totp secret err: %v", err)
		return err
	}
	return nil
}

func GetTOTPSecret(ctx context.Context, userID string) (string, error) {
	secret, err := rUser.Get(ctx, constants.TOTPSecret+userID).Result()
	if errors.Is(err, redis.Nil) {
		logger.Infof("totp secret has expired")
		return "", nil
	}
	if err != nil {
		logger.Errorf("redis get totp secret err: %v", err)
		return "", err
	}
	return secret, nil
}
