package kafka

import (
	"fmt"
	"time"
)

const (
	// RandomIDRange 用于生成事件ID的随机数范围
	// 值为 1000 表示生成 0-999 的随机数（3位数）
	RandomIDRange = 1000
)

// EventType 事件类型
type EventType string

const (
	EventTypeLike EventType = "like"
)

// BaseEvent 基础事件结构
type BaseEvent struct {
	EventID   string    `json:"event_id"`   // 事件唯一ID
	EventType EventType `json:"event_type"` // 事件类型
	Timestamp int64     `json:"timestamp"`  // 事件时间戳（毫秒）
	UserID    int64     `json:"user_id"`    // 操作用户ID
}

// LikeEvent 点赞事件
type LikeEvent struct {
	BaseEvent
	VideoID   *int64 `json:"video_id,omitempty"`   // 视频ID（可选）
	CommentID *int64 `json:"comment_id,omitempty"` // 评论ID（可选）
	Action    int64  `json:"action"`               // 1=点赞, 2=取消点赞
}

// GenerateEventID 生成事件ID
// 格式: {userID}_{timestamp}_{random}
func GenerateEventID(userID, timestamp int64) string {
	return fmt.Sprintf("%d_%d_%d", userID, timestamp, time.Now().UnixNano()%RandomIDRange)
}
