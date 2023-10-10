package api

import (
	"goweb/app/account/service"
	"goweb/app/account/types"
	"goweb/common/result"
	"goweb/common/validx"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Tags user2group
// @Summary GroupsListByUser
// @Param userId query int true "userId"
// @Success 200 {object} types.GroupsResp
// @Router /user/groupsList [get]
func GroupsListByUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("userId"), 10, 64)
	if err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	ugs := service.NewUser2GroupService()
	resp, err := ugs.GroupsListByUser(c.Request.Context(), userId)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}
	result.SuccessWithData(c, resp)
	return
}

// @Tags user2group
// @Summary GroupsAddToUser
// @Param request body types.GroupsInUserReq true "GroupsInUserReq"
// @Success 200 {object} types.ModifyRowsResp
// @Router /user/groupsAdd [post]
func GroupsAddToUser(c *gin.Context) {
	var req types.GroupsInUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}
	if err := validx.GetValidator().Struct(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	ugs := service.NewUser2GroupService()
	rowsAffected, err := ugs.GroupsAddToUser(c, req.UserId, req.GroupIds)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}

	result.SuccessWithData(c, types.ModifyRowsResp{
		RowsAffected: rowsAffected,
	})
	return
}

// @Tags user2group
// @Summary GroupsDeleteFromUser
// @Param request body types.GroupsInUserReq true "GroupsInUserReq"
// @Success 200 {object} types.ModifyRowsResp
// @Router /user/groupsDelete [post]
func GroupsDeleteFromUser(c *gin.Context) {
	var req types.GroupsInUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}
	if err := validx.GetValidator().Struct(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	ugs := service.NewUser2GroupService()
	rowsAffected, err := ugs.GroupsDeleteFromUser(c, req.UserId, req.GroupIds)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}
	result.SuccessWithData(c, types.ModifyRowsResp{
		RowsAffected: rowsAffected,
	})
	return
}

// @Tags user2group
// @Summary UsersListByGroup
// @Param groupId query int true "groupId"
// @Success 200 {object} types.UsersResp
// @Router /group/usersList [get]
func UsersListByGroup(c *gin.Context) {
	groupId, err := strconv.ParseInt(c.Query("groupId"), 10, 64)
	if err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	ugs := service.NewUser2GroupService()
	resp, err := ugs.UsersListByGroup(c.Request.Context(), groupId)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}
	result.SuccessWithData(c, resp)
	return
}

// @Tags user2group
// @Summary UsersAddToGroup
// @Param request body types.UsersInGroupReq true "UsersInGroupReq"
// @Success 200 {object} types.ModifyRowsResp
// @Router /group/usersAdd [post]
func UsersAddToGroup(c *gin.Context) {
	var req types.UsersInGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}
	if err := validx.GetValidator().Struct(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	ugs := service.NewUser2GroupService()
	rowsAffected, err := ugs.UsersAddToGroup(c, req.GroupId, req.UserIds)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}

	result.SuccessWithData(c, types.ModifyRowsResp{
		RowsAffected: rowsAffected,
	})
	return
}

// @Tags user2group
// @Summary UsersDeleteFromGroup
// @Param request body types.UsersInGroupReq true "UsersInGroupReq"
// @Success 200 {object} types.ModifyRowsResp
// @Router /group/usersDelete [post]
func UsersDeleteFromGroup(c *gin.Context) {
	var req types.UsersInGroupReq
	if err := c.ShouldBindJSON(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}
	if err := validx.GetValidator().Struct(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	ugs := service.NewUser2GroupService()
	rowsAffected, err := ugs.UsersDeleteFromGroup(c, req.GroupId, req.UserIds)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}

	result.SuccessWithData(c, types.ModifyRowsResp{
		RowsAffected: rowsAffected,
	})
	return
}
