package service

import (
	"github.com/nnieie/golanglab5/cmd/video/dal/db"
	"github.com/nnieie/golanglab5/cmd/video/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *VideoService) SearchVideo(keywords string, pageNum, pageSize int64, fromDate, toDate *int64, username *string) (
	[]*base.Video, int64, error) {
	videos, total, err := db.SearchVideos(s.ctx, keywords, pageNum, pageSize, fromDate, toDate, username)
	if err != nil {
		return nil, 0, err
	}
	return pack.DBVideoToBaseVideos(videos), total, nil
}
