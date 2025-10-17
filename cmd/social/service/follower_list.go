package service

import (
	"github.com/nnieie/golanglab5/cmd/social/dal/db"
	"github.com/nnieie/golanglab5/cmd/social/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *socialService) GetFollowerList(userID, pageNum, pageSize int64) ([]*base.User, int64, error) {
	followersID, err := db.QueryFollowerList(s.ctx, userID, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := db.QueryFollowerCount(s.ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	followers, err := rpc.QueryUsersByIDs(s.ctx, followersID)
	if err != nil {
		return nil, 0, err
	}
	return followers, count, nil
}
