package interactive

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gl-app/api/internal/constants"
	interactiverepo "gl-app/api/internal/repo/interactive_repo"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type MutateLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMutateLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MutateLikeLogic {
	return &MutateLikeLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *MutateLikeLogic) Mutate(req *types.LikeResourceReq, liked bool) (*types.LikeResourceResp, error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID <= 0 {
		return nil, fmt.Errorf("请先登录")
	}
	if req == nil || req.ResourceId <= 0 {
		return nil, fmt.Errorf("资源ID不正确")
	}
	bizType, err := parseBizType(req.ResourceType)
	if err != nil {
		return nil, err
	}

	mutation, err := l.svcCtx.InteractiveRepo.MutateLike(l.ctx, userID, req.ResourceId, bizType, liked)
	if err != nil {
		if errors.Is(err, interactiverepo.ErrResourceNotFound) {
			return nil, fmt.Errorf("资源不存在")
		}
		l.Errorf("更新Redis点赞状态失败: %v", err)
		return nil, fmt.Errorf("点赞操作失败，请稍后重试")
	}
	// 重复请求也会补发同版本事件：Kafka/网络瞬时故障后，客户端重试可触发消费端对账；
	// 消费端通过 user_like.version 保证重复消息不会重复累计。
	if mutation.Version > 0 {
		event := interactiverepo.LikeEvent{
			EventID: uuid.NewString(), UserID: userID, ResourceID: req.ResourceId,
			ResourceType: bizType, Liked: liked, Version: mutation.Version, OccurredAt: time.Now(),
		}
		if err = l.svcCtx.LikeQueue.Publish(l.ctx, event); err != nil {
			l.Errorf("投递点赞事件失败: %v", err)
			return nil, fmt.Errorf("点赞操作提交失败，请稍后重试")
		}
	}
	return &types.LikeResourceResp{ResourceType: strings.ToLower(req.ResourceType), ResourceId: req.ResourceId, Liked: liked, Count: mutation.Count}, nil
}

func parseBizType(value string) (constants.BizType, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "work":
		return constants.WorkType, nil
	case "comment":
		return constants.CommentType, nil
	case "avatar":
		return constants.AvatarType, nil
	default:
		return "", fmt.Errorf("资源类型仅支持 work、comment、avatar")
	}
}
