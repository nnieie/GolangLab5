package service

import (
	"github.com/nnieie/golanglab5/cmd/user/dal/cache"
	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func (s *UserService) BindMFA(code, secret string, userID int64) (bool, error) {
	// 先检查此用户15分钟内有没有获取过 TOTPSecret
	cacheSecret, err := cache.GetTOTPSecret(s.ctx, userID)
	if err != nil {
		return false, err
	}
	if cacheSecret == "" {
		return false, errno.NotGenerateTotpErr
	}
	if cacheSecret != secret {
		return false, errno.MFAInvalidCodeErr
	}

	ok := utils.CheckTotp(code, secret)
	if !ok {
		return false, errno.MFAInvalidCodeErr
	}

	return true, nil
}
