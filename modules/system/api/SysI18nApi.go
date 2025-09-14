package api

import (
	"common/api"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	"system/model"
	"system/service"
	"system/vo"
)

type SysI18nApi struct {
	api.BaseApi
}

//	========================================================================
//
// api
// =========================================================================

func (w SysI18nApi) GetI18nInfo(c *gin.Context) {
	id := c.Param("id")
	role := model.SysI18n{}
	i18n, err := role.FindById(cast.ToInt(id))
	lv_err.HasErrAndPanic(err)
	util.Success(c, i18n, "Success")
}

// ListSysI18n 查询列表
func (w SysI18nApi) ListSysI18n(c *gin.Context) {
	req := vo.SysI18nReq{}
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	req.BeginTime = c.DefaultQuery("params[beginTime]", "")
	req.EndTime = c.DefaultQuery("params[endTime]", "")

	var svc = service.GetSysI18nServiceInstance()
	result, total, _ := svc.ListByPage(&req)
	util.SuccessPage(c, result, total)
}

// CreateSysI18n 新增页面保存
func (w SysI18nApi) CreateSysI18n(c *gin.Context) {
	form := model.SysI18n{}
	if err := c.ShouldBind(&form); err != nil {
		util.Fail(c, err.Error())
		return
	}
	w.FillInCreate(c, &form.BaseModel)
	var svc = service.GetSysI18nServiceInstance()
	po, err := svc.AddSave(&form)
	lv_err.HasErrAndPanic(err)
	util.Success(c, po, "Success")
}

// SaveSysI18n 修改页面保存
func (w SysI18nApi) UpdateSysI18n(c *gin.Context) {
	form := model.SysI18n{}
	if err := c.ShouldBind(&form); err != nil {
		util.Fail(c, err.Error())
		return
	}
	w.FillInUpdate(c, &form.BaseModel)
	var svc = service.GetSysI18nServiceInstance()
	po, err := svc.EditSave(&form)
	lv_err.HasErrAndPanic(err)
	util.Success(c, po, "Success")
}

// RemoveSysI18n 删除数据
func (w SysI18nApi) DeleteSysI18n(c *gin.Context) {
	var ids = c.Param("ids")
	rows, err := service.GetSysI18nServiceInstance().DeleteByIds(ids)
	lv_err.HasErrAndPanic(err)
	util.Success(c, rows, "Success")
}

// ExportSysI18n 导出
func (w SysI18nApi) ExportSysI18n(c *gin.Context) {
	req := vo.SysI18nReq{}
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	req.BeginTime = c.DefaultQuery("params[beginTime]", "")
	req.EndTime = c.DefaultQuery("params[endTime]", "")

	var svc = service.GetSysI18nServiceInstance()
	headerMap, listMap, err := svc.ExportAll(&req)
	lv_err.HasErrAndPanic(err)
	ex := util.NewMyExcel()
	ex.ExportToWeb(c, *headerMap, *listMap)
}
