package service

import (
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/nnieie/golanglab5/internal/social/dal/db"
	"github.com/nnieie/golanglab5/internal/social/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *socialService) GetFollowingList(userID string, pageNum, pageSize int64) ([]*base.User, int64, error) {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, 0, err
	}
	var followings []*base.User
	var count int64

	eg, ctx := errgroup.WithContext(s.ctx)
	eg.Go(func() error {
		followingsID, err := db.QueryFollowingList(ctx, intUserID, pageNum, pageSize)
		if err != nil {
			return err
		}
		followings, err = rpc.QueryUsersByIDs(ctx, followingsID)
		return err
	})
	eg.Go(func() error {
		var err error
		count, err = db.QueryFollowingCount(ctx, intUserID)
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}
	return followings, count, nil
}
