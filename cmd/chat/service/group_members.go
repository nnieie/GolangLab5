package service

import (
	"github.com/nnieie/golanglab5/cmd/chat/dal/db"
)

func (s *ChatService) GetGroupMembers(groupID string) ([]int64, error) {
	return db.QueryGroupMemberIDList(s.ctx, groupID)
}

func (s *ChatService) CheckUserExistInGroup(userID string, groupID string) (bool, error) {
	return db.CheckUserExistInGroup(s.ctx, userID, groupID)
}
