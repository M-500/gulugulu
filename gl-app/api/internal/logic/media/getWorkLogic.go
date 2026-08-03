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

type GetWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取作品详情
func NewGetWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkLogic {
	return &GetWorkLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetWorkLogic) GetWork(req *types.WorkIdReq) (*types.WorkDetailResp, error) {
	row, err := l.svcCtx.WorkRepo.FindOwnerDetail(l.ctx, req.WorkId, ctxdata.GetUidFromCtx(l.ctx))
	if err != nil {
		return nil, fmt.Errorf("作品不存在: %w", err)
	}
	assets, err := l.svcCtx.WorkRepo.ListReadyAssets(l.ctx, req.WorkId)
	if err != nil {
		return nil, err
	}
	topics, err := l.svcCtx.WorkRepo.ListTopics(l.ctx, req.WorkId)
	if err != nil {
		return nil, err
	}
	resp := &types.WorkDetailResp{
		WorkId: row.ID, Type: row.Type, Title: row.Title, Content: row.Content, Visibility: row.Visibility,
		CollectionId: row.CollectionID, Original: row.Original, ProcessStatus: row.ProcessStatus,
		ReviewStatus: row.ReviewStatus, PublishStatus: row.PublishStatus, ReviewReason: row.ReviewReason,
		CreatedAt: row.CreatedAt.Format(time.RFC3339), Assets: []types.WorkAssetItem{}, Topics: []types.WorkTopicItem{},
	}
	if row.ScheduledAt != nil {
		resp.ScheduledAt = row.ScheduledAt.Format(time.RFC3339)
	}
	if row.PublishedAt != nil {
		resp.PublishedAt = row.PublishedAt.Format(time.RFC3339)
	}
	for _, asset := range assets {
		url := l.presignWorkAsset(asset.Bucket, asset.ObjectKey)
		switch asset.Role {
		case "cover":
			resp.CoverUrl = url
		case "video":
			resp.DurationMs = asset.DurationMs
			resp.VideoPlaylist = fmt.Sprintf("/api/v1/works/%d/playlist", row.ID)
		}
		resp.Assets = append(resp.Assets, types.WorkAssetItem{MediaId: asset.MediaID, Role: asset.Role, Sort: asset.Sort, Url: url})
	}
	for _, topic := range topics {
		resp.Topics = append(resp.Topics, types.WorkTopicItem{TopicId: topic.ID, Name: topic.Name})
	}
	return resp, nil
}

func (l *GetWorkLogic) presignWorkAsset(bucket, objectKey string) string {
	if bucket == "" || objectKey == "" {
		return ""
	}
	value, err := l.svcCtx.MinioClient.PresignedGetObject(l.ctx, bucket, objectKey, 30*time.Minute, nil)
	if err != nil {
		l.Errorf("生成作品资源地址失败: %v", err)
		return ""
	}
	return value.String()
}
