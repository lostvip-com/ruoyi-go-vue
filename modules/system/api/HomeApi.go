package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"system/model"
	"system/service"
)

type HomeApi struct{}

func (w *HomeApi) Home(c *gin.Context) {
	util.BuildTpl(c, "home").WriteTpl()
}

func (w *HomeApi) GetUserInfo(c *gin.Context) {
	var req *lv_dto.IdsReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Del).SetMsg(err.Error()).WriteJsonExit()
	}
	user := service.GetUserService().GetProfile(c)

	data := make(map[string]any)
	data["roles"] = ""
	data["permissions"] = ""
	data["user"] = user
	data["dept"] = ""
	util.Success(c, data, "success")
}

// GetRouters 后台框架菜单
func (w *HomeApi) GetRouters(c *gin.Context) {
	userService := service.GetUserService()
	user := userService.GetProfile(c)
	var menus *[]model.SysMenu
	menuService := service.MenuService{}
	if userService.IsAdmin(user.UserId) {
		menus, _ = menuService.SelectMenuNormalAll("")
	} else {
		menus, _ = menuService.SelectMenusByUserId(user.UserId, "")
	}
	//获取配置数
	util.Success(c, menus, "success")
}
