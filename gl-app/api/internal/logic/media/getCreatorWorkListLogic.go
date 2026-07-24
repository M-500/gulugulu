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

type GetCreatorWorkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询当前用户所有处理完成的作品
func NewGetCreatorWorkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCreatorWorkListLogic {
	return &GetCreatorWorkListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCreatorWorkListLogic) GetCreatorWorkList(req *types.CreatorWorkListReq) (resp *types.CreatorWorkListResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID <= 0 {
		return nil, fmt.Errorf("用户未登录")
	}

	page, pageSize := normalizeCreatorWorkPage(req.Page, req.PageSize)
	var total int64
	if err = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &total, `SELECT COUNT(1)
		FROM work
		WHERE user_id=? AND process_status='succeeded' AND deleted_at IS NULL`, userID); err != nil {
		return nil, fmt.Errorf("查询作品数量失败: %w", err)
	}

	resp = &types.CreatorWorkListResp{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     make([]types.CreatorWorkListItem, 0),
	}
	if total == 0 {
		return resp, nil
	}

	// 作品只有在媒体处理全部成功后才会进入管理列表。审核中、审核失败、
	// 审核通过、定时发布和已发布状态都保留，供CMS按状态展示。
	var rows []creatorWorkRow
	offset := (page - 1) * pageSize
	if err = l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &rows, `SELECT
			w.id,w.type,w.title,w.review_status,w.publish_status,w.review_reason,
			w.scheduled_at,w.published_at,w.created_at,
			COALESCE(ma.formal_bucket,'') AS cover_bucket,
			COALESCE(ma.formal_object_key,'') AS cover_object_key
		FROM work w
		LEFT JOIN media_asset ma ON ma.id=w.cover_asset_id
		WHERE w.user_id=? AND w.process_status='succeeded' AND w.deleted_at IS NULL
		ORDER BY w.created_at DESC,w.id DESC
		LIMIT ? OFFSET ?`, userID, pageSize, offset); err != nil {
		return nil, fmt.Errorf("查询作品列表失败: %w", err)
	}

	for _, row := range rows {
		item := types.CreatorWorkListItem{
			WorkId:        row.ID,
			Type:          row.Type,
			Title:         row.Title,
			ReviewStatus:  row.ReviewStatus,
			PublishStatus: row.PublishStatus,
			ReviewReason:  row.ReviewReason,
			CreatedAt:     row.CreatedAt.Format(time.RFC3339),
		}
		if row.ScheduledAt.Valid {
			item.ScheduledAt = row.ScheduledAt.Time.Format(time.RFC3339)
		}
		if row.PublishedAt.Valid {
			item.PublishedAt = row.PublishedAt.Time.Format(time.RFC3339)
		}
		if row.CoverBucket != "" && row.CoverObjectKey != "" {
			coverURL, signErr := l.svcCtx.MinioClient.PresignedGetObject(
				l.ctx, row.CoverBucket, row.CoverObjectKey, 15*time.Minute, nil,
			)
			if signErr != nil {
				l.Errorf("生成作品%d封面地址失败: %v", row.ID, signErr)
			} else {
				item.CoverUrl = coverURL.String()
			}
		}
		resp.List = append(resp.List, item)
	}

	return resp, nil
}

type creatorWorkRow struct {
	ID             int64        `db:"id"`
	Type           string       `db:"type"`
	Title          string       `db:"title"`
	ReviewStatus   string       `db:"review_status"`
	PublishStatus  string       `db:"publish_status"`
	ReviewReason   string       `db:"review_reason"`
	ScheduledAt    sql.NullTime `db:"scheduled_at"`
	PublishedAt    sql.NullTime `db:"published_at"`
	CreatedAt      time.Time    `db:"created_at"`
	CoverBucket    string       `db:"cover_bucket"`
	CoverObjectKey string       `db:"cover_object_key"`
}

func normalizeCreatorWorkPage(page, pageSize int64) (int64, int64) {
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
