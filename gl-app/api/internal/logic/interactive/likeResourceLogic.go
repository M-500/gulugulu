package interactive

import (
	"context"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeResourceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 点赞作品、评论或头像
func NewLikeResourceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeResourceLogic {
	return &LikeResourceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeResourceLogic) LikeResource(req *types.LikeResourceReq) (resp *types.LikeResourceResp, err error) {
	// todo: add your logic here and delete this line

	return
}
