package api

import (
	"common/global"
	"common/models"
	"common/util"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/spf13/cast"
	"net/http"
	"system/model"
	"system/service"
	"time"
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
	po.UpdateTime = time.Now()
}
func (api *BaseApi) FillInCreate(c *gin.Context, po *models.BaseModel) {
	po.CreateBy = api.GetCurrUsername(c)
	po.CreateTime = time.Now()
}

/************************************************************************************
 * 返回数据并记录日志 相关的公用方法
 *
 ***********************************************************************************/

// LogParamIn 临时保存入参，用于日志记录
func (api *BaseApi) LogParamIn(c *gin.Context, req any) *BaseApi {
	c.Set(global.KEY_GIN_IN_PARAM, req)
	return api
}

// SuccessData 成功数据处理，并记录日志
func (api *BaseApi) SuccessData(c *gin.Context, data any) {
	api.Success(c, data, "Success")
}
func (api *BaseApi) SuccessMsg(c *gin.Context, msg string) {
	api.Success(c, nil, "Success")
}

// Success 成功数据处理，并记录日志
func (api *BaseApi) Success(c *gin.Context, data any, msg string) {
	ret := &lv_dto.Resp{Code: 200, Data: data, Msg: msg}
	inContent, _ := c.Get(global.KEY_GIN_IN_PARAM)
	bizCode := c.Request.Method
	service.GetOperLogServiceInstance().SaveLog(c, bizCode, inContent, ret)
	if data != nil {
		local := c.GetHeader("Accept-Language")
		util.TranslateI18nTagAll(local, data) //只翻译指针类型如：结构体指针 或 结构体切片指针
	}
	c.AbortWithStatusJSON(http.StatusOK, ret)
}

// Error 失败数据处理
func (api *BaseApi) Error(c *gin.Context, err error) {
	var msg string
	if err != nil {
		msg = err.Error()
	}
	util.Fail(c, msg)
}
func (api *BaseApi) ErrResp(c *gin.Context, res lv_dto.Resp) {
	c.AbortWithStatusJSON(http.StatusOK, res)
}

// SuccessPage 分页数据处理 ， 自动翻译 Tag locale标记的字段
func (api *BaseApi) SuccessPage(c *gin.Context, rows any, total any) {
	if rows != nil {
		local := c.GetHeader("Accept-Language")
		util.TranslateI18nTagAll(local, rows) //只翻译指针类型如：结构体指针 或 结构体切片指针
	}
	ret := &lv_dto.RespPage{Code: 200, Rows: rows, Total: total, Msg: "Success"}
	inContent, _ := c.Get(global.KEY_GIN_IN_PARAM)
	bizCode := c.Request.Method
	service.GetOperLogServiceInstance().SaveLog(c, bizCode, inContent, ret)
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "rows": rows, "total": total})
}
