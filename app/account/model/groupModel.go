package model

import (
	"context"
	"goweb/app/account/types"
	"goweb/common/baseModel"
	mctx "goweb/common/context"
	"goweb/common/jsonx"
	"goweb/common/stringx"
)

type GroupModel interface {
	groupModel
	ListByIds(c context.Context, ids []int64) (groups []Group, _ error)
	PageList(c context.Context, limit int, offset int, query types.GroupQuery) ([]Group, error)
	Search(c context.Context, req types.GroupSearchReq) ([]Group, error)
}

type customGroupModel struct {
	*defaultGroupModel
}

var _ GroupModel = (*customGroupModel)(nil)

func NewGroupModel(modelCtx mctx.ModelContext) GroupModel {
	return &customGroupModel{
		defaultGroupModel: newGroupModel(modelCtx),
	}
}

func (g *customGroupModel) ListByIds(c context.Context, ids []int64) (groups []Group, _ error) {
	b := new(stringx.SqlBuilder)
	stmt, params, err :=
		b.Select(groupSqlFields...).
			From(g.table).
			Where(b.In("id", ids)).
			Result()
	if err != nil {
		return nil, err
	}
	if err := g.ModelContext.SelectContext(c, &groups, stmt, params); err != nil {
		return nil, err
	}
	return groups, nil
}

func (g *customGroupModel) PageList(c context.Context, limit int, offset int, query types.GroupQuery) ([]Group, error) {
	kvs, err := jsonx.Struct2MapWithRule(&query, jsonx.RuleWithoutZero)
	if err != nil {
		return nil, err
	}

	b := new(stringx.SqlBuilder)
	b = b.Select(groupSqlFields...).
		From(g.table).
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

	groups := []Group{}
	if err = g.SelectContext(c, &groups, stmt, params...); err != nil {
		return nil, err
	}
	return groups, nil
}

func (g *customGroupModel) Search(c context.Context, req types.GroupSearchReq) ([]Group, error) {
	b := new(stringx.SqlBuilder)
	b = b.Select(groupSqlFields...).
		From(g.table).
		Where(b.Eq(baseModel.IsDeleted, 0))

	for _, cond := range req.Query {
		b = stringx.WithCond(b, cond)
	}
	stmt, params, err := b.Result()
	if err != nil {
		return nil, err
	}

	var groups []Group
	if err := g.SelectContext(c, &groups, stmt, params...); err != nil {
		return nil, err
	}

	return groups, nil
}
