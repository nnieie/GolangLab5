package kafka

import (
	"context"
	"encoding/json"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/kafka"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/tracer"
)

func SendLikeEvent(ctx context.Context, userID int64, videoID, commentID *int64, action int64) {
	logger.Debugf("SendLikeEvent called: userID=%d, videoID=%v, commentID=%v, action=%d", userID, videoID, commentID, action)
	event := NewLikeEvent(userID, videoID, commentID, action)
	v, err := json.Marshal(event)
	if err != nil {
		logger.Errorf("marshal like event failed: %v", err)
		return
	}
	logger.Debugf("Sending like event to Kafka: %s", string(v))
	err = KafkaInstance.Send(context.Background(), constants.LikeTopic, []*kafka.Message{
		{
			K: []byte(strconv.FormatInt(userID, 10)),
			V: v,
		},
	})
	if err != nil {
		logger.Errorf("kafka send msg failed: %v", err)
		tracer.MQProduceCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("topic", constants.LikeTopic),
			attribute.String("status", "fail"),
		))
	} else {
		logger.Infof("Successfully sent like event to Kafka")
		tracer.MQProduceCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("topic", constants.LikeTopic),
			attribute.String("status", "success"),
		))
	}
}
