package comment_repo

import "context"

// CommentRepo 是评论业务唯一依赖的持久化接口，Service 层不直接查询数据库。
type CommentRepo interface {
	Create(ctx context.Context, input CreateInput) (*CommentModel, error)
	FindByID(ctx context.Context, commentID int64) (*CommentView, error)
	ListRoots(ctx context.Context, workID, limit, offset int64) (int64, []CommentView, error)
	ListReplies(ctx context.Context, rootID, limit, offset int64) (int64, []CommentView, error)
}

type commentRepoImpl struct{ dao CommentDAO }

func NewCommentRepo(dao CommentDAO) CommentRepo { return &commentRepoImpl{dao: dao} }

func (r *commentRepoImpl) Create(ctx context.Context, input CreateInput) (*CommentModel, error) {
	return r.dao.Create(ctx, input)
}

func (r *commentRepoImpl) FindByID(ctx context.Context, commentID int64) (*CommentView, error) {
	return r.dao.FindByID(ctx, commentID)
}

func (r *commentRepoImpl) ListRoots(ctx context.Context, workID, limit, offset int64) (int64, []CommentView, error) {
	return r.dao.ListRoots(ctx, workID, limit, offset)
}

func (r *commentRepoImpl) ListReplies(ctx context.Context, rootID, limit, offset int64) (int64, []CommentView, error) {
	return r.dao.ListReplies(ctx, rootID, limit, offset)
}
