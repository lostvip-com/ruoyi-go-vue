package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"system/model"
	"system/service"
	"system/vo"
)

type HomeApi struct{}

func (w *HomeApi) GetUserInfo(c *gin.Context) {
	u, ok := c.Get("user")
	if !ok {
		util.Fail(c, "fail!")
	}
	user := u.(*model.SysUser)
	roles := user.GetRoleKeys()
	isAdmin := service.GetUserService().IsAdmin(user.UserId)
	var permissions []string
	if isAdmin {
		permissions = []string{"*:*:*"}
	} else {
		permissions = service.GetPermissionServiceInstance().FindPerms(roles)
	}
	//获取参数
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "user": user, "roles": roles, "permissions": permissions})
}

// GetRouters 后台框架菜单
func (w *HomeApi) GetRouters(c *gin.Context) {
	userService := service.GetUserService()
	user := userService.GetProfile(c)
	var menus []vo.RouterVO
	menuService := service.MenuService{}
	if userService.IsAdmin(user.UserId) {
		menus, _ = menuService.SelectMenuNormalAll(0, "")
	} else {
		menus, _ = menuService.SelectMenuNormalAll(user.UserId, "")
	}
	//获取配置数
	util.Success(c, menus, "success")
}
