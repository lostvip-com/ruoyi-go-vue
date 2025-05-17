package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	sysService "system/service"
	"things/internal/iot_driver/service"
	"things/internal/iot_driver/vo"
)

type SysJobLogApi struct {
}

//	========================================================================
//
// api
// =========================================================================

// ListSysJobLog 新增页面保存
func (w SysJobLogApi) ListSysJobLog(c *gin.Context) {
	req := new(vo.SysJobLogReq)
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	//if req.DeptId == 0 {
	//	req.DeptId = session.GetLoginInfo(c).DeptId
	//}
	var svc service.SysJobLogService
	result, total, _ := svc.ListByPage(req)
	util.SucessPage(c, result, total)
}

// CreateSysJobLog 新增页面保存
func (w SysJobLogApi) AddSysJobLog(c *gin.Context) {
	req := new(vo.AddSysJobLogReq)
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var svc service.SysJobLogService
	var userService sysService.UserService
	user := userService.GetProfile(c)
	req.CreateBy = user.UserName
	id, err := svc.AddSave(req)
	lv_err.HasErrAndPanic(err)
	util.SucessData(c, id)
}

// SaveSysJobLog 修改页面保存
func (w SysJobLogApi) SaveSysJobLog(c *gin.Context) {
	req := new(vo.EditSysJobLogReq)
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var svc service.SysJobLogService
	var userService sysService.UserService
	user := userService.GetProfile(c)
	req.UpdateBy = user.UserName
	err = svc.EditSave(req)
	lv_err.HasErrAndPanic(err)
	util.Success(c, nil, "success")
}

// RemoveSysJobLog 删除数据
func (w SysJobLogApi) RemoveSysJobLog(c *gin.Context) {
	req := new(lv_dto.IdsReq)
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var svc service.SysJobLogService
	rs := svc.DeleteByIds(req.Ids)
	util.SuccessData(c, rs)
}

// 导出
func (w SysJobLogApi) ExportSysJobLog(c *gin.Context) {
	req := new(vo.SysJobLogReq)
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var svc service.SysJobLogService
	url, err := svc.ExportAll(req)
	lv_err.HasErrAndPanic(err)
	util.SucessDataMsg(c, url, url)
}

