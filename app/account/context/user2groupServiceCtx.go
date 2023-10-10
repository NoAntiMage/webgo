package context

import (
	"goweb/app/account/model"
	ctx "goweb/common/context"
)

type User2GroupServiceContext struct {
	User2GroupModel model.User2GroupModel
}

func NewUser2GroupServiceContext() User2GroupServiceContext {
	mctx := ctx.NewModelCtx()
	u2gMo := model.NewUser2GroupModel(mctx)
	return User2GroupServiceContext{
		User2GroupModel: u2gMo,
	}
}
