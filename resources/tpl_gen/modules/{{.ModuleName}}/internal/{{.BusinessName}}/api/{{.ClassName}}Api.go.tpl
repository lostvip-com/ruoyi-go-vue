package api

import (
  "common/util"
  "common/api"
  "{{.ModuleName}}/internal/{{.BusinessName}}/service"
  "{{.ModuleName}}/internal/{{.BusinessName}}/model"
  "{{.ModuleName}}/internal/{{.BusinessName}}/vo"
  "github.com/spf13/cast"
  "github.com/gin-gonic/gin"
  "github.com/lostvip-com/lv_framework/utils/lv_err"
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
	role := model.{{.ClassName}}{}
	{{.BusinessName}}, err := role.FindById(cast.ToInt64(id))
	lv_err.HasErrAndPanic(err)
	util.Success(c, {{.BusinessName}})
}
// List{{.ClassName}} 查询列表
func (w {{.ClassName}}Api) List{{.ClassName}}(c *gin.Context) {
	req := vo.{{.ClassName}}Req{}
    if err := c.ShouldBind(&req); err != nil {
        util.Fail(c, err.Error())
        return
    }
    req.BeginTime = c.DefaultQuery("params[beginTime]", "")
    req.EndTime = c.DefaultQuery("params[endTime]", "")

	var svc = service.Get{{.ClassName}}ServiceInstance()
	result, total, _ := svc.ListByPage(&req)
	util.SuccessPage(c, result, total)
}

// Create{{.ClassName}} 新增页面保存
func (w {{.ClassName}}Api) Create{{.ClassName}}(c *gin.Context) {
	form := model.{{.ClassName}}{}
    if err := c.ShouldBind(&form); err != nil {
        util.Fail(c, err.Error())
        return
    }
    w.FillInCreate(c, &form.BaseModel)
	var svc = service.Get{{.ClassName}}ServiceInstance()
	po,err := svc.AddSave(&form)
	lv_err.HasErrAndPanic(err)
	util.Success(c, po)
}

// Save{{.ClassName}} 修改页面保存
func (w {{.ClassName}}Api) Update{{.ClassName}}(c *gin.Context) {
	form := model.{{.ClassName}}{}
    if err := c.ShouldBind(&form); err != nil {
        util.Fail(c, err.Error())
        return
    }
    w.FillInUpdate(c, &form.BaseModel)
	var svc = service.Get{{.ClassName}}ServiceInstance()
	po,err := svc.EditSave(&form)
	lv_err.HasErrAndPanic(err)
	util.Success(c, po)
}

// Remove{{.ClassName}} 删除数据
func (w {{.ClassName}}Api) Delete{{.ClassName}}(c *gin.Context) {
    var ids = c.Param("ids")
	rows,err := service.Get{{.ClassName}}ServiceInstance().DeleteByIds(ids)
	lv_err.HasErrAndPanic(err)
	util.Success(c, rows)
}

//Export{{.ClassName}} 导出
func (w {{.ClassName}}Api) Export{{.ClassName}}(c *gin.Context) {
	req := vo.{{.ClassName}}Req{}
    if err := c.ShouldBind(&req); err != nil {
        util.Fail(c, err.Error())
        return
    }
    req.BeginTime = c.DefaultQuery("params[beginTime]", "")
    req.EndTime = c.DefaultQuery("params[endTime]", "")

	var svc = service.Get{{.ClassName}}ServiceInstance()
    headerMap, listMap, err := svc.ExportAll(&req)
    lv_err.HasErrAndPanic(err)
    ex := util.NewMyExcel()
    ex.ExportToWeb(c, *headerMap, *listMap)
}

