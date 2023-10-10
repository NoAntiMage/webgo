package types

import (
	"goweb/common/condition"
	"goweb/common/condition/pagination"
)

type GroupCreateReq struct {
	Name string
}

type GroupCreateResp struct {
	LastId int64
}

type GroupResp struct {
	Id   int64
	Name string
}

type GroupsResp struct {
	Ids []int64
}

type GroupUpdateReq struct {
	Id   int64 `validate:"gt=0"`
	Name string
}

type GroupListReq struct {
	pagination.Page
	GroupQuery
}

type GroupQuery struct {
	Id   int64  `form:"id"`
	Name string `form:"name"`
}

//restful api where
type GroupSearchReq struct {
	// pagination.Page
	Query []condition.QueryCondition
}
