package service

import (
	"context"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/oss"
	"github.com/nnieie/golanglab5/pkg/utils"
)

type VideoService struct {
	ctx         context.Context
	videoBucket *oss.VideoOSSCli
}

func NewVideoService(ctx context.Context, snowflake *utils.Snowflake) *VideoService {
	return &VideoService{
		ctx:         ctx,
		videoBucket: oss.NewVideoOSSCli(constants.VideoBucketName, constants.VideoPublicDomain, snowflake),
	}
}
