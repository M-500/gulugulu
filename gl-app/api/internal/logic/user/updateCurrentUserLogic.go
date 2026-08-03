package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	userrepo "gl-app/api/internal/repo/user_repo"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCurrentUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改当前用户资料
func NewUpdateCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCurrentUserLogic {
	return &UpdateCurrentUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCurrentUserLogic) UpdateCurrentUser(req *types.UpdateUserProfileReq) (resp *types.UserProfileResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID <= 0 {
		return nil, fmt.Errorf("用户未登录")
	}
	if req == nil {
		return nil, fmt.Errorf("用户资料不能为空")
	}

	updateMap, err := buildProfileUpdate(req, time.Now())
	if err != nil {
		return nil, err
	}
	if _, err = l.svcCtx.UserRepo.FindOneByID(l.ctx, userID); err != nil {
		if errors.Is(err, userrepo.ErrNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户资料失败: %w", err)
	}

	if err = l.svcCtx.UserRepo.UpdateByMap(l.ctx, userID, updateMap); err != nil {
		return nil, fmt.Errorf("更新用户资料失败: %w", err)
	}
	return QueryUserProfile(l.ctx, l.svcCtx, userID)
}

// buildProfileUpdate 集中完成资料校验和字段转换，避免把字符串直接写入数值或日期列。
func buildProfileUpdate(req *types.UpdateUserProfileReq, now time.Time) (map[string]any, error) {
	nickname := strings.TrimSpace(req.NickName)
	length := utf8.RuneCountInString(nickname)
	if length < 2 || length > 32 {
		return nil, fmt.Errorf("昵称长度必须是2到32个字符")
	}

	bio := strings.TrimSpace(req.Bio)
	if utf8.RuneCountInString(bio) > 200 {
		return nil, fmt.Errorf("个性签名不能超过200个字符")
	}
	if req.Sex < 0 || req.Sex > 2 {
		return nil, fmt.Errorf("性别参数不正确")
	}

	var bothDay *time.Time
	if value := strings.TrimSpace(req.BothDay); value != "" {
		parsed, err := time.ParseInLocation(time.DateOnly, value, now.Location())
		if err != nil {
			return nil, fmt.Errorf("生日格式必须为YYYY-MM-DD")
		}
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if parsed.After(today) {
			return nil, fmt.Errorf("生日不能晚于今天")
		}
		bothDay = &parsed
	}

	// 使用 map 更新以确保 sex=0、bio="" 和 both_day=NULL 等零值也能正确落库。
	return map[string]any{
		"nickname": nickname,
		"bio":      bio,
		"sex":      req.Sex,
		"both_day": bothDay,
	}, nil
}
