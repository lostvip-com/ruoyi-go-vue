package service

import (
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_db"
	"system/model"
	"system/vo"
	"time"
)

type OperLogService struct {
}

var operLogService *OperLogService

func GetOperLogServiceInstance() *OperLogService {
	if operLogService == nil {
		operLogService = &OperLogService{}
	}
	return operLogService
}

// FindPage 根据条件分页查询用户列表
func (svc OperLogService) FindPage(param *vo.OperLogPageReq) (*[]model.SysOperLog, int, error) {
	db := lv_db.GetOrmDefault()
	tb := db.Table("sys_oper_log")
	if param != nil {
		if param.Title != "" {
			tb.Where("title like ?", "%"+param.Title+"%")
		}

		if param.OperName != "" {
			tb.Where("oper_name like ?", "%"+param.OperName+"%")
		}

		if param.Status != "" {
			tb.Where("status = ?", param.Status)
		}

		if param.BusinessTypes != "" {
			tb.Where("business_type = ?", param.BusinessTypes)
		}

		if param.BeginTime != "" {
			tb.Where("oper_time >= ? ", param.BeginTime)
		}

		if param.EndTime != "" {
			tb.Where("oper_time <= ? ", param.EndTime)
		}
	}
	var result []model.SysOperLog
	var total int64
	err := tb.Count(&total).Offset(param.GetStartNum()).Limit(param.GetPageSize()).Order("oper_id desc").Find(&result).Error
	return &result, int(total), err
}

func (svc OperLogService) TruncateLogTable() error {
	err := lv_db.GetOrmDefault().Exec("truncate table sys_oper_log").Error
	return err
}

// SaveLog  新增记录
func (svc OperLogService) SaveLog(c *gin.Context, status int, inContent, outContent string, user *model.SysUser) error {
	var operLog model.SysOperLog
	operLog.Title = c.Request.RequestURI
	if c.Request.Method == "POST" { // 0其它 1新增 2修改 3删除
		operLog.BusinessType = 1
	} else if c.Request.Method == "PUT" {
		operLog.BusinessType = 2
	} else if c.Request.Method == "DELETE" {
		operLog.BusinessType = 3
	}
	operLog.OperIp = c.ClientIP()
	//operLog.OperLocation = GetCityByIp(operLog.OperIp)
	if len(inContent) > 64 {
		inContent = inContent[0:64]
	}
	operLog.OperParam = inContent + "..."
	if len(outContent) > 64 {
		outContent = outContent[0:64] + "..."
	}
	operLog.JsonResult = outContent
	//操作类别（0其它 1后台用户 2手机端用户）
	operLog.OperatorType = 1
	//操作状态（0正常 1异常）
	operLog.Status = status
	operLog.OperName = user.UserName
	operLog.RequestMethod = c.Request.Method
	operLog.DeptName = user.Dept.DeptName
	path := c.Request.URL.Path
	if len(path) > 64 {
		path = path[0:64] + "..."
	}
	operLog.OperUrl = path
	operLog.Method = c.Request.Method
	operLog.OperTime = time.Now()
	return operLog.Save()
}
