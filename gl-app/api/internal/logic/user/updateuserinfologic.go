package user

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改当前用户基本信息
func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(
	nickname string,
	avatarFile multipart.File,
	avatarHeader *multipart.FileHeader,
) (resp *types.UserProfileResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID <= 0 {
		return nil, fmt.Errorf("用户未登录")
	}
	hasNickname := strings.TrimSpace(nickname) != ""
	hasAvatar := avatarFile != nil && avatarHeader != nil
	if !hasNickname && !hasAvatar {
		return nil, fmt.Errorf("请至少修改昵称或头像")
	}

	if hasNickname {
		if err = updateNickname(l.ctx, l.svcCtx, userID, nickname); err != nil {
			l.Errorf("修改用户昵称失败：%v", err)
			return nil, err
		}
	}

	if hasAvatar {
		objectKey, uploadErr := uploadAvatarToPublicBucket(l.ctx, l.svcCtx, userID, avatarFile, avatarHeader)
		if uploadErr != nil {
			l.Errorf("上传用户头像失败：%v", uploadErr)
			return nil, uploadErr
		}
		if err = saveAvatarObjectKey(l.ctx, l.svcCtx, userID, objectKey); err != nil {
			l.Errorf("保存用户头像失败：%v", err)
			return nil, err
		}
	}

	return QueryUserProfile(l.ctx, l.svcCtx, userID)
}
