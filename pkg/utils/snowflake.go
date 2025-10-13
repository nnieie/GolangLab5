package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	// 起始时间戳
	Epoch int64 = 1145141919
	// workerIDBits 工作节点 ID 的位数
	workerIDBits uint8 = 10
	// sequenceBits 每个节点每毫秒生成的序列号的位数
	sequenceBits uint8 = 12
	// workerIDMax workerIDBits 位数所能表示的最大值 (1023)
	workerIDMax int64 = -1 ^ (-1 << workerIDBits)
	// sequenceMask 获取序列号的低 12 位
	sequenceMask int64 = -1 ^ (-1 << sequenceBits)
	// workerIDShift 工作节点 ID 左移的位数
	workerIDShift uint8 = sequenceBits
	// timestampShift 时间戳左移的位数
	timestampShift uint8 = sequenceBits + workerIDBits
)

type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	workerID      int64
	sequence      int64
}

func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > workerIDMax {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", workerIDMax-1)
	}
	return &Snowflake{
		workerID: workerID,
	}, nil
}

// Generate 生成一个唯一的雪花 ID
func (n *Snowflake) Generate() (int64, error) {
	// 加锁，保证并发安全
	n.mu.Lock()
	defer n.mu.Unlock()

	// 获取当前时间的毫秒数
	now := time.Now().UnixMilli()

	// 如果当前时间小于上次生成 ID 的时间，说明发生了时钟回拨
	if now < n.lastTimestamp {
		return 0, errors.New("clock moved backwards, refusing to generate id")
	}

	// 如果是同一毫秒内生成的
	if n.lastTimestamp == now {
		// 序列号加 1，然后与 sequenceMask 进行位与操作，防止溢出
		n.sequence = (n.sequence + 1) & sequenceMask
		// 如果序列号溢出（回到 0），则表示当前毫秒的 4096 个 ID 已用完
		if n.sequence == 0 {
			// 等待下一毫秒
			for now <= n.lastTimestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 如果是新的毫秒，则序列号重置为 0
		n.sequence = 0
	}

	// 更新最后的时间戳
	n.lastTimestamp = now

	id := ((now - Epoch) << timestampShift) |
		(n.workerID << workerIDShift) |
		n.sequence

	return id, nil
}
