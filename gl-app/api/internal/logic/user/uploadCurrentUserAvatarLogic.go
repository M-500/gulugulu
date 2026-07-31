package user

import (
	"context"
	"fmt"
	"mime/multipart"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadCurrentUserAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 上传当前用户头像
func NewUploadCurrentUserAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCurrentUserAvatarLogic {
	return &UploadCurrentUserAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCurrentUserAvatarLogic) UploadCurrentUserAvatar(
	file multipart.File, header *multipart.FileHeader,
) (resp *types.UserProfileResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID <= 0 {
		return nil, fmt.Errorf("用户未登录")
	}
	objectKey, err := uploadAvatarToPublicBucket(l.ctx, l.svcCtx, userID, file, header)
	if err != nil {
		l.Errorf("上传当前用户头像失败：%v", err)
		return nil, err
	}
	if err = saveAvatarObjectKey(l.ctx, l.svcCtx, userID, objectKey); err != nil {
		l.Errorf("保存当前用户头像失败：%v", err)
		return nil, err
	}
	return QueryUserProfile(l.ctx, l.svcCtx, userID)
}
