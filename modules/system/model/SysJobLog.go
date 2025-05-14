package model

import (
	"common/util"
	"github.com/lostvip-com/lv_framework/lv_db"
	"time"
)

type SysJobLog struct {
	JobLogId      int       `json:"jobLogId" gorm:"column:job_log_id;primaryKey"` //表示主键
	JobName       string    `json:"jobName" gorm:"job_name"`
	JobGroup      string    `json:"jobGroup" gorm:"job_group"`
	InvokeTarget  string    `json:"invokeTarget" gorm:"invoke_target"`
	JobMessage    string    `json:"jobMessage" gorm:"job_message"`
	ExceptionInfo string    `json:"exceptionInfo"`
	StartTime     string    `json:"startTime"`
	StopTime      string    `json:"stopTime"`
	Status        string    `json:"status" gorm:"status"`
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time;type:datetime"`
}

// TableName 指定数据库表名称
func (e *SysJobLog) TableName() string {
	return "sys_job_log"
}

func (e *SysJobLog) FindJobLogById(id int64) (*SysJobLog, error) {
	err := lv_db.GetMasterGorm().Where("id = ?", id).First(e).Error
	return e, err
}

func (e *SysJobLog) Save() (*SysJobLog, error) {
	err := lv_db.GetMasterGorm().Create(e).Error
	return e, err
}

func (e *SysJobLog) DetectJobLog(ids string) error {
	arr := util.SplitToInt(ids, ",")
	err := lv_db.GetMasterGorm().Where("id in ( ? )", arr).Delete(&SysJobLog{}).Error
	return err
}

func (e *SysJobLog) ClearJobLog() error {
	err := lv_db.GetMasterGorm().Exec("truncate table sys_job_log").Error
	return err
}
