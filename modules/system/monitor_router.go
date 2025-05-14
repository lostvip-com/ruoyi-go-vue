package system

import (
	auth2 "common/middleware/auth"
	"github.com/lostvip-com/lv_framework/web/router"
	"system/api"
)

// 加载路由
func init() {
	//登录日志
	g2 := router.New("/monitorApi/logininfor", auth2.TokenCheck(), auth2.PermitCheck)
	loginInfoApi := api.LogininfoApi{}
	g2.GET("/list", "monitorApi:logininfor:list", loginInfoApi.ListAjax)
	g2.POST("/export", "monitorApi:logininfor:export", loginInfoApi.Export)
	g2.POST("/clean", "monitorApi:logininfor:remove", loginInfoApi.Clean)
	g2.DELETE("/:infoIds", "monitorApi:logininfor:remove", loginInfoApi.Remove)
	g2.POST("/unlock", "monitorApi:logininfor:unlock", loginInfoApi.Unlock)

	//操作日志
	g3 := router.New("/monitorApi/operlog", auth2.TokenCheck(), auth2.PermitCheck)
	operApi := api.OperateLogApi{}
	g3.GET("/list", "monitorApi:operlog:list", operApi.ListAjax)
	g3.DELETE("/:operId", "monitorApi:operlog:Remove", operApi.DelectOperlog)
	g3.DELETE("/clean", "monitorApi:operlog:export", operApi.Clean)
	// 监控
	monitorApi := new(api.MonitorApi)
	monitor := router.New("/monitorApi", auth2.TokenCheck(), auth2.PermitCheck)
	monitor.GET("/cache", "", monitorApi.CacheHandler)
	monitor.GET("/cache/getNames", "", monitorApi.CacheHandler)
	monitor.GET("/cache/getKeys/:cacheName", "", monitorApi.GetCacheKeysHandler)
	monitor.GET("/cache/getValue/:cacheName/:cacheKey", "", monitorApi.GetCacheValueHandler)
	monitor.DELETE("/cache/clearCacheName/:cacheName", "", monitorApi.ClearCacheNameHandler)
	monitor.DELETE("/cache/clearCacheKey/:cacheKey", "", monitorApi.ClearCacheKeyHandler)
	monitor.DELETE("/cache/clearCacheAll", "", monitorApi.ClearCacheAllHandler)
	monitor.GET("/server", "", monitorApi.ServerInfo)
	//online
	monitor.GET("/online/list", "", monitorApi.ListOnLine)
	monitor.DELETE("online/:tokenId", "", monitorApi.DetectOnLine)
	//job
	monitor.POST("/job", "", monitorApi.SaveJob)
	monitor.PUT("/job", "", monitorApi.UploadJob)
	monitor.GET("/job/:jobId", "", monitorApi.GetJobById)
	monitor.GET("/job/list", "", monitorApi.ListJob)
	monitor.POST("/job/export", "", monitorApi.ExportJob)
	monitor.PUT("/job/changeStatus", "", monitorApi.ChangeStatus)
	monitor.DELETE("/job/:jobIds", "", monitorApi.DelectJob)
	monitor.PUT("/job/run/:jobId", "", monitorApi.RunJob)
	//job log
	monitor.GET("/jobLog/list", "", monitorApi.ListJobLog)
	monitor.POST("/jobLog/export", "", monitorApi.ExportJobLog)
	monitor.GET("/jobLog/:logId", "", monitorApi.GetJobLog)
	monitor.DELETE("/jobLog/:jobLogIds", "", monitorApi.DetectJobLog)
	monitor.DELETE("/jobLog/clean", "", monitorApi.ClearJobLog)

}
