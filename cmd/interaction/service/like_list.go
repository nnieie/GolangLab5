package service

import (
	"strconv"

	"github.com/nnieie/golanglab5/cmd/interaction/dal/db"
	"github.com/nnieie/golanglab5/cmd/interaction/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *interactionService) GetLikeList(userID string, pageNum, pageSize int64) ([]*base.Video, error) {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, err
	}
	likes, err := db.QueryLikeVideoListByUserID(s.ctx, intUserID, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	videos, err := rpc.QueryVideosByIDs(s.ctx, likes)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
