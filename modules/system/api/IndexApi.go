package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_global"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_net"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"io"
	"net/http"
	"net/url"
	"os"
	"system/service"
)

type IndexApi struct{}

// 注销
func (w *IndexApi) Index(c *gin.Context) {
	var user service.SessionService
	tokenStr := lv_net.GetParam(c, "token")
	err := user.SignOut(tokenStr)
	if err != nil {
		return
	}
	util.Success(c, nil, "welcome!")
}

// 下载 public/upload 文件头像之类
func (w *IndexApi) Download(c *gin.Context) {
	fileName := c.Query("fileName")
	//delete := c.Query("delete")
	if fileName == "" {
		util.BuildTpl(c, lv_dto.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "参数错误",
		})
		return
	}
	curDir, err := os.Getwd()
	filepath := curDir + "/static/upload/" + fileName
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		util.BuildTpl(c, lv_dto.ERROR_PAGE).WriteTpl(gin.H{
			"desc": "参数错误",
		})
		return
	}
	b, _ := io.ReadAll(file)
	c.Writer.Header().Add("Content-Disposition", "attachment")
	c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Writer.Write(b)
	c.Abort()
}

// CaptchaImage 图形验证码生成逻辑,使用其它服务的公共接口生成，不再单独维护
func (w *IndexApi) CaptchaImage(c *gin.Context) {
	//传参数
	clientId := c.PostForm("uuid")
	//返回值
	params := url.Values{}
	params.Set("uuid", clientId)
	url_captcha := lv_global.Config().GetValueStr("url_captcha")
	if url_captcha == "" {
		util.Err(c, "请配置验证码地址url_captcha")
		return
	}
	json, err := lv_net.PostForm(url_captcha, params)
	lv_err.HasErrAndPanic(err)
	dataMap := lv_conv.ToMap(json)
	base64stringC := dataMap["data"]
	c.JSON(http.StatusOK, lv_dto.CaptchaRes{
		Code:           200,
		CaptchaEnabled: true,
		Uuid:           clientId,
		Img:            base64stringC,
		Msg:            "操作成功",
	})
}
