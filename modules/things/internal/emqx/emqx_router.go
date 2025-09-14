// ==========================================================================
// LV自动生成路由代码,只生成一次,按需修改,再次生成不会覆盖.
// 生成日期:2025-03-13 07:55:45 &#43;0000 UTC
// 生成人:lv
// ==========================================================================
package emqx

import (
	"common/middleware/auth"
	"github.com/lostvip-com/lv_framework/web/router"
	"things/internal/emqx/controller"
)

func init() {
	group_emqx := router.New("/iot/emqx", auth.TokenCheck())
	emqx := controller.EmqxHookApi{}
	group_emqx.POST("/authV5", "", emqx.AuthV5)
	group_emqx.POST("/online", "", emqx.Online)
	group_emqx.POST("/offline", "", emqx.Offline)
}
