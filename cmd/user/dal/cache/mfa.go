package cache

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func SetTOTPSecret(ctx context.Context, secret string, userID int64) error {
	err := rUser.Set(ctx, constants.TOTPSecret+utils.I64ToStr(userID), secret, constants.TOTPSecretExpTime).Err()
	if err != nil {
		logger.Errorf("redis set totp secret err: %v", err)
		return err
	}
	return nil
}

func GetTOTPSecret(ctx context.Context, userID int64) (string, error) {
	secret, err := rUser.Get(ctx, constants.TOTPSecret+utils.I64ToStr(userID)).Result()
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
