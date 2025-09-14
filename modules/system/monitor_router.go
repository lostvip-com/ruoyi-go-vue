package system

import (
	auth2 "common/middleware/auth"
	"github.com/lostvip-com/lv_framework/web/router"
	"system/api"
)

// 加载路由
func init() {
	monitorApi := new(api.MonitorApi)
	monitor := router.New("/monitor", auth2.TokenCheck(), auth2.PermitCheck)
	//登录日志
	loginInfoApi := api.LogininfoApi{}
	monitor.GET("/logininfor/list", "monitor:logininfor:list", loginInfoApi.ListAjax)
	monitor.POST("/logininfor/export", "monitor:logininfor:export", loginInfoApi.Export)
	monitor.DELETE("/logininfor/clean", "monitor:logininfor:remove", loginInfoApi.Clean)
	monitor.DELETE("/logininfor/:infoIds", "monitor:logininfor:remove", loginInfoApi.Remove)
	monitor.POST("/logininfor/unlock", "monitor:logininfor:unlock", loginInfoApi.Unlock)

	//操作日志
	operApi := api.OperateLogApi{}
	monitor.GET("/operlog/list", "monitor:operlog:list", operApi.ListAjax)
	monitor.DELETE("/operlog/:operIds", "monitor:operlog:remove", operApi.DelLogs)
	monitor.DELETE("/operlog/clean", "monitor:operlog:remove", operApi.Clean)
	monitor.POST("/operlog/export", "monitor:operlog:export", operApi.Export)
	// 监控

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
	monitor.DELETE("/job/:jobIds", "", monitorApi.DelJobs)
	monitor.PUT("/job/run", "", monitorApi.RunJob)
	//job log
	monitor.GET("/jobLog/list", "", monitorApi.ListJobLog)
	monitor.POST("/jobLog/export", "", monitorApi.ExportJobLog)
	monitor.GET("/jobLog/:logId", "", monitorApi.GetJobLog)
	monitor.DELETE("/jobLog/:jobLogIds", "", monitorApi.DelLogs)
	monitor.DELETE("/jobLog/clean", "", monitorApi.ClearJobLog)

}
