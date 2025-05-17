package api

import (
	"common/myconf"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_log"
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
	util.Success(c, nil)
}

// 下载 public/upload 文件头像之类
func (w *IndexApi) Download(c *gin.Context) {
	fileName := c.Query("fileName")
	if fileName == "" {
		util.Fail(c, "参数错误")
		return
	}
	curDir, err := os.Getwd()
	filepath := curDir + "/static/upload/" + fileName
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		util.Fail(c, "参数错误")
		return
	}
	b, _ := io.ReadAll(file)
	c.Writer.Header().Add("Content-Disposition", "attachment")
	c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetmxls.sheet")
	c.Writer.Write(b)
	c.Abort()
}

// CaptchaImage 图形验证码生成逻辑,使用其它服务的公共接口生成，不再单独维护
func (w *IndexApi) CaptchaImage(c *gin.Context) {
	//传参数
	clientId := c.Query("uuid")
	//返回值
	params := url.Values{}
	params.Set("uuid", clientId)

	captcha := lv_dto.CaptchaRes{Code: 200, Uuid: clientId}
	captcha.CaptchaEnabled = myconf.GetConfigInstance().GetBool("sys.account.captchaEnabled")
	url_captcha := myconf.GetConfigInstance().GetValueStr("sys.url.captcha")
	if !captcha.CaptchaEnabled || url_captcha == "" {
		lv_log.Warn("未配置验证码地址url-captcha 或 未启用开关")
	} else {
		json, err := lv_net.Get(url_captcha)
		lv_err.HasErrAndPanic(err)
		dataMap := lv_conv.ToMap(json)
		captcha.Img = dataMap["img"]
	}
	c.JSON(http.StatusOK, captcha)
}
