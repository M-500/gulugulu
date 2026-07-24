package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/zeromicro/go-queue/kq"
)

const ProcessWorkMessageType = "process_work"

type Message struct {
	Type    string `json:"type"`
	WorkID  int64  `json:"workId"`
	Attempt int    `json:"attempt"`
}

type MediaQueue struct {
	pusher *kq.Pusher
}

func NewMediaQueue(brokers []string, topic string) *MediaQueue {
	return &MediaQueue{
		pusher: kq.NewPusher(
			brokers,
			topic,
			kq.WithSyncPush(),
			kq.WithAllowAutoTopicCreation(),
		),
	}
}

func (q *MediaQueue) PublishProcessWork(ctx context.Context, workID int64, attempt int) error {
	message, err := json.Marshal(Message{
		Type:    ProcessWorkMessageType,
		WorkID:  workID,
		Attempt: attempt,
	})
	if err != nil {
		return fmt.Errorf("序列化媒体处理消息失败: %w", err)
	}
	if err = q.pusher.PushWithKey(ctx, strconv.FormatInt(workID, 10), string(message)); err != nil {
		return fmt.Errorf("投递Kafka媒体处理任务失败: %w", err)
	}
	return nil
}

func DecodeMessage(value string) (Message, error) {
	var message Message
	if err := json.Unmarshal([]byte(value), &message); err != nil {
		return Message{}, fmt.Errorf("解析Kafka媒体处理消息失败: %w", err)
	}
	if message.Type != ProcessWorkMessageType || message.WorkID <= 0 || message.Attempt < 0 {
		return Message{}, fmt.Errorf("Kafka媒体处理消息字段不合法")
	}
	return message, nil
}

func (q *MediaQueue) Close() error {
	return q.pusher.Close()
}
