package service

import (
	"context"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/oss"
	"github.com/nnieie/golanglab5/pkg/utils"
)

type UserService struct {
	ctx          context.Context
	avatarBucket *oss.AvatarOSSCli
}

func NewUserService(ctx context.Context, snowflake *utils.Snowflake) *UserService {
	return &UserService{
		ctx:          ctx,
		avatarBucket: oss.NewAvatarOSSCli(constants.AvatarBucketName, constants.AvatarPublicDomain, snowflake),
	}
}
