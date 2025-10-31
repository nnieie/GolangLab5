package service

import (
	"golang.org/x/sync/errgroup"

	"github.com/nnieie/golanglab5/cmd/social/dal/db"
	"github.com/nnieie/golanglab5/cmd/social/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *socialService) GetFollowingList(userID, pageNum, pageSize int64) ([]*base.User, int64, error) {
	var followings []*base.User
	var count int64

	eg, ctx := errgroup.WithContext(s.ctx)
	eg.Go(func() error {
		followingsID, err := db.QueryFollowingList(ctx, userID, pageNum, pageSize)
		if err != nil {
			return err
		}
		followings, err = rpc.QueryUsersByIDs(ctx, followingsID)
		return err
	})
	eg.Go(func() error {
		var err error
		count, err = db.QueryFollowingCount(ctx, userID)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}
	return followings, count, nil
}
