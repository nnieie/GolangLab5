package service

import (
	"github.com/nnieie/golanglab5/cmd/user/dal/cache"
	"github.com/nnieie/golanglab5/cmd/user/dal/db"
	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func (s *UserService) BindMFA(code, secret string, userID int64) (bool, error) {
	// 先检查此用户15分钟内有没有获取过 TOTPSecret
	cacheEncryptedSecret, err := cache.GetTOTPSecret(s.ctx, userID)
	if err != nil {
		return false, err
	}
	if cacheEncryptedSecret == "" {
		logger.Infof("user %d bind mfa err: not generate totp", userID)
		return false, errno.NotGenerateTotpErr
	}
	// 解密 redis 里的 secret 再验证
	cacheSecret, err := utils.Decrypt(cacheEncryptedSecret)
	if err != nil {
		return false, err
	}
	if cacheSecret != secret {
		logger.Infof("user %d bind mfa err: secret not match %s and %s", userID, secret, cacheSecret)
		return false, errno.MFAInvalidCodeErr
	}

	ok := utils.CheckTotp(code, cacheEncryptedSecret)
	if !ok {
		logger.Infof("user %d bind mfa err: invalid code, %s %s", userID, code, cacheSecret)
		return false, errno.MFAInvalidCodeErr
	}

	err = db.UpdateMFA(s.ctx, secret, userID)
	if err != nil {
		return false, err
	}

	return true, nil
}
