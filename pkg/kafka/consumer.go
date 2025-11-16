package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

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
func (c *Consumer) ConsumeMessages(ctx context.Context, handler func([]byte) error) error {
	for {
		select {
		case <-ctx.Done():
			logger.Infof("Consumer context canceled, stopping...")
			return ctx.Err()
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				logger.Errorf("Failed to read message: %v", err)
				continue
			}

			logger.Infof("Received message: topic=%s, partition=%d, offset=%d",
				msg.Topic, msg.Partition, msg.Offset)

			// 处理消息
			if err := handler(msg.Value); err != nil {
				logger.Errorf("Failed to handle message: %v", err)
				// 这里可以实现重试逻辑或将消息发送到死信队列
				continue
			}
		}
	}
}

// ConsumeLikeEvents 消费点赞事件
func (c *Consumer) ConsumeLikeEvents(ctx context.Context, handler func(*LikeEvent) error) error {
	return c.ConsumeMessages(ctx, func(data []byte) error {
		var event LikeEvent
		if err := json.Unmarshal(data, &event); err != nil {
			return fmt.Errorf("unmarshal like event failed: %w", err)
		}
		return handler(&event)
	})
}
