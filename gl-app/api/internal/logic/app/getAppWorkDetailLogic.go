package app

import (
	"context"
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
	return &GetAppWorkDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAppWorkDetailLogic) GetAppWorkDetail(req *types.AppWorkDetailReq) (*types.AppWorkDetailResp, error) {
	if req == nil || req.WorkId <= 0 {
		return nil, fmt.Errorf("作品ID不正确")
	}
	row, err := l.svcCtx.WorkRepo.FindPublicDetail(l.ctx, req.WorkId)
	if err != nil {
		l.Errorf("查询App公开作品%d失败: %v", req.WorkId, err)
		return nil, fmt.Errorf("作品不存在或暂不可见")
	}
	assets, err := l.svcCtx.WorkRepo.ListReadyAssets(l.ctx, req.WorkId)
	if err != nil {
		return nil, fmt.Errorf("查询作品资源失败")
	}
	topics, err := l.svcCtx.WorkRepo.ListTopics(l.ctx, req.WorkId)
	if err != nil {
		return nil, fmt.Errorf("查询作品话题失败")
	}
	resp := &types.AppWorkDetailResp{
		WorkId: row.ID, Type: row.Type, Title: row.Title, Content: row.Content,
		Author: types.RecommendAuthorItem{UserId: row.AuthorID, NickName: normalizeAuthorName(row.AuthorName),
			AvatarUrl: buildAppPublicAvatarURL(l.svcCtx, row.AuthorAvatar)},
		Like: mockRecommendLike(row.ID), FavoriteCount: 19 + (row.ID*53)%5000,
		CommentCount: 8 + (row.ID*31)%1000, ShareCount: 3 + (row.ID*17)%400,
		Assets: make([]types.AppWorkAssetItem, 0, len(assets)), Topics: make([]types.AppWorkTopicItem, 0, len(topics)),
	}
	if row.PublishedAt != nil {
		resp.PublishedAt = row.PublishedAt.Format(time.RFC3339)
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
			MediaId: asset.MediaID, Role: asset.Role, Sort: asset.Sort, Url: assetURL,
			Width: asset.Width, Height: asset.Height,
		})
	}
	for _, topic := range topics {
		resp.Topics = append(resp.Topics, types.AppWorkTopicItem{TopicId: topic.ID, Name: topic.Name})
	}
	return resp, nil
}

func (l *GetAppWorkDetailLogic) presignAppWorkAsset(bucket, objectKey string) string {
	if bucket == "" || objectKey == "" {
		return ""
	}
	value, err := l.svcCtx.MinioClient.PresignedGetObject(l.ctx, bucket, objectKey, 30*time.Minute, nil)
	if err != nil {
		l.Errorf("生成App作品资源签名地址失败: %v", err)
		return ""
	}
	return value.String()
}
