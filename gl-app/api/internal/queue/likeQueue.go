package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	interactiverepo "gl-app/api/internal/repo/interactive_repo"

	"github.com/segmentio/kafka-go"
)

// LikeQueue 使用 kafka-go 异步 Writer。HTTP 请求只把事件交给生产者缓冲区，
// 不等待 Broker 确认，降低高频点赞操作的响应延迟。
type LikeQueue struct{ writer *kafka.Writer }

func NewLikeQueue(brokers []string, topic string) *LikeQueue {
	return &LikeQueue{writer: &kafka.Writer{
		Addr: kafka.TCP(brokers...), Topic: topic, Async: true,
		Balancer: &kafka.Hash{}, BatchTimeout: 10 * time.Millisecond,
		AllowAutoTopicCreation: true,
	}}
}

func (q *LikeQueue) Publish(ctx context.Context, event interactiverepo.LikeEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("序列化点赞事件失败: %w", err)
	}
	key := strconv.FormatInt(event.UserID, 10) + ":" + string(event.ResourceType) + ":" + strconv.FormatInt(event.ResourceID, 10)
	if err = q.writer.WriteMessages(ctx, kafka.Message{Key: []byte(key), Value: value}); err != nil {
		return fmt.Errorf("异步投递点赞事件失败: %w", err)
	}
	return nil
}

func DecodeLikeEvent(value string) (interactiverepo.LikeEvent, error) {
	var event interactiverepo.LikeEvent
	if err := json.Unmarshal([]byte(value), &event); err != nil {
		return event, fmt.Errorf("解析点赞事件失败: %w", err)
	}
	if event.EventID == "" || event.UserID <= 0 || event.ResourceID <= 0 || !event.ResourceType.Valid() || event.Version <= 0 {
		return event, fmt.Errorf("点赞事件字段不合法")
	}
	return event, nil
}

func (q *LikeQueue) Close() error { return q.writer.Close() }
