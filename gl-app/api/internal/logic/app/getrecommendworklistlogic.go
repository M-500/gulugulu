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

type GetRecommendWorkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// App首页推荐作品列表
func NewGetRecommendWorkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendWorkListLogic {
	return &GetRecommendWorkListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRecommendWorkListLogic) GetRecommendWorkList(req *types.RecommendWorkListReq) (resp *types.RecommendWorkListResp, err error) {
	page, pageSize := normalizeRecommendPage(req.Page, req.PageSize)
	limit := pageSize + 1
	offset := (page - 1) * pageSize

	// App推荐流是公开接口，只展示已经处理成功、审核通过、正式发布且公开可见的作品。
	var rows []recommendWorkRow
	query := `SELECT
			w.id,
			w.type,
			w.title,
			LEFT(w.content, 120) AS content_excerpt,
			w.published_at,
			COALESCE(cover.formal_bucket, '') AS cover_bucket,
			COALESCE(cover.formal_object_key, '') AS cover_object_key,
			COALESCE(video.duration_ms, 0) AS duration_ms,
			u.id AS author_id,
			COALESCE(NULLIF(u.nickname, ''), u.email) AS author_name
		FROM work w
		INNER JOIN user u ON u.id = w.user_id AND u.deleted_at IS NULL
		LEFT JOIN media_asset cover ON cover.id = w.cover_asset_id AND cover.deleted_at IS NULL
		LEFT JOIN work_asset video_relation ON video_relation.work_id = w.id AND video_relation.role = 'video' AND video_relation.sort = 0
		LEFT JOIN media_asset video ON video.id = video_relation.media_asset_id AND video.deleted_at IS NULL
		WHERE w.deleted_at IS NULL
			AND w.process_status = 'succeeded'
			AND w.review_status = 'approved'
			AND w.publish_status = 'published'
			AND w.visibility = 'public'
		ORDER BY w.published_at DESC, w.id DESC
		LIMIT ? OFFSET ?`
	if err = l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &rows, query, limit, offset); err != nil {
		l.Errorf("查询App推荐作品列表失败: %v", err)
		return nil, fmt.Errorf("查询推荐作品失败")
	}

	hasMore := int64(len(rows)) > pageSize
	if hasMore {
		rows = rows[:pageSize]
	}

	resp = &types.RecommendWorkListResp{
		Page:     page,
		PageSize: pageSize,
		HasMore:  hasMore,
		List:     make([]types.RecommendWorkItem, 0, len(rows)),
	}

	for _, row := range rows {
		item := types.RecommendWorkItem{
			WorkId:         row.ID,
			Type:           row.Type,
			Title:          row.Title,
			ContentExcerpt: strings.TrimSpace(row.ContentExcerpt),
			DurationMs:     row.DurationMs,
			Author: types.RecommendAuthorItem{
				UserId:   row.AuthorID,
				NickName: normalizeAuthorName(row.AuthorName),
			},
			Like: mockRecommendLike(row.ID),
		}
		if row.PublishedAt.Valid {
			item.PublishedAt = row.PublishedAt.Time.Format(time.RFC3339)
		}
		if row.CoverBucket != "" && row.CoverObjectKey != "" {
			coverURL, signErr := l.svcCtx.MinioClient.PresignedGetObject(
				l.ctx, row.CoverBucket, row.CoverObjectKey, 30*time.Minute, nil,
			)
			if signErr != nil {
				l.Errorf("生成App推荐作品%d封面预签名地址失败: %v", row.ID, signErr)
			} else {
				item.CoverUrl = coverURL.String()
			}
		}
		resp.List = append(resp.List, item)
	}

	return resp, nil
}

type recommendWorkRow struct {
	ID             int64        `db:"id"`
	Type           string       `db:"type"`
	Title          string       `db:"title"`
	ContentExcerpt string       `db:"content_excerpt"`
	PublishedAt    sql.NullTime `db:"published_at"`
	CoverBucket    string       `db:"cover_bucket"`
	CoverObjectKey string       `db:"cover_object_key"`
	DurationMs     int64        `db:"duration_ms"`
	AuthorID       int64        `db:"author_id"`
	AuthorName     string       `db:"author_name"`
}

func normalizeRecommendPage(page, pageSize int64) (int64, int64) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}
	return page, pageSize
}

func normalizeAuthorName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "咕噜用户"
	}
	return name
}

func mockRecommendLike(workID int64) types.RecommendLikeInfo {
	// 点赞表还没有接入，先用作品ID生成稳定假数据，避免每次刷新跳变。
	count := int64(37) + (workID*97)%23800
	return types.RecommendLikeInfo{
		Liked: false,
		Count: count,
		Text:  formatLikeText(count),
	}
}

func formatLikeText(count int64) string {
	if count >= 10000 {
		value := float64(count) / 10000
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", value), "0"), ".") + "万"
	}
	return fmt.Sprintf("%d", count)
}
