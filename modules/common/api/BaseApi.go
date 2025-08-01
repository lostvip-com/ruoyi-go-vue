package api

import (
	"common/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	"system/model"
	"system/service"
	"time"
)

type BaseApi struct {
}

func (w *BaseApi) GetCurrUser(c *gin.Context) *model.SysUser {
	user := service.GetUserServiceInstance().GetCurrUser(c)
	return user
}
func (w *BaseApi) GetCurrUserId(c *gin.Context) int {
	userId, exist := c.Get("userId") //获取用户ID，中间件已处理
	if !exist {
		lv_err.HasError1(errors.New("not login error! userId" + cast.ToString(userId)))
	}
	return cast.ToInt(userId)
}
func (w *BaseApi) GetDeptId(c *gin.Context) int {
	deptId := c.GetInt("deptId") //获取用户ID，中间件已处理
	if deptId == 0 {
		lv_err.HasError1(errors.New("not login error! deptId"))
	}
	return deptId
}

func (w *BaseApi) GetCurrUsername(c *gin.Context) string {
	userName := c.GetString("userName") //获取用户ID，中间件已处理
	if userName == "" {
		lv_err.HasError1(errors.New("not login error! userName is empty "))
	}
	return userName
}
func (svc *BaseApi) IsAdmin(userId int) bool {
	if userId == 1 {
		return true
	} else {
		return false
	}
}

func (svc *BaseApi) FillInUpdate(c *gin.Context, po *models.BaseModel) {
	po.UpdateBy = svc.GetCurrUsername(c)
	po.UpdateTime = time.Now()
}
func (svc *BaseApi) FillInCreate(c *gin.Context, po *models.BaseModel) {
	po.CreateBy = svc.GetCurrUsername(c)
	po.CreateTime = time.Now()
}
