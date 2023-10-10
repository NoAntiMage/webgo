package types

type GroupsInUserReq struct {
	UserId   int64   `validate:"required,gt=0"`
	GroupIds []int64 `validate:"required,gt=0"`
}

type UsersInGroupReq struct {
	GroupId int64   `validate:"required,gt=0"`
	UserIds []int64 `validate:"required,gt=0"`
}

type GroupsInUserResp struct {
	UserId   int64
	GroupIds []int64
}

type UsersInGroupResp struct {
	GroupId int64
	UserIds []int64
}
