package db

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/cmd/video/rpc"
	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/logger"
)

type Video struct {
	UserID       int64
	VideoURL     string
	CoverURL     string
	Title        string
	Description  string
	VisitCount   int64
	LikeCount    int64
	CommentCount int64
	gorm.Model
}

const (
	maxSearchPageSize = 100
)

func CreateVideo(ctx context.Context, video *Video) (int64, error) {
	err := DB.WithContext(ctx).Create(video).Error
	if err != nil {
		return 0, err
	}
	return int64(video.ID), nil
}

func QueryVideoByID(ctx context.Context, videoID int64) (*Video, error) {
	var video Video
	err := DB.WithContext(ctx).Where("id = ?", videoID).Find(&video).Error
	if err != nil {
		return nil, err
	}
	if video == (Video{}) {
		return nil, errno.VideoIsNotExistErr
	}
	return &video, nil
}

func QueryVideoByIDs(ctx context.Context, videoIDs []int64) ([]*Video, error) {
	var videos []*Video
	// 使用 WHERE id IN (?) 保持传入 ID 的顺序
	err := DB.WithContext(ctx).Where("id IN (?)", videoIDs).Order(gorm.Expr("FIELD(id, ?)", videoIDs)).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func QueryVideoByLatestTime(ctx context.Context, latestTime time.Time) ([]*Video, error) {
	var videos []*Video
	err := DB.WithContext(ctx).Where("created_at > ?", latestTime).Order("created_at DESC").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return nil, errno.VideoIsNotExistErr
	}
	return videos, nil
}

func QueryVideoByUserID(ctx context.Context, userID int64, pageNum, pageSize int64) ([]*Video, error) {
	var videos []*Video
	err := DB.WithContext(ctx).Where("user_id = ?", userID).Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return nil, errno.VideoIsNotExistErr
	}
	return videos, nil
}

func QueryVideoCountByUserID(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&Video{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func QueryVideoByPopular(ctx context.Context, pageNum, pageSize int64) ([]*Video, error) {
	var videos []*Video
	err := DB.WithContext(ctx).Order("visit_count DESC").Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	if len(videos) == 0 {
		return nil, errno.VideoIsNotExistErr
	}
	return videos, nil
}

func SearchVideos(
	ctx context.Context,
	keywords string,
	pageNum, pageSize int64,
	fromDate, toDate *int64,
	username *string,
) ([]*Video, int64, error) {
	var videos []*Video
	logger.Debugf("SearchVideos keywords: %s, pageNum: %d, pageSize: %d, fromDate: %v, toDate: %v, username: %v",
		keywords, pageNum, pageSize, fromDate, toDate, username)
	query := DB.WithContext(ctx)

	if strings.TrimSpace(keywords) != "" {
		likePattern := "%" + keywords + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", likePattern, likePattern)
	}

	if fromDate != nil {
		query = query.Where("created_at >= ?", time.Unix(*fromDate, 0))
	}
	if toDate != nil {
		query = query.Where("created_at <= ?", time.Unix(*toDate, 0))
	}

	if username != nil && strings.TrimSpace(*username) != "" {
		userIds, err := rpc.SearchUserIds(ctx, *username, 1, maxSearchPageSize)
		if err != nil {
			return nil, 0, err
		}
		query = query.Where("user_id IN ?", userIds)
	}

	err := query.Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&videos).Error
	if err != nil {
		return nil, 0, err
	}
	if len(videos) == 0 {
		return nil, 0, errno.VideoIsNotExistErr
	}
	var count int64
	err = query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return videos, count, nil
}

func IncVideoVisitCount(ctx context.Context, videoID int64) error {
	return DB.WithContext(ctx).Where("id = ?", videoID).UpdateColumn("visit_count", gorm.Expr("visit_count + 1")).Error
}
