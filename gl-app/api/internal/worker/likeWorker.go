package worker

import (
	"context"
	"sync"

	"gl-app/api/internal/queue"
	"gl-app/api/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	corequeue "github.com/zeromicro/go-zero/core/queue"
	"github.com/zeromicro/go-zero/core/service"
)

type LikeWorker struct {
	svcCtx   *svc.ServiceContext
	queue    corequeue.MessageQueue
	stopOnce sync.Once
}

var _ service.Service = (*LikeWorker)(nil)

func NewLikeWorker(svcCtx *svc.ServiceContext) *LikeWorker {
	w := &LikeWorker{svcCtx: svcCtx}
	w.queue = kq.MustNewQueue(svcCtx.Config.LikeQueue.KqConf, w)
	return w
}

func (w *LikeWorker) Start() {
	logx.Info("Kafka点赞消费者已启动")
	w.queue.Start()
}

func (w *LikeWorker) Stop() { w.stopOnce.Do(w.queue.Stop) }

func (w *LikeWorker) Consume(ctx context.Context, _ string, value string) error {
	event, err := queue.DecodeLikeEvent(value)
	if err != nil {
		// 非法消息重试也无法恢复，记录后提交 offset，避免阻塞整个分区。
		logx.WithContext(ctx).Errorf("忽略非法Kafka点赞消息: %v", err)
		return nil
	}
	if err = w.svcCtx.InteractiveRepo.ApplyLikeEvent(ctx, event); err != nil {
		logx.WithContext(ctx).Errorf("消费点赞事件失败(eventId=%s): %v", event.EventID, err)
		return err
	}
	return nil
}
