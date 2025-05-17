// ==========================================================================
// LV自动生成路由代码,只生成一次,按需修改,再次生成不会覆盖.
// 生成日期:2025-05-16 05:37:58 &#43;0000 UTC
// 生成人:lv
// ==========================================================================
package iot_driver

import (
        "github.com/lostvip-com/lv_framework/web/router"
        "common/middleware/auth"
        "things/internal/iot_driver/api"
)

func init() {
	group_log := router.New( "/iot_driver/log", auth.TokenCheck())

	log := api.SysJobLogApi{}
	group_log.GET("/", "iot_driver:log:view", log.PreListSysJobLog)
	group_log.GET("/preAddSysJobLog", "iot_driver:log:add", log.PreAddSysJobLog)
	group_log.GET("/preEditSysJobLog", "iot_driver:log:edit", log.PreEditSysJobLog)
	// api
	group_log.POST("/listSysJobLog", "iot_driver:log:list", log.ListSysJobLog)
	group_log.POST("/addSysJobLog", "iot_driver:log:add", log.AddSysJobLog)
	group_log.POST("/removeSysJobLog", "iot_driver:log:remove", log.RemoveSysJobLog)
	group_log.POST("/editSysJobLog", "iot_driver:log:edit",log.SaveSysJobLog)
	group_log.POST("/exportSysJobLog", "iot_driver:log:export", log.ExportSysJobLog)
}
