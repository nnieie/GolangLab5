package service

import (
	"github.com/nnieie/golanglab5/cmd/video/dal/db"
	"github.com/nnieie/golanglab5/cmd/video/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *VideoService) GetPopularVideo(pageNum, pageSize int64) ([]*base.Video, error) {
	videos, err := db.QueryVideoByPopular(s.ctx, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return pack.DBVideoToBaseVideos(videos), nil
}
