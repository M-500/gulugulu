package media

import (
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gl-app/api/internal/models/media"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUploadPresignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// create upload presigned url
func NewCreateUploadPresignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUploadPresignLogic {
	return &CreateUploadPresignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUploadPresignLogic) CreateUploadPresign(req *types.CreateUploadPresignReq) (resp *types.CreateUploadPresignResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	if userId <= 0 {
		return nil, fmt.Errorf("用户未登录")
	}

	resourceType, err := normalizeResourceType(req.ResourceType)
	if err != nil {
		return nil, err
	}

	originName, ext, err := normalizeFileName(req.FileName)
	if err != nil {
		return nil, err
	}

	if !isAllowedExt(resourceType, ext) {
		return nil, fmt.Errorf("不支持的%s格式: %s", resourceType, ext)
	}

	contentType := strings.TrimSpace(req.ContentType)
	if contentType == "" {
		contentType = mime.TypeByExtension(ext)
	}
	if contentType == "" {
		contentType = defaultContentType(resourceType)
	}

	expiresIn := l.svcCtx.Config.Minio.PresignExpire
	if expiresIn <= 0 {
		expiresIn = 900
	}

	bucket := l.svcCtx.Config.Minio.TempBucket
	if bucket == "" {
		return nil, fmt.Errorf("未配置MinIO临时上传桶")
	}

	if err := l.ensureBucket(bucket); err != nil {
		return nil, err
	}

	resourceName := fmt.Sprintf("%s-%s", uuid.NewString(), sanitizeResourceName(originName))
	objectKey := fmt.Sprintf("%d/%s/%s", userId, time.Now().Format("20060102"), resourceName)
	presignClient, err := l.presignClient()
	if err != nil {
		return nil, err
	}

	uploadUrl, err := presignClient.PresignedPutObject(l.ctx, bucket, objectKey, time.Duration(expiresIn)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("生成预签上传地址失败: %w", err)
	}
	previewUrl, err := presignClient.PresignedGetObject(l.ctx, bucket, objectKey, time.Duration(expiresIn)*time.Second, nil)
	if err != nil {
		return nil, fmt.Errorf("生成预览地址失败: %w", err)
	}

	result, err := l.svcCtx.MediaAssetRepo.Insert(l.ctx, &media.MediaAsset{
		UserId:       userId,
		ResourceType: resourceType,
		Bucket:       bucket,
		ObjectKey:    objectKey,
		OriginName:   originName,
		ContentType:  contentType,
		Ext:          strings.TrimPrefix(ext, "."),
		Status:       "uploading",
	})
	if err != nil {
		return nil, fmt.Errorf("创建媒体素材记录失败: %w", err)
	}

	mediaId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("获取媒体素材ID失败: %w", err)
	}

	return &types.CreateUploadPresignResp{
		MediaId:    mediaId,
		Bucket:     bucket,
		ObjectKey:  objectKey,
		UploadUrl:  uploadUrl.String(),
		PreviewUrl: previewUrl.String(),
		Method:     "PUT",
		ExpiresIn:  expiresIn,
		Headers: map[string]string{
			"Content-Type": contentType,
		},
	}, nil
}

func (l *CreateUploadPresignLogic) presignClient() (*minio.Client, error) {
	return minio.New(l.svcCtx.Config.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(l.svcCtx.Config.Minio.AccessKeyID, l.svcCtx.Config.Minio.SecretAccessKey, ""),
		Secure: l.svcCtx.Config.Minio.UseSSL,
	})
}

func (l *CreateUploadPresignLogic) ensureBucket(bucket string) error {
	exists, err := l.svcCtx.MinioClient.BucketExists(l.ctx, bucket)
	if err != nil {
		return fmt.Errorf("检查MinIO临时桶失败: %w", err)
	}

	if exists {
		return nil
	}

	if err := l.svcCtx.MinioClient.MakeBucket(l.ctx, bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("创建MinIO临时桶失败: %w", err)
	}

	return nil
}

func normalizeResourceType(resourceType string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(resourceType)) {
	case "image":
		return "image", nil
	case "video":
		return "video", nil
	default:
		return "", fmt.Errorf("resourceType必须是Image或Video")
	}
}

func normalizeFileName(fileName string) (originName, ext string, err error) {
	originName = filepath.Base(strings.TrimSpace(fileName))
	if originName == "." || originName == "/" || originName == "" {
		return "", "", fmt.Errorf("文件名不能为空")
	}

	ext = strings.ToLower(filepath.Ext(originName))
	if ext == "" {
		return "", "", fmt.Errorf("文件扩展名不能为空")
	}

	return originName, ext, nil
}

func sanitizeResourceName(fileName string) string {
	fileName = strings.ReplaceAll(fileName, " ", "_")
	var builder strings.Builder

	for _, item := range fileName {
		if item >= 'a' && item <= 'z' || item >= 'A' && item <= 'Z' || item >= '0' && item <= '9' || item == '.' || item == '_' || item == '-' {
			builder.WriteRune(item)
		}
	}

	if builder.Len() == 0 {
		return "resource"
	}

	return builder.String()
}

func isAllowedExt(resourceType, ext string) bool {
	videoExts := map[string]struct{}{
		".mp4": {}, ".m4v": {}, ".mov": {}, ".avi": {}, ".wmv": {}, ".flv": {}, ".mkv": {},
		".webm": {}, ".mpeg": {}, ".mpg": {}, ".3gp": {}, ".3g2": {}, ".ts": {}, ".mts": {},
		".m2ts": {}, ".m3u8": {}, ".rm": {}, ".rmvb": {}, ".asf": {},
	}
	imageExts := map[string]struct{}{
		".jpg": {}, ".jpeg": {}, ".png": {}, ".gif": {}, ".webp": {}, ".bmp": {}, ".tif": {},
		".tiff": {}, ".heic": {}, ".heif": {}, ".avif": {}, ".svg": {}, ".ico": {},
	}

	if resourceType == "video" {
		_, ok := videoExts[ext]
		return ok
	}

	_, ok := imageExts[ext]
	return ok
}

func defaultContentType(resourceType string) string {
	if resourceType == "video" {
		return "application/octet-stream"
	}

	return "application/octet-stream"
}
