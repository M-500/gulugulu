package work

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkModel = (*customWorkModel)(nil)

type (
	// WorkModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkModel.
	WorkModel interface {
		workModel
	}

	customWorkModel struct {
		*defaultWorkModel
	}
)

// NewWorkModel returns a model for the database table.
func NewWorkModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WorkModel {
	return &customWorkModel{
		defaultWorkModel: newWorkModel(conn, c, opts...),
	}
}
