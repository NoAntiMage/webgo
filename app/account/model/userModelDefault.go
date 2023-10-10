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
	UserTable = "user"
)

var (
	userFields              = stringx.StructFieldNames(&User{})
	userSqlFields           = stringx.FieldsFromStructToSql(userFields)
	userSqlFieldsAutoSet    = stringx.Remove(userSqlFields, baseModel.BaseSqlFields...)
	userRows                = strings.Join(userSqlFields, ",")
	userRowsAutoSet         = strings.Join(userSqlFieldsAutoSet, ",")
	userRowsWithPlaceHolder = stringx.RowsWithPlaceHolder(userSqlFieldsAutoSet)
)

type (
	userModel interface {
		Insert(ctx context.Context, data *User) (sql.Result, error)
		Select(ctx context.Context, id int64) (*User, error)
		Update(ctx context.Context, data *User) (sql.Result, error)
		Delete(ctx context.Context, id int64) (sql.Result, error)
	}

	defaultUserModel struct {
		ctx.ModelContext
		table string
		m2m   map[string]string
	}

	User struct {
		baseModel.BaseModel
		Username string
		Realname string
		Email    string `default:"devops@qq.com"`
		Age      int    `default:"18"`
		Gender   int
	}
)

var _ userModel = (*defaultUserModel)(nil)

func newUserModel(modelCtx ctx.ModelContext) *defaultUserModel {
	return &defaultUserModel{
		ModelContext: modelCtx,
		table:        UserTable,
		m2m: map[string]string{
			GroupTable: "user_groupx_m2m",
		},
	}
}

func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
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

func (m *defaultUserModel) Select(ctx context.Context, id int64) (*User, error) {
	b := stringx.NewSqlBuilder()
	stmt, err := b.
		Select(userSqlFields...).
		From(m.table).
		Where(b.Eq(baseModel.IsDeleted, 0)).
		And(b.Eq("id", id)).
		Limit(1).
		RawString()
	if err != nil {
		return nil, err
	}

	u := new(User)

	err = m.GetContext(ctx, u, stmt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (m *defaultUserModel) Update(ctx context.Context, data *User) (sql.Result, error) {
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

func (m *defaultUserModel) Delete(ctx context.Context, id int64) (sql.Result, error) {
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
