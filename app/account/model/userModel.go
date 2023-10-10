package model

import (
	"context"
	"goweb/app/account/types"
	"goweb/common/baseModel"
	ctx "goweb/common/context"
	"goweb/common/jsonx"
	"goweb/common/stringx"
)

type UserModel interface {
	userModel
	ListByIds(c context.Context, ids []int64) ([]User, error)
	PageList(c context.Context, limit int, offset int, query types.UserQuery) ([]User, error)
	Search(c context.Context, req types.UserSearchReq) ([]User, error)
}

type customUserModel struct {
	*defaultUserModel
}

var _ UserModel = (*customUserModel)(nil)

func NewUserModel(modelCtx ctx.ModelContext) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(modelCtx),
	}
}

func (u *customUserModel) ListByIds(c context.Context, ids []int64) ([]User, error) {
	b := new(stringx.SqlBuilder)
	stmt, params, err :=
		b.Select(userSqlFields...).
			From(u.table).
			Where(b.In("id", ids)).
			Result()
	if err != nil {
		return nil, err
	}
	users := []User{}
	if err := u.ModelContext.SelectContext(c, &users, stmt, params...); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *customUserModel) PageList(c context.Context, limit int, offset int, query types.UserQuery) ([]User, error) {
	kvs, err := jsonx.Struct2MapWithRule(&query, jsonx.RuleWithoutZero)
	if err != nil {
		return nil, err
	}

	b := new(stringx.SqlBuilder)
	b = b.Select(userSqlFields...).
		From(u.table).
		Where(b.Eq(baseModel.IsDeleted, 0))

	for k, v := range kvs {
		b = b.And(b.Eq(k, v))
	}

	b = b.Limit(limit).
		Offset(offset)

	stmt, params, err := b.Result()
	if err != nil {
		return nil, err
	}

	users := []User{}
	if err = u.SelectContext(c, &users, stmt, params...); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *customUserModel) Search(c context.Context, req types.UserSearchReq) ([]User, error) {
	b := new(stringx.SqlBuilder)
	b = b.Select(userSqlFields...).
		From(u.table).
		Where(b.Eq(baseModel.IsDeleted, 0))

	for _, cond := range req.Query {
		b = stringx.WithCond(b, cond)
	}
	stmt, params, err := b.Result()
	if err != nil {
		return nil, err
	}

	var users []User
	if err := u.SelectContext(c, &users, stmt, params...); err != nil {
		return nil, err
	}

	return users, nil
}
