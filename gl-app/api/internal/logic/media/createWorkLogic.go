package media

import (
	"context"
	"database/sql"
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
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	mediaModel "gl-app/api/internal/models/media"
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

type lockedMediaAsset struct {
	ID           int64  `db:"id"`
	UserID       int64  `db:"user_id"`
	ResourceType string `db:"resource_type"`
	Status       string `db:"status"`
	BoundWorkID  int64  `db:"bound_work_id"`
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

	if existed, findErr := l.svcCtx.WorkRepo.FindOneByUserIdIdempotencyKey(l.ctx, userID, payload.IdempotencyKey); findErr == nil {
		_ = l.svcCtx.MediaQueue.PublishProcessWork(l.ctx, existed.Id, 0)
		return &types.CreateWorkResp{
			WorkId:        existed.Id,
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
			_ = l.svcCtx.MediaAssetRepo.Delete(context.Background(), coverAsset.Id)
		}
	}()

	var workID int64
	err = l.svcCtx.SqlConn.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if payload.CollectionID > 0 {
			var collectionOwnerID int64
			if queryErr := session.QueryRowCtx(ctx, &collectionOwnerID,
				"SELECT user_id FROM collection WHERE id=? AND deleted_at IS NULL LIMIT 1",
				payload.CollectionID); queryErr != nil || collectionOwnerID != userID {
				return fmt.Errorf("合集不存在或不属于当前用户")
			}
		}
		for _, item := range payload.Assets {
			var asset lockedMediaAsset
			if queryErr := session.QueryRowCtx(ctx, &asset,
				"SELECT id,user_id,resource_type,status,bound_work_id FROM media_asset WHERE id=? FOR UPDATE",
				item.MediaID); queryErr != nil {
				return fmt.Errorf("素材%d不存在: %w", item.MediaID, queryErr)
			}
			if asset.UserID != userID {
				return fmt.Errorf("素材%d不属于当前用户", item.MediaID)
			}
			if asset.Status != "uploaded" || asset.BoundWorkID != 0 {
				return fmt.Errorf("素材%d状态不可发布", item.MediaID)
			}
			if asset.ResourceType != payload.Type {
				return fmt.Errorf("作品类型与素材%d类型不一致", item.MediaID)
			}
		}

		visibilityUsers, _ := json.Marshal(payload.Visibility.UserIDs)
		var scheduledAt sql.NullTime
		if payload.ScheduledAt != "" {
			parsed, parseErr := time.Parse(time.RFC3339, payload.ScheduledAt)
			if parseErr != nil {
				return fmt.Errorf("scheduledAt必须是RFC3339时间: %w", parseErr)
			}
			scheduledAt = sql.NullTime{Time: parsed.UTC(), Valid: true}
		}

		result, insertErr := session.ExecCtx(ctx, `INSERT INTO work
			(user_id,type,title,content,visibility,visibility_user_ids,collection_id,original,cover_asset_id,
			 process_status,review_status,publish_status,scheduled_at,idempotency_key)
			VALUES(?,?,?,?,?,?,?,?,?,'pending','waiting_process','pending',?,?)`,
			userID, payload.Type, strings.TrimSpace(payload.Title), payload.Content, payload.Visibility.Type,
			string(visibilityUsers), payload.CollectionID, payload.Original, coverAsset.Id, scheduledAt, payload.IdempotencyKey)
		if insertErr != nil {
			return fmt.Errorf("创建作品失败: %w", insertErr)
		}
		workID, insertErr = result.LastInsertId()
		if insertErr != nil {
			return insertErr
		}

		for _, item := range payload.Assets {
			role := payload.Type
			if _, insertErr = session.ExecCtx(ctx,
				"INSERT INTO work_asset(work_id,media_asset_id,role,sort) VALUES(?,?,?,?)",
				workID, item.MediaID, role, item.Sort); insertErr != nil {
				return insertErr
			}
			if _, insertErr = session.ExecCtx(ctx,
				"UPDATE media_asset SET status='bound',bound_work_id=? WHERE id=?",
				workID, item.MediaID); insertErr != nil {
				return insertErr
			}
		}
		if _, insertErr = session.ExecCtx(ctx,
			"INSERT INTO work_asset(work_id,media_asset_id,role,sort) VALUES(?,?, 'cover',0)",
			workID, coverAsset.Id); insertErr != nil {
			return insertErr
		}
		if _, insertErr = session.ExecCtx(ctx,
			"UPDATE media_asset SET status='bound',bound_work_id=? WHERE id=?",
			workID, coverAsset.Id); insertErr != nil {
			return insertErr
		}

		for index, topic := range payload.Topics {
			name := strings.TrimSpace(strings.TrimPrefix(topic.Name, "#"))
			if name == "" {
				continue
			}
			normalized := strings.ToLower(name)
			topicResult, topicErr := session.ExecCtx(ctx,
				"INSERT INTO topic(name,normalized_name) VALUES(?,?) ON DUPLICATE KEY UPDATE id=LAST_INSERT_ID(id)",
				name, normalized)
			if topicErr != nil {
				return topicErr
			}
			topicID, topicErr := topicResult.LastInsertId()
			if topicErr != nil {
				return topicErr
			}
			if _, topicErr = session.ExecCtx(ctx,
				"INSERT IGNORE INTO work_topic(work_id,topic_id,sort) VALUES(?,?,?)",
				workID, topicID, index); topicErr != nil {
				return topicErr
			}
		}

		_, insertErr = session.ExecCtx(ctx,
			"INSERT INTO media_process_task(work_id,task_type,status,stage,progress) VALUES(?,'process_work','pending','queued',0)",
			workID)
		return insertErr
	})
	if err != nil {
		return nil, err
	}
	cleanupCover = false

	if queueErr := l.svcCtx.MediaQueue.PublishProcessWork(l.ctx, workID, 0); queueErr != nil {
		l.Errorf("work %d queue publish failed, dispatcher will retry: %v", workID, queueErr)
	}

	return &types.CreateWorkResp{
		WorkId:        workID,
		ProcessStatus: "pending",
		ReviewStatus:  "waiting_process",
		PublishStatus: "pending",
	}, nil
}

func (l *CreateWorkLogic) storeCover(userID int64, file multipart.File, header *multipart.FileHeader) (*mediaModel.MediaAsset, error) {
	if header.Size <= 0 || header.Size > 5<<20 {
		return nil, fmt.Errorf("封面大小必须在5MB以内")
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

	asset := &mediaModel.MediaAsset{
		UserId:       userID,
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
	result, err := l.svcCtx.MediaAssetRepo.Insert(l.ctx, asset)
	if err != nil {
		_ = l.svcCtx.MinioClient.RemoveObject(l.ctx, bucket, objectKey, minio.RemoveObjectOptions{})
		return nil, fmt.Errorf("保存封面记录失败: %w", err)
	}
	asset.Id, err = result.LastInsertId()
	return asset, err
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
	if payload.Type == "image" && (len(payload.Assets) < 1 || len(payload.Assets) > 19) {
		return fmt.Errorf("图片作品必须包含1到19张图片")
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
