package db

import (
	"context"
	"strconv"
	"time"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/errno"
)

func QueryLastLogoutTime(ctx context.Context, userID string) (time.Time, error) {
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return time.Time{}, errno.ParamErr
	}
	var logoutTime time.Time
	err = DB.WithContext(ctx).Table(constants.LastLogoutTimeTableName).Where("user_id = ?", id).Select("logout_time").Scan(&logoutTime).Error
	return logoutTime, err
}

func UpdateLastLogoutTime(ctx context.Context, userID string, logoutTime time.Time) error {
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return errno.ParamErr
	}
	return DB.WithContext(ctx).Table(constants.LastLogoutTimeTableName).
		Where("user_id = ?", id).
		Assign(map[string]interface{}{"logout_time": logoutTime}).
		FirstOrCreate(&struct {
			UserID     int64     `gorm:"column:user_id"`
			LogoutTime time.Time `gorm:"column:logout_time"`
		}{UserID: id, LogoutTime: logoutTime}).Error
}
