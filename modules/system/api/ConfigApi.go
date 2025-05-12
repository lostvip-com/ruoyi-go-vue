package api

import (
	"common/common_vo"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	dao2 "system/dao"
	"system/model"
	"system/service"
)

type ConfigApi struct {
	BaseApi
}

func (w *ConfigApi) GetConfigInfo(c *gin.Context) {
	configId := c.Param("configId")
	var configService = service.GetConfigServiceInstance()
	value, err := configService.FindConfigById(cast.ToInt64(configId))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, value)
}

func (w *ConfigApi) GetConfigKey(c *gin.Context) {
	configKey := c.Param("configKey")
	var configService = service.GetConfigServiceInstance()
	value := configService.GetValueFromCache(configKey)
	util.Success(c, value)
}

func (w *ConfigApi) ListAjax(c *gin.Context) {
	req := new(common_vo.SelectConfigPageReq)
	var configService service.ConfigService
	if err := c.ShouldBind(req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	rows, total, err := configService.FindPage(req)
	lv_err.HasErrAndPanic(err)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessPage(c, rows, total)

}

func (w *ConfigApi) AddSave(c *gin.Context) {
	entity := new(model.SysConfig)
	//获取参数
	err := c.ShouldBind(entity)
	lv_err.HasErrAndPanic(err)
	var config dao2.ConfigDao
	count, err := config.Count(entity.ConfigKey)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	if count > 0 {
		util.Fail(c, "此编号已经存在！请更换编号")
		return
	}
	entity.CreateBy = w.GetCurrUser(c).UserName
	err = entity.Save()
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, "")
}

func (w *ConfigApi) EditSave(c *gin.Context) {
	req := new(common_vo.EditConfigReq)
	//获取参数
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	var configService = service.GetConfigServiceInstance()
	configService.EditSave(req, c)
	util.Success(c, "")
}

// 删除数据
func (w *ConfigApi) Remove(c *gin.Context) {
	var configIds = c.Param("configIds")
	var configService service.ConfigService
	configService.DeleteByIds(configIds)
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
