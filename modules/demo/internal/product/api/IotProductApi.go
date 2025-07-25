package api

import (
  "common/util"
  "common/api"
  "demo/internal/product/service"
  "demo/internal/product/model"
  "demo/internal/product/vo"
  "github.com/spf13/cast"
  "github.com/gin-gonic/gin"
  "github.com/lostvip-com/lv_framework/utils/lv_err"
)

type IotProductApi struct {
    api.BaseApi
}

//	========================================================================
//
// api
// =========================================================================

func (w IotProductApi) GetRoleInfo(c *gin.Context) {
	id := c.Param("id")
	role := model.IotProduct{}
	product, err := role.FindById(cast.ToInt64(id))
	lv_err.HasErrAndPanic(err)
	util.Success(c, product)
}
// ListIotProduct 查询列表
func (w IotProductApi) ListIotProduct(c *gin.Context) {
	req := vo.IotProductReq{}
    if err := c.ShouldBind(&req); err != nil {
        util.Fail(c, err.Error())
        return
    }
    req.BeginTime = c.DefaultQuery("params[beginTime]", "")
    req.EndTime = c.DefaultQuery("params[endTime]", "")

	var svc = service.GetIotProductServiceInstance()
	result, total, _ := svc.ListByPage(&req)
	util.SuccessPage(c, result, total)
}

// CreateIotProduct 新增页面保存
func (w IotProductApi) CreateIotProduct(c *gin.Context) {
	form := model.IotProduct{}
    if err := c.ShouldBind(&form); err != nil {
        util.Fail(c, err.Error())
        return
    }
    w.FillInCreate(c, &form.BaseModel)
	var svc = service.GetIotProductServiceInstance()
	po,err := svc.AddSave(&form)
	lv_err.HasErrAndPanic(err)
	util.Success(c, po)
}

// SaveIotProduct 修改页面保存
func (w IotProductApi) UpdateIotProduct(c *gin.Context) {
	form := model.IotProduct{}
    if err := c.ShouldBind(&form); err != nil {
        util.Fail(c, err.Error())
        return
    }
    w.FillInUpdate(c, &form.BaseModel)
	var svc = service.GetIotProductServiceInstance()
	po,err := svc.EditSave(&form)
	lv_err.HasErrAndPanic(err)
	util.Success(c, po)
}

// RemoveIotProduct 删除数据
func (w IotProductApi) DeleteIotProduct(c *gin.Context) {
    var ids = c.Param("ids")
	rows,err := service.GetIotProductServiceInstance().DeleteByIds(ids)
	lv_err.HasErrAndPanic(err)
	util.Success(c, rows)
}

//ExportIotProduct 导出
func (w IotProductApi) ExportIotProduct(c *gin.Context) {
	req := vo.IotProductReq{}
    if err := c.ShouldBind(&req); err != nil {
        util.Fail(c, err.Error())
        return
    }
    req.BeginTime = c.DefaultQuery("params[beginTime]", "")
    req.EndTime = c.DefaultQuery("params[endTime]", "")

	var svc = service.GetIotProductServiceInstance()
    headerMap, listMap, err := svc.ExportAll(&req)
    lv_err.HasErrAndPanic(err)
    ex := util.NewMyExcel()
    ex.ExportToWeb(c, *headerMap, *listMap)
}

