package media

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuditWorkDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询审核作品详情
func NewGetAuditWorkDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuditWorkDetailLogic {
	return &GetAuditWorkDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAuditWorkDetailLogic) GetAuditWorkDetail(req *types.WorkIdReq) (resp *types.AuditWorkDetailResp, err error) {
	if err = ensureAuditPermission(l.svcCtx, ctxdata.GetUidFromCtx(l.ctx)); err != nil {
		return nil, err
	}

	var row struct {
		ID            int64        `db:"id"`
		Type          string       `db:"type"`
		Title         string       `db:"title"`
		Content       string       `db:"content"`
		Visibility    string       `db:"visibility"`
		Original      int64        `db:"original"`
		AuthorID      int64        `db:"user_id"`
		AuthorName    string       `db:"author_name"`
		AuthorEmail   string       `db:"author_email"`
		ReviewStatus  string       `db:"review_status"`
		PublishStatus string       `db:"publish_status"`
		ReviewReason  string       `db:"review_reason"`
		ScheduledAt   sql.NullTime `db:"scheduled_at"`
		CreatedAt     time.Time    `db:"created_at"`
		ReviewedAt    sql.NullTime `db:"reviewed_at"`
		CoverBucket   string       `db:"cover_bucket"`
		CoverObject   string       `db:"cover_object_key"`
	}
	if err = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &row, `SELECT
			w.id,w.type,w.title,w.content,w.visibility,w.original,w.user_id,
			COALESCE(NULLIF(u.nickname,''),u.email) AS author_name,u.email AS author_email,
			w.review_status,w.publish_status,w.review_reason,w.scheduled_at,w.created_at,w.reviewed_at,
			COALESCE(cover.formal_bucket,'') AS cover_bucket,
			COALESCE(cover.formal_object_key,'') AS cover_object_key
		FROM work w JOIN user u ON u.id=w.user_id
		LEFT JOIN media_asset cover ON cover.id=w.cover_asset_id
		WHERE w.id=? AND w.process_status='succeeded' AND w.deleted_at IS NULL LIMIT 1`, req.WorkId); err != nil {
		return nil, fmt.Errorf("审核作品不存在: %w", err)
	}

	resp = &types.AuditWorkDetailResp{
		WorkId: row.ID, Type: row.Type, Title: row.Title, Content: row.Content,
		Visibility: row.Visibility, Original: row.Original == 1, AuthorId: row.AuthorID,
		AuthorName: row.AuthorName, AuthorEmail: row.AuthorEmail, ReviewStatus: row.ReviewStatus,
		PublishStatus: row.PublishStatus, ReviewReason: row.ReviewReason,
		SubmittedAt: row.CreatedAt.Format(time.RFC3339),
		Assets:      make([]types.WorkAssetItem, 0), Topics: make([]types.WorkTopicItem, 0),
	}
	resp.CoverUrl = l.presignAuditAsset(row.CoverBucket, row.CoverObject)
	if row.ScheduledAt.Valid {
		resp.ScheduledAt = row.ScheduledAt.Time.Format(time.RFC3339)
	}
	if row.ReviewedAt.Valid {
		resp.ReviewedAt = row.ReviewedAt.Time.Format(time.RFC3339)
	}

	var assets []struct {
		MediaID   int64  `db:"media_asset_id"`
		Role      string `db:"role"`
		Sort      int64  `db:"sort"`
		Bucket    string `db:"formal_bucket"`
		ObjectKey string `db:"formal_object_key"`
		Duration  int64  `db:"duration_ms"`
	}
	if err = l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &assets, `SELECT wa.media_asset_id,wa.role,wa.sort,
			ma.formal_bucket,ma.formal_object_key,ma.duration_ms
		FROM work_asset wa JOIN media_asset ma ON ma.id=wa.media_asset_id
		WHERE wa.work_id=? AND ma.status='ready'
		ORDER BY CASE wa.role WHEN 'video' THEN 0 WHEN 'image' THEN 1 ELSE 2 END,wa.sort`, req.WorkId); err != nil {
		return nil, fmt.Errorf("查询审核素材失败: %w", err)
	}
	for _, asset := range assets {
		resp.Assets = append(resp.Assets, types.WorkAssetItem{
			MediaId: asset.MediaID, Role: asset.Role, Sort: asset.Sort,
			Url: l.presignAuditAsset(asset.Bucket, asset.ObjectKey),
		})
		if asset.Role == "video" {
			resp.DurationMs = asset.Duration
			resp.VideoPlaylist = l.buildSignedPlaylist(asset.Bucket, asset.ObjectKey)
		}
	}

	var topics []struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	if err = l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &topics, `SELECT t.id,t.name
		FROM work_topic wt JOIN topic t ON t.id=wt.topic_id
		WHERE wt.work_id=? ORDER BY wt.sort`, req.WorkId); err != nil {
		return nil, fmt.Errorf("查询作品话题失败: %w", err)
	}
	for _, topic := range topics {
		resp.Topics = append(resp.Topics, types.WorkTopicItem{TopicId: topic.ID, Name: topic.Name})
	}

	return resp, nil
}

// 将HLS分片改写为短期签名地址，审核页面无需公开正式桶即可播放视频。
func (l *GetAuditWorkDetailLogic) buildSignedPlaylist(bucket, objectKey string) string {
	if bucket == "" || objectKey == "" {
		return ""
	}
	object, err := l.svcCtx.MinioClient.GetObject(l.ctx, bucket, objectKey, minio.GetObjectOptions{})
	if err != nil {
		l.Errorf("读取审核视频清单失败: %v", err)
		return ""
	}
	defer object.Close()
	content, err := io.ReadAll(io.LimitReader(object, 2<<20))
	if err != nil {
		l.Errorf("读取审核视频清单失败: %v", err)
		return ""
	}
	lines := strings.Split(string(content), "\n")
	basePath := path.Dir(objectKey)
	for index, line := range lines {
		segment := strings.TrimSpace(line)
		if segment == "" || strings.HasPrefix(segment, "#") {
			continue
		}
		lines[index] = l.presignAuditAsset(bucket, path.Join(basePath, segment))
	}
	return strings.Join(lines, "\n")
}

func (l *GetAuditWorkDetailLogic) presignAuditAsset(bucket, objectKey string) string {
	if bucket == "" || objectKey == "" {
		return ""
	}
	value, err := l.svcCtx.MinioClient.PresignedGetObject(l.ctx, bucket, objectKey, 30*time.Minute, nil)
	if err != nil {
		l.Errorf("生成审核资源地址失败: %v", err)
		return ""
	}
	return value.String()
}
