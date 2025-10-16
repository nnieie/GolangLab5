package pack

import (
	"github.com/nnieie/golanglab5/cmd/video/dal/db"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func DBVideoToBaseVideo(video *db.Video) *base.Video {
	return &base.Video{
		Id:          int64(video.ID),
		Title:       video.Title,
		Description: video.Description,
		CreatedAt:   video.CreatedAt.String(),
		UpdatedAt:   video.UpdatedAt.String(),
	}
}

func DBVideoToBaseVideos(videos []*db.Video) []*base.Video {
	baseVideos := make([]*base.Video, 0, len(videos))
	for _, video := range videos {
		baseVideos = append(baseVideos, DBVideoToBaseVideo(video))
	}
	return baseVideos
}
