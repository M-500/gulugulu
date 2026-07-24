package work

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WorkAssetModel = (*customWorkAssetModel)(nil)

type (
	// WorkAssetModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWorkAssetModel.
	WorkAssetModel interface {
		workAssetModel
	}

	customWorkAssetModel struct {
		*defaultWorkAssetModel
	}
)

// NewWorkAssetModel returns a model for the database table.
func NewWorkAssetModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) WorkAssetModel {
	return &customWorkAssetModel{
		defaultWorkAssetModel: newWorkAssetModel(conn, c, opts...),
	}
}
