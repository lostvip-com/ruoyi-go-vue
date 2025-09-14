package service

import (
	"common/schedule"
	"system/vo"
)

type JobService struct {
	BaseService
}

var jobService *JobService

func GetJobServiceInstance() *JobService {
	if jobService == nil {
		jobService = &JobService{}
	}
	return jobService
}

func (e *JobService) FindJobList(params *vo.JobReq) (*[]schedule.SysJob, int, error) {
	var total int64
	db := e.GetDb().Table("sys_job")
	if params.JobName != "" {
		db.Where("job_name like ?", "%"+params.JobName+"%")
	}
	if params.JobGroup != "" {
		db.Where("job_group = ?", params.JobGroup)
	}
	if params.Status != "" {
		db.Where("status = ?", params.Status)
	}
	if params.InvokeTarget != "" {
		db.Where("invoke_target like concat('%', ?, '%')", params.InvokeTarget)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db.Order("job_id DESC")
	var list []schedule.SysJob
	if params.PageNum >= 1 && params.PageSize > 0 {
		offset := (params.PageNum - 1) * params.PageSize
		db.Offset(offset).Limit(params.PageSize)
	}
	err := db.Find(&list).Error
	return &list, int(total), err
}
func (e *JobService) FindJobLogList(params *vo.JobReq) (*[]schedule.SysJobLog, int, error) {
	var total int64
	db := e.GetDb().Table("sys_job_log")
	if params.JobName != "" {
		db.Where("job_name like ?", "%"+params.JobName+"%")
	}
	if params.JobGroup != "" {
		db.Where("job_group = ?", params.JobGroup)
	}

	if params.Status != "" {
		db.Where("status = ?", params.Status)
	}
	if params.BeginTime != "" {
		db.Where("create_time >= ?", params.BeginTime)
		db.Where("create_time <= ?", params.EndTime+" 23:59:59")
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	db.Order("job_log_id DESC")

	var list []schedule.SysJobLog
	if params.PageNum >= 1 && params.PageSize > 0 {
		offset := (params.PageNum - 1) * params.PageSize
		db.Offset(offset).Limit(params.PageSize)
	}
	err := db.Find(&list).Error
	return &list, int(total), err
}

func (e *JobService) ListActive() (*[]schedule.SysJob, error) {
	db := e.GetDb().Table("sys_job").Where("status = 0 ")
	var list []schedule.SysJob
	err := db.Find(&list).Error
	return &list, err
}
