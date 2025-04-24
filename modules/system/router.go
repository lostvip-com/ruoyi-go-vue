package system

import (
	"common/middleware/auth"
	"github.com/lostvip-com/lv_framework/web/router"
	"system/api"
)

func init() {
	// 加载登录路由
	g0 := router.New("/")
	index := api.IndexApi{}
	login := api.LoginApi{}
	g0.GET("/", "", index.Index)
	g0.GET("/captchaImage", "", index.CaptchaImage)
	g0.POST("/login", "", login.Login)
	//下在要检测是否登录
	g1 := router.New("/", auth.TokenCheck())
	home := api.HomeApi{}
	g1.GET("/logout", "", login.Logout)
	g1.POST("/getInfo", "system:user:view", home.GetUserInfo)
}
