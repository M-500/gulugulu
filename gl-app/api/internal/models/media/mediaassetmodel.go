package media

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MediaAssetModel = (*customMediaAssetModel)(nil)

type (
	// MediaAssetModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMediaAssetModel.
	MediaAssetModel interface {
		mediaAssetModel
	}

	customMediaAssetModel struct {
		*defaultMediaAssetModel
	}
)

// NewMediaAssetModel returns a model for the database table.
func NewMediaAssetModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MediaAssetModel {
	return &customMediaAssetModel{
		defaultMediaAssetModel: newMediaAssetModel(conn, c, opts...),
	}
}
