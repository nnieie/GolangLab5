package service

import (
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/nnieie/golanglab5/cmd/video/dal/db"
	"github.com/nnieie/golanglab5/cmd/video/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *VideoService) GetVideoList(userID string, pageNum, pageSize int64) ([]*base.Video, int64, error) {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, 0, err
	}

	var videos []*db.Video
	var count int64

	// 创建一个 errgroup 管理 goroutine 和错误处理
	eg, ctx := errgroup.WithContext(s.ctx)

	// 获取视频列表
	eg.Go(func() error {
		var err error
		// 传入的是 ctx 而不是 s.ctx 有 goroutine 出错时会取消其他操作
		videos, err = db.QueryVideoByUserID(ctx, intUserID, pageNum, pageSize)
		return err // 如果出错，errgroup 会捕获
	})

	// 获取视频总数
	eg.Go(func() error {
		var err error
		count, err = db.QueryVideoCountByUserID(ctx, intUserID)
		return err
	})

	// 等待两个 goroutine 完成
	if err := eg.Wait(); err != nil {
		// 只要有一个出错，就返回错误
		return nil, 0, err
	}

	return pack.DBVideosToBaseVideos(videos), count, nil
}
