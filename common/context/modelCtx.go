package context

import (
	"goweb/common/sqlc"
)

type ModelContext struct {
	*sqlc.SqlConn
}

//modelCtx contain resource used in model operation, [db, cache, redis, file, ...]
func NewModelCtx() ModelContext {
	return ModelContext{
		SqlConn: sqlc.NewSqlConnect(),
	}
}
