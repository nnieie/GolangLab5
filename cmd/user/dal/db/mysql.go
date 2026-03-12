package db

import (
	"context"
	"errors"
	"strconv"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/pkg/constants"
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

func (User) TableName() string {
	return constants.UserTableName
}

func parseUserID(userID string) (int64, error) {
	return strconv.ParseInt(userID, 10, 64)
}

func parseUserIDs(userIDs []string) ([]int64, error) {
	ids := make([]int64, 0, len(userIDs))
	for _, userID := range userIDs {
		id, err := parseUserID(userID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func CreateUser(ctx context.Context, user *User) (string, error) {
	err := DB.WithContext(ctx).Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", errno.UserAlreadyExistErr
		}
		logger.Errorf("create user err: %v", err)
		return "", err
	}
	logger.Infof("CreateUser: created id=%d user=%s", user.ID, user.UserName)
	return strconv.FormatUint(uint64(user.ID), 10), nil
}

func QueryUserByID(ctx context.Context, userID string) (*User, error) {
	id, err := parseUserID(userID)
	if err != nil {
		return nil, errno.ParamErr
	}
	var user User
	err = DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.UserIsNotExistErr
		}
		return nil, err
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
	return &user, nil
}

func QueryUserByIDList(ctx context.Context, userIds []string) ([]*User, error) {
	ids, err := parseUserIDs(userIds)
	if err != nil {
		return nil, errno.ParamErr
	}
	users := make([]*User, 0, len(ids))

	if err := DB.WithContext(ctx).Where("id IN ?", ids).Order(gorm.Expr("FIELD(id, ?)", ids)).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func UpdateMFA(ctx context.Context, secret string, userID string) error {
	id, err := parseUserID(userID)
	if err != nil {
		return errno.ParamErr
	}
	err = DB.WithContext(ctx).Model(User{}).Where("id = ?", id).Update("totp", secret).Error
	if err != nil {
		logger.Errorf("update totp secret err: %v", err)
		return err
	}
	return nil
}

func UpdateAvatar(ctx context.Context, userID string, avatar string) (*User, error) {
	id, err := parseUserID(userID)
	if err != nil {
		return nil, errno.ParamErr
	}

	err = DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("avatar", avatar).Error
	if err != nil {
		logger.Errorf("update avatar err: %v", err)
		return nil, err
	}

	var user User
	err = DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func SearchUserIdsByName(ctx context.Context, pattern string, page, pageSize int64) ([]string, error) {
	ids := make([]int64, 0)
	err := DB.WithContext(ctx).Model(&User{}).
		Where("user_name LIKE ?", "%"+pattern+"%").
		Limit(int(pageSize)).Offset(int((page-1)*pageSize)).
		Pluck("id", &ids).Error
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(ids))
	for _, id := range ids {
		userIds = append(userIds, strconv.FormatInt(id, 10))
	}
	return userIds, nil
}
