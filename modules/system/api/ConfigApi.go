package api

import (
	"common/common_vo"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	dao2 "system/dao"
	"system/service"
)

type ConfigApi struct {
}

// 修改页面保存
func (w *ConfigApi) GetConfigKey(c *gin.Context) {
	configKey := c.Param("configKey")
	//获取参数
	var configService = service.GetConfigServiceInstance()
	value := configService.GetValueFromCache(configKey)
	util.Success(c, value)
}

// 列表分页数据
func (w *ConfigApi) ListAjax(c *gin.Context) {
	req := new(common_vo.SelectConfigPageReq)
	var configService service.ConfigService
	//获取参数
	if err := c.ShouldBind(req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("参数管理", req).WriteJsonExit()
		return
	}
	result, total, err := configService.FindPage(req)
	lv_err.HasErrAndPanic(err)
	util.BuildTable(c, total, result).WriteJsonExit()
}

// 新增页面保存
func (w *ConfigApi) AddSave(c *gin.Context) {
	req := new(common_vo.AddConfigReq)
	//获取参数
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var config dao2.ConfigDao
	count, err := config.Count(req.ConfigKey)
	lv_err.HasErrAndPanic(err)
	if count > 0 {
		util.Fail(c, "此编号已经存在！请更换编号")
	} else {
		util.Success(c, "")
	}
}

// 修改页面保存
func (w *ConfigApi) EditSave(c *gin.Context) {
	req := new(common_vo.EditConfigReq)
	//获取参数
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var configService service.ConfigService
	configService.EditSave(req, c)
	util.Success(c, "")
}

// 删除数据
func (w *ConfigApi) Remove(c *gin.Context) {
	var configIds = c.Param("configIds")
	var configService service.ConfigService
	configService.DeleteRecordByIds(configIds)
	util.Success(c, "")
}

// 导出
func (w *ConfigApi) Export(c *gin.Context) {
	req := new(common_vo.SelectConfigPageReq)
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var configService service.ConfigService
	url, err := configService.Export(req)
	lv_err.HasErrAndPanic(err)

	util.Success(c, url)
}
