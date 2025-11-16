package kafka

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/nnieie/golanglab5/cmd/interaction/dal/db"
	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/pkg/kafka"
	"github.com/nnieie/golanglab5/pkg/logger"
)

const (
	bufferSize    = 100
	flushInterval = 5 * time.Second
)

func StartInteractionConsumer(ctx context.Context) {
	// 创建批量处理器
	// 参数：bufferSize=100（累积100条），flushInterval=5s（最多等待5秒）
	batchProcessor := NewBatchProcessor(bufferSize, flushInterval)
	defer batchProcessor.Close()

	// 创建 Kafka 消费者
	consumer, err := kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers: config.Kafka.Brokers,
		Topic:   config.Kafka.Topic,
		GroupID: "like-persistence-group",
	})
	if err != nil {
		logger.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// 使用批量处理器
	go func() {
		logger.Infof("Starting Kafka consumer with batch processing...")
		err := consumer.ConsumeLikeEvents(ctx, handleLikeEventBatch(batchProcessor))
		if err != nil && !errors.Is(err, context.Canceled) {
			logger.Errorf("Consumer error: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Infof("Shutting down...")
}

// BatchProcessor 批量处理器
type BatchProcessor struct {
	buffer        []*kafka.LikeEvent // 消息缓冲区
	bufferSize    int                // 批量大小
	flushInterval time.Duration      // 刷新间隔
	mu            sync.Mutex         // 并发保护
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewBatchProcessor 创建批量处理器
func NewBatchProcessor(bufferSize int, flushInterval time.Duration) *BatchProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	bp := &BatchProcessor{
		buffer:        make([]*kafka.LikeEvent, 0, bufferSize),
		bufferSize:    bufferSize,
		flushInterval: flushInterval,
		ctx:           ctx,
		cancel:        cancel,
	}

	// 启动定时刷新 goroutine
	go bp.periodicFlush()

	return bp
}

// AddEvent 添加事件到缓冲区
func (bp *BatchProcessor) AddEvent(event *kafka.LikeEvent) error {
	bp.mu.Lock()
	bp.buffer = append(bp.buffer, event)
	shouldFlush := len(bp.buffer) >= bp.bufferSize
	bp.mu.Unlock()

	// 达到批量大小，立即刷新
	if shouldFlush {
		return bp.Flush()
	}

	return nil
}

// periodicFlush 定时刷新（防止消息积压）
func (bp *BatchProcessor) periodicFlush() {
	ticker := time.NewTicker(bp.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-bp.ctx.Done():
			// 退出前刷新剩余消息
			bp.Flush()
			return
		case <-ticker.C:
			bp.Flush()
		}
	}
}

// Flush 批量刷新到数据库
func (bp *BatchProcessor) Flush() error {
	bp.mu.Lock()
	if len(bp.buffer) == 0 {
		bp.mu.Unlock()
		return nil
	}

	// 复制缓冲区，快速释放锁
	batch := make([]*kafka.LikeEvent, len(bp.buffer))
	copy(batch, bp.buffer)
	bp.buffer = bp.buffer[:0] // 清空缓冲区
	bp.mu.Unlock()

	logger.Infof("Flushing %d like events to database", len(batch))

	// 分组处理：点赞 vs 取消点赞
	var likesToInsert []db.Like
	var likesToDelete []struct {
		UserID   int64
		TargetID int64
		Type     int64
	}

	for _, event := range batch {
		var targetID, likeType int64

		switch {
		case event.VideoID != nil:
			targetID = *event.VideoID
			likeType = db.VideoLikeType
		case event.CommentID != nil:
			targetID = *event.CommentID
			likeType = db.CommentLikeType
		default:
			continue
		}

		switch event.Action {
		case 1: // 点赞
			likesToInsert = append(likesToInsert, db.Like{
				UserID:   event.UserID,
				TargetID: targetID,
				Type:     likeType,
			})
		case 2: // 取消点赞
			likesToDelete = append(likesToDelete, struct {
				UserID   int64
				TargetID int64
				Type     int64
			}{event.UserID, targetID, likeType})
		}
	}

	ctx := context.Background()

	// 批量插入点赞
	if len(likesToInsert) > 0 {
		if err := db.BatchLikeAction(ctx, likesToInsert); err != nil {
			logger.Errorf("Failed to batch insert likes: %v", err)
			return err
		}
		logger.Infof("Successfully inserted %d likes", len(likesToInsert))
	}

	// 批量删除点赞
	if len(likesToDelete) > 0 {
		if err := db.BatchUnlikeAction(ctx, likesToDelete); err != nil {
			logger.Errorf("Failed to batch delete likes: %v", err)
			return err
		}
		logger.Infof("Successfully deleted %d likes", len(likesToDelete))
	}

	return nil
}

// Close 关闭处理器
func (bp *BatchProcessor) Close() error {
	bp.cancel()
	return bp.Flush() // 最后一次刷新
}

// handleLikeEventBatch 批量处理
func handleLikeEventBatch(processor *BatchProcessor) func(*kafka.LikeEvent) error {
	return func(event *kafka.LikeEvent) error {
		return processor.AddEvent(event)
	}
}
