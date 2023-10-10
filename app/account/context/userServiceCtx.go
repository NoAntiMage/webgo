package context

import (
	"goweb/app/account/model"
	ctx "goweb/common/context"
)

type UserServiceContext struct {
	UserModel model.UserModel
}

func NewUserServiceContext() UserServiceContext {
	mctx := ctx.NewModelCtx()
	um := model.NewUserModel(mctx)
	return UserServiceContext{
		UserModel: um,
	}
}
