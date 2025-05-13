package system

import (
	auth2 "common/middleware/auth"
	"github.com/lostvip-com/lv_framework/web/router"
	"system/api"
)

// 加载路由
func init() {
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
	operApi := api.OperateLogApi{}
	g3.GET("/list", "monitor:operlog:list", operApi.ListAjax)
	g3.DELETE("/:operId", "monitor:operlog:Remove", operApi.DelectOperlog)
	g3.DELETE("/clean", "monitor:operlog:export", operApi.Clean)

	//在线用户
	g4 := router.New("/monitor/online", auth2.TokenCheck(), auth2.PermitCheck)
	onlineController := api.OnlineApi{}
	g4.POST("/list", "monitor:online:list", onlineController.ListAjax)
	g4.POST("/forceLogout", "monitor:online:forceLogout", onlineController.ForceLogout)
	g4.POST("/batchForceLogout", "monitor:online:batchForceLogout", onlineController.BatchForceLogout)

	//
	monitor := new(api.MonitorApi)
	monitorGroup := router.New("/monitor", auth2.TokenCheck(), auth2.PermitCheck)
	monitorGroup.GET("/cache", "", monitor.CacheHandler)
	monitorGroup.GET("/cache/getNames", "", monitor.CacheHandler)
	monitorGroup.GET("/cache/getKeys/:cacheName", "", monitor.GetCacheKeysHandler)
	monitorGroup.GET("/cache/getValue/:cacheName/:cacheKey", "", monitor.GetCacheValueHandler)
	monitorGroup.DELETE("/cache/clearCacheName/:cacheName", "", monitor.ClearCacheNameHandler)
	monitorGroup.DELETE("/cache/clearCacheKey/:cacheKey", "", monitor.ClearCacheKeyHandler)
	monitorGroup.DELETE("/cache/clearCacheAll", "", monitor.ClearCacheAllHandler)
	monitorGroup.GET("/server", "", monitor.ServerInfo)
}
