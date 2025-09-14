package api

import (
	"common/global"
	"common/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	"system/model"
	"system/service"
)

type BaseApi struct {
}

func (api *BaseApi) GetCurrUser(c *gin.Context) *model.SysUser {
	user := service.GetUserServiceInstance().GetCurrUser(c)
	return user
}
func (api *BaseApi) GetCurrUserId(c *gin.Context) int {
	userId, exist := c.Get("userId") //获取用户ID，中间件已处理
	if !exist {
		lv_err.HasError1(errors.New("not login error! userId" + cast.ToString(userId)))
	}
	return cast.ToInt(userId)
}
func (api *BaseApi) GetDeptId(c *gin.Context) int {
	deptId := c.GetInt("deptId") //获取用户ID，中间件已处理
	if deptId == 0 {
		lv_err.HasError1(errors.New("not login error! deptId"))
	}
	return deptId
}

func (api *BaseApi) GetCurrUsername(c *gin.Context) string {
	userName := c.GetString(global.KEY_GIN_USERNAME) //获取用户ID，中间件已处理
	if userName == "" {
		lv_err.HasError1(errors.New("not login error! userName is empty "))
	}
	return userName
}
func (api *BaseApi) IsAdmin(userId int) bool {
	if userId == 1 {
		return true
	} else {
		return false
	}
}

func (api *BaseApi) FillInUpdate(c *gin.Context, po *models.BaseModel) {
	po.UpdateBy = api.GetCurrUsername(c)
}
func (api *BaseApi) FillInCreate(c *gin.Context, po *models.BaseModel) {
	po.CreateBy = api.GetCurrUsername(c)
}
