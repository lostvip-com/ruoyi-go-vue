package system

import (
	auth2 "common/middleware/auth"
	"github.com/lostvip-com/lv_framework/web/router"
	"system/api"
)

// 加载路由
func init() {
	serverController := api.ServiceApi{}
	g0 := router.New("/health")
	g0.GET("/", "", serverController.Health)
	// 服务监控
	g1 := router.New("/monitor/server", auth2.TokenCheck(), auth2.PermitCheck)
	g1.GET("/", "", serverController.Server)
	//登录日志
	g2 := router.New("/monitor/logininfor", auth2.TokenCheck(), auth2.PermitCheck)
	loginInfoApi := api.LogininfoApi{}
	g2.GET("/list", "monitor:logininfor:list", loginInfoApi.ListAjax)
	g2.POST("/export", "monitor:logininfor:export", loginInfoApi.Export)
	g2.POST("/clean", "monitor:logininfor:remove", loginInfoApi.Clean)
	g2.DELETE("/:infoIds", "monitor:logininfor:remove", loginInfoApi.Remove)
	g2.POST("/unlock", "monitor:logininfor:unlock", loginInfoApi.Unlock)

	//操作日志
	g3 := router.New("/monitor/operlog", auth2.TokenCheck(), auth2.PermitCheck)
	operController := api.OperateLogApi{}
	g3.POST("/list", "monitor:operlog:list", operController.ListAjax)
	g3.POST("/remove", "monitor:operlog:export", operController.Remove)
	g3.POST("/clean", "monitor:operlog:export", operController.Clean)

	//在线用户
	g4 := router.New("/monitor/online", auth2.TokenCheck(), auth2.PermitCheck)
	onlineController := api.OnlineApi{}
	g4.POST("/list", "monitor:online:list", onlineController.ListAjax)
	g4.POST("/forceLogout", "monitor:online:forceLogout", onlineController.ForceLogout)
	g4.POST("/batchForceLogout", "monitor:online:batchForceLogout", onlineController.BatchForceLogout)
}
