package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"net/http"
	"system/dao"
	"system/model"
	"system/service"
	"system/vo"
)

type MenuApi struct {
}

// ListAjax 列表分页数据
func (w *MenuApi) ListAjax(c *gin.Context) {
	var req = new(vo.SelectMenuPageReq)
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("菜单管理", req).WriteJsonExit()
		return
	}
	rows := make([]model.SysMenu, 0)
	result, err := dao.GetMenuDaoInstance().SelectListAll(req)

	if err == nil && len(result) > 0 {
		rows = result
	}
	c.JSON(http.StatusOK, rows)
}

// AddSave 新增页面保存
func (w *MenuApi) AddSave(c *gin.Context) {
	var req = new(model.SysMenu)
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg(err.Error()).Log("菜单管理", req).WriteJsonExit()
		return
	}
	user := service.GetUserService().GetProfile(c)
	if user != nil {
		req.CreateBy = user.UserName
	}
	id, err := service.GetMenuServiceInstance().AddSave(req)

	if err != nil || id <= 0 {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg(err.Error()).Log("菜单管理", req).WriteJsonExit()
		return
	}
	util.SucessResp(c).SetBtype(lv_dto.Buniss_Add).SetData(id).Log("菜单管理", req).WriteJsonExit()
}

// EditSave 修改页面保存
func (w *MenuApi) EditSave(c *gin.Context) {
	var req = new(model.SysMenu)
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).Log("菜单管理", req).WriteJsonExit()
		return
	}
	user := service.GetUserService().GetProfile(c)
	if user != nil {
		req.UpdateBy = user.UserName
	}
	rs, err := service.GetMenuServiceInstance().Edit(req)

	if err != nil || rs <= 0 {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).Log("菜单管理", req).WriteJsonExit()
		return
	}
	util.SuccessData(c, req)
}

// Remove 删除数据
func (w *MenuApi) Remove(c *gin.Context) {
	id := lv_conv.Int64(c.Query("id"))
	err := service.GetMenuServiceInstance().DeleteRecordById(id)
	if err == nil {
		util.SuccessData(c, gin.H{"id": id})
	} else {
		util.Fail(c, err.Error())
	}
}

// MenuTreeData 加载所有菜单列表树
func (w *MenuApi) MenuTreeData(c *gin.Context) {
	user := service.GetUserService().GetProfile(c)
	if user == nil {
		util.ErrorResp(c).SetMsg("登录超时").Log("菜单管理", gin.H{"userId": user.UserId}).WriteJsonExit()
		return
	}
	ztrees, err := service.GetMenuServiceInstance().MenuTreeData(user.UserId)
	if err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("菜单管理", gin.H{"userId": user.UserId}).WriteJsonExit()
		return
	}
	c.JSON(http.StatusOK, ztrees)
}

// RoleMenuTreeData 加载角色菜单列表树
func (w *MenuApi) RoleMenuTreeData(c *gin.Context) {
	var userService service.UserService
	roleId := lv_conv.Int64(c.Query("roleId"))
	user := userService.GetProfile(c)
	if user == nil || user.UserId <= 0 {
		util.ErrorResp(c).SetMsg("登录超时").Log("菜单管理", gin.H{"roleId": roleId}).WriteJsonExit()
		return
	}
	var service service.MenuService
	result, err := service.RoleMenuTreeData(roleId, user.UserId, "")

	if err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("菜单管理", gin.H{"roleId": roleId}).WriteJsonExit()
		return
	}

	c.JSON(http.StatusOK, result)
}
