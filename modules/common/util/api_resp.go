package util

import (
	"common/session"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"net/http"
	"reflect"
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
	c.AbortWithStatusJSON(http.StatusOK, &lv_dto.CommonRes{Code: 500, Msg: msg})
}

func FailNoAuth(c *gin.Context, msg string) {
	ret := lv_dto.CommonRes{
		Code: http.StatusUnauthorized,
		Msg:  msg,
	}
	c.AbortWithStatusJSON(http.StatusOK, &ret)
}

func Success(c *gin.Context, data any, msg string) {
	if data != nil {
		local := c.GetHeader("Accept-Language")
		TranslateI18nTagAll(local, data) //只翻译指针类型如：结构体指针 或 结构体切片指针
	}
	c.AbortWithStatusJSON(http.StatusOK, &lv_dto.CommonRes{
		Code: 200,
		Data: data,
		Msg:  msg,
	})
}

// TranslateI18nTagAll 递归翻译 data 内部的所有结构体
// 支持的类型：
// 1. 结构体指针 (*struct)
// 2. 结构体 (struct)
// 3. 切片指针 (*[]struct)
// 4. 切片 ([]struct)
func TranslateI18nTagAll(local string, data any) {
	if data == nil {
		return
	}
	local = local[0:2] //如 zh-CN -> zh
	val := reflect.ValueOf(data)

	// 处理指针类型
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		// 解引用指针
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		// 结构体：取地址后翻译
		if val.CanAddr() {
			TranslateI18nTag(local, val.Addr().Interface())
		}
	case reflect.Slice:
		// 切片：遍历每个元素进行翻译
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i)
			// 如果元素本身是指针，则直接传入翻译
			if item.Kind() == reflect.Ptr && !item.IsNil() && item.Elem().Kind() == reflect.Struct {
				TranslateI18nTag(local, item.Interface())
			} else if item.Kind() == reflect.Struct {
				// 如果元素是结构体，则取地址后翻译
				if item.CanAddr() {
					TranslateI18nTag(local, item.Addr().Interface())
				}
			}
			// 忽略非结构体元素，继续处理下一个
		}
	}
}

// SuccessData 通常成功数据处理
func SuccessData(c *gin.Context, data any) {
	Success(c, data, "Success")
}

func SuccessMsg(c *gin.Context, msg string) {

	c.AbortWithStatusJSON(http.StatusOK, &lv_dto.CommonRes{
		Code: 200,
		Data: nil,
		Msg:  msg,
	})
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
func SuccessPage(c *gin.Context, rows any, total any) {
	if rows != nil {
		local := c.GetHeader("Accept-Language")
		TranslateI18nTagAll(local, rows) //只翻译指针类型如：结构体指针 或 结构体切片指针
	}
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
