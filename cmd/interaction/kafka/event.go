package kafka

// LikeEvent 点赞事件
type LikeEvent struct {
	UserID    int64  `json:"user_id"`              // 操作用户ID
	VideoID   *int64 `json:"video_id,omitempty"`   // 视频ID
	CommentID *int64 `json:"comment_id,omitempty"` // 评论ID
	Action    int64  `json:"action"`               // 操作类型：1=点赞, 2=取消点赞
}

// NewLikeEvent 创建点赞事件
func NewLikeEvent(userID int64, videoID, commentID *int64, action int64) *LikeEvent {
	return &LikeEvent{
		UserID:    userID,
		VideoID:   videoID,
		CommentID: commentID,
		Action:    action,
	}
}
