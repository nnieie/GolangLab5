package service

import (
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/nnieie/golanglab5/cmd/social/dal/db"
	"github.com/nnieie/golanglab5/cmd/social/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *socialService) GetFriendList(userID string, pageNum, pageSize int64) ([]*base.User, int64, error) {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, 0, err
	}
	var friends []*base.User
	var count int64

	eg, ctx := errgroup.WithContext(s.ctx)
	eg.Go(func() error {
		friendsID, err := db.QueryFriendList(ctx, intUserID, pageNum, pageSize)
		if err != nil {
			return err
		}
		friends, err = rpc.QueryUsersByIDs(ctx, friendsID)
		return err
	})
	eg.Go(func() error {
		var err error
		count, err = db.QueryFriendCount(ctx, intUserID)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}
	return friends, count, nil
}
