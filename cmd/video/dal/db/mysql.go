package db

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/cmd/video/rpc"
	"github.com/nnieie/golanglab5/pkg/constants"
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

func (Video) TableName() string {
	return constants.VideoTableName
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

func parseID(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}

func parseIDs(ids []string) ([]int64, error) {
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		parsedID, err := parseID(id)
		if err != nil {
			return nil, err
		}
		result = append(result, parsedID)
	}
	return result, nil
}

func QueryVideoByID(ctx context.Context, videoID string) (*Video, error) {
	parsedVideoID, err := parseID(videoID)
	if err != nil {
		return nil, errno.ParamErr
	}
	var video Video
	err = DB.WithContext(ctx).Where("id = ?", parsedVideoID).First(&video).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.VideoIsNotExistErr
		}
		return nil, err
	}
	return &video, nil
}

func QueryVideosByIDs(ctx context.Context, videoIDs []string) ([]*Video, error) {
	parsedVideoIDs, err := parseIDs(videoIDs)
	if err != nil {
		return nil, errno.ParamErr
	}
	var videos []*Video
	// 使用 FIELD 保持传入 ID 的顺序
	err = DB.WithContext(ctx).Where("id IN (?)", parsedVideoIDs).Order(gorm.Expr("FIELD(id, ?)", parsedVideoIDs)).Find(&videos).Error
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
	return videos, nil
}

func QueryVideoByUserID(ctx context.Context, userID int64, pageNum, pageSize int64) ([]*Video, error) {
	var videos []*Video
	err := DB.WithContext(ctx).Where("user_id = ?", userID).Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&videos).Error
	if err != nil {
		return nil, err
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
	query := DB.WithContext(ctx).Model(&Video{})

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
	var count int64
	// 先 Count，再 Limit/Offset
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&videos).Error
	if err != nil {
		return nil, 0, err
	}
	return videos, count, nil
}

func IncVideoVisitCount(ctx context.Context, videoID string) error {
	parsedVideoID, err := parseID(videoID)
	if err != nil {
		return errno.ParamErr
	}
	return DB.WithContext(ctx).Model(&Video{}).Where("id = ?", parsedVideoID).UpdateColumn("visit_count", gorm.Expr("visit_count + 1")).Error
}

func QueryVideoLikeCount(ctx context.Context, videoID string) (int64, error) {
	parsedVideoID, err := parseID(videoID)
	if err != nil {
		return 0, errno.ParamErr
	}
	var video Video
	err = DB.WithContext(ctx).Where("id = ?", parsedVideoID).First(&video).Error
	if err != nil {
		logger.Errorf("QueryVideoLikeCount err: %v", err)
		return 0, err
	}
	return video.LikeCount, nil
}

func SetVideoLikeCount(ctx context.Context, videoID string, likeCount int64) error {
	parsedVideoID, err := parseID(videoID)
	if err != nil {
		return errno.ParamErr
	}
	return DB.WithContext(ctx).Model(&Video{}).Where("id = ?", parsedVideoID).Update("like_count", likeCount).Error
}

func UpdateVideoLikeCount(ctx context.Context, videoID string, delta int64) error {
	parsedVideoID, err := parseID(videoID)
	if err != nil {
		return errno.ParamErr
	}
	return DB.WithContext(ctx).Model(&Video{}).Where("id = ?", parsedVideoID).UpdateColumn("like_count", gorm.Expr("like_count + ?", delta)).Error
}

func BatchUpdateVideoLikeCount(ctx context.Context, counts map[int64]int64) error {
	if len(counts) == 0 {
		return nil
	}

	query := "UPDATE videos SET like_count = like_count + CASE id "
	var args []interface{}
	var ids []int64

	for id, count := range counts {
		query += "WHEN ? THEN ? "
		args = append(args, id, count)
		ids = append(ids, id)
	}
	query += "END WHERE id IN ?"
	args = append(args, ids)

	if err := DB.WithContext(ctx).Exec(query, args...).Error; err != nil {
		return err
	}
	return nil
}
