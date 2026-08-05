package comment

import (
	"context"
	"fmt"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentReplyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentReplyListLogic {
	return &GetCommentReplyListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetCommentReplyListLogic) GetCommentReplyList(req *types.CommentReplyListReq) (*types.CommentListResp, error) {
	if req == nil || req.CommentId <= 0 {
		return nil, fmt.Errorf("评论ID不正确")
	}
	page, pageSize := normalizeCommentPage(req.Page, req.PageSize)
	total, rows, err := l.svcCtx.CommentRepo.ListReplies(l.ctx, req.CommentId, pageSize, (page-1)*pageSize)
	if err != nil {
		l.Errorf("查询评论回复失败: %v", err)
		return nil, fmt.Errorf("查询回复失败")
	}
	resp := &types.CommentListResp{Total: total, Page: page, PageSize: pageSize, HasMore: page*pageSize < total, List: make([]types.CommentItem, 0, len(rows))}
	for _, row := range rows {
		resp.List = append(resp.List, buildCommentItem(l.ctx, l.svcCtx, row))
	}
	return resp, nil
}
