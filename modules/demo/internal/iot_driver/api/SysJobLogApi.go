package api

import (
	"common/util"
	"common/api"
    "demo/internal/iot_driver/service"
    "demo/internal/iot_driver/model"
    "demo/internal/iot_driver/vo"
    sysService "system/service"
	"github.com/spf13/cast"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"

)

type SysJobLogApi struct {
    api.BaseApi
}

//	========================================================================
//
// api
// =========================================================================

func (w SysJobLogApi) GetRoleInfo(c *gin.Context) {
	id := c.Param("id")
	role := new(model.SysJobLog)
	log, err := role.FindById(cast.ToInt64(id))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, log)
}
// ListSysJobLog 新增页面保存
func (w SysJobLogApi) ListSysJobLog(c *gin.Context) {
	req := new(vo.SysJobLogReq)
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var svc = service.GetSysJobLogServiceInstance()
	result, total, _ := svc.ListByPage(req)
	util.SuccessPage(c, result, total)
}

// CreateSysJobLog 新增页面保存
func (w SysJobLogApi) CreateSysJobLog(c *gin.Context) {
	form := new(model.SysJobLog)
	err := c.ShouldBind(form)
	lv_err.HasErrAndPanic(err)
    w.FillInCreate(c, &form.BaseModel)
	var svc = service.GetSysJobLogServiceInstance()
	id, err := svc.AddSave(form)
	lv_err.HasErrAndPanic(err)
	util.Success(c, id)
}

// SaveSysJobLog 修改页面保存
func (w SysJobLogApi) UpdateSysJobLog(c *gin.Context) {
	form := new(model.SysJobLog)
	err := c.ShouldBind(form)
	lv_err.HasErrAndPanic(err)
    w.FillInUpdate(c, &form.BaseModel)
	var svc = service.GetSysJobLogServiceInstance()
	err = svc.EditSave(form)
	lv_err.HasErrAndPanic(err)
	util.Success(c, nil)
}

// RemoveSysJobLog 删除数据
func (w SysJobLogApi) DeleteSysJobLog(c *gin.Context) {
    var ids = c.Param("ids")
	err := service.GetSysJobLogServiceInstance().DeleteByIds(ids)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

// 导出
func (w SysJobLogApi) ExportSysJobLog(c *gin.Context) {
	req := new(vo.SysJobLogReq)
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var svc = service.GetSysJobLogServiceInstance()
	url, err := svc.ExportAll(req)
	lv_err.HasErrAndPanic(err)
	util.Success(c, url)
}

