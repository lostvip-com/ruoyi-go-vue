// ==========================================================================
// LV自动生成数据库操作代码，无需手动修改，重新生成会自动覆盖.
// 生成日期: 2025-05-17 13:18:47 &#43;0000 UTC
// 生成人: lv
// ==========================================================================
package model

import (
    "common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/namedsql"
	"time"
)

// SysJobLog 测试
type SysJobLog struct {
    CreateTime time.Time `gorm:"type:datetime;comment:创建时间;" json:"createTime" time_format:"2006-01-02 15:04:05"`
    ExceptionInfo  string  `gorm:"type:varchar(2000);comment:异常信息;" json:"exceptionInfo"`
    Status  string  `gorm:"type:char(1);comment:执行状态（0正常 1失败）;" json:"status"`
    JobMessage  string  `gorm:"type:varchar(500);comment:日志信息;" json:"jobMessage"`
    InvokeTarget  string  `gorm:"type:varchar(500);comment:调用目标字符串;" json:"invokeTarget"`
    JobGroup  string  `gorm:"type:varchar(64);comment:任务组名;" json:"jobGroup"`
    JobName  string  `gorm:"type:varchar(64);comment:任务名称;" json:"jobName"`
    JobLogId  int64  `gorm:"type:bigint;primary_key;auto_increment;任务日志ID;" json:"jobLogId"`
   models.BaseModel
}

func (e *SysJobLog) TableName() string {
	return "sys_job_log"
}

func (e *SysJobLog) Save() error {
	return lv_db.GetMasterGorm().Save(e).Error
}

func (e *SysJobLog) FindById(id int64) (*SysJobLog,error) {
	err := lv_db.GetMasterGorm().Take(e,id).Error
	return e,err
}

func (e *SysJobLog) FindOne() (*SysJobLog,error) {
    tb := lv_db.GetMasterGorm().Table(e.TableName())

    if e.ExceptionInfo != "" {
         tb = tb.Where("exception_info=?", e.ExceptionInfo)
    }
    if e.Status != "" {
         tb = tb.Where("status=?", e.Status)
    }
    if e.JobMessage != "" {
         tb = tb.Where("job_message=?", e.JobMessage)
    }
    if e.InvokeTarget != "" {
         tb = tb.Where("invoke_target=?", e.InvokeTarget)
    }
    err := tb.First(e).Error
    return e,err
}

func (e *SysJobLog) Updates() error {
	return lv_db.GetMasterGorm().Table(e.TableName()).Updates(e).Error
}

func (e *SysJobLog) Delete() error {
	return lv_db.GetMasterGorm().Delete(e).Error
}

func (e *SysJobLog) Count() (int64, error) {
	sql := " select count(*) from sys_job_log where del_flag = 0 "

	
         if e.ExceptionInfo != "" {
            sql += " and exception_info = @ExceptionInfo "
         }
         if e.Status != "" {
            sql += " and status = @Status "
         }
         if e.JobMessage != "" {
            sql += " and job_message = @JobMessage "
         }
         if e.InvokeTarget != "" {
            sql += " and invoke_target = @InvokeTarget "
         }

	return namedsql.Count(lv_db.GetMasterGorm(), sql, e)
}