package service

import (
	"context"
	svcCtx "goweb/app/account/context"
	"goweb/app/account/types"
)

type user2groupService interface {
	GroupsListByUser(ctx context.Context, userId int64) (*types.GroupsResp, error)
	GroupsAddToUser(ctx context.Context, userId int64, groupIds []int64) (int64, error)
	GroupsDeleteFromUser(ctx context.Context, userId int64, groupIds []int64) (int64, error)
	UsersListByGroup(ctx context.Context, groupId int64) (*types.UsersResp, error)
	UsersAddToGroup(ctx context.Context, groupId int64, userIds []int64) (int64, error)
	UsersDeleteFromGroup(ctx context.Context, groupId int64, userIds []int64) (int64, error)
}

type User2GroupService struct {
	SvcCtx svcCtx.User2GroupServiceContext
}

var _ user2groupService = (*User2GroupService)(nil)

func NewUser2GroupService() user2groupService {
	return &User2GroupService{
		SvcCtx: svcCtx.NewUser2GroupServiceContext(),
	}
}

func (ugs *User2GroupService) GroupsListByUser(ctx context.Context, userId int64) (*types.GroupsResp, error) {
	groupIds, err := ugs.SvcCtx.User2GroupModel.GroupsListByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	resp := &types.GroupsResp{
		Ids: groupIds,
	}

	return resp, nil
}

func (ugs *User2GroupService) GroupsAddToUser(ctx context.Context, userId int64, groupIds []int64) (int64, error) {
	return ugs.SvcCtx.User2GroupModel.GroupsAddToUser(ctx, userId, groupIds)
}

func (ugs *User2GroupService) GroupsDeleteFromUser(ctx context.Context, userId int64, groupIds []int64) (int64, error) {
	result, err := ugs.SvcCtx.User2GroupModel.GroupsDeleteFromUser(ctx, userId, groupIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (ugs *User2GroupService) UsersListByGroup(ctx context.Context, groupId int64) (*types.UsersResp, error) {
	userIds, err := ugs.SvcCtx.User2GroupModel.UserListByGroup(ctx, groupId)
	if err != nil {
		return nil, err
	}

	resp := &types.UsersResp{
		Ids: userIds,
	}

	return resp, nil
}

func (ugs *User2GroupService) UsersAddToGroup(ctx context.Context, groupId int64, userIds []int64) (int64, error) {
	return ugs.SvcCtx.User2GroupModel.UsersAddToGroup(ctx, groupId, userIds)
}
func (ugs *User2GroupService) UsersDeleteFromGroup(ctx context.Context, groupId int64, userIds []int64) (int64, error) {
	result, err := ugs.SvcCtx.User2GroupModel.UsersDeleteFromGroup(ctx, groupId, userIds)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
