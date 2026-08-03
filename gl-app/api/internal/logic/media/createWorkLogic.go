package media

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	mediarepo "gl-app/api/internal/repo/media_repo"
	workrepo "gl-app/api/internal/repo/work_repo"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

const maxImageAssets = 500

// 发布作品
func NewCreateWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkLogic {
	return &CreateWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type createWorkPayload struct {
	IdempotencyKey string `json:"idempotencyKey"`
	Type           string `json:"type"`
	Title          string `json:"title"`
	Content        string `json:"content"`
	Visibility     struct {
		Type    string  `json:"type"`
		UserIDs []int64 `json:"userIds"`
	} `json:"visibility"`
	Topics []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"topics"`
	Assets []struct {
		MediaID int64 `json:"mediaId"`
		Sort    int   `json:"sort"`
	} `json:"assets"`
	CollectionID int64  `json:"collectionId"`
	Original     bool   `json:"original"`
	ScheduledAt  string `json:"scheduledAt"`
}

func (l *CreateWorkLogic) CreateWork(req *types.CreateWorkReq, cover multipart.File, coverHeader *multipart.FileHeader, idempotencyHeader string) (resp *types.CreateWorkResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID <= 0 {
		return nil, fmt.Errorf("用户未登录")
	}

	var payload createWorkPayload
	if err := json.Unmarshal([]byte(req.Payload), &payload); err != nil {
		return nil, fmt.Errorf("payload不是合法JSON: %w", err)
	}
	if strings.TrimSpace(idempotencyHeader) != "" {
		payload.IdempotencyKey = strings.TrimSpace(idempotencyHeader)
	}
	if err := validateCreateWorkPayload(&payload); err != nil {
		return nil, err
	}

	if existed, findErr := l.svcCtx.WorkRepo.FindByIdempotencyKey(l.ctx, userID, payload.IdempotencyKey); findErr == nil {
		if existed.Type == "image" && existed.ProcessStatus != "succeeded" {
			if promoteErr := promoteImageWork(l.ctx, l.svcCtx, int64(existed.ID)); promoteErr != nil {
				markImagePromotionFailed(l.ctx, l.svcCtx, int64(existed.ID), promoteErr)
				return nil, promoteErr
			}
			return &types.CreateWorkResp{
				WorkId:        int64(existed.ID),
				ProcessStatus: "succeeded",
				ReviewStatus:  "pending_review",
				PublishStatus: "pending",
			}, nil
		}
		// 数据库事务已成功但Kafka曾短暂不可用时，前端可使用同一个
		// Idempotency-Key安全重试。只有消息成功写入Kafka才返回成功。
		if existed.ProcessStatus == "pending" {
			if queueErr := l.svcCtx.MediaQueue.PublishProcessWork(l.ctx, int64(existed.ID), 0); queueErr != nil {
				return nil, queueErr
			}
		}
		return &types.CreateWorkResp{
			WorkId:        int64(existed.ID),
			ProcessStatus: existed.ProcessStatus,
			ReviewStatus:  existed.ReviewStatus,
			PublishStatus: existed.PublishStatus,
		}, nil
	}

	coverAsset, err := l.storeCover(userID, cover, coverHeader)
	if err != nil {
		return nil, err
	}
	cleanupCover := true
	defer func() {
		if cleanupCover {
			_ = l.svcCtx.MinioClient.RemoveObject(context.Background(), coverAsset.Bucket, coverAsset.ObjectKey, minio.RemoveObjectOptions{})
			_ = l.svcCtx.MediaRepo.DeleteAsset(context.Background(), int64(coverAsset.ID))
		}
	}()

	visibilityUsers := ""
	if len(payload.Visibility.UserIDs) > 0 {
		value, marshalErr := json.Marshal(payload.Visibility.UserIDs)
		if marshalErr != nil {
			return nil, fmt.Errorf("序列化可见用户失败: %w", marshalErr)
		}
		visibilityUsers = string(value)
	}
	var scheduledAt *time.Time
	if payload.ScheduledAt != "" {
		parsed, parseErr := time.Parse(time.RFC3339, payload.ScheduledAt)
		if parseErr != nil {
			return nil, fmt.Errorf("scheduledAt必须是RFC3339时间: %w", parseErr)
		}
		value := parsed.UTC()
		scheduledAt = &value
	}
	input := workrepo.CreateWorkInput{
		UserID: userID, Type: payload.Type, Title: strings.TrimSpace(payload.Title), Content: payload.Content,
		Visibility: payload.Visibility.Type, VisibilityUserIDs: visibilityUsers, CollectionID: payload.CollectionID,
		Original: payload.Original, CoverAssetID: int64(coverAsset.ID), ScheduledAt: scheduledAt, IdempotencyKey: payload.IdempotencyKey,
		Assets: make([]workrepo.CreateAssetInput, 0, len(payload.Assets)),
		Topics: make([]workrepo.CreateTopicInput, 0, len(payload.Topics)),
	}
	for _, item := range payload.Assets {
		input.Assets = append(input.Assets, workrepo.CreateAssetInput{MediaID: item.MediaID, Sort: int64(item.Sort)})
	}
	for index, item := range payload.Topics {
		input.Topics = append(input.Topics, workrepo.CreateTopicInput{Name: item.Name, Sort: int64(index)})
	}
	workID, err := l.svcCtx.WorkRepo.Create(l.ctx, input)
	if err != nil {
		return nil, err
	}
	cleanupCover = false

	if payload.Type == "image" {
		if promoteErr := promoteImageWork(l.ctx, l.svcCtx, workID); promoteErr != nil {
			markImagePromotionFailed(l.ctx, l.svcCtx, workID, promoteErr)
			return nil, promoteErr
		}
		return &types.CreateWorkResp{
			WorkId:        workID,
			ProcessStatus: "succeeded",
			ReviewStatus:  "pending_review",
			PublishStatus: "pending",
		}, nil
	}

	// 只有视频需要进入Kafka，由独立Worker下载、探测并转码为HLS。
	if queueErr := l.svcCtx.MediaQueue.PublishProcessWork(l.ctx, workID, 0); queueErr != nil {
		return nil, queueErr
	}

	return &types.CreateWorkResp{
		WorkId:        workID,
		ProcessStatus: "pending",
		ReviewStatus:  "waiting_process",
		PublishStatus: "pending",
	}, nil
}

func (l *CreateWorkLogic) storeCover(userID int64, file multipart.File, header *multipart.FileHeader) (*mediarepo.MediaAsset, error) {
	if header.Size <= 0 || header.Size > 100<<20 {
		return nil, fmt.Errorf("封面大小必须在100MB以内")
	}
	config, format, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("封面不是有效图片: %w", err)
	}
	if format != "jpeg" && format != "png" {
		return nil, fmt.Errorf("封面只支持jpg/jpeg/png")
	}
	if config.Width < 1 || config.Height < 1 {
		return nil, fmt.Errorf("封面尺寸无效")
	}
	if _, err = file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("读取封面失败: %w", err)
	}

	bucket := l.svcCtx.Config.Minio.TempBucket
	objectKey := fmt.Sprintf("%d/%s/cover-%s%s", userID, time.Now().Format("20060102"), uuid.NewString(), strings.ToLower(filepath.Ext(header.Filename)))
	contentType := "image/jpeg"
	if format == "png" {
		contentType = "image/png"
	}
	if _, err = l.svcCtx.MinioClient.PutObject(l.ctx, bucket, objectKey, file, header.Size, minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return nil, fmt.Errorf("上传封面失败: %w", err)
	}

	asset := &mediarepo.MediaAsset{
		UserID:       userID,
		ResourceType: "image",
		Bucket:       bucket,
		ObjectKey:    objectKey,
		OriginName:   filepath.Base(header.Filename),
		ContentType:  contentType,
		Ext:          strings.TrimPrefix(strings.ToLower(filepath.Ext(header.Filename)), "."),
		FileSize:     header.Size,
		Status:       "uploaded",
		Width:        int64(config.Width),
		Height:       int64(config.Height),
	}
	err = l.svcCtx.MediaRepo.CreateAsset(l.ctx, asset)
	if err != nil {
		_ = l.svcCtx.MinioClient.RemoveObject(l.ctx, bucket, objectKey, minio.RemoveObjectOptions{})
		return nil, fmt.Errorf("保存封面记录失败: %w", err)
	}
	return asset, nil
}

func validateCreateWorkPayload(payload *createWorkPayload) error {
	payload.Type = strings.ToLower(strings.TrimSpace(payload.Type))
	payload.Title = strings.TrimSpace(payload.Title)
	payload.IdempotencyKey = strings.TrimSpace(payload.IdempotencyKey)
	if payload.IdempotencyKey == "" || len(payload.IdempotencyKey) > 64 {
		return fmt.Errorf("Idempotency-Key不能为空且最多64字符")
	}
	if payload.Type != "image" && payload.Type != "video" {
		return fmt.Errorf("type必须是image或video")
	}
	if payload.Title == "" || len([]rune(payload.Title)) > 50 {
		return fmt.Errorf("title不能为空且最多50个字符")
	}
	if len([]rune(payload.Content)) > 1000 {
		return fmt.Errorf("content最多1000个字符")
	}
	if payload.Type == "video" && len(payload.Assets) != 1 {
		return fmt.Errorf("视频作品必须且只能包含一个视频")
	}
	if payload.Type == "image" && (len(payload.Assets) < 1 || len(payload.Assets) > maxImageAssets) {
		return fmt.Errorf("图片作品必须包含1到%d张图片", maxImageAssets)
	}
	seen := make(map[int64]struct{}, len(payload.Assets))
	for _, asset := range payload.Assets {
		if asset.MediaID <= 0 {
			return fmt.Errorf("mediaId必须大于0")
		}
		if _, ok := seen[asset.MediaID]; ok {
			return fmt.Errorf("mediaId不能重复")
		}
		seen[asset.MediaID] = struct{}{}
	}
	switch payload.Visibility.Type {
	case "public", "private", "mutual":
		payload.Visibility.UserIDs = nil
	case "selected", "excluded":
		if len(payload.Visibility.UserIDs) == 0 {
			return fmt.Errorf("%s可见性必须指定用户", payload.Visibility.Type)
		}
	default:
		return fmt.Errorf("visibility.type不合法")
	}
	if len(payload.Topics) > 10 {
		return fmt.Errorf("话题最多10个")
	}
	if payload.ScheduledAt != "" {
		value, err := time.Parse(time.RFC3339, payload.ScheduledAt)
		if err != nil || value.Before(time.Now().Add(5*time.Minute)) {
			return fmt.Errorf("scheduledAt必须是至少晚于当前时间5分钟的RFC3339时间")
		}
	}
	return nil
}
