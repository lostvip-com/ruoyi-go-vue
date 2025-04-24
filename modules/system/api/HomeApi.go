package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
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

// 后台框架首页
//func (w *HomeApi) Index(c *gin.Context) {
//	var userService service.UserService
//	user := userService.GetProfile(c)
//	loginname := user.LoginName
//	username := user.UserName
//	avatar := user.Avatar
//	if avatar == "" {
//		avatar = "/resource/img/profile.jpg"
//	}
//	var menus *[]model.SysMenu
//	//获取菜单数据
//	menuService := service.MenuService{}
//	if userService.IsAdmin(user.UserId) {
//		tmp, err := menuService.SelectMenuNormalAll("")
//		if err == nil {
//			menus = tmp
//		}
//	} else {
//		tmp, err := menuService.SelectMenusByUserId(user.UserId, "")
//		if err == nil {
//			menus = tmp
//		}
//	}
//
//	//获取配置数据
//	var configService service.ConfigService
//	sideTheme := configService.GetValueFromRam("sys.index.sideTheme")
//	skinName := configService.GetValueFromRam("sys.index.skinName")
//	//设置首页风格
//	menuStyle := c.Query("menuStyle")
//	cookie, _ := c.Request.Cookie("menuStyle")
//	if cookie == nil {
//		cookie = &http.Cookie{
//			Name:     "menuStyle",
//			Value:    menuStyle,
//			HttpOnly: true,
//		}
//		http.SetCookie(c.Writer, cookie)
//	}
//	if menuStyle == "" { //未指定则从cookie中取
//		menuStyle = cookie.Value
//	}
//	var targetIndex string         //默认首页
//	if menuStyle == "index_left" { //指定了左侧风格,
//		targetIndex = "index_left"
//	} else { //否则默认风格
//		targetIndex = "index"
//	}
//	//"menuStyle", cookie.Value, 1000, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly
//	c.SetCookie(cookie.Name, menuStyle, cookie.MaxAge, cookie.Path, cookie.Domain, cookie.Secure, cookie.HttpOnly)
//	util.BuildTpl(c, targetIndex).WriteTpl(gin.H{
//		"avatar":    avatar,
//		"loginname": loginname,
//		"username":  username,
//		"menus":     menus,
//		"sideTheme": sideTheme,
//		"skinName":  skinName,
//	})
//	c.Abort()
//}
