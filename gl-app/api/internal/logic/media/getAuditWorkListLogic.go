package media

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuditWorkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询审核中心作品列表
func NewGetAuditWorkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuditWorkListLogic {
	return &GetAuditWorkListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAuditWorkListLogic) GetAuditWorkList(req *types.AuditWorkListReq) (resp *types.AuditWorkListResp, err error) {
	if err = ensureAuditPermission(l.svcCtx, ctxdata.GetUidFromCtx(l.ctx)); err != nil {
		return nil, err
	}

	page, pageSize := normalizeCreatorWorkPage(req.Page, req.PageSize)
	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "" {
		status = "pending"
	}
	if status != "all" && status != "pending" && status != "approved" && status != "rejected" {
		return nil, fmt.Errorf("status必须是all、pending、approved或rejected")
	}
	workType := strings.ToLower(strings.TrimSpace(req.Type))
	if workType != "" && workType != "image" && workType != "video" {
		return nil, fmt.Errorf("type必须是image或video")
	}

	var summary auditWorkSummary
	if err = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &summary, `SELECT
			COALESCE(SUM(review_status='pending_review'),0) AS pending_count,
			COALESCE(SUM(review_status='approved'),0) AS approved_count,
			COALESCE(SUM(review_status='rejected'),0) AS rejected_count
		FROM work WHERE process_status='succeeded' AND deleted_at IS NULL`); err != nil {
		return nil, fmt.Errorf("查询审核统计失败: %w", err)
	}

	whereClause, args := buildAuditWorkFilter(status, workType, strings.TrimSpace(req.Keyword))
	var total int64
	if err = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &total,
		"SELECT COUNT(1) FROM work w JOIN user u ON u.id=w.user_id "+whereClause, args...); err != nil {
		return nil, fmt.Errorf("查询审核作品数量失败: %w", err)
	}

	resp = &types.AuditWorkListResp{
		Total: total, Page: page, PageSize: pageSize,
		PendingCount: summary.PendingCount, ApprovedCount: summary.ApprovedCount,
		RejectedCount: summary.RejectedCount, List: make([]types.AuditWorkListItem, 0),
	}
	if total == 0 {
		return resp, nil
	}

	var rows []auditWorkListRow
	query := `SELECT w.id,w.type,w.title,LEFT(w.content,160) AS content_excerpt,w.visibility,w.user_id,
			COALESCE(NULLIF(u.nickname,''),u.email) AS author_name,w.review_status,w.publish_status,
			w.review_reason,w.created_at,w.reviewed_at,
			COALESCE(video_asset.duration_ms,0) AS duration_ms,
			COALESCE(cover.formal_bucket,'') AS cover_bucket,
			COALESCE(cover.formal_object_key,'') AS cover_object_key
		FROM work w
		JOIN user u ON u.id=w.user_id
		LEFT JOIN media_asset cover ON cover.id=w.cover_asset_id
		LEFT JOIN work_asset video_relation ON video_relation.work_id=w.id AND video_relation.role='video'
		LEFT JOIN media_asset video_asset ON video_asset.id=video_relation.media_asset_id
		` + whereClause + `
		ORDER BY CASE w.review_status WHEN 'pending_review' THEN 0 ELSE 1 END,w.created_at ASC,w.id ASC
		LIMIT ? OFFSET ?`
	args = append(args, pageSize, (page-1)*pageSize)
	if err = l.svcCtx.SqlConn.QueryRowsCtx(l.ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("查询审核作品列表失败: %w", err)
	}

	for _, row := range rows {
		item := types.AuditWorkListItem{
			WorkId: row.ID, Type: row.Type, Title: row.Title, ContentExcerpt: row.ContentExcerpt,
			DurationMs: row.DurationMs, Visibility: row.Visibility, AuthorId: row.AuthorID,
			AuthorName: row.AuthorName, ReviewStatus: row.ReviewStatus, PublishStatus: row.PublishStatus,
			ReviewReason: row.ReviewReason, SubmittedAt: row.CreatedAt.Format(time.RFC3339),
		}
		if row.ReviewedAt.Valid {
			item.ReviewedAt = row.ReviewedAt.Time.Format(time.RFC3339)
		}
		item.CoverUrl = l.presignAuditAsset(row.CoverBucket, row.CoverObjectKey)
		resp.List = append(resp.List, item)
	}

	return resp, nil
}

type auditWorkSummary struct {
	PendingCount  int64 `db:"pending_count"`
	ApprovedCount int64 `db:"approved_count"`
	RejectedCount int64 `db:"rejected_count"`
}

type auditWorkListRow struct {
	ID             int64        `db:"id"`
	Type           string       `db:"type"`
	Title          string       `db:"title"`
	ContentExcerpt string       `db:"content_excerpt"`
	DurationMs     int64        `db:"duration_ms"`
	Visibility     string       `db:"visibility"`
	AuthorID       int64        `db:"user_id"`
	AuthorName     string       `db:"author_name"`
	ReviewStatus   string       `db:"review_status"`
	PublishStatus  string       `db:"publish_status"`
	ReviewReason   string       `db:"review_reason"`
	CreatedAt      time.Time    `db:"created_at"`
	ReviewedAt     sql.NullTime `db:"reviewed_at"`
	CoverBucket    string       `db:"cover_bucket"`
	CoverObjectKey string       `db:"cover_object_key"`
}

func buildAuditWorkFilter(status, workType, keyword string) (string, []any) {
	conditions := []string{"w.process_status='succeeded'", "w.deleted_at IS NULL"}
	args := make([]any, 0, 4)
	switch status {
	case "pending":
		conditions = append(conditions, "w.review_status='pending_review'")
	case "approved":
		conditions = append(conditions, "w.review_status='approved'")
	case "rejected":
		conditions = append(conditions, "w.review_status='rejected'")
	default:
		conditions = append(conditions, "w.review_status IN ('pending_review','approved','rejected')")
	}
	if workType != "" {
		conditions = append(conditions, "w.type=?")
		args = append(args, workType)
	}
	if keyword != "" {
		escaped := "%" + escapeLikeKeyword(keyword) + "%"
		conditions = append(conditions, "(w.title LIKE ? ESCAPE '\\\\' OR u.nickname LIKE ? ESCAPE '\\\\' OR u.email LIKE ? ESCAPE '\\\\')")
		args = append(args, escaped, escaped, escaped)
	}
	return "WHERE " + strings.Join(conditions, " AND "), args
}

func (l *GetAuditWorkListLogic) presignAuditAsset(bucket, objectKey string) string {
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
