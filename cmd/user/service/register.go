package service

import (
	"github.com/nnieie/golanglab5/cmd/user/dal/db"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func (s *UserService) Register(username, password string) (int64, error) {
	pwd, err := utils.Crypt(password)
	if err != nil {
		return 0, err
	}
	userID, err := db.CreateUser(s.ctx, &db.User{
		UserName: username,
		Password: pwd,
	})
	if err != nil {
		return 0, err
	}
	return userID, nil
}
