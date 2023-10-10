package model

import (
	"context"
	"database/sql"
	"goweb/common/baseModel"
	ctx "goweb/common/context"
	"goweb/common/jsonx"
	"goweb/common/stringx"
	"strings"
)

const (
	GroupTable = "groupx"
)

var (
	groupFields              = stringx.StructFieldNames(&Group{})
	groupSqlFields           = stringx.FieldsFromStructToSql(groupFields)
	groupSqlFieldsAutoSet    = stringx.Remove(groupSqlFields, baseModel.BaseSqlFields...)
	groupRows                = strings.Join(groupSqlFields, ",")
	groupRowsAutoSet         = strings.Join(groupSqlFieldsAutoSet, ",")
	groupRowsWithPlaceHolder = stringx.RowsWithPlaceHolder(groupSqlFieldsAutoSet)
)

type (
	groupModel interface {
		Insert(ctx context.Context, data *Group) (sql.Result, error)
		Select(ctx context.Context, id int64) (*Group, error)
		Update(ctx context.Context, data *Group) (sql.Result, error)
		Delete(ctx context.Context, id int64) (sql.Result, error)
	}

	defaultGroupModel struct {
		ctx.ModelContext
		table string
		m2m   map[string]string
	}

	Group struct {
		baseModel.BaseModel
		Name string
		//	Comment string
	}
)

var _ groupModel = (*defaultGroupModel)(nil)

func newGroupModel(modelCtx ctx.ModelContext) *defaultGroupModel {
	return &defaultGroupModel{
		ModelContext: modelCtx,
		table:        GroupTable,
		m2m: map[string]string{
			UserTable: "user_groupx_m2m",
		},
	}
}

func (m *defaultGroupModel) Insert(ctx context.Context, data *Group) (sql.Result, error) {
	b := stringx.NewSqlBuilder()
	kvs, err := jsonx.Struct2Map(data)
	if err != nil {
		return nil, err
	}

	stmt, params, err := b.Insert(m.table, kvs).Result()
	if err != nil {
		return nil, err
	}

	return m.ExecContextWithLog(ctx, stmt, params...)
}

func (m *defaultGroupModel) Select(ctx context.Context, id int64) (*Group, error) {
	b := stringx.NewSqlBuilder()
	stmt, err := b.
		Select(groupSqlFields...).
		From(m.table).
		Where(b.Eq(baseModel.IsDeleted, 0)).
		And(b.Eq("id", id)).
		Limit(1).
		RawString()
	if err != nil {
		return nil, err
	}

	u := new(Group)

	err = m.GetContext(ctx, u, stmt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (m *defaultGroupModel) Update(ctx context.Context, data *Group) (sql.Result, error) {
	b := stringx.NewSqlBuilder()
	kvs, err := jsonx.Struct2Map(data)
	if err != nil {
		return nil, err
	}

	stmt, params, err := b.
		Update(m.table, kvs).
		Where(b.Eq("id", data.Id)).
		Result()
	if err != nil {
		return nil, err
	}

	return m.ExecContextWithLog(ctx, stmt, params...)

}
func (m *defaultGroupModel) Delete(ctx context.Context, id int64) (sql.Result, error) {
	b := stringx.NewSqlBuilder()
	kvs := map[string]any{
		baseModel.IsDeletedKey: 1,
	}
	stmt, err := b.
		Update(m.table, kvs).
		Where(b.Eq("id", id)).
		RawString()
	if err != nil {
		return nil, err
	}
	return m.ExecContextWithLog(ctx, stmt)
}
