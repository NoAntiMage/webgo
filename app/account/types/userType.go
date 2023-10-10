package types

import (
	"goweb/common/condition"
	"goweb/common/condition/pagination"
)

type UserCreateReq struct {
	Username string
	Realname string
	Email    string
	Age      int
	Gender   int
}

type UserCreateResp struct {
	LastId int64
}

type UserResp struct {
	Id       int64
	Username string
	Realname string
	Email    string
	Age      int
	Gender   int
}

type UsersResp struct {
	Ids []int64
}

type UserUpdateReq struct {
	Id       int64
	Username string
	Realname string
	Email    string
	Age      int
	Gender   int
}

type UserListReq struct {
	pagination.Page
	UserQuery
}

type UserQuery struct {
	Id       int64  `form:"id"`
	Username string `form:"username"`
	Realname string `form:"realname"`
	Email    string `form:"email"`
	Age      int    `form:"age"`
	Gender   int    `form:"gender"`
}

//restful api where
type UserSearchReq struct {
	// pagination.Page
	Query []condition.QueryCondition
}
