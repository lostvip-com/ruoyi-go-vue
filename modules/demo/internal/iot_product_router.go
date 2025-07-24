// ==========================================================================
// LV自动生成路由代码,只生成一次,按需修改,再次生成不会覆盖.
// 生成日期:2025-07-24 02:41:35
// 生成人:lv
// ==========================================================================
package internal

import (
	"common/middleware/auth"
	"demo/internal/product/api"
	"github.com/lostvip-com/lv_framework/web/router"
)

func init() {
	product := router.New("/product/product", auth.TokenCheck())
	productApi := api.IotProductApi{}
	product.GET("/:id", "product:product:info", productApi.GetRoleInfo)
	product.GET("/listIotProduct", "product:product:list", productApi.ListIotProduct)
	product.POST("", "product:product:new", productApi.CreateIotProduct)
	product.PUT("", "product:product:edit", productApi.UpdateIotProduct)
	product.DELETE("/ids", "product:product:del", productApi.DeleteIotProduct)
	product.POST("/exportIotProduct", "product:product:export", productApi.ExportIotProduct)
}
