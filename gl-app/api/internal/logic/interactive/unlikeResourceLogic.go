package interactive

import (
	"context"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnlikeResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 取消点赞作品、评论或头像
func NewUnlikeResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikeResourceLogic {
	return &UnlikeResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnlikeResourceLogic) UnlikeResource(req *types.LikeResourceReq) (resp *types.LikeResourceResp, err error) {
	// todo: add your logic here and delete this line

	return
}
