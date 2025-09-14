package schedule

import (
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/robfig/cron/v3"
	"sync"
)

// SchedulerService 定时任务调度器服务
type SchedulerService struct {
	cron       *cron.Cron
	jobs       map[int]cron.EntryID    // 存储任务ID与cron EntryID的映射
	jobFuncMap map[string]func() error // 存储执行函数的映射
	mutex      sync.RWMutex
}

var schedulerService *SchedulerService
var once sync.Once

// GetSchedulerServiceInstance 获取调度器服务单例实例
func GetSchedulerServiceInstance() *SchedulerService {
	once.Do(func() {
		schedulerService = &SchedulerService{
			cron:       cron.New(cron.WithSeconds()), // 支持秒级调度
			jobs:       make(map[int]cron.EntryID),
			jobFuncMap: make(map[string]func() error),
		}
		// 启动调度器
		schedulerService.cron.Start()
		// 初始化时加载所有已存在的任务
		schedulerService.loadAllJobs()
	})
	return schedulerService
}

// loadAllJobs 初始化时加载所有数据库中的任务
func (s *SchedulerService) loadAllJobs() {
	var jobs []SysJob
	err := lv_db.GetOrmDefault().Table("sys_job").Where("status = ?", "0").Find(&jobs).Error
	if err != nil {
		lv_log.Error("加载定时任务失败:", err)
		return
	}

	for _, job := range jobs {
		s.AddJob(&job)
	}
}

// AddJob 添加定时任务
func (s *SchedulerService) AddJob(job *SysJob) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 如果任务已存在，先移除
	if entryID, exists := s.jobs[job.JobId]; exists {
		s.cron.Remove(entryID)
		delete(s.jobs, job.JobId)
	}

	// 只有启用状态的任务才添加到调度器中
	if job.Status == "0" {
		// 创建任务执行函数
		jobFunc := s.createJobFunc(job)

		// 添加任务到cron调度器
		entryID, err := s.cron.AddFunc(job.CronExpression, jobFunc)
		if err != nil {
			lv_log.Error("添加定时任务失败:", err)
			return err
		}

		// 记录任务ID与EntryID的映射关系
		s.jobs[job.JobId] = entryID
		lv_log.Info("添加定时任务成功:", job.JobName, " cron:", job.CronExpression)
	}

	return nil
}

// UpdateJob 更新定时任务
func (s *SchedulerService) UpdateJob(job *SysJob) error {
	// 更新任务相当于先移除再添加
	s.removeFromCache(job.JobId)
	return s.AddJob(job)
}

// RemoveJob 删除定时任务
func (s *SchedulerService) removeFromCache(jobId int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if entryID, exists := s.jobs[jobId]; exists {
		s.cron.Remove(entryID)
		delete(s.jobs, jobId)
		lv_log.Info("删除定时任务成功:", jobId)
	}

	return nil
}

// StartJob 启动定时任务
func (s *SchedulerService) StartJob(job *SysJob) error {
	// 更新任务状态为启用
	job.Status = "0"
	_, err := job.Update()
	if err != nil {
		return err
	}

	// 添加任务到调度器
	return s.AddJob(job)
}

// StopJob 停止定时任务
func (s *SchedulerService) StopJob(job *SysJob) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 更新任务状态为停止
	job.Status = "1"
	_, err := job.Update()
	if err != nil {
		return err
	}

	// 从调度器中移除任务
	if entryID, exists := s.jobs[job.JobId]; exists {
		s.cron.Remove(entryID)
		delete(s.jobs, job.JobId)
		lv_log.Info("停止定时任务成功:", job.JobId)
	}

	return nil
}

// RunJobOnce 立即执行一次任务
func (s *SchedulerService) RunJobOnce(job *SysJob) {
	jobFunc := s.createJobFunc(job)
	go jobFunc() // 在goroutine中执行，避免阻塞
}

// createJobFunc 创建任务执行函数
func (s *SchedulerService) createJobFunc(job *SysJob) func() {
	return func() {
		lv_log.Info("开始执行定时任务:", job.JobName)

		// 记录任务执行日志
		jobLog := &SysJobLog{
			JobName:      job.JobName,
			JobGroup:     job.JobGroup,
			InvokeTarget: job.InvokeTarget,
			JobMessage:   "开始执行任务",
			Status:       "0", // 默认执行成功
		}

		// 根据 invokeTarget 调用具体的方法
		if jobFunc, exists := s.jobFuncMap[job.InvokeTarget]; exists {
			jobFunc()
		} else {
			lv_log.Error("请实现函数:", job.InvokeTarget)
			jobLog.Status = "1" // 执行失败
			jobLog.JobMessage = "未找到指定的函数"
		}
		// 保存执行日志
		_, err := jobLog.Save()
		if err != nil {
			lv_log.Error("保存任务执行日志失败:", err)
		}

		lv_log.Info("定时任务执行完成:", job.JobName)
	}
}

// GetEntryID 获取任务的EntryID
func (s *SchedulerService) GetEntryID(jobId int) (cron.EntryID, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entryID, exists := s.jobs[jobId]
	return entryID, exists
}

// ListEntries 获取所有任务条目
func (s *SchedulerService) ListEntries() []cron.Entry {
	return s.cron.Entries()
}

func (s *SchedulerService) RegFunc(jobName string, callFunc func() error) {
	s.jobFuncMap[jobName] = callFunc
}
