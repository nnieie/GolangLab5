package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/nnieie/golanglab5/pkg/logger"
)

// Consumer Kafka 消费者
type Consumer struct {
	reader *kafka.Reader
}

// ConsumerConfig Kafka 消费者配置
type ConsumerConfig struct {
	Brokers  []string // Kafka broker 地址列表
	Topic    string   // 主题名称
	GroupID  string   // 消费者组ID
	MinBytes int      // 最小读取字节数
	MaxBytes int      // 最大读取字节数
}

// NewConsumer 创建 Kafka 消费者
func NewConsumer(config ConsumerConfig) (*Consumer, error) {
	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("kafka brokers cannot be empty")
	}
	if config.Topic == "" {
		return nil, fmt.Errorf("kafka topic cannot be empty")
	}
	if config.GroupID == "" {
		return nil, fmt.Errorf("kafka group id cannot be empty")
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          config.Topic,
		GroupID:        config.GroupID,
		MinBytes:       config.MinBytes,
		MaxBytes:       config.MaxBytes,
		CommitInterval: time.Second,      // 自动提交间隔
		StartOffset:    kafka.LastOffset, // 从最新位置开始消费
	})

	logger.Infof("Kafka consumer initialized, brokers: %v, topic: %s, group: %s",
		config.Brokers, config.Topic, config.GroupID)

	return &Consumer{
		reader: reader,
	}, nil
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}

// ConsumeMessages 消费消息
func (c *Consumer) ConsumeMessages(ctx context.Context, handler func(context.Context, []byte) error) error {
	for {
		select {
		case <-ctx.Done():
			logger.Infof("Consumer context canceled, stopping...")
			return ctx.Err()
		default:
			// 读取消息
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				logger.Errorf("Failed to read message: %v", err)
				continue
			}

			// 提取 TraceID 把 Kafka Headers 转换回 MapCarrier
			carrier := propagation.MapCarrier{}
			for _, h := range msg.Headers {
				carrier[h.Key] = string(h.Value)
			}

			// 从 carrier 里提取出父级 Context
			parentCtx := otel.GetTextMapPropagator().Extract(context.Background(), carrier)

			// 开启一个新的 Consumer Span
			tr := otel.Tracer("kafka-consumer")
			// 使用提取出来的 parentCtx 启动 Span，这样就能连上 Producer 了
			spanName := fmt.Sprintf("consume %s", msg.Topic)
			spanCtx, span := tr.Start(parentCtx, spanName, trace.WithSpanKind(trace.SpanKindConsumer))

			logger.Infof("Received message: topic=%s, partition=%d, offset=%d",
				msg.Topic, msg.Partition, msg.Offset)

			// 处理消息 (把带有 TraceID 的 spanCtx 传给业务逻辑)
			if err := handler(spanCtx, msg.Value); err != nil {
				span.RecordError(err) // 记录错误到 Jaeger
				logger.Errorf("Failed to handle message: %v", err)
				span.End() // 结束 Span
				continue
			}

			span.End() // 处理成功，结束 Span
		}
	}
}

// ConsumeLikeEvents 消费点赞事件
func (c *Consumer) ConsumeLikeEvents(ctx context.Context, handler func(context.Context, *LikeEvent) error) error {
	// 这里的回调函数接收到了 spanCtx
	return c.ConsumeMessages(ctx, func(spanCtx context.Context, data []byte) error {
		var event LikeEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("unmarshal like event failed: %w", err)
		}
		// 把 spanCtx 继续透传给具体的业务 Handler
		return handler(spanCtx, &event)
	})
}
