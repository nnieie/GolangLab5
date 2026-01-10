package kafka

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
)

const (
	batchSize    = 100                   // 批量发送大小
	batchTimeout = 10 * time.Millisecond // 批量发送超时时间
)

type Kafka struct {
	readers      []*kafka.Reader
	writers      map[string]*kafka.Writer
	consumeChans map[string]chan *Message
}

type ConsumerConfig struct {
	Brokers  []string // Kafka broker 地址列表
	Topic    string   // 主题名称
	GroupID  string   // 消费者组ID
	MinBytes int      // 最小读取字节数
	MaxBytes int      // 最大读取字节数
}

type ProducerConfig struct {
	Brokers []string // Kafka broker 地址列表
	Topic   string   // 主题名称
}

type Message struct {
	K, V []byte
}

// NewKafkaInstance 返回一个新的 kafka 实例
func NewKafkaInstance() *Kafka {
	return &Kafka{
		readers:      make([]*kafka.Reader, 0),
		writers:      make(map[string]*kafka.Writer),
		consumeChans: make(map[string]chan *Message),
	}
}

// Consume 获取消息并将消息通过 channel 传递出去
func (k *Kafka) Consume(ctx context.Context, topic string, groupID string, chanCap ...int) <-chan *Message {
	if k.consumeChans[topic] != nil {
		return k.consumeChans[topic]
	}

	chCap := constants.DefaultConsumerChanCap
	if chanCap != nil {
		chCap = chanCap[0]
	}
	ch := make(chan *Message, chCap)
	k.consumeChans[topic] = ch

	readers := getNewReader(ConsumerConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	k.readers = append(k.readers, readers)
	go k.consume(ctx, topic, readers)
	return ch
}

func (k *Kafka) consume(ctx context.Context, topic string, r *kafka.Reader) {
	ch := k.consumeChans[topic]
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			logger.Errorf("read message from kafka reader failed,err: %v", err.Error())
			return
		}

		ch <- &Message{K: msg.Key, V: msg.Value}
	}
}

// Send 发送消息到指定的 topic
func (k *Kafka) Send(ctx context.Context, topic string, messages []*Message) error {
	if k.writers[topic] == nil {
		if err := k.SetWriter(topic); err != nil {
			return err
		}
	}

	return k.send(ctx, topic, messages)
}

// SetWriter 生成 writer 并存储到 map 中
func (k *Kafka) SetWriter(topic string) error {
	w := getNewWriter(ProducerConfig{
		Topic:   topic,
		Brokers: config.Kafka.Brokers,
	})
	k.writers[topic] = w
	return nil
}

func (k *Kafka) send(ctx context.Context, topic string, messages []*Message) error {
	msgs := make([]kafka.Message, 0, len(messages))
	for _, m := range messages {
		msgs = append(msgs, kafka.Message{
			Key:   m.K,
			Value: m.V,
		})
	}

	err := k.writers[topic].WriteMessages(ctx, msgs...)
	return err
}

func (k *Kafka) Close() {
	for _, reader := range k.readers {
		if err := reader.Close(); err != nil {
			logger.Errorf("close kafka reader failed, err: %v", err)
		}
	}

	for _, writer := range k.writers {
		if err := writer.Close(); err != nil {
			logger.Errorf("close kafka writer failed, err: %v", err)
		}
	}

	for _, ch := range k.consumeChans {
		close(ch)
	}
}

func getNewReader(config ConsumerConfig) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  config.Brokers,
		Topic:    config.Topic,
		GroupID:  config.GroupID,
		MinBytes: config.MinBytes,
		MaxBytes: config.MaxBytes,
	})
	return reader
}

func getNewWriter(config ProducerConfig) *kafka.Writer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{}, // 负载均衡策略
		BatchSize:    batchSize,           // 批量发送大小
		BatchTimeout: batchTimeout,
		RequiredAcks: kafka.RequireOne, // 至少一个副本确认
		Async:        false,            // 同步发送
		Compression:  kafka.Snappy,     // 压缩算法
	}
	return writer
}
