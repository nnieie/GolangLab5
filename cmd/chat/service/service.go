package service

import "context"

type ChatService struct {
	ctx context.Context
}

func NewChatService(ctx context.Context) *ChatService {
	return &ChatService{
		ctx: ctx,
	}
}
