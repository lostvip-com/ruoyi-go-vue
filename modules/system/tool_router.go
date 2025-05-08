package system

import (
	auth2 "common/middleware/auth"
	"github.com/lostvip-com/lv_framework/web/router"
	"system/api"
)

// 加载路由
func init() {
	// 服务监控
	tool := new(api.GenApi)
	g1 := router.New("/tool", auth2.TokenCheck(), auth2.PermitCheck)
	g1.GET("/build", "", tool.Build)
	g1.GET("/swagger", "", tool.Swagger)
	g1.GET("/gen", "", tool.Gen)
	g1.POST("/gen/list", "tool:gen:list", tool.GenList)
	g1.POST("/gen/remove", "tool:gen:remove", tool.Remove)
	g1.GET("/gen/importTable", "", tool.ImportTable)
	g1.POST("/gen/db/list", "tool:gen:list", tool.DataList)
	g1.POST("/gen/importTable", "tool:gen:list", tool.ImportTableSave)
	g1.GET("/gen/edit", "", tool.Edit)
	g1.POST("/gen/edit", "tool:gen:edit", tool.EditSave)
	g1.POST("/gen/column/list", "tool:gen:list", tool.ColumnList)
	g1.GET("/gen/preview", "", tool.Preview)
	g1.GET("/gen/genCode", "", tool.GenCode)
	//执行sql文件
	g1.GET("/gen/execSqlFile", "", tool.ExecSqlFile)

}
