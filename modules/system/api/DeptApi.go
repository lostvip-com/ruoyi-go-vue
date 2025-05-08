package api

import (
	"common/common_vo"
	"common/session"
	util2 "common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"net/http"
	"system/service"
)

type DeptApi struct {
}

// ListAjax 列表分页数据
func (w *DeptApi) ListAjax(c *gin.Context) {
	var service service.DeptService
	var req = common_vo.DeptPageReq{}

	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetMsg(err.Error()).Log("部门管理", req).WriteJsonExit()
		return
	}
	result, err := service.SelectListAll(&req)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, result)
}

// AddSave 新增页面保存
func (w *DeptApi) AddSave(c *gin.Context) {
	var req *common_vo.AddDeptReq
	var service service.DeptService

	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	rid, err := service.AddSave(req, c)
	if err != nil || rid <= 0 {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).WriteJsonExit()
		return
	}
	util2.SucessResp(c).SetBtype(lv_dto.Buniss_Add).WriteJsonExit()
}

// EditSave 修改页面保存
func (w *DeptApi) EditSave(c *gin.Context) {
	var service service.DeptService

	var req *common_vo.EditDeptReq

	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).Log("部门管理", req).WriteJsonExit()
		return
	}

	rs, err := service.EditSave(req, c)

	if err != nil || rs <= 0 {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).Log("部门管理", req).WriteJsonExit()
		return
	}
	util2.SucessResp(c).SetData(rs).SetBtype(lv_dto.Buniss_Edit).Log("部门管理", req).WriteJsonExit()
}

// 删除数据
func (w *DeptApi) Remove(c *gin.Context) {
	id := lv_conv.Int64(c.Query("id"))
	service := service.DeptService{}
	err := service.DeleteDeptById(id)
	if err != nil {
		util2.Fail(c, err.Error())
	} else {
		util2.Success(c, id)
	}
}

// 加载部门列表树结构的数据
func (w *DeptApi) TreeData(c *gin.Context) {
	var service service.DeptService
	tenantId := session.GetTenantId(c)
	result, _ := service.SelectDeptTree(0, "", "", tenantId)
	c.JSON(http.StatusOK, result)
}

// 加载角色部门（数据权限）列表树
func (w *DeptApi) RoleDeptTreeData(c *gin.Context) {
	var service service.DeptService
	tenantId := session.GetTenantId(c)
	roleId := lv_conv.Int64(c.Query("roleId"))
	result, err := service.RoleDeptTreeData(roleId, tenantId)

	if err != nil {
		util2.ErrorResp(c).SetMsg(err.Error()).Log("菜单树", gin.H{"roleId": roleId})
		return
	}

	c.JSON(http.StatusOK, result)
}
