package context

import (
	"goweb/app/account/model"
	ctx "goweb/common/context"
)

type GroupServiceContext struct {
	GroupModel model.GroupModel
}

func NewGroupServiceContext() GroupServiceContext {
	mctx := ctx.NewModelCtx()
	gm := model.NewGroupModel(mctx)
	return GroupServiceContext{
		GroupModel: gm,
	}
}
