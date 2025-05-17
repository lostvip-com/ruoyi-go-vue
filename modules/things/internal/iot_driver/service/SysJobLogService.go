// ==========================================================================
// LV自动生成业务逻辑层相关代码: 只生成一次,按需修改,再次生成不会覆盖.
// date  : 2025-05-16 05:37:58 &#43;0000 UTC
// author: lv
// ==========================================================================
package service

import (
	"time"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_office"
	"github.com/lostvip-com/lv_framework/utils/lv_reflect"
    "things/internal/iot_driver/model"
    "things/internal/iot_driver/dao"
    "things/internal/iot_driver/vo"
)
type SysJobLogService struct{}
// FindById 根据主键查询数据
func (svc SysJobLogService) FindById(id int64) (*model.SysJobLog, error) {
	var po = new(model.SysJobLog)
    po,err := po.FindById(id)
	return po, err
}

// DeleteById 根据主键删除数据
func (svc SysJobLogService) DeleteById(id int64) error {
	err := (&model.SysJobLog{ JobLogId: id}).Delete()
	return err
}

// DeleteByIds 批量删除数据记录
func (svc SysJobLogService) DeleteByIds(ids string) int64 {
	ida := lv_conv.ToInt64Array(ids, ",")
    var d dao.SysJobLogDao
	rowsAffected := d.DeleteByIds(ida)
	return rowsAffected
}

// AddSave 添加数据
func (svc SysJobLogService) AddSave(req *vo.AddSysJobLogReq) (int64, error) {
	var entity = new(model.SysJobLog)
	lv_reflect.CopyProperties(req, entity)
	entity.CreateTime = time.Now()
	entity.UpdateTime = entity.CreateTime
	entity.CreateBy = req.CreateBy
	err := entity.Save()
	lv_err.HasErrAndPanic(err)
	return entity.JobLogId, err
}

// EditSave 修改数据
func (svc SysJobLogService) EditSave(req *vo.EditSysJobLogReq) error {
	var po = new(model.SysJobLog)
	po,err := po.FindById(req.JobLogId)
    lv_err.HasErrAndPanic(err)
	lv_reflect.CopyProperties(req, po)
	po.UpdateTime = time.Now()
	po.UpdateBy = req.UpdateBy
	return po.Updates()
}

// ListByPage 根据条件分页查询数据
func (svc SysJobLogService) ListByPage(params *vo.SysJobLogReq) (*[]vo.SysJobLogResp,int64, error) {
	var d dao.SysJobLogDao
	return d.ListByPage(params)
}

// ExportAll 导出excel
func (svc SysJobLogService) ExportAll(param *vo.SysJobLogReq) (string, error) {
    var d dao.SysJobLogDao
    var err error
    var listMap *[]map[string]any
    if param.PageNum > 0 { //分页导出
        listMap, _, err = d.ListMapByPage(param)
    } else { //全部导出
        listMap, err = d.ListAll(param, true)
    }
    lv_err.HasErrAndPanic(err)
	heads := []string{  "创建时间" ,"异常信息" ,"执行状态（0正常 1失败）" ,"日志信息" ,"调用目标字符串" ,"任务组名" ,"任务名称" ,"任务日志ID"}
	keys  := []string{  "createTime" ,"exceptionInfo" ,"status" ,"jobMessage" ,"invokeTarget" ,"jobGroup" ,"jobName" ,"jobLogId"}
	url, err := lv_office.DownlaodExcelByListMap(&heads, &keys, listMap)
	return url, err
}