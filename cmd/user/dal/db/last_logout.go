package db

import (
	"context"
	"time"

	"github.com/nnieie/golanglab5/pkg/constants"
)

func QueryLastLogoutTime(ctx context.Context, userID int64) (time.Time, error) {
	var logoutTime time.Time
	err := DB.WithContext(ctx).Table(constants.LastLogoutTimeTableName).Where("user_id = ?", userID).Select("logout_time").Scan(&logoutTime).Error
	return logoutTime, err
}

func UpdateLastLogoutTime(ctx context.Context, userID int64, logoutTime time.Time) error {
	return DB.WithContext(ctx).Table(constants.LastLogoutTimeTableName).
		Where("user_id = ?", userID).
		Assign(map[string]interface{}{"logout_time": logoutTime}).
		FirstOrCreate(&struct {
			UserID     int64     `gorm:"column:user_id"`
			LogoutTime time.Time `gorm:"column:logout_time"`
		}{UserID: userID, LogoutTime: logoutTime}).Error
}
