package service

import (
	"errors"

	"github.com/nnieie/golanglab5/cmd/user/dal/db"
	"github.com/nnieie/golanglab5/cmd/user/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func (s *UserService) Login(username, password, code string) (*base.User, error) {
	user, err := db.QueryUserByNameWithPassword(s.ctx, username)
	if errors.Is(err, errno.UserIsNotExistErr) || user == nil || !utils.VerifyPassword(password, user.Password) {
		return nil, errno.UserIsNotExistOrPasswordErr
	}
	if err != nil {
		return nil, err
	}
	return pack.DBUserTobaseUser(user), nil
}
