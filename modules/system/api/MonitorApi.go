package api

import (
	api2 "common/api"
	"common/global"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_cache/lv_redis"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_net"
	"github.com/spf13/cast"
	"strings"
	"system/model"
	"system/service"
	"system/vo"
)

type MonitorApi struct {
	api2.BaseApi
}

func (m MonitorApi) CacheHandler(c *gin.Context) {
	var list []vo.CacheVO
	list = append(list, vo.CacheVO{CacheName: global.LoginCacheKey, Remark: "用户信息"})
	list = append(list, vo.CacheVO{CacheName: "sys_config:", Remark: "配置信息"})
	list = append(list, vo.CacheVO{CacheName: global.SysDictCacheKey, Remark: "数据字典"})
	list = append(list, vo.CacheVO{CacheName: "captcha_codes:", Remark: "验证码"})
	list = append(list, vo.CacheVO{CacheName: "repeat_submit:", Remark: "防重提交"})
	list = append(list, vo.CacheVO{CacheName: "rate_limit:", Remark: "限流处理"})
	list = append(list, vo.CacheVO{CacheName: "pwd_err_cnt:", Remark: "密码错误次数"})
	util.Success(c, list)
}
func (m MonitorApi) GetCacheKeysHandler(c *gin.Context) {
	cacheName := c.Param("cacheName")
	redisCache := lv_redis.GetInstance(0)
	keys, _, err := redisCache.Scan(0, cacheName+"*", global.ScanCountMax)
	if err != nil {
		util.Fail(c, err.Error())
	}
	util.Success(c, keys)
}

func (m MonitorApi) GetCacheValueHandler(c *gin.Context) {
	cacheName := c.Param("cacheName")
	cacheKey := c.Param("cacheKey")
	redisCache := lv_redis.GetInstance(0)
	value, err := redisCache.Get(cacheKey)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	var cache = vo.CacheVO{
		CacheName:  cacheName,
		CacheKey:   cacheKey,
		CacheValue: value,
	}
	util.Success(c, cache)
}

func (m MonitorApi) ClearCacheNameHandler(c *gin.Context) {
	cacheName := c.Param("cacheName")
	redisCache := lv_redis.GetInstance(0)
	keys, _, err := redisCache.Scan(0, cacheName+"*", global.ScanCountMax)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	for i := range keys {
		err := redisCache.Del(keys[i])
		if err != nil {
			lv_log.Error("ClearCacheNameHandler error:", err.Error())
		}
	}
	util.Success(c, nil)
}

func (m MonitorApi) ClearCacheKeyHandler(c *gin.Context) {
	cacheKey := c.Param("cacheKey")
	redisCache := lv_redis.GetInstance(0)
	err := redisCache.Del(cacheKey)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

func (m MonitorApi) ClearCacheAllHandler(c *gin.Context) {
	redisCache := lv_redis.GetInstance(0)
	keys, _, err := redisCache.Scan(0, "*", global.ScanCountMax)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	for i := range keys {
		err := redisCache.Del(keys[i])
		if err != nil {
			return
		}
	}
	util.Success(c, nil)
}

func (m MonitorApi) ServerInfo(c *gin.Context) {
	ip := lv_net.GetRemoteClientIp(c.Request)
	var server service.MonitorService
	var info = vo.ServerInfo{}
	info.Cpu = server.GetCpu()
	info.GoInfo = server.GetGoInfo()
	info.Mem = server.GetMem()
	info.Sys = server.GetSys(ip)
	info.SysFile = server.GetSysFile()
	util.Success(c, info)
}

func (m MonitorApi) ListOnLine(c *gin.Context) {
	var ipaddr = c.DefaultQuery("ipaddr", "")
	var userName = c.DefaultQuery("userName", "")
	key := global.LoginCacheKey + "*"
	redisCache := lv_redis.GetInstance(0)
	keyList, _, _ := redisCache.Scan(0, key, 0)
	var rows []map[string]string
	for i := range keyList {
		keyString := keyList[i]
		result, _ := redisCache.HGetAll(keyString)
		rows = append(rows, result)
	}
	if userName != "" || ipaddr != "" { //按条件搜索
		rows = *m.FindSearchTarget(rows, userName, ipaddr)
	}
	util.SuccessPage(c, rows, int64(len(rows)))
}

func (m MonitorApi) FindSearchTarget(rows []map[string]string, userName string, ipaddr string) *[]map[string]string {
	var search = make([]map[string]string, 0)
	for i := range rows {
		row := rows[i]
		if userName != "" || row["userName"] == userName {
			if strings.Contains(userName, row["userName"]) {
				search = append(search, row)
			}
		}
		if ipaddr != "" || row["ipaddr"] == ipaddr {
			if strings.Contains(ipaddr, row["ipaddr"]) {
				search = append(search, row)
			}
		}
	}
	return &search
}

func (m MonitorApi) DetectOnLine(c *gin.Context) {
	var tokenId = c.Param("tokenId")
	var key = global.LoginCacheKey + tokenId
	redisCache := lv_redis.GetInstance(0)
	error := redisCache.Del(key)
	if error != nil {
		util.Fail(c, error.Error())
		return
	}
	util.Success(c, nil)
}

func (m MonitorApi) ListJob(c *gin.Context) {
	var req *vo.JobReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	list, total, err := service.GetJobServiceInstance().FindJobList(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessPage(c, list, total)
}

func (m MonitorApi) ExportJob(c *gin.Context) {
	var req *vo.JobReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	svc := service.GetJobServiceInstance()
	listPtr, _, _ := svc.FindJobList(req)
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":   "jobId",
		"title": "任务序号",
		"width": "10",
	})
	dataKey = append(dataKey, map[string]string{
		"key":   "jobName",
		"title": "任务名称",
		"width": "10",
	})
	dataKey = append(dataKey, map[string]string{
		"key":   "jobGroup",
		"title": "任务组名",
		"width": "10",
	})
	dataKey = append(dataKey, map[string]string{
		"key":   "invokeTarget",
		"title": "调用目标字符串",
		"width": "10",
	})
	dataKey = append(dataKey, map[string]string{
		"key":   "cronExpression",
		"title": "执行表达式",
		"width": "10",
	})
	dataKey = append(dataKey, map[string]string{
		"key":   "misfirePolicy",
		"title": "计划策略",
		"width": "10",
	})
	dataKey = append(dataKey, map[string]string{
		"key":   "concurrent",
		"title": "并发执行",
		"width": "10",
	})
	dataKey = append(dataKey, map[string]string{
		"key":   "status",
		"title": "任务状态",
		"width": "10",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)
	list := *listPtr
	if len(list) > 0 {
		for _, v := range list {
			misfirePolicyKey := v.MisfirePolicy
			var misfirePolicy = ""
			if "0" == misfirePolicyKey {
				misfirePolicy = "默认"
			}
			if "1" == misfirePolicyKey {
				misfirePolicy = "立即触发执行"
			}
			if "2" == misfirePolicyKey {
				misfirePolicy = "触发一次执行"
			}
			if "3" == misfirePolicyKey {
				misfirePolicy = "不触发立即执行"
			}
			concurrentKey := v.Concurrent
			var concurrent = ""
			if "0" == concurrentKey {
				concurrent = "允许"
			}
			if "1" == concurrentKey {
				concurrent = "禁止"
			}
			statusKey := v.Concurrent
			var status = ""
			if "0" == statusKey {
				status = "正常"
			}
			if "1" == statusKey {
				status = "暂停"
			}
			data = append(data, map[string]interface{}{
				"jobId":          v.JobId,
				"jobName":        v.JobName,
				"jobGroup":       v.JobGroup,
				"cronExpression": v.CronExpression,
				"invokeTarget":   v.InvokeTarget,
				"misfirePolicy":  misfirePolicy,
				"concurrent":     concurrent,
				"status":         status,
			})
		}
	}
	ex := util.NewMyExcel()
	ex.ExportToWeb(c, dataKey, data)
}

func (m MonitorApi) GetJobById(c *gin.Context) {
	jobId := c.Param("jobId")
	job := new(model.SysJob)
	result, err := job.FindById(cast.ToInt64(jobId))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, result)
}

func (m MonitorApi) SaveJob(c *gin.Context) {
	job := new(model.SysJob)
	if err := c.ShouldBind(job); err != nil {
		util.Fail(c, err.Error())
		return
	}
	result, err := job.Save()
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, result)
}

func (m MonitorApi) UploadJob(c *gin.Context) {
	job := new(model.SysJob)
	if err := c.ShouldBind(job); err != nil {
		util.Fail(c, err.Error())
		return
	}
	m.FillInUpdate(c, &job.BaseModel)
	result, err := job.Update()
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, result)
}

func (m MonitorApi) ChangeStatus(c *gin.Context) {
	js := new(vo.JobStatus)
	if err := c.ShouldBind(js); err != nil {
		util.Fail(c, err.Error())
		return
	}
	err := lv_db.GetOrmDefault().Table("sys_job").Where("job_id=?", js.JobId).UpdateColumn("status", js.Status).Error
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

func (m MonitorApi) RunJob(c *gin.Context) {
	js := new(vo.JobStatus)
	if err := c.ShouldBind(js); err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Fail(c, "未实现！")
}

func (m MonitorApi) DelectJob(c *gin.Context) {
	jobIds := c.Param("jobIds")
	arr := strings.Split(jobIds, ",")
	err := lv_db.GetOrmDefault().Where("job_id in ( ? )", arr).Delete(&model.SysJob{}).Error
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

func (m MonitorApi) ListJobLog(c *gin.Context) {
	var req *vo.JobReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	list, total, err := service.GetJobServiceInstance().FindJobLogList(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessPage(c, list, total)
}

func (m MonitorApi) ExportJobLog(c *gin.Context) {
	var req *vo.JobReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	//定义首行标题
	dataKey := make([]map[string]string, 0)
	dataKey = append(dataKey, map[string]string{
		"key":    "jobLogId",
		"title":  "日志序号",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "jobName",
		"title":  "任务名称",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "jobGroup",
		"title":  "任务组名",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "invokeTarget",
		"title":  "调用目标字符串",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "jobMessage",
		"title":  "日志信息",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "status",
		"title":  "执行状态",
		"width":  "10",
		"is_num": "0",
	})
	dataKey = append(dataKey, map[string]string{
		"key":    "exceptionInfo",
		"title":  "异常信息",
		"width":  "10",
		"is_num": "0",
	})
	//填充数据
	data := make([]map[string]interface{}, 0)
	svc := service.GetJobServiceInstance()
	listPtr, _, _ := svc.FindJobLogList(req)
	list := *listPtr
	if len(list) > 0 {
		for _, v := range list {
			statusStr := ""
			status := v.Status
			if status == "0" {
				statusStr = "正常"
			}
			if status == "1" {
				statusStr = "失败"
			}
			data = append(data, map[string]interface{}{
				"jobLogId":      v.JobLogId,
				"jobName":       v.JobName,
				"jobGroup":      v.JobGroup,
				"invokeTarget":  v.InvokeTarget,
				"jobMessage":    v.JobMessage,
				"status":        statusStr,
				"exceptionInfo": v.ExceptionInfo,
			})
		}
	}

	ex := util.NewMyExcel()
	ex.ExportToWeb(c, dataKey, data)
}

func (m MonitorApi) GetJobLog(c *gin.Context) {
	var logId = c.Param("logId")
	log := new(model.SysJobLog)
	log, err := log.FindJobLogById(cast.ToInt64(logId))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

func (m MonitorApi) DetectJobLog(c *gin.Context) {
	var ids = c.Param("jobLogIds")
	err := lv_db.GetOrmDefault().Where("id in ( ? )", util.SplitToInt(ids, ",")).Delete(&model.SysJobLog{}).Error
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

func (m MonitorApi) ClearJobLog(c *gin.Context) {
	err := lv_db.GetOrmDefault().Exec("truncate table sys_job_log").Error
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}
