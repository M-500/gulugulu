package work

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkTopicModel = (*customWorkTopicModel)(nil)

type (
	// WorkTopicModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkTopicModel.
	WorkTopicModel interface {
		workTopicModel
	}

	customWorkTopicModel struct {
		*defaultWorkTopicModel
	}
)

// NewWorkTopicModel returns a model for the database table.
func NewWorkTopicModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WorkTopicModel {
	return &customWorkTopicModel{
		defaultWorkTopicModel: newWorkTopicModel(conn, c, opts...),
	}
}
