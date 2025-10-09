package main

import (
	"context"

	chat "github.com/nnieie/golanglab5/kitex_gen/chat"
)

// ChatServiceImpl implements the last service interface defined in the IDL.
type ChatServiceImpl struct{}

// SendMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) SendMessage(ctx context.Context, req *chat.SendMessageRequest,
) (resp *chat.SendMessageResponse, err error) {
	// TODO: Your code here...
	return
}

// QueryPrivateOfflineMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) QueryPrivateOfflineMessage(ctx context.Context, req *chat.QueryPrivateOfflineMessageRequest,
) (resp *chat.QueryPrivateOfflineMessageResponse, err error) {
	// TODO: Your code here...
	return
}

// QueryPrivateHistoryMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) QueryPrivateHistoryMessage(ctx context.Context, req *chat.QueryPrivateHistoryMessageRequest,
) (resp *chat.QueryPrivateHistoryMessageResponse, err error) {
	// TODO: Your code here...
	return
}

// QueryGroupHistoryMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) QueryGroupHistoryMessage(ctx context.Context, req *chat.QueryGroupHistoryMessageRequest,
) (resp *chat.QueryGroupHistoryMessageResponse, err error) {
	// TODO: Your code here...
	return
}
