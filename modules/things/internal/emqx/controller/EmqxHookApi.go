//	=========================================================================
//
// LV自动生成控制器相关代码，只生成一次，按需修改,默认再次生成不会覆盖.
// date：2024-08-29 09:08:30 +0000 UTC
// author：lv
// ==========================================================================
package controller

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	service2 "things/internal/emqx/service"
	"things/internal/emqx/vo_emqx"
	"things/internal/iot_product/service"
)

type EmqxHookApi struct{}

// ========================================================================
//
//	emqx 回调使用，包括 授权、上线、下线
//
// =========================================================================

//	{
//		 "clientId": "${clientid}",
//		 "password": "${password}",
//		 "username": "${username}"
//	}
//
// AuthMqtt mqtt 登录鉴权
func (w EmqxHookApi) AuthV5(c *gin.Context) {
	req := new(vo_emqx.EmqxLoginVO)
	err := c.ShouldBindJSON(req)
	lv_err.HasErrAndPanic(err)
	lv_log.Infof("<==========%v", req)
	err = service2.GetClientService().LoginClient(req)
	lv_err.HasErrAndPanic(err)
	util.SuccessData(c, req)
}

// AddSave 新增页面保存
func (w EmqxHookApi) Online(c *gin.Context) {
	req := new(vo_emqx.EmqxOnlineVO)
	err := c.ShouldBindJSON(req)
	lv_err.HasErrAndPanic(err)
	lv_log.Infof("<==========%v", req)
	deviceId := cast.ToInt(req.Clientid) //client 全部使用数字，且是device表的主键
	_ = service.GetDeviceCacheService().Online(deviceId)
	lv_err.HasErrAndPanic(err)
	util.SuccessData(c, req)
}

func (w EmqxHookApi) Offline(c *gin.Context) {
	req := new(vo_emqx.EmqxOnlineVO)
	err := c.ShouldBindJSON(req)
	lv_err.HasErrAndPanic(err)
	lv_log.Infof("<==========%v", req)
	deviceId := cast.ToInt(req.Clientid) //client 全部使用数字，且是device表的主键
	_ = service.GetDeviceCacheService().Offline(deviceId)
	util.SuccessData(c, req)
}
