package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	"system/model"
)

type BaseApi struct {
}

func (w *BaseApi) GetCurrUser(c *gin.Context) model.SysUser {
	user, exist := c.Get("user") //获取用户信息，中间件已处理
	if !exist {
		lv_err.HasError1(errors.New("not login error! "))
	}
	return user.(model.SysUser)
}
func (w *BaseApi) GetCurrUserId(c *gin.Context) int64 {
	userId, exist := c.Get("userId") //获取用户ID，中间件已处理
	if !exist {
		lv_err.HasError1(errors.New("not login error! userId" + cast.ToString(userId)))
	}
	return cast.ToInt64(userId)
}
func (svc *BaseApi) IsAdmin(userId int64) bool {
	if userId == 1 {
		return true
	} else {
		return false
	}
}
