package service

import (
	"github.com/nnieie/golanglab5/cmd/social/dal/db"
	"github.com/nnieie/golanglab5/cmd/social/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *socialService) GetFollowingList(userID, pageNum, pageSize int64) ([]*base.User, int64, error) {
	followingsID, err := db.QueryFollowingList(s.ctx, userID, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := db.QueryFollowingCount(s.ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	followings, err := rpc.QueryUsersByIDs(s.ctx, followingsID)
	if err != nil {
		return nil, 0, err
	}
	return followings, count, nil
}
