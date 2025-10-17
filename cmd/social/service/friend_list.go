package service

import (
	"github.com/nnieie/golanglab5/cmd/social/dal/db"
	"github.com/nnieie/golanglab5/cmd/social/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *socialService) GetFriendList(userID, pageNum, pageSize int64) ([]*base.User, int64, error) {
	friendsID, err := db.QueryFriendList(s.ctx, userID, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := db.QueryFriendCount(s.ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	friends, err := rpc.QueryUsersByIDs(s.ctx, friendsID)
	if err != nil {
		return nil, 0, err
	}
	return friends, count, nil
}
