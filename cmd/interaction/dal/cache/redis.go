package cache

import (
	"context"
	"fmt"

	"github.com/nnieie/golanglab5/pkg/constants"
)

func CheckVideoLikeExists(ctx context.Context, userID, videoID int64) (bool, error) {
	key := fmt.Sprintf("%s:%d:%d", constants.VideoLikePrefix, userID, videoID)
	exists, err := rInteraction.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func CheckCommentLikeExists(ctx context.Context, userID, commentID int64) (bool, error) {
	key := fmt.Sprintf("%s:%d:%d", constants.CommentLikePrefix, userID, commentID)
	exists, err := rInteraction.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func IncrVideoLikeCount(ctx context.Context, videoID int64, delta int64) error {
	key := fmt.Sprintf("%s:%d", constants.VideoLikeCountPrefix, videoID)
	return rInteraction.IncrBy(ctx, key, delta).Err()
}

func IncrCommentLikeCount(ctx context.Context, commentID int64, delta int64) error {
	key := fmt.Sprintf("%s:%d", constants.CommentLikeCountPrefix, commentID)
	return rInteraction.IncrBy(ctx, key, delta).Err()
}

func SetVideoLike(ctx context.Context, userID, videoID int64) error {
	key := fmt.Sprintf("%s:%d:%d", constants.VideoLikePrefix, userID, videoID)
	return rInteraction.Set(ctx, key, "1", 0).Err()
}

func SetCommentLike(ctx context.Context, userID, commentID int64) error {
	key := fmt.Sprintf("%s:%d:%d", constants.CommentLikePrefix, userID, commentID)
	return rInteraction.Set(ctx, key, "1", 0).Err()
}

func DelVideoLike(ctx context.Context, userID, videoID int64) error {
	key := fmt.Sprintf("%s:%d:%d", constants.VideoLikePrefix, userID, videoID)
	return rInteraction.Del(ctx, key).Err()
}

func DelCommentLike(ctx context.Context, userID, commentID int64) error {
	key := fmt.Sprintf("%s:%d:%d", constants.CommentLikePrefix, userID, commentID)
	return rInteraction.Del(ctx, key).Err()
}
