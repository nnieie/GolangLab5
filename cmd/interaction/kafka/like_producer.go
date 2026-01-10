package kafka

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/kafka"
	"github.com/nnieie/golanglab5/pkg/logger"
)

func SendLikeEvent(ctx context.Context, userID int64, videoID, commentID *int64, action int64) {
	event := NewLikeEvent(userID, videoID, commentID, action)
	v, err := json.Marshal(event)
	if err != nil {
		logger.Errorf("marshal like event failed: %v", err)
		return
	}
	err = KafkaInstance.Send(ctx, constants.LikeTopic, []*kafka.Message{
		{
			K: []byte(strconv.FormatInt(userID, 10)),
			V: v,
		},
	})
	if err != nil {
		logger.Errorf("kafka send msg failed: %v", err)
	}
}
