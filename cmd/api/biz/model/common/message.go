package common

import "github.com/nnieie/golanglab5/cmd/api/biz/model/base"

// Message ws消息结构
type Message struct {
	Type       int64  // 消息类型：1-私聊发送 2-私聊历史 3-私聊离线 4-群聊发送 5-群聊历史 6-群聊离线
	FromUserID int64  `json:"from_user_id"`
	ToUserID   int64  `json:"to_user_id"`
	GroupID    int64  `json:"group_id"`
	Content    string `json:"content"`
	PageNum    int64  `json:"page_num"`
	PageSize   int64  `json:"page_size"`
	CreatedAt  int64  `json:"created_at"`
}

type SendPrivateMessageRequest struct {
	ToUserID int64
	Content  string
}

type SendPrivateMessageResponse struct {
	Base *base.BaseResp `json:"base"`
}

type QueryPrivateHistoryMessageRequest struct {
	ToUserID int64
	PageNum  int64
	PageSize int64
}

type QueryPrivateHistoryMessageResponse struct {
	Base     *base.BaseResp         `json:"base"`
	Messages []*base.PrivateMessage `json:"messages"`
}

type QueryPrivateOfflineMessageRequest struct {
	ToUserID int64
	PageNum  int64
	PageSize int64
}

type QueryPrivateOfflineMessageResponse struct {
	Base     *base.BaseResp         `json:"base"`
	Messages []*base.PrivateMessage `json:"messages"`
}

type SendGroupMessageRequest struct {
	GroupID int64
	Content string
}

type SendGroupMessageResponse struct {
	Base *base.BaseResp `json:"base"`
}

type QueryGroupHistoryMessageRequest struct {
	GroupID  int64
	PageNum  int64
	PageSize int64
}

type QueryGroupHistoryMessageResponse struct {
	Base     *base.BaseResp       `json:"base"`
	Messages []*base.GroupMessage `json:"messages"`
}

type QueryGroupOfflineMessageRequest struct {
	GroupID  int64
	PageNum  int64
	PageSize int64
}

type QueryGroupOfflineMessageResponse struct {
	Base     *base.BaseResp       `json:"base"`
	Messages []*base.GroupMessage `json:"messages"`
}
