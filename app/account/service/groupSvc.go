package service

import (
	"context"
	"fmt"
	svcCtx "goweb/app/account/context"
	"goweb/app/account/model"
	"goweb/app/account/types"
	"goweb/common/copier"
	"goweb/common/singleflight"
)

type groupService interface {
	GetGroupById(ctx context.Context, id int64) (*types.GroupResp, error)
	GroupCreate(ctx context.Context, req *types.GroupCreateReq) (id int64, _ error)
	GroupUpdate(ctx context.Context, req *types.GroupUpdateReq) (int64, error)
	GroupDelete(ctx context.Context, id int64) (rowsAffected int64, _ error)
	GroupList(ctx context.Context, req types.GroupListReq) (groupsResp []types.GroupResp, _ error)
	GroupSearch(ctx context.Context, req types.GroupSearchReq) (groupResp []types.GroupResp, _ error)
}

type GroupService struct {
	SvcCtx svcCtx.GroupServiceContext
}

var _ groupService = (*GroupService)(nil)

func NewGroupService() groupService {
	return &GroupService{
		SvcCtx: svcCtx.NewGroupServiceContext(),
	}
}

func (gs *GroupService) GetGroupById(ctx context.Context, id int64) (*types.GroupResp, error) {
	u, err := gs.SvcCtx.GroupModel.Select(ctx, id)
	if err != nil {
		return nil, err
	}

	groupResp := types.GroupResp{}

	err = copier.Copy(&groupResp, u)
	if err != nil {
		return nil, err
	}

	return &groupResp, err
}

func (gs *GroupService) GroupCreate(ctx context.Context, req *types.GroupCreateReq) (id int64, _ error) {
	u := model.Group{}
	copier.Copy(&u, req)
	result, err := gs.SvcCtx.GroupModel.Insert(ctx, &u)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (gs *GroupService) GroupUpdate(ctx context.Context, req *types.GroupUpdateReq) (int64, error) {
	u := model.Group{}
	copier.Copy(&u, req)
	result, err := gs.SvcCtx.GroupModel.Update(ctx, &u)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (gs *GroupService) GroupDelete(ctx context.Context, id int64) (int64, error) {
	result, err := gs.SvcCtx.GroupModel.Delete(ctx, id)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (gs *GroupService) GroupList(ctx context.Context, req types.GroupListReq) (groupsResp []types.GroupResp, _ error) {
	limit, offset := req.Page.ToSqlLimit()
	v, err, _ := singleflight.SfGroup.Do(fmt.Sprintf("%+v", req), func() (any, error) {
		return gs.SvcCtx.GroupModel.PageList(ctx, limit, offset, req.GroupQuery)
	})
	if err != nil {
		return nil, err
	}

	groups := v.([]model.Group)
	if err := copier.Copy(&groupsResp, &groups); err != nil {
		return nil, err
	}
	return groupsResp, nil
}

func (gs *GroupService) GroupSearch(ctx context.Context, req types.GroupSearchReq) (groupsResp []types.GroupResp, _ error) {
	v, err, _ := singleflight.SfGroup.Do(fmt.Sprintf("%+v", req), func() (any, error) {
		return gs.SvcCtx.GroupModel.Search(ctx, req)
	})
	if err != nil {
		return nil, err
	}

	groups := v.([]model.Group)
	if err := copier.Copy(&groupsResp, &groups); err != nil {
		return nil, err
	}
	return groupsResp, nil
}
