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

type userService interface {
	GetUserById(ctx context.Context, id int64) (*types.UserResp, error)
	UserCreate(ctx context.Context, req *types.UserCreateReq) (id int64, _ error)
	UserUpdate(ctx context.Context, req *types.UserUpdateReq) (int64, error)
	UserDelete(ctx context.Context, id int64) (rowsAffected int64, _ error)
	UserList(ctx context.Context, req types.UserListReq) (usersResp []types.UserResp, _ error)
	UserSearch(ctx context.Context, req types.UserSearchReq) (userResp []types.UserResp, _ error)
}

type UserService struct {
	SvcCtx svcCtx.UserServiceContext
}

var _ userService = (*UserService)(nil)

func NewUserService() userService {
	return &UserService{
		SvcCtx: svcCtx.NewUserServiceContext(),
	}
}

func (us *UserService) GetUserById(ctx context.Context, id int64) (*types.UserResp, error) {
	u, err := us.SvcCtx.UserModel.Select(ctx, id)
	if err != nil {
		return nil, err
	}

	userResp := types.UserResp{}

	err = copier.Copy(&userResp, u)
	if err != nil {
		return nil, err
	}

	return &userResp, err
}

func (us *UserService) UserCreate(ctx context.Context, req *types.UserCreateReq) (id int64, _ error) {
	u := model.User{}
	copier.Copy(&u, req)
	result, err := us.SvcCtx.UserModel.Insert(ctx, &u)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (us *UserService) UserUpdate(ctx context.Context, req *types.UserUpdateReq) (int64, error) {
	u := model.User{}
	copier.Copy(&u, req)
	result, err := us.SvcCtx.UserModel.Update(ctx, &u)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (us *UserService) UserDelete(ctx context.Context, id int64) (int64, error) {
	result, err := us.SvcCtx.UserModel.Delete(ctx, id)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (us *UserService) UserList(ctx context.Context, req types.UserListReq) (usersResp []types.UserResp, _ error) {
	limit, offset := req.Page.ToSqlLimit()
	v, err, _ := singleflight.SfGroup.Do(fmt.Sprintf("%+v", req), func() (any, error) {
		return us.SvcCtx.UserModel.PageList(ctx, limit, offset, req.UserQuery)
	})
	if err != nil {
		return nil, err
	}

	users := v.([]model.User)
	if err := copier.Copy(&usersResp, &users); err != nil {
		return nil, err
	}
	return usersResp, nil
}

func (us *UserService) UserSearch(ctx context.Context, req types.UserSearchReq) (usersResp []types.UserResp, _ error) {
	v, err, _ := singleflight.SfGroup.Do(fmt.Sprintf("%+v", req), func() (any, error) {
		return us.SvcCtx.UserModel.Search(ctx, req)
	})
	if err != nil {
		return nil, err
	}

	users := v.([]model.User)
	if err := copier.Copy(&usersResp, &users); err != nil {
		return nil, err
	}
	return usersResp, nil
}
