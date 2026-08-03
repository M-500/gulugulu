package comment_repo

import "context"

type CommentDao interface {
	FindOneByResourceID(ctx context.Context, resourceID int64) ([]*CommentModel, error)
}
