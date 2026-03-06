package service

import (
	"github.com/nnieie/golanglab5/cmd/user/dal/db"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func (s *UserService) Register(username, password string) (string, error) {
	pwd, err := utils.Crypt(password)
	if err != nil {
		return "", err
	}
	userID, err := db.CreateUser(s.ctx, &db.User{
		UserName: username,
		Password: pwd,
	})
	if err != nil {
		return "", err
	}
	return userID, nil
}
