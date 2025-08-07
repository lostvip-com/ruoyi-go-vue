package service

import (
	global2 "common/global"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
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

// SaveLog  新增记录
func (svc OperLogService) SaveLog(c *gin.Context, bizCode string, inContent any, outContent lv_dto.R) error {
	var operLog model.SysOperLog
	//outJson, _ := json.Marshal(outContent)
	//outJsonStr := string(outJson)
	operLog.Title = bizCode
	if inContent != nil {
		operLog.OperParam = lv_conv.ToJsonStr(inContent)
	}
	//operLog.JsonResult = outJsonStr
	operLog.BusinessType = c.GetInt(global2.KEY_GIN_BIZ_TYPE)
	//操作类别（0其它 1后台用户 2手机端用户）
	operLog.OperatorType = 1
	//操作状态（0正常 1异常）
	if outContent.GetCode() == 0 {
		operLog.Status = 0
	} else {
		operLog.Status = 1
	}
	u, _ := c.Get(global2.KEY_GIN_USER_PTR)
	user := u.(*model.SysUser)
	operLog.OperName = user.UserName
	operLog.RequestMethod = c.Request.Method
	//获取用户部门
	operLog.DeptName = user.Dept.DeptName
	operLog.OperUrl = c.Request.URL.Path
	operLog.Method = c.Request.Method
	operLog.JsonResult = outContent.GetMsg()
	operLog.OperIp = c.ClientIP()
	operLog.OperLocation = util.GetCityByIp(operLog.OperIp)
	operLog.OperTime = time.Now()
	return operLog.Save()
}

// 根据条件分页查询用户列表
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

		if param.BusinessTypes >= 0 {
			tb.Where("business_type = ?", param.BusinessTypes)
		}

		if param.BeginTime != "" {
			tb.Where("date_format(oper_time,'%y%m%d') >= date_format(?,'%y%m%d')", param.BeginTime)
		}

		if param.EndTime != "" {
			tb.Where("date_format(oper_time,'%y%m%d') <= date_format(?,'%y%m%d')", param.EndTime)
		}
	}
	var result []model.SysOperLog
	var total int64
	err := tb.Count(&total).Offset(param.GetStartNum()).Limit(param.GetPageSize()).Order("oper_id desc").Find(&result).Error
	return &result, int(total), err
}

// 根据主键查询用户信息
func (svc OperLogService) FindById(id int) (*model.SysOperLog, error) {
	entity := &model.SysOperLog{OperId: id}
	_, err := entity.FindOne()
	return entity, err
}

func (svc OperLogService) DeleteRecordAll() error {
	err := lv_db.GetOrmDefault().Exec("truncate table sys_oper_log").Error
	return err
}
