// ==========================================================================
// LV自动生成业务逻辑层相关代码: 只生成一次,按需修改,再次生成不会覆盖.
// date  : 2025-05-17 13:18:47 &#43;0000 UTC
// author: lv
// ==========================================================================
package service

import (
	"time"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_office"
	"github.com/lostvip-com/lv_framework/utils/lv_reflect"
    "demo/internal/iot_driver/model"
    "demo/internal/iot_driver/dao"
    "demo/internal/iot_driver/vo"
)
type SysJobLogService struct{}

var logService *SysJobLogService

func GetSysJobLogServiceInstance() *SysJobLogService {
	if logService == nil {
		logService = &SysJobLogService{}
	}
	return logService
}

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
func (svc SysJobLogService) DeleteByIds(ids string) (int64,error) {
	ida := lv_conv.ToInt64Array(ids, ",")
    var logDao = dao.GetSysJobLogDaoInstance()
	rows,err := logDao.DeleteByIds(ida)
	return rows,err
}

// AddSave 添加数据
func (svc SysJobLogService) AddSave(form *model.SysJobLog) (*model.SysJobLog, error) {
	err := form.Save()
	lv_err.HasErrAndPanic(err)
	return form, err
}

// EditSave 修改数据
func (svc SysJobLogService) EditSave(form *model.SysJobLog)  (*model.SysJobLog, error) {
	var po = new(model.SysJobLog)
	po,err := po.FindById(form.JobLogId)
    if err!=nil{
        return nil,err
    }
	_ = lv_reflect.CopyProperties(form, po)
	err = po.Updates()
	return po,err
}

// ListByPage 根据条件分页查询数据
func (svc SysJobLogService) ListByPage(params *vo.SysJobLogReq) (*[]vo.SysJobLogResp,int64, error) {
	var logDao = dao.GetSysJobLogDaoInstance()
	return logDao.ListByPage(params)
}

// ExportAll 导出excel
func (svc SysJobLogService) ExportAll(param *vo.SysJobLogReq) (string, error) {
    var logDao = dao.GetSysJobLogDaoInstance()
    var err error
    var listMap *[]map[string]any
    if param.PageNum > 0 { //分页导出
        listMap, _, err = logDao.ListMapByPage(param)
    } else { //全部导出
        listMap, err = logDao.ListAll(param, true)
    }
    lv_err.HasErrAndPanic(err)
	heads := []string{  "创建时间" ,"异常信息" ,"执行状态（0正常 1失败）" ,"日志信息" ,"调用目标字符串" ,"任务组名" ,"任务名称" ,"任务日志ID"}
	keys  := []string{  "createTime" ,"exceptionInfo" ,"status" ,"jobMessage" ,"invokeTarget" ,"jobGroup" ,"jobName" ,"jobLogId"}
	url, err := lv_office.DownlaodExcelByListMap(&heads, &keys, listMap)
	return url, err
}