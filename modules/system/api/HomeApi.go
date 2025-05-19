package api

import (
	api2 "common/api"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"net/http"
	"system/model"
	"system/service"
)

type HomeApi struct {
	api2.BaseApi
}

func (w *HomeApi) GetUserInfo(c *gin.Context) {
	u, ok := c.Get("user")
	if !ok {
		util.Fail(c, "fail!")
	}
	user := u.(*model.SysUser)
	roles := user.GetRoleKeys()
	isAdmin := service.GetUserServiceInstance().IsAdmin(user.UserId)
	var permissions []string
	if isAdmin {
		permissions = []string{"*:*:*"}
	} else {
		permissions = service.GetPermissionServiceInstance().FindPerms(roles)
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "user": user, "roles": roles, "permissions": permissions})
}

// GetRouters 后台框架菜单
func (w *HomeApi) GetRouters(c *gin.Context) {
	userId := w.GetCurrUserId(c)
	var menus []model.SysMenu
	var err error
	menuService := service.GetMenuServiceInstance()
	if w.IsAdmin(userId) {
		menus, err = menuService.FindRouterTreeAll()
	} else {
		menus, err = menuService.FindRouterTreeAllByUserId(userId)
	}
	lv_err.HasErrAndPanic(err)
	var data = menuService.BuildMenus(menus) //获取配置数
	util.Success(c, data)
}
