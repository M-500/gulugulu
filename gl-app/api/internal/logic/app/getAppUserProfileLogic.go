package app

import (
	"context"
	"fmt"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppUserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取App用户公开资料
func NewGetAppUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppUserProfileLogic {
	return &GetAppUserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppUserProfileLogic) GetAppUserProfile(req *types.AppUserProfileReq) (resp *types.AppUserProfileResp, err error) {
	if req == nil || req.UserId <= 0 {
		return nil, fmt.Errorf("用户ID不正确")
	}

	userInfo, err := findAppPublicUser(l.ctx, l.svcCtx, req.UserId)
	if err != nil {
		return nil, err
	}
	res := types.AppUserProfileResp{
		UserId:    userInfo.ID,
		NickName:  normalizeAuthorName(userInfo.Nickname),
		Bio:       userInfo.Bio,
		Sex:       userInfo.Sex,
		IpAddress: userInfo.IPAddress,
		AvatarUrl: buildAppPublicAvatarURL(l.svcCtx, userInfo.Avatar),
	}
	return &res, nil
}
