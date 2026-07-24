package media

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取作品详情
func NewGetWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkLogic {
	return &GetWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWorkLogic) GetWork(req *types.WorkIdReq) (resp *types.WorkDetailResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	var row struct {
		ID            int64        `db:"id"`
		Type          string       `db:"type"`
		Title         string       `db:"title"`
		Content       string       `db:"content"`
		Visibility    string       `db:"visibility"`
		CollectionID  int64        `db:"collection_id"`
		Original      int64        `db:"original"`
		ProcessStatus string       `db:"process_status"`
		ReviewStatus  string       `db:"review_status"`
		PublishStatus string       `db:"publish_status"`
		ScheduledAt   sql.NullTime `db:"scheduled_at"`
		PublishedAt   sql.NullTime `db:"published_at"`
	}
	if err = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &row, `SELECT id,type,title,content,visibility,collection_id,original,
		process_status,review_status,publish_status,scheduled_at,published_at
		FROM work WHERE id=? AND user_id=? AND deleted_at IS NULL LIMIT 1`, req.WorkId, userID); err != nil {
		return nil, fmt.Errorf("作品不存在: %w", err)
	}
	resp = &types.WorkDetailResp{
		WorkId: row.ID, Type: row.Type, Title: row.Title, Content: row.Content, Visibility: row.Visibility,
		CollectionId: row.CollectionID, Original: row.Original == 1, ProcessStatus: row.ProcessStatus,
		ReviewStatus: row.ReviewStatus, PublishStatus: row.PublishStatus,
		Assets: []types.WorkAssetItem{}, Topics: []types.WorkTopicItem{},
	}
	if row.ScheduledAt.Valid {
		resp.ScheduledAt = row.ScheduledAt.Time.Format(time.RFC3339)
	}
	if row.PublishedAt.Valid {
		resp.PublishedAt = row.PublishedAt.Time.Format(time.RFC3339)
	}

	var assets []struct {
		MediaID   int64  `db:"media_asset_id"`
		Role      string `db:"role"`
		Sort      int64  `db:"sort"`
		Bucket    string `db:"formal_bucket"`
		ObjectKey string `db:"formal_object_key"`
	}
	if err = l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &assets, `SELECT wa.media_asset_id,wa.role,wa.sort,
		ma.formal_bucket,ma.formal_object_key FROM work_asset wa JOIN media_asset ma ON ma.id=wa.media_asset_id
		WHERE wa.work_id=? ORDER BY wa.role,wa.sort`, req.WorkId); err != nil {
		return nil, err
	}
	for _, asset := range assets {
		url := ""
		if asset.Bucket != "" && asset.ObjectKey != "" {
			presigned, signErr := l.svcCtx.MinioClient.PresignedGetObject(l.ctx, asset.Bucket, asset.ObjectKey, 15*time.Minute, nil)
			if signErr == nil {
				url = presigned.String()
			}
		}
		resp.Assets = append(resp.Assets, types.WorkAssetItem{MediaId: asset.MediaID, Role: asset.Role, Sort: asset.Sort, Url: url})
	}
	var topics []struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	if err = l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &topics, `SELECT t.id,t.name FROM work_topic wt
		JOIN topic t ON t.id=wt.topic_id WHERE wt.work_id=? ORDER BY wt.sort`, req.WorkId); err != nil {
		return nil, err
	}
	for _, topic := range topics {
		resp.Topics = append(resp.Topics, types.WorkTopicItem{TopicId: topic.ID, Name: topic.Name})
	}

	return resp, nil
}
