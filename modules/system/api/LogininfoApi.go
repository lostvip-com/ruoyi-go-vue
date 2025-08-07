package api

import (
	api2 "common/api"
	"common/util"
	"github.com/gin-gonic/gin"
	"system/model"
	"system/service"
	"system/vo"
)

type LogininfoApi struct {
	api2.BaseApi
}

// ListAjax 用户列表分页数据
func (w *LogininfoApi) ListAjax(c *gin.Context) {
	var req *vo.LoginInfoPageReq

	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}

	rows := make([]model.SysLoginInfo, 0)
	result, total, err := service.GetLoginServiceInstance().FindPage(req)

	if err == nil && len(*result) > 0 {
		rows = *result
	}
	util.BuildTable(c, total, rows).WriteJsonExit()
}

func (w *LogininfoApi) Remove(c *gin.Context) {
	var ids = c.Param("infoIds")
	//获取参数
	err := service.GetLoginServiceInstance().DeleteByIds(ids)
	if err == nil {
		util.SuccessData(c, nil)
	} else {
		util.Fail(c, "fail")
	}
}

// 清空记录
func (w *LogininfoApi) Clean(c *gin.Context) {
	err := service.GetLoginServiceInstance().DeleteRecordAll()

	if err == nil {
		util.SuccessData(c, nil)
	} else {
		util.Fail(c, err.Error())
	}
}

// 导出
func (w *LogininfoApi) Export(c *gin.Context) {
	var req *vo.LoginInfoPageReq

	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	url, err := service.GetLoginServiceInstance().Export(req)
	if err != nil {
		util.Fail(c, err.Error())
	} else {
		util.SuccessData(c, url)
	}
}

// Unlock 解锁账号
func (w *LogininfoApi) Unlock(c *gin.Context) {
	UserName := c.Query("UserName")
	if UserName == "" {
		util.ErrorResp(c).SetMsg("参数错误").Log("解锁账号", "UserName="+UserName).WriteJsonExit()
	} else {
		var loginService = service.GetLoginServiceInstance()
		loginService.RemovePasswordCounts(UserName)
		loginService.Unlock(UserName)
		util.SuccessData(c, nil)
	}

}
