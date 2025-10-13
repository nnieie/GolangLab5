package service

import (
	"errors"

	"github.com/nnieie/golanglab5/cmd/user/dal/db"
	"github.com/nnieie/golanglab5/cmd/user/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func (s *UserService) Login(username, password string, code *string) (*base.User, error) {
	logger.Debugf("Login request(rpc): username=%s, password=%s, code=%v", username, password, code)

	user, err := db.QueryUserByNameWithPassword(s.ctx, username)
	if errors.Is(err, errno.UserIsNotExistErr) || user == nil || !utils.VerifyPassword(password, user.Password) {
		return nil, errno.UserIsNotExistOrPasswordErr
	}
	if err != nil {
		return nil, err
	}

	if user.TOTP != "" {
		if code == nil || !utils.CheckTotp(*code, user.TOTP) {
			logger.Infof("user %d login err: invalid mfa code, %v %s", user.ID, code, user.TOTP)
			return nil, errno.MFAInvalidCodeErr
		}
	}

	return pack.DBUserTobaseUser(user), nil
}
