package work

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MediaProcessTaskModel = (*customMediaProcessTaskModel)(nil)

type (
	// MediaProcessTaskModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMediaProcessTaskModel.
	MediaProcessTaskModel interface {
		mediaProcessTaskModel
	}

	customMediaProcessTaskModel struct {
		*defaultMediaProcessTaskModel
	}
)

// NewMediaProcessTaskModel returns a model for the database table.
func NewMediaProcessTaskModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MediaProcessTaskModel {
	return &customMediaProcessTaskModel{
		defaultMediaProcessTaskModel: newMediaProcessTaskModel(conn, c, opts...),
	}
}
