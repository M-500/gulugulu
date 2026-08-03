package app

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	userrepo "gl-app/api/internal/repo/user_repo"
	"gl-app/api/internal/svc"
)

// findAppPublicUser 查询公开主页需要的用户记录，并统一屏蔽已删除用户。
func findAppPublicUser(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) (*userrepo.User, error) {
	userInfo, err := svcCtx.UserRepo.FindOneByID(ctx, userID)
	if errors.Is(err, userrepo.ErrNotFound) {
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
