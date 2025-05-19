// ==========================================================================
// LV自动生成路由代码,只生成一次,按需修改,再次生成不会覆盖.
// 生成日期:{{.CreateTime}}
// 生成人:{{.FunctionAuthor}}
// ==========================================================================
package internal

import (
        "github.com/lostvip-com/lv_framework/web/router"
        "common/middleware/auth"
        "{{.ModuleName}}/internal/{{.PackageName}}/api"

)

func init() {
	{{.BusinessName}} := router.New( "/{{.PackageName}}/{{.BusinessName}}", auth.TokenCheck())

	{{.BusinessName}}Api := api.{{.ClassName}}Api{}
    {{.BusinessName}}.GET("/:id", "{{.PackageName}}:{{.BusinessName}}:info", {{.BusinessName}}Api.GetRoleInfo)
    {{.BusinessName}}.GET("/list{{.ClassName}}", "{{.PackageName}}:{{.BusinessName}}:list", {{.BusinessName}}Api.List{{.ClassName}})
	{{.BusinessName}}.GET("/list{{.ClassName}}", "{{.PackageName}}:{{.BusinessName}}:list", {{.BusinessName}}Api.List{{.ClassName}})
	{{.BusinessName}}.POST("", "{{.PackageName}}:{{.BusinessName}}:new", {{.BusinessName}}Api.Create{{.ClassName}})
	{{.BusinessName}}.PUT("", "{{.PackageName}}:{{.BusinessName}}:edit",{{.BusinessName}}Api.Update{{.ClassName}})
    {{.BusinessName}}.DELETE("/ids", "{{.PackageName}}:{{.BusinessName}}:del", {{.BusinessName}}Api.Delete{{.ClassName}})
	{{.BusinessName}}.POST("/export{{.ClassName}}", "{{.PackageName}}:{{.BusinessName}}:export", {{.BusinessName}}Api.Export{{.ClassName}})
}
