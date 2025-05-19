// ==========================================================================
// LV自动生成路由代码,只生成一次,按需修改,再次生成不会覆盖.
// 生成日期:2025-05-17 13:18:47 &#43;0000 UTC
// 生成人:lv
// ==========================================================================
package iot_driver

import (
        "github.com/lostvip-com/lv_framework/web/router"
        "common/middleware/auth"
        "demo/internal/iot_driver/api"

)

func init() {
	log := router.New( "/iot_driver/log", auth.TokenCheck())

	logApi := api.SysJobLogApi{}
    log.GET("/:id", "iot_driver:log:info", logApi.GetRoleInfo)
    log.GET("/listSysJobLog", "iot_driver:log:list", logApi.ListSysJobLog)
	log.GET("/listSysJobLog", "iot_driver:log:list", logApi.ListSysJobLog)
	log.POST("", "iot_driver:log:new", logApi.CreateSysJobLog)
	log.PUT("", "iot_driver:log:edit",logApi.UpdateSysJobLog)
    log.DELETE("/ids", "iot_driver:log:del", logApi.DeleteSysJobLog)
	log.POST("/exportSysJobLog", "iot_driver:log:export", logApi.ExportSysJobLog)
}
