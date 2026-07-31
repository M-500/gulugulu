package user

import (
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
)

type profileRow struct {
	UserID   int64  `db:"id"`
	Email    string `db:"email"`
	NickName string `db:"nickname"`
	Avatar   string `db:"avatar"`
}

// QueryUserProfile 统一查询用户基础资料，登录态初始化和个人资料页都走这一套返回结构。
func QueryUserProfile(ctx context.Context, svcCtx *svc.ServiceContext, userID int64) (*types.UserProfileResp, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("用户未登录")
	}
	var row profileRow
	if err := svcCtx.SqlConn.QueryRowCtx(ctx, &row,
		"SELECT id,email,nickname,avatar FROM user WHERE id=? AND deleted_at IS NULL LIMIT 1", userID); err != nil {
		return nil, fmt.Errorf("查询用户资料失败: %w", err)
	}
	resp := &types.UserProfileResp{
		UserId: row.UserID, Email: row.Email, NickName: row.NickName,
	}
	if row.Avatar != "" {
		resp.AvatarUrl = publicObjectURL(svcCtx, publicAvatarBucket(svcCtx), row.Avatar)
	}
	return resp, nil
}

func updateNickname(ctx context.Context, svcCtx *svc.ServiceContext, userID int64, nickname string) error {
	nickname = strings.TrimSpace(nickname)
	length := utf8.RuneCountInString(nickname)
	if length < 2 || length > 32 {
		return fmt.Errorf("昵称长度必须是2到32个字符")
	}
	if _, err := svcCtx.SqlConn.ExecCtx(ctx,
		"UPDATE user SET nickname=? WHERE id=? AND deleted_at IS NULL", nickname, userID); err != nil {
		return fmt.Errorf("更新用户昵称失败: %w", err)
	}
	return nil
}

func uploadAvatarToPublicBucket(
	ctx context.Context,
	svcCtx *svc.ServiceContext,
	userID int64,
	file multipart.File,
	header *multipart.FileHeader,
) (string, error) {
	if file == nil || header == nil {
		return "", nil
	}
	if userID <= 0 {
		return "", fmt.Errorf("用户未登录")
	}
	if header.Size <= 0 {
		return "", fmt.Errorf("头像文件不能为空")
	}
	ext, contentType, err := normalizeAvatarFile(header.Filename)
	if err != nil {
		return "", err
	}
	bucket := publicAvatarBucket(svcCtx)
	if err = ensurePublicBucket(ctx, svcCtx, bucket); err != nil {
		return "", err
	}
	objectKey := fmt.Sprintf("users/%d/avatar/%s%s", userID, uuid.NewString(), ext)
	if _, err = svcCtx.MinioClient.PutObject(
		ctx, bucket, objectKey, file, header.Size, minio.PutObjectOptions{ContentType: contentType},
	); err != nil {
		return "", fmt.Errorf("上传头像到公共桶失败: %w", err)
	}
	return objectKey, nil
}

func saveAvatarObjectKey(ctx context.Context, svcCtx *svc.ServiceContext, userID int64, objectKey string) error {
	if objectKey == "" {
		return nil
	}
	if _, err := svcCtx.SqlConn.ExecCtx(ctx,
		"UPDATE user SET avatar=? WHERE id=? AND deleted_at IS NULL", objectKey, userID); err != nil {
		return fmt.Errorf("保存用户头像失败: %w", err)
	}
	return nil
}

func normalizeAvatarFile(fileName string) (ext string, contentType string, err error) {
	ext = strings.ToLower(filepath.Ext(fileName))
	allowed := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".avif": "image/avif",
		".heic": "image/heic",
		".heif": "image/heif",
	}
	contentType, ok := allowed[ext]
	if !ok {
		return "", "", fmt.Errorf("头像仅支持常见图片格式")
	}
	if value := mime.TypeByExtension(ext); value != "" {
		contentType = value
	}
	return ext, contentType, nil
}

func ensurePublicBucket(ctx context.Context, svcCtx *svc.ServiceContext, bucket string) error {
	if svcCtx.MinioClient == nil {
		return fmt.Errorf("MinIO客户端未初始化")
	}
	exists, err := svcCtx.MinioClient.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("检查MinIO公共桶失败: %w", err)
	}
	if !exists {
		if err = svcCtx.MinioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("创建MinIO公共桶失败: %w", err)
		}
	}
	// 头像桶允许公开读取，前端可以直接用稳定URL渲染头像，不需要每次生成预签名地址。
	policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}`, bucket)
	if err = svcCtx.MinioClient.SetBucketPolicy(ctx, bucket, policy); err != nil {
		return fmt.Errorf("设置MinIO公共桶读权限失败: %w", err)
	}
	return nil
}

func publicAvatarBucket(svcCtx *svc.ServiceContext) string {
	if svcCtx.Config.Minio.PublicBucket != "" {
		return svcCtx.Config.Minio.PublicBucket
	}
	return svcCtx.Config.Minio.FormalBucket
}

func publicObjectURL(svcCtx *svc.ServiceContext, bucket, objectKey string) string {
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
