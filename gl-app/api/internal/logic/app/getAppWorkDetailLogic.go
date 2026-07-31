package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppWorkDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取App公开作品详情
func NewGetAppWorkDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppWorkDetailLogic {
	return &GetAppWorkDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppWorkDetailLogic) GetAppWorkDetail(req *types.AppWorkDetailReq) (resp *types.AppWorkDetailResp, err error) {
	if req == nil || req.WorkId <= 0 {
		return nil, fmt.Errorf("作品ID不正确")
	}

	var row appWorkDetailRow
	query := `SELECT
			w.id,w.type,w.title,w.content,w.published_at,
			u.id AS author_id,COALESCE(NULLIF(u.nickname,''),u.email) AS author_name,u.avatar AS author_avatar
		FROM work w
		INNER JOIN user u ON u.id=w.user_id AND u.deleted_at IS NULL
		WHERE w.id=?
			AND w.deleted_at IS NULL
			AND w.process_status='succeeded'
			AND w.review_status='approved'
			AND w.publish_status='published'
			AND w.visibility='public'
		LIMIT 1`
	if err = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &row, query, req.WorkId); err != nil {
		l.Errorf("查询App公开作品%d失败: %v", req.WorkId, err)
		return nil, fmt.Errorf("作品不存在或暂不可见")
	}

	resp = &types.AppWorkDetailResp{
		WorkId:  row.ID,
		Type:    row.Type,
		Title:   row.Title,
		Content: row.Content,
		Author: types.RecommendAuthorItem{
			UserId:    row.AuthorID,
			NickName:  normalizeAuthorName(row.AuthorName),
			AvatarUrl: buildAppPublicAvatarURL(l.svcCtx, row.AuthorAvatar),
		},
		Like:          mockRecommendLike(row.ID),
		FavoriteCount: int64(19) + (row.ID*53)%5000,
		CommentCount:  int64(8) + (row.ID*31)%1000,
		ShareCount:    int64(3) + (row.ID*17)%400,
		Assets:        make([]types.AppWorkAssetItem, 0),
		Topics:        make([]types.AppWorkTopicItem, 0),
	}
	if row.PublishedAt.Valid {
		resp.PublishedAt = row.PublishedAt.Time.Format(time.RFC3339)
	}

	var assets []appWorkAssetRow
	assetQuery := `SELECT
			wa.media_asset_id,wa.role,wa.sort,
			ma.formal_bucket,ma.formal_object_key,ma.duration_ms,ma.width,ma.height
		FROM work_asset wa
		INNER JOIN media_asset ma ON ma.id=wa.media_asset_id AND ma.deleted_at IS NULL
		WHERE wa.work_id=? AND ma.status='ready'
		ORDER BY CASE wa.role WHEN 'image' THEN 0 WHEN 'video' THEN 1 ELSE 2 END,wa.sort,wa.id`
	if err = l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &assets, assetQuery, req.WorkId); err != nil {
		l.Errorf("查询App公开作品%d资源失败: %v", req.WorkId, err)
		return nil, fmt.Errorf("查询作品资源失败")
	}

	for _, asset := range assets {
		assetURL := ""
		if asset.Role == "video" {
			resp.DurationMs = asset.DurationMs
			resp.VideoPlaylist = fmt.Sprintf("/app/v1/works/%d/playlist", row.ID)
			assetURL = resp.VideoPlaylist
		} else {
			assetURL = l.presignAppWorkAsset(asset.Bucket, asset.ObjectKey)
		}
		if asset.Role == "cover" {
			resp.CoverUrl = assetURL
		}
		resp.Assets = append(resp.Assets, types.AppWorkAssetItem{
			MediaId: asset.MediaID,
			Role:    asset.Role,
			Sort:    asset.Sort,
			Url:     assetURL,
			Width:   asset.Width,
			Height:  asset.Height,
		})
	}

	var topics []appWorkTopicRow
	topicQuery := `SELECT t.id,t.name
		FROM work_topic wt
		INNER JOIN topic t ON t.id=wt.topic_id
		WHERE wt.work_id=?
		ORDER BY wt.sort,wt.id`
	if err = l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &topics, topicQuery, req.WorkId); err != nil {
		l.Errorf("查询App公开作品%d话题失败: %v", req.WorkId, err)
		return nil, fmt.Errorf("查询作品话题失败")
	}
	for _, topic := range topics {
		resp.Topics = append(resp.Topics, types.AppWorkTopicItem{
			TopicId: topic.ID,
			Name:    topic.Name,
		})
	}

	return resp, nil
}

func (l *GetAppWorkDetailLogic) presignAppWorkAsset(bucket, objectKey string) string {
	if bucket == "" || objectKey == "" {
		return ""
	}
	presigned, err := l.svcCtx.MinioClient.PresignedGetObject(
		l.ctx, bucket, objectKey, 30*time.Minute, nil,
	)
	if err != nil {
		l.Errorf("生成App作品资源签名地址失败: %v", err)
		return ""
	}
	return presigned.String()
}

type appWorkDetailRow struct {
	ID           int64        `db:"id"`
	Type         string       `db:"type"`
	Title        string       `db:"title"`
	Content      string       `db:"content"`
	PublishedAt  sql.NullTime `db:"published_at"`
	AuthorID     int64        `db:"author_id"`
	AuthorName   string       `db:"author_name"`
	AuthorAvatar string       `db:"author_avatar"`
}

type appWorkAssetRow struct {
	MediaID    int64  `db:"media_asset_id"`
	Role       string `db:"role"`
	Sort       int64  `db:"sort"`
	Bucket     string `db:"formal_bucket"`
	ObjectKey  string `db:"formal_object_key"`
	DurationMs int64  `db:"duration_ms"`
	Width      int64  `db:"width"`
	Height     int64  `db:"height"`
}

type appWorkTopicRow struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
