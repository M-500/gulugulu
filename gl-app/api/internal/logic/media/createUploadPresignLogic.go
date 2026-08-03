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
	mediarepo "gl-app/api/internal/repo/media_repo"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUploadPresignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	// MinIO客户端由ServiceContext统一创建并注入，避免在业务逻辑中重复初始化连接配置。
	minioClient *minio.Client
}

// create upload presigned url
func NewCreateUploadPresignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUploadPresignLogic {
	return &CreateUploadPresignLogic{
		Logger:      logx.WithContext(ctx),
		ctx:         ctx,
		svcCtx:      svcCtx,
		minioClient: svcCtx.MinioClient,
	}
}

func (l *CreateUploadPresignLogic) CreateUploadPresign(req *types.CreateUploadPresignReq) (resp *types.CreateUploadPresignResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	if userId <= 0 {
		return nil, l.logAndReturn("创建上传预签名失败", fmt.Errorf("用户未登录"))
	}

	resourceType, err := normalizeResourceType(req.ResourceType)
	if err != nil {
		return nil, l.logAndReturn("创建上传预签名失败", err)
	}

	originName, ext, err := normalizeFileName(req.FileName)
	if err != nil {
		return nil, l.logAndReturn("创建上传预签名失败", err)
	}

	if !isAllowedExt(resourceType, ext) {
		return nil, l.logAndReturn("创建上传预签名失败", fmt.Errorf("不支持的%s格式: %s", resourceType, ext))
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
		return nil, l.logAndReturn("创建上传预签名失败", fmt.Errorf("未配置MinIO临时上传桶"))
	}

	// 临时桶用于接收前端直传的原始素材，后续发布成功后再迁移到正式桶。
	if err := l.ensureBucket(bucket); err != nil {
		return nil, l.logAndReturn("创建上传预签名失败", err)
	}

	resourceName := fmt.Sprintf("%s-%s", uuid.NewString(), sanitizeResourceName(originName))
	// 对象路径按用户和日期分区，便于后续排查、生命周期清理和迁移处理。
	objectKey := fmt.Sprintf("%d/%s/%s", userId, time.Now().Format("20060102"), resourceName)

	uploadUrl, err := l.minioClient.PresignedPutObject(l.ctx, bucket, objectKey, time.Duration(expiresIn)*time.Second)
	if err != nil {
		return nil, l.logAndReturn("生成预签上传地址失败", fmt.Errorf("调用MinIO生成PUT预签名失败: %w", err))
	}
	previewUrl, err := l.minioClient.PresignedGetObject(l.ctx, bucket, objectKey, time.Duration(expiresIn)*time.Second, nil)
	if err != nil {
		return nil, l.logAndReturn("生成预览地址失败", fmt.Errorf("调用MinIO生成GET预签名失败: %w", err))
	}

	asset := &mediarepo.MediaAsset{
		UserID:       userId,
		ResourceType: resourceType,
		Bucket:       bucket,
		ObjectKey:    objectKey,
		OriginName:   originName,
		ContentType:  contentType,
		Ext:          strings.TrimPrefix(ext, "."),
		Status:       "uploading",
	}
	err = l.svcCtx.MediaRepo.CreateAsset(l.ctx, asset)
	if err != nil {
		return nil, l.logAndReturn("创建媒体素材记录失败", fmt.Errorf("写入媒体素材记录失败: %w", err))
	}

	return &types.CreateUploadPresignResp{
		MediaId:    asset.ID,
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

func (l *CreateUploadPresignLogic) ensureBucket(bucket string) error {
	if l.minioClient == nil {
		return fmt.Errorf("MinIO客户端未初始化")
	}

	// 本地开发或新环境首次启动时桶可能不存在，这里做一次幂等兜底。
	exists, err := l.minioClient.BucketExists(l.ctx, bucket)
	if err != nil {
		return fmt.Errorf("检查MinIO临时桶失败: %w", err)
	}

	if exists {
		return nil
	}

	if err := l.minioClient.MakeBucket(l.ctx, bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("创建MinIO临时桶失败: %w", err)
	}

	return nil
}

func (l *CreateUploadPresignLogic) logAndReturn(message string, err error) error {
	l.Errorf("%s：%v", message, err)
	return err
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
