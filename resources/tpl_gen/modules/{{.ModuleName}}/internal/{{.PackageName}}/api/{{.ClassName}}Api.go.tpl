package api

import (
	"common/util"
	"common/api"
    "{{.ModuleName}}/internal/{{.PackageName}}/service"
    "{{.ModuleName}}/internal/{{.PackageName}}/model"
    "{{.ModuleName}}/internal/{{.PackageName}}/vo"
    sysService "system/service"
	"github.com/spf13/cast"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"

)

type {{.ClassName}}Api struct {
    api.BaseApi
}

//	========================================================================
//
// api
// =========================================================================

func (w {{.ClassName}}Api) GetRoleInfo(c *gin.Context) {
	id := c.Param("id")
	role := new(model.{{.ClassName}})
	{{.BusinessName}}, err := role.FindById(cast.ToInt64(id))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, {{.BusinessName}})
}
// List{{.ClassName}} 新增页面保存
func (w {{.ClassName}}Api) List{{.ClassName}}(c *gin.Context) {
	req := new(vo.{{.ClassName}}Req)
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var svc = service.Get{{.ClassName}}ServiceInstance()
	result, total, _ := svc.ListByPage(req)
	util.SuccessPage(c, result, total)
}

// Create{{.ClassName}} 新增页面保存
func (w {{.ClassName}}Api) Create{{.ClassName}}(c *gin.Context) {
	form := new(model.{{.ClassName}})
	err := c.ShouldBind(form)
	lv_err.HasErrAndPanic(err)
    w.FillInCreate(c, &form.BaseModel)
	var svc = service.Get{{.ClassName}}ServiceInstance()
	id, err := svc.AddSave(form)
	lv_err.HasErrAndPanic(err)
	util.Success(c, id)
}

// Save{{.ClassName}} 修改页面保存
func (w {{.ClassName}}Api) Update{{.ClassName}}(c *gin.Context) {
	form := new(model.{{.ClassName}})
	err := c.ShouldBind(form)
	lv_err.HasErrAndPanic(err)
    w.FillInUpdate(c, &form.BaseModel)
	var svc = service.Get{{.ClassName}}ServiceInstance()
	err = svc.EditSave(form)
	lv_err.HasErrAndPanic(err)
	util.Success(c, nil)
}

// Remove{{.ClassName}} 删除数据
func (w {{.ClassName}}Api) Delete{{.ClassName}}(c *gin.Context) {
    var ids = c.Param("ids")
	err := service.Get{{.ClassName}}ServiceInstance().DeleteByIds(ids)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

// 导出
func (w {{.ClassName}}Api) Export{{.ClassName}}(c *gin.Context) {
	req := new(vo.{{.ClassName}}Req)
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var svc = service.Get{{.ClassName}}ServiceInstance()
	url, err := svc.ExportAll(req)
	lv_err.HasErrAndPanic(err)
	util.Success(c, url)
}

