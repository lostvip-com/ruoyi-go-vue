package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
)

type SysJobParam struct {
	JobId          int64
	Concurrent     string
	CronExpression string
	InvokeTarget   string
	JobGroup       string
	JobName        string
	MisfirePolicy  string
	Status         string
}
type SysJob struct {
	JobId          int64  `json:"jobId" form:"jobId" gorm:"column:job_id;primaryKey"` //表示主键
	JobName        string `json:"jobName" form:"jobName" gorm:"job_name"`
	JobGroup       string `json:"jobGroup" form:"jobGroup" gorm:"job_group"`
	InvokeTarget   string `json:"invokeTarget" form:"invokeTarget" gorm:"invoke_target"`
	CronExpression string `json:"cronExpression" form:"cronExpression" gorm:"cron_expression"`
	MisfirePolicy  int    `json:"misfirePolicy" form:"misfirePolicy" gorm:"misfire_policy"`
	Concurrent     int    `json:"concurrent" form:"concurrent"  gorm:"concurrent"`
	Status         string `json:"status" form:"status" gorm:"status"`
	Remark         string `json:"remark" form:"remark" gorm:"remark"`
	models.BaseModel
}

// TableName 指定数据库表名称
func (e *SysJob) TableName() string {
	return "sys_job"
}

func (e *SysJob) FindById(jobId int64) (*SysJob, error) {
	var job SysJob
	err := lv_db.GetMasterGorm().Where("job_id = ?", jobId).First(&job).Error
	return &job, err
}

func (e *SysJob) Save() (*SysJob, error) {
	err := lv_db.GetMasterGorm().Create(e).Error
	return e, err
}

func (e *SysJob) Update() (*SysJob, error) {
	err := lv_db.GetMasterGorm().Updates(e).Error
	return e, err
}
