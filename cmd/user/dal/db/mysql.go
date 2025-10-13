package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/logger"
)

type User struct {
	UserName string
	Password string
	Avatar   string
	TOTP     string
	gorm.Model
}

func CreateUser(ctx context.Context, user *User) (int64, error) {
	err := DB.WithContext(ctx).Create(user).Error
	if err != nil {
		logger.Errorf("create user err: %v", err)
		return 0, err
	}
	logger.Infof("CreateUser: created id=%d user=%s", user.ID, user.UserName)
	return int64(user.ID), nil
}

func QueryUserByID(ctx context.Context, userID int64) (*User, error) {
	var user User
	err := DB.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.UserIsNotExistErr
		}
		return nil, err
	}
	if user == (User{}) {
		return nil, errno.UserIsNotExistErr
	}
	user.Password = ""
	return &user, nil
}

func QueryUserByName(ctx context.Context, username string) (*User, error) {
	var user User
	err := DB.WithContext(ctx).Where("user_name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.UserIsNotExistErr
		}
		return nil, err
	}
	if user == (User{}) {
		return nil, errno.UserIsNotExistErr
	}
	user.Password = ""
	return &user, nil
}

func QueryUserByNameWithPassword(ctx context.Context, username string) (*User, error) {
	var user User
	err := DB.WithContext(ctx).Where("user_name = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.UserIsNotExistErr
		}
		return nil, err
	}
	if user == (User{}) {
		return nil, errno.UserIsNotExistErr
	}
	return &user, nil
}

func QueryUserByIDList(ctx context.Context, userIds []int64) (*[]User, error) {
	users := new([]User)
	if len(userIds) == 0 {
		return users, errno.UserIsNotExistErr
	}

	// 构建 ORDER BY FIELD 字符串，保证返回顺序
	var builder strings.Builder
	builder.WriteString("FIELD(id")
	for _, id := range userIds {
		builder.WriteString(fmt.Sprintf(", %d", id))
	}
	builder.WriteString(")")
	orderClause := builder.String()

	if err := DB.WithContext(ctx).Where("id IN ?", userIds).Order(orderClause).Find(users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func UpdateMFA(ctx context.Context, secret string, userID int64) error {
	err := DB.WithContext(ctx).Model(User{}).Where("id = ?", userID).Update("totp", secret).Error
	if err != nil {
		logger.Errorf("update totp secret err: %v", err)
		return err
	}
	return nil
}

func UpdateAvatar(ctx context.Context, userID int64, avatar string) (*User, error) {
	// 先执行更新
	if err := DB.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("avatar", avatar).Error; err != nil {
		logger.Errorf("update avatar err: %v", err)
		return nil, err
	}

	// 再单独查询最新记录
	user := new(User)
	if err := DB.WithContext(ctx).Where("id = ?", userID).First(user).Error; err != nil {
		logger.Errorf("query updated user err: %v", err)
		return nil, err
	}
	return user, nil
}
