package cache

import "time"

const (
	privateMessageTTL  = 7 * 24 * time.Hour
	groupMessageTTL    = 7 * 24 * time.Hour
	maxMessageListSize = 1000
)

// 缓存中的私聊消息结构
type CachedPrivateMessage struct {
	ID         string
	FromUserID int64
	ToUserID   int64
	Content    string
	CreatedAt  time.Time
}

// 缓存中的群聊消息结构
type CachedGroupMessage struct {
	ID         string
	GroupID    int64
	FromUserID int64
	Content    string
	CreatedAt  time.Time
}

// Redis键名常量
const (
	PrivateMessagePrefix   = "msg:private:"
	GroupMessagePrefix     = "msg:group:"
	PrivateMessageQueueKey = "msg:queue:private"
	GroupMessageQueueKey   = "msg:queue:group"
)
