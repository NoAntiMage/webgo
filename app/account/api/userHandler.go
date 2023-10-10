package api

import (
	svc "goweb/app/account/service"
	"goweb/app/account/types"
	"goweb/common/jsonx"
	"goweb/common/result"
	"goweb/common/validx"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Tags User
// @Summary Fetch UserInfo profile
// @Param userId path int true "userId"
// @Success 200 {object} types.UserResp
// @Router /user/{userId} [get]
func UserInfo(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	userSvc := svc.NewUserService()
	userResp, err := userSvc.GetUserById(c.Request.Context(), userId)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}

	result.SuccessWithData(c, userResp)
	return
}

// @Tags User
// @Summary Create User
// @Accept json
// @Param request body types.UserCreateReq true "UserCreate Param"
// @Success 200 {object} types.UserCreateResp
// @Router /user/create [post]
func UserCreate(c *gin.Context) {
	req := types.UserCreateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	userSvc := svc.NewUserService()
	id, err := userSvc.UserCreate(c.Request.Context(), &req)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}

	result.SuccessWithData(c,
		types.UserCreateResp{
			LastId: id,
		})
	return
}

// @Tags User
// @Summary userInfo update
// @Param userId path int true "userId"
// @Param request body types.UserUpdateReq true "UserUpdate Param"
// @Success 200 {object} types.ModifyRowsResp
// @Router /user/{userId}/update [post]
func UserUpdate(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
	}
	req := types.UserUpdateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}
	req.Id = userId

	userSvc := svc.NewUserService()
	rowsAffected, err := userSvc.UserUpdate(c.Request.Context(), &req)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}

	result.SuccessWithData(c,
		types.ModifyRowsResp{
			RowsAffected: rowsAffected,
		})
	return
}

// @Tags User
// @Summary user delete
// @Param userId path int true "userId"
// @Success 200 {object} types.ModifyRowsResp
// @Router /user/{userId}/delete [get]
func UserDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}
	userSvc := svc.NewUserService()
	rowsAffected, err := userSvc.UserDelete(c.Request.Context(), id)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}

	result.SuccessWithData(c,
		types.ModifyRowsResp{
			RowsAffected: rowsAffected,
		})
	return
}

// @Tags User
// @Summary user list
// @Param request query types.UserListReq true "UserList Param"
// @Success 200 {array} types.UserResp
// @Router /user/list [get]
/*
pagination
*/
func UserList(c *gin.Context) {
	var req types.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}
	if err := jsonx.Struct2StructWithRule(&req, jsonx.RuleDefault); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}
	valid := validx.GetValidator()
	if err := valid.Struct(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	userSvc := svc.NewUserService()
	resp, err := userSvc.UserList(c, req)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}
	result.SuccessWithData(c, resp)
	return
}

// @Tags User
// @Summary user search
// @Param request body types.UserSearchReq true "UserSearch Param"
// @Success 200 {object} types.UserResp
// @Router /user/search [post]
/*
complex search
multiple conditions
*/
func UserSearch(c *gin.Context) {
	var req types.UserSearchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	valid := validx.GetValidator()
	for _, cond := range req.Query {
		if err := valid.Struct(&cond); err != nil {
			result.FailWithErr(c, http.StatusBadRequest, err)
			return
		}
	}

	userSvc := svc.NewUserService()
	resp, err := userSvc.UserSearch(c, req)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}
	result.SuccessWithData(c, resp)
	return
}
