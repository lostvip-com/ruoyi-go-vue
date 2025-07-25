// ==========================================================================
// LV自动生成路由代码,只生成一次,按需修改,再次生成不会覆盖.
// 生成日期:2025-07-25 02:44:32
// 生成人:lv
// ==========================================================================
package internal

import (
	"common/middleware/auth"
	"demo/internal/product/api"
	"github.com/lostvip-com/lv_framework/web/router"
)

func init() {
	product := router.New("/demo/product", auth.TokenCheck())
	productApi := api.IotProductApi{}
	product.GET("/:id", "demo:product:info", productApi.GetRoleInfo)
	product.GET("/listIotProduct", "demo:product:list", productApi.ListIotProduct)
	product.POST("", "demo:product:new", productApi.CreateIotProduct)
	product.PUT("", "demo:product:edit", productApi.UpdateIotProduct)
	product.DELETE("/ids", "demo:product:del", productApi.DeleteIotProduct)
	product.POST("/exportIotProduct", "product:product:export", productApi.ExportIotProduct)
}
