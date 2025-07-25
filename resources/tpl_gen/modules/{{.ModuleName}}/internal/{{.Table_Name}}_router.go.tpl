// ==========================================================================
// LV自动生成路由代码,只生成一次,按需修改,再次生成不会覆盖.
// 生成日期:{{FmtTime .CreateTime}}
// 生成人:{{.FunctionAuthor}}
// ==========================================================================
package internal

import (
        "github.com/lostvip-com/lv_framework/web/router"
        "common/middleware/auth"
        "{{.ModuleName}}/internal/{{.BusinessName}}/api"

)

func init() {
    {{.BusinessName}} := router.New( "/{{.ModuleName}}/{{.BusinessName}}", auth.TokenCheck())
    {{.BusinessName}}Api := api.{{.ClassName}}Api{}
    {{.BusinessName}}.GET("/:id", "{{.ModuleName}}:{{.BusinessName}}:info", {{.BusinessName}}Api.GetRoleInfo)
    {{.BusinessName}}.GET("/list{{.ClassName}}", "{{.ModuleName}}:{{.BusinessName}}:list", {{.BusinessName}}Api.List{{.ClassName}})
    {{.BusinessName}}.POST("", "{{.ModuleName}}:{{.BusinessName}}:new", {{.BusinessName}}Api.Create{{.ClassName}})
    {{.BusinessName}}.PUT("", "{{.ModuleName}}:{{.BusinessName}}:edit",{{.BusinessName}}Api.Update{{.ClassName}})
    {{.BusinessName}}.DELETE("/ids", "{{.ModuleName}}:{{.BusinessName}}:del", {{.BusinessName}}Api.Delete{{.ClassName}})
    {{.BusinessName}}.POST("/export{{.ClassName}}", "{{.BusinessName}}:{{.BusinessName}}:export", {{.BusinessName}}Api.Export{{.ClassName}})
}
