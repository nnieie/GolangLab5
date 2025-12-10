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

const (
	batchSize    = 100                   // 批量发送大小
	batchTimeout = 10 * time.Millisecond // 批量发送超时时间
)

// Producer Kafka 生产者
type Producer struct {
	writer *kafka.Writer
}

// ProducerConfig Kafka 生产者配置
type ProducerConfig struct {
	Brokers []string // Kafka broker 地址列表
	Topic   string   // 主题名称
}

var defaultProducer *Producer

// InitProducer 初始化 Kafka 生产者
func InitProducer(config ProducerConfig) error {
	if len(config.Brokers) == 0 {
		return fmt.Errorf("kafka brokers cannot be empty")
	}
	if config.Topic == "" {
		return fmt.Errorf("kafka topic cannot be empty")
	}

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

	defaultProducer = &Producer{
		writer: writer,
	}

	logger.Infof("Kafka producer initialized successfully, brokers: %v, topic: %s", config.Brokers, config.Topic)
	return nil
}

// GetProducer 获取默认生产者实例
func GetProducer() *Producer {
	return defaultProducer
}

// Close 关闭生产者
func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

// PublishLikeEvent 发布点赞事件
func (p *Producer) PublishLikeEvent(ctx context.Context, event *LikeEvent) error {
	return p.publishEvent(ctx, event)
}

// publishEvent 发布事件到 Kafka
func (p *Producer) publishEvent(ctx context.Context, event interface{}) error {
	if p.writer == nil {
		return fmt.Errorf("kafka producer not initialized")
	}

	// 开启一个 Span，记录 发送消息 这个动作
	tr := otel.Tracer("kafka-producer")
	ctx, span := tr.Start(ctx, "publish_event", trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	// 序列化事件
	data, err := json.Marshal(event)
	if err != nil {
		logger.Errorf("Failed to marshal event: %v", err)
		return fmt.Errorf("marshal event failed: %w", err)
	}

	// 注入 TraceID 到 Kafka Headers
	// 创建一个 map 来承载 Trace 信息
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	// 把 map 转换成 Kafka 的 Headers 格式
	var headers []kafka.Header
	for k, v := range carrier {
		headers = append(headers, kafka.Header{Key: k, Value: []byte(v)})
	}

	// 构造 Kafka 消息
	msg := kafka.Message{
		Value:   data,
		Time:    time.Now(),
		Headers: headers,
	}

	// 发送消息
	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		// 记录 Error 到 Span
		span.RecordError(err)
		logger.Errorf("Failed to write message to kafka: %v", err)
		return fmt.Errorf("write message failed: %w", err)
	}

	return nil
}

// PublishEventWithRetry 带重试机制的发布事件
func (p *Producer) PublishEventWithRetry(ctx context.Context, event interface{}, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = p.publishEvent(ctx, event)
		if err == nil {
			return nil
		}
		logger.Warnf("Failed to publish event (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond) // 指数退避
	}
	return fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}
