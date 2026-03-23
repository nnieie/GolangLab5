package pack

import (
	"strconv"

	"github.com/nnieie/golanglab5/internal/chat/dal/cache"
	"github.com/nnieie/golanglab5/internal/chat/dal/db"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func DBPrivateMessageToChatPrivateMessage(dbMsg *db.PrivateMessage) *base.PrivateMessage {
	if dbMsg == nil {
		return nil
	}
	return &base.PrivateMessage{
		FromUserId: strconv.FormatInt(dbMsg.FromUserID, 10),
		ToUserId:   strconv.FormatInt(dbMsg.ToUserID, 10),
		Content:    dbMsg.Content,
		CreatedAt:  dbMsg.CreatedAt.Unix(),
	}
}

func DBPrivateMessagesToChatPrivateMessages(dbMsgs []*db.PrivateMessage) []*base.PrivateMessage {
	chatMsgs := make([]*base.PrivateMessage, 0, len(dbMsgs))
	for _, msg := range dbMsgs {
		chatMsgs = append(chatMsgs, DBPrivateMessageToChatPrivateMessage(msg))
	}
	return chatMsgs
}
func DBGroupMessageToChatGroupMessage(dbMsg *db.GroupMessage) *base.GroupMessage {
	if dbMsg == nil {
		return nil
	}
	return &base.GroupMessage{
		FromUserId: strconv.FormatInt(dbMsg.FromUserID, 10),
		GroupId:    strconv.FormatInt(dbMsg.GroupID, 10),
		Content:    dbMsg.Content,
		CreatedAt:  dbMsg.CreatedAt.Unix(),
	}
}
func DBGroupMessagesToChatGroupMessages(dbMsgs []*db.GroupMessage) []*base.GroupMessage {
	chatMsgs := make([]*base.GroupMessage, 0, len(dbMsgs))
	for _, msg := range dbMsgs {
		chatMsgs = append(chatMsgs, DBGroupMessageToChatGroupMessage(msg))
	}
	return chatMsgs
}

func ConvertCachedToBaseMessages(cachedMsgs []*cache.CachedPrivateMessage) []*base.PrivateMessage {
	baseMsgs := make([]*base.PrivateMessage, 0, len(cachedMsgs))
	for _, cachedMsg := range cachedMsgs {
		baseMsgs = append(baseMsgs, &base.PrivateMessage{
			FromUserId: strconv.FormatInt(cachedMsg.FromUserID, 10),
			ToUserId:   strconv.FormatInt(cachedMsg.ToUserID, 10),
			Content:    cachedMsg.Content,
			CreatedAt:  cachedMsg.CreatedAt.Unix(),
		})
	}
	return baseMsgs
}

// 转换缓存的群聊消息为基础消息
func ConvertCachedGroupToBaseMessages(cachedMsgs []*cache.CachedGroupMessage) []*base.GroupMessage {
	baseMsgs := make([]*base.GroupMessage, 0, len(cachedMsgs))
	for _, cachedMsg := range cachedMsgs {
		baseMsgs = append(baseMsgs, &base.GroupMessage{
			FromUserId: strconv.FormatInt(cachedMsg.FromUserID, 10),
			GroupId:    strconv.FormatInt(cachedMsg.GroupID, 10),
			Content:    cachedMsg.Content,
			CreatedAt:  cachedMsg.CreatedAt.Unix(),
		})
	}
	return baseMsgs
}
