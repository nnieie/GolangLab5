package pack

import (
	"strconv"

	"github.com/nnieie/golanglab5/cmd/video/dal/db"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func DBVideoToBaseVideo(video *db.Video) *base.Video {
	return &base.Video{
		Id:           strconv.FormatUint(uint64(video.ID), 10),
		UserId:       strconv.FormatInt(video.UserID, 10),
		Title:        video.Title,
		Description:  video.Description,
		VideoUrl:     video.VideoURL,
		CoverUrl:     video.CoverURL,
		VisitCount:   video.VisitCount,
		LikeCount:    video.LikeCount,
		CommentCount: video.CommentCount,
		CreatedAt:    video.CreatedAt.String(),
		UpdatedAt:    video.UpdatedAt.String(),
	}
}

func DBVideosToBaseVideos(videos []*db.Video) []*base.Video {
	baseVideos := make([]*base.Video, 0, len(videos))
	for _, video := range videos {
		baseVideos = append(baseVideos, DBVideoToBaseVideo(video))
	}
	return baseVideos
}
