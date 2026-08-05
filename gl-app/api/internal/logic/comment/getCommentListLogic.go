package comment

import (
	"context"
	"fmt"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetCommentListLogic) GetCommentList(req *types.CommentListReq) (*types.CommentListResp, error) {
	if req == nil || req.WorkId <= 0 {
		return nil, fmt.Errorf("作品ID不正确")
	}
	if _, err := l.svcCtx.WorkRepo.FindPublicDetail(l.ctx, req.WorkId); err != nil {
		return nil, fmt.Errorf("作品不存在或暂不可见")
	}
	page, pageSize := normalizeCommentPage(req.Page, req.PageSize)
	total, rows, err := l.svcCtx.CommentRepo.ListRoots(l.ctx, req.WorkId, pageSize, (page-1)*pageSize)
	if err != nil {
		l.Errorf("查询作品评论失败: %v", err)
		return nil, fmt.Errorf("查询评论失败")
	}
	resp := &types.CommentListResp{Total: total, Page: page, PageSize: pageSize, HasMore: page*pageSize < total, List: make([]types.CommentItem, 0, len(rows))}
	for _, row := range rows {
		item := buildCommentItem(l.ctx, l.svcCtx, row)
		if row.ReplyCount > 0 {
			_, replies, replyErr := l.svcCtx.CommentRepo.ListReplies(l.ctx, row.ID, 2, 0)
			if replyErr == nil {
				for _, reply := range replies {
					item.Replies = append(item.Replies, buildCommentItem(l.ctx, l.svcCtx, reply))
				}
				item.ReplyHasMore = int64(len(replies)) < row.ReplyCount
			}
		}
		resp.List = append(resp.List, item)
	}
	return resp, nil
}
