package app

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPublishedWorkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取用户已发布的公开作品
func NewGetUserPublishedWorkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPublishedWorkListLogic {
	return &GetUserPublishedWorkListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPublishedWorkListLogic) GetUserPublishedWorkList(req *types.UserPublishedWorkListReq) (resp *types.UserPublishedWorkListResp, err error) {
	if req == nil || req.UserId <= 0 {
		return nil, fmt.Errorf("用户ID不正确")
	}

	// 先确认用户存在，避免把“不存在的用户”和“用户暂时没有作品”混为一谈。
	userInfo, err := findAppPublicUser(l.ctx, l.svcCtx, req.UserId)
	if err != nil {
		return nil, err
	}

	page, pageSize := normalizeRecommendPage(req.Page, req.PageSize)
	var total int64
	countQuery := `SELECT COUNT(1)
		FROM work w
		WHERE w.user_id=?
			AND w.deleted_at IS NULL
			AND w.process_status='succeeded'
			AND w.review_status='approved'
			AND w.publish_status='published'
			AND w.visibility='public'`
	if err = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &total, countQuery, req.UserId); err != nil {
		l.Errorf("查询用户%d已发布作品数量失败: %v", req.UserId, err)
		return nil, fmt.Errorf("查询用户作品失败")
	}

	resp = &types.UserPublishedWorkListResp{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     make([]types.RecommendWorkItem, 0),
	}
	if total == 0 {
		return resp, nil
	}

	// 用户公开主页只允许展示处理成功、审核通过、正式发布且公开可见的作品。
	var rows []userPublishedWorkRow
	offset := (page - 1) * pageSize
	listQuery := `SELECT
			w.id,
			w.type,
			w.title,
			LEFT(w.content,120) AS content_excerpt,
			w.published_at,
			COALESCE(cover.formal_bucket,'') AS cover_bucket,
			COALESCE(cover.formal_object_key,'') AS cover_object_key,
			COALESCE(video.duration_ms,0) AS duration_ms
		FROM work w
		LEFT JOIN media_asset cover
			ON cover.id=w.cover_asset_id AND cover.deleted_at IS NULL
		LEFT JOIN work_asset video_relation
			ON video_relation.work_id=w.id AND video_relation.role='video' AND video_relation.sort=0
		LEFT JOIN media_asset video
			ON video.id=video_relation.media_asset_id AND video.deleted_at IS NULL
		WHERE w.user_id=?
			AND w.deleted_at IS NULL
			AND w.process_status='succeeded'
			AND w.review_status='approved'
			AND w.publish_status='published'
			AND w.visibility='public'
		ORDER BY w.published_at DESC,w.id DESC
		LIMIT ? OFFSET ?`
	if err = l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &rows, listQuery, req.UserId, pageSize, offset); err != nil {
		l.Errorf("查询用户%d已发布作品列表失败: %v", req.UserId, err)
		return nil, fmt.Errorf("查询用户作品失败")
	}

	author := types.RecommendAuthorItem{
		UserId:    userInfo.Id,
		NickName:  normalizeAuthorName(userInfo.Nickname),
		AvatarUrl: buildAppPublicAvatarURL(l.svcCtx, userInfo.Avatar),
	}
	resp.List = make([]types.RecommendWorkItem, 0, len(rows))
	for _, row := range rows {
		item := types.RecommendWorkItem{
			WorkId:         row.ID,
			Type:           row.Type,
			Title:          row.Title,
			ContentExcerpt: strings.TrimSpace(row.ContentExcerpt),
			DurationMs:     row.DurationMs,
			Author:         author,
			Like:           mockRecommendLike(row.ID),
		}
		if row.PublishedAt.Valid {
			item.PublishedAt = row.PublishedAt.Time.Format(time.RFC3339)
		}
		if row.CoverBucket != "" && row.CoverObjectKey != "" {
			coverURL, signErr := l.svcCtx.MinioClient.PresignedGetObject(
				l.ctx, row.CoverBucket, row.CoverObjectKey, 30*time.Minute, nil,
			)
			if signErr != nil {
				l.Errorf("生成用户作品%d封面预签名地址失败: %v", row.ID, signErr)
			} else {
				item.CoverUrl = coverURL.String()
			}
		}
		resp.List = append(resp.List, item)
	}

	resp.HasMore = offset+int64(len(resp.List)) < total
	return resp, nil
}

type userPublishedWorkRow struct {
	ID             int64        `db:"id"`
	Type           string       `db:"type"`
	Title          string       `db:"title"`
	ContentExcerpt string       `db:"content_excerpt"`
	PublishedAt    sql.NullTime `db:"published_at"`
	CoverBucket    string       `db:"cover_bucket"`
	CoverObjectKey string       `db:"cover_object_key"`
	DurationMs     int64        `db:"duration_ms"`
}
