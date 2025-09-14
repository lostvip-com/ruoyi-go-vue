package schedule

import (
	"common/models"
	"common/util"
	"github.com/lostvip-com/lv_framework/lv_db"
)

type SysJobLog struct {
	JobLogId      int    `json:"jobLogId" gorm:"column:job_log_id;primary_key;auto_increment;主键;" ` //表示主键
	JobName       string `json:"jobName"  gorm:"size:64;comment:任务名称;"`
	JobGroup      string `json:"jobGroup" gorm:"size:64;comment:任务组名;"`
	InvokeTarget  string `json:"invokeTarget" gorm:"size:64;comment:调用目标函数;"`
	JobMessage    string `json:"jobMessage" gorm:"type:text;comment:异常信息;"`
	ExceptionInfo string `json:"exceptionInfo" gorm:"type:text;comment:异常信息;"`
	StartTime     string `json:"startTime" gorm:"size:11;comment:开始时间;"`
	StopTime      string `json:"stopTime" gorm:"size:11;comment:结束时间;"`
	Status        string `json:"status" gorm:"size:1;comment:执行状态（0正常 1失败）;" `
	models.BaseModel
}

// TableName 指定数据库表名称
func (e *SysJobLog) TableName() string {
	return "sys_job_log"
}

func (e *SysJobLog) FindJobLogById(id int) (*SysJobLog, error) {
	err := lv_db.GetOrmDefault().Where("id = ?", id).First(e).Error
	return e, err
}

func (e *SysJobLog) Save() (*SysJobLog, error) {
	err := lv_db.GetOrmDefault().Create(e).Error
	return e, err
}

func (e *SysJobLog) DetectJobLog(ids string) error {
	arr := util.SplitToInt(ids, ",")
	err := lv_db.GetOrmDefault().Where("id in ( ? )", arr).Delete(&SysJobLog{}).Error
	return err
}

func (e *SysJobLog) ClearJobLog() error {
	err := lv_db.GetOrmDefault().Exec("truncate table sys_job_log").Error
	return err
}
