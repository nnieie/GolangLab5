package pack

import (
	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/cmd/user/dal/db"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func DBUserTobaseUser(user *db.User) *base.User {
	if user == nil {
		return nil
	}

	return &base.User{
		Id:        int64(user.ID),
		Username:  user.UserName,
		AvatarUrl: user.Avatar,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt: formatDeletedAt(user.DeletedAt),
	}
}

func DBUserTobaseUsers(users []*db.User) []*base.User {
	if users == nil {
		return nil
	}

	baseUsers := make([]*base.User, len(users))
	for i, user := range users {
		baseUsers[i] = DBUserTobaseUser(user)
	}
	return baseUsers
}

func formatDeletedAt(deletedAt gorm.DeletedAt) string {
	if deletedAt.Valid {
		return deletedAt.Time.Format("2006-01-02 15:04:05")
	}
	return "1970-01-01 08:00:00"
}
