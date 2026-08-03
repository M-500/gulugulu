package media

import (
	"context"
	"fmt"
	"time"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuditWorkDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询审核作品详情
func NewGetAuditWorkDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuditWorkDetailLogic {
	return &GetAuditWorkDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAuditWorkDetailLogic) GetAuditWorkDetail(req *types.WorkIdReq) (*types.AuditWorkDetailResp, error) {
	if err := ensureAuditPermission(l.svcCtx, ctxdata.GetUidFromCtx(l.ctx)); err != nil {
		return nil, err
	}
	row, err := l.svcCtx.WorkRepo.FindAuditDetail(l.ctx, req.WorkId)
	if err != nil {
		return nil, fmt.Errorf("审核作品不存在: %w", err)
	}
	assets, err := l.svcCtx.WorkRepo.ListReadyAssets(l.ctx, req.WorkId)
	if err != nil {
		return nil, fmt.Errorf("查询审核素材失败: %w", err)
	}
	topics, err := l.svcCtx.WorkRepo.ListTopics(l.ctx, req.WorkId)
	if err != nil {
		return nil, fmt.Errorf("查询作品话题失败: %w", err)
	}
	resp := &types.AuditWorkDetailResp{
		WorkId: row.ID, Type: row.Type, Title: row.Title, Content: row.Content,
		Visibility: row.Visibility, Original: row.Original, AuthorId: row.AuthorID,
		AuthorName: row.AuthorName, AuthorEmail: row.AuthorEmail, ReviewStatus: row.ReviewStatus,
		PublishStatus: row.PublishStatus, ReviewReason: row.ReviewReason,
		SubmittedAt: row.CreatedAt.Format(time.RFC3339), CoverUrl: l.presignAuditDetailAsset(row.CoverBucket, row.CoverObject),
		Assets: make([]types.WorkAssetItem, 0, len(assets)), Topics: make([]types.WorkTopicItem, 0, len(topics)),
	}
	if row.ScheduledAt != nil {
		resp.ScheduledAt = row.ScheduledAt.Format(time.RFC3339)
	}
	if row.ReviewedAt != nil {
		resp.ReviewedAt = row.ReviewedAt.Format(time.RFC3339)
	}
	for _, asset := range assets {
		resp.Assets = append(resp.Assets, types.WorkAssetItem{
			MediaId: asset.MediaID, Role: asset.Role, Sort: asset.Sort,
			Url: l.presignAuditDetailAsset(asset.Bucket, asset.ObjectKey),
		})
		if asset.Role == "video" {
			resp.DurationMs = asset.DurationMs
			resp.VideoPlaylist = fmt.Sprintf("/api/v1/audit/works/%d/playlist", row.ID)
		}
	}
	for _, topic := range topics {
		resp.Topics = append(resp.Topics, types.WorkTopicItem{TopicId: topic.ID, Name: topic.Name})
	}
	return resp, nil
}

func (l *GetAuditWorkDetailLogic) presignAuditDetailAsset(bucket, objectKey string) string {
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
