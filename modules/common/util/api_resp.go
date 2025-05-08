package util

import (
	"common/session"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"net/http"
)

func WriteTpl(c *gin.Context, tpl string, params ...gin.H) {
	var data gin.H
	login := session.GetLoginInfo(c)
	if params == nil || len(params) == 0 {
		data = gin.H{}
	} else {
		data = params[0]
	}
	for k, v := range c.Keys {
		data[k] = v
	}
	data["uid"] = login.UserId
	c.HTML(http.StatusOK, tpl, data)
	c.Abort()
}

// Fail 返回一个成功的消息体
func Fail(c *gin.Context, msg string) {
	ret := lv_dto.CommonRes{
		Code: 500,
		Msg:  msg,
	}
	c.AbortWithStatusJSON(http.StatusOK, &ret)
}

// Success 通常成功数据处理
func Success(c *gin.Context, data any) {
	//if data!=nil{
	//	if reflect.TypeOf(data).Kind() == reflect.Ptr {
	//		util.TranslateByTag(data)
	//	}
	//}
	msg := lv_dto.CommonRes{
		Code: 200,
		Data: data,
		Msg:  "Success",
	}
	c.AbortWithStatusJSON(http.StatusOK, &msg)
	//c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 200, "data": data, "msg": "success"})
}

// Error 失败数据处理
func Error(c *gin.Context, err error) {
	var msg string
	if err != nil {
		msg = err.Error()
	}
	Fail(c, msg)
}
func ErrResp(c *gin.Context, res lv_dto.Resp) {
	c.AbortWithStatusJSON(http.StatusOK, res)
}

// SuccessPage 分页数据处理 ， 自动翻译 Tag locale标记的字段
func SuccessPage(c *gin.Context, rows any, total int64) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "rows": rows, "total": total})
}

///**
// * 自动翻译 Tag locale标记的字段
// */
//func PageOKLocale(c *gin.Context, result []interface{}, count int, pageIndex int, pageSize int, msg string) {
//	var res Page
//	util.TranslateSliceByTag(result)
//	// result.([]interface{})
//	res.List = result
//	res.Count = count
//	res.PageIndex = pageIndex
//	res.PageSize = pageSize
//	Ok(c, res, msg)
//}

// 兼容函数
func Custum(c *gin.Context, data gin.H) {
	c.AbortWithStatusJSON(http.StatusOK, data)
}
