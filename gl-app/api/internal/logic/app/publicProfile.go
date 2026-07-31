package app

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"gl-app/api/internal/models/user"
	"gl-app/api/internal/svc"
)

// findAppPublicUser 查询公开主页需要的用户记录，并统一屏蔽已删除用户。
func findAppPublicUser(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) (*user.User, error) {
	userInfo, err := svcCtx.UserRepo.FindOne(ctx, userID)
	if errors.Is(err, user.ErrNotFound) || (err == nil && userInfo.DeletedAt.Valid) {
		return nil, fmt.Errorf("用户不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户资料失败")
	}
	return userInfo, nil
}

// buildAppPublicAvatarURL 将数据库中的头像对象Key转换为前端可直接访问的公共地址。
func buildAppPublicAvatarURL(svcCtx *svc.ServiceContext, objectKey string) string {
	objectKey = strings.TrimSpace(objectKey)
	if objectKey == "" {
		return ""
	}

	bucket := svcCtx.Config.Minio.PublicBucket
	if bucket == "" {
		bucket = svcCtx.Config.Minio.FormalBucket
	}
	scheme := "http"
	if svcCtx.Config.Minio.UseSSL {
		scheme = "https"
	}
	endpoint := strings.TrimPrefix(strings.TrimPrefix(svcCtx.Config.Minio.Endpoint, "http://"), "https://")
	value := url.URL{
		Scheme: scheme,
		Host:   endpoint,
		Path:   path.Join(bucket, objectKey),
	}
	return value.String()
}
