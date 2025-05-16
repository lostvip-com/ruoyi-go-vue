package system

import (
	auth2 "common/middleware/auth"
	"github.com/lostvip-com/lv_framework/web/router"
	"system/api"
)

// 加载路由
func init() {
	codeGenApi := new(api.GenCodeApi)
	tool := router.New("/tool", auth2.TokenCheck(), auth2.PermitCheck)
	tool.GET("/build", "", codeGenApi.Build)
	tool.GET("/swagger", "", codeGenApi.Swagger)
	//code gen
	tool.GET("/gen/list", "tool:gen:list", codeGenApi.GenList)
	tool.GET("/gen/:tableId", "tool:gen:list", codeGenApi.GetGenTableInfo)
	tool.GET("/gen/db/list", "tool:gen:list", codeGenApi.DataList)
	tool.DELETE("/gen/:tableId", "tool:gen:remove", codeGenApi.RemoveByTableId)
	tool.POST("/gen/importTable", "tool:gen:list", codeGenApi.ImportTableSave)
	tool.PUT("/gen", "tool:gen:edit", codeGenApi.EditSave)
	//column
	tool.GET("/gen/column/list", "tool:gen:list", codeGenApi.ColumnList)
	tool.GET("/gen/preview", "", codeGenApi.Preview)
	tool.POST("/gen/genCode", "", codeGenApi.GenCode)
	//执行sql文件
	tool.POST("/gen/execSqlFile", "", codeGenApi.ExecSqlFile)
}
