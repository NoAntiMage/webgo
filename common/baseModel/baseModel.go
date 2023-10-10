package baseModel

import (
	"goweb/common/stringx"
)

const (
	DefaultTag string = "default"
)

const (
	IsDeletedKey = "IsDeleted"
)

var (
	IsDeleted = stringx.SnakeFormat(IsDeletedKey)

	BaseFields    = stringx.StructFieldNames(&BaseModel{})
	BaseSqlFields = stringx.FieldsFromStructToSql(BaseFields)
)

type BaseModel struct {
	Id        int64 `json:"-"`
	IsDeleted int   `json:"-"`
}
