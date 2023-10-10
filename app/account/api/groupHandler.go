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

// @Tags Group
// @Summary GroupInfo
// @Param groupId path int true "groupId"
// @Success 200 {object} types.GroupResp
// @Router /group/{groupId} [get]
func GroupInfo(c *gin.Context) {
	groupId, err := strconv.ParseInt(c.Param("groupId"), 10, 64)
	if err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	groupSvc := svc.NewGroupService()
	groupResp, err := groupSvc.GetGroupById(c.Request.Context(), groupId)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}

	result.SuccessWithData(c, groupResp)
	return
}

// @Tags Group
// @Summary GroupInfo
// @Param request body types.GroupCreateReq true "GroupCreateReq"
// @Success 200 {object}  types.GroupCreateResp
// @Router /group/create [post]
func GroupCreate(c *gin.Context) {
	req := types.GroupCreateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}

	groupSvc := svc.NewGroupService()
	id, err := groupSvc.GroupCreate(c.Request.Context(), &req)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}

	result.SuccessWithData(c,
		types.GroupCreateResp{
			LastId: id,
		})
}

// @Tags Group
// @Summary GroupInfo
// @Param groupId path int true "groupId"
// @Param request body types.GroupUpdateReq true "GroupUpdateReq"
// @Success 200 {object} types.ModifyRowsResp
// @Router /group/{groupId}/update [post]
func GroupUpdate(c *gin.Context) {
	groupId, err := strconv.ParseInt(c.Param("groupId"), 10, 64)
	if err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
	}
	req := types.GroupUpdateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}
	req.Id = groupId

	groupSvc := svc.NewGroupService()
	rowsAffected, err := groupSvc.GroupUpdate(c.Request.Context(), &req)
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

// @Tags Group
// @Summary GroupInfo
// @Param groupId path int true "groupId"
// @Success 200 {object} types.ModifyRowsResp
// @Router /group/{groupId}/delete [get]
func GroupDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("groupId"), 10, 64)
	if err != nil {
		result.FailWithErr(c, http.StatusBadRequest, err)
		return
	}
	groupSvc := svc.NewGroupService()
	rowsAffected, err := groupSvc.GroupDelete(c.Request.Context(), id)
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

// @Tags Group
// @Summary GroupInfo
// @Param request query types.GroupListReq true "GroupListReq"
// @Success 200 {array} types.GroupResp
// @Router /group/list [get]
/*
pagination
*/
func GroupList(c *gin.Context) {
	var req types.GroupListReq
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

	groupSvc := svc.NewGroupService()
	resp, err := groupSvc.GroupList(c, req)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}
	result.SuccessWithData(c, resp)
	return
}

// @Tags Group
// @Summary GroupInfo
// @Param request body types.GroupSearchReq true "GroupSearchReq"
// @Success 200
// @Router /group/search [post]
/*
complex search
multiple conditions
*/
func GroupSearch(c *gin.Context) {
	var req types.GroupSearchReq
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

	groupSvc := svc.NewGroupService()
	resp, err := groupSvc.GroupSearch(c, req)
	if err != nil {
		result.FailWithErr(c, http.StatusBadGateway, err)
		return
	}
	result.SuccessWithData(c, resp)
	return
}
