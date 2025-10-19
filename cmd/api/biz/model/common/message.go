package common

import "github.com/nnieie/golanglab5/cmd/api/biz/model/base"

// Message ws消息结构
type Message struct {
	Type       int64 // 消息类型：1-私聊发送 2-私聊历史 3-私聊离线 4-群聊发送 5-群聊历史 6-群聊离线
	FromUserID int64
	ToUserID   int64
	GroupID    int64
	Content    string
	PageNum    int64
	PageSize   int64
	CreatedAt  int64
}

type SendPrivateMessageRequest struct {
	ToUserID int64
	Content  string
}

type SendPrivateMessageResponse struct {
	Base *base.BaseResp
}

type QueryPrivateHistoryMessageRequest struct {
	ToUserID int64
	PageNum  int64
	PageSize int64
}

type QueryPrivateHistoryMessageResponse struct {
	Base     *base.BaseResp
	Messages []*base.PrivateMessage
}

type QueryPrivateOfflineMessageRequest struct {
	ToUserID int64
	PageNum  int64
	PageSize int64
}

type QueryPrivateOfflineMessageResponse struct {
	Base     *base.BaseResp
	Messages []*base.PrivateMessage
}

type SendGroupMessageRequest struct {
	GroupID int64
	Content string
}

type SendGroupMessageResponse struct {
	Base *base.BaseResp
}

type QueryGroupHistoryMessageRequest struct {
	GroupID  int64
	PageNum  int64
	PageSize int64
}

type QueryGroupHistoryMessageResponse struct {
	Base     *base.BaseResp
	Messages []*base.GroupMessage
}

type QueryGroupOfflineMessageRequest struct {
	GroupID  int64
	PageNum  int64
	PageSize int64
}

type QueryGroupOfflineMessageResponse struct {
	Base     *base.BaseResp
	Messages []*base.GroupMessage
}
