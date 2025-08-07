package schedule

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
)

type SysJobParam struct {
	JobId          int
	Concurrent     string
	CronExpression string
	InvokeTarget   string
	JobGroup       string
	JobName        string
	MisfirePolicy  string
	Status         string
}
type SysJob struct {
	JobId          int    `json:"jobId" form:"jobId" gorm:"column:job_id;primaryKey"` //表示主键
	JobName        string `json:"jobName"  gorm:"size:64;comment:任务名称;"`
	JobGroup       string `json:"jobGroup" gorm:"size:64;comment:任务组名;"`
	InvokeTarget   string `json:"invokeTarget" gorm:"size:64;comment:调用目标函数;"`
	CronExpression string `json:"cronExpression" form:"cronExpression" gorm:"size:64;comment:cron执行表达式;"`
	MisfirePolicy  string `json:"misfirePolicy" form:"misfirePolicy" gorm:"size:64;comment:计划执行错误策略（1立即执行 2执行一次 3放弃执行）;"`
	Concurrent     string `json:"concurrent" form:"concurrent"  gorm:"size:64;comment:是否并发执行（0允许 1禁止）;"`
	Status         string `json:"status" form:"status" gorm:"size:64;comment:调用目标函数;1正常 0暂停"`
	Remark         string `json:"remark" form:"remark" gorm:"size:64;comment:备注信息;"`
	models.BaseModel
}

// TableName 指定数据库表名称
func (e *SysJob) TableName() string {
	return "sys_job"
}

func (e *SysJob) FindById(jobId int) (*SysJob, error) {
	var job SysJob
	err := lv_db.GetOrmDefault().Where("job_id = ?", jobId).First(&job).Error
	return &job, err
}

func (e *SysJob) Save() (*SysJob, error) {
	err := lv_db.GetOrmDefault().Create(e).Error
	return e, err
}

func (e *SysJob) Update() (*SysJob, error) {
	err := lv_db.GetOrmDefault().Updates(e).Error
	return e, err
}
