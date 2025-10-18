package service

import (
	"github.com/nnieie/golanglab5/cmd/interaction/dal/db"
	"github.com/nnieie/golanglab5/cmd/interaction/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *interactionService) GetLikeList(userID, pageNum, pageSize int64) ([]*base.Video, error) {
	likes, err := db.QueryLikeVideoListByUserID(s.ctx, userID, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	videos, err := rpc.QueryVideosByIDs(s.ctx, likes)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
