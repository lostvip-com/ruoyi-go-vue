package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/spf13/cast"
	"net/http"
	"system/dao"
	"system/model"
	"system/service"
	"system/vo"
)

type MenuApi struct {
	BaseApi
}

func (w *MenuApi) GetMenuInfo(c *gin.Context) {
	var menuId = c.Param("menuId")
	menu := new(model.SysMenu)
	menu, err := menu.FindById(cast.ToInt64(menuId))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, menu)
}

// ListMenu 列表分页数据
func (w *MenuApi) ListMenu(c *gin.Context) {
	var req = new(vo.SelectMenuPageReq)
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	userId := w.GetCurrUserId(c)
	rows := make([]model.SysMenu, 0)
	var err error
	if w.IsAdmin(userId) {
		rows, err = dao.GetMenuDaoInstance().SelectListAll(req)
	} else {
		rows, err = dao.GetMenuDaoInstance().FindMenusByUserId(userId, req)
	}
	if err != nil {
		util.Fail(c, err.Error())
	}
	util.Success(c, rows)
}

// AddSave 新增页面保存
func (w *MenuApi) AddSave(c *gin.Context) {
	var req = new(model.SysMenu)

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
	util.Success(c, req)
}

// Remove 删除数据
func (w *MenuApi) Remove(c *gin.Context) {
	id := lv_conv.Int64(c.Query("menuId"))
	err := service.GetMenuServiceInstance().DeleteRecordById(id)
	if err == nil {
		util.Success(c, gin.H{"id": id})
	} else {
		util.Fail(c, err.Error())
	}
}

func (w *MenuApi) GetTreeSelect(c *gin.Context) {
	userId := w.GetCurrUserId(c)
	var menuParm model.SysMenu
	if err := c.ShouldBind(&menuParm); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	svc := service.GetMenuServiceInstance()
	menus, err := svc.SelectMenuTree(userId, &menuParm)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	var arrTree = svc.BuildMenuTreeSelect(menus)
	c.JSON(http.StatusOK, gin.H{
		"msg":   "操作成功",
		"code":  http.StatusOK,
		"menus": arrTree,
	})
}

func (w *MenuApi) TreeSelectByRole(c *gin.Context) {
	userId := w.GetCurrUserId(c)
	var roleId = c.Param("roleId")
	var menuParm model.SysMenu
	if err := c.ShouldBind(&menuParm); err != nil {
		c.JSON(http.StatusOK, err.Error())
		return
	}
	svc := service.GetMenuServiceInstance()
	menus, err := svc.SelectMenuTree(userId, &menuParm)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	role := new(model.SysRole)
	role, err = role.FindById(cast.ToInt64(roleId))
	checkedKeys, _ := svc.SelectMenuListByRoleId(roleId, role.MenuCheckStrictly)
	var arrTree = svc.BuildMenuTreeSelect(menus)
	c.JSON(http.StatusOK, gin.H{
		"msg":         "操作成功",
		"code":        http.StatusOK,
		"menus":       arrTree,
		"checkedKeys": checkedKeys,
	})
}
