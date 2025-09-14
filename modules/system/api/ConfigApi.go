package api

import (
	api2 "common/api"
	"common/common_vo"
	"common/global"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_cache/lv_redis"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	"system/model"
	"system/service"
)

type ConfigApi struct {
	api2.BaseApi
}

func (w *ConfigApi) GetConfigInfo(c *gin.Context) {
	configId := c.Param("configId")
	var configService = service.GetConfigServiceInstance()
	value, err := configService.FindConfigById(cast.ToInt(configId))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessData(c, value)
}

func (w *ConfigApi) GetConfigValueByKey(c *gin.Context) {
	configKey := c.Param("configKey")
	var configService = service.GetConfigServiceInstance()
	value := configService.GetValue(configKey)
	util.SuccessData(c, value)
}

func (w *ConfigApi) ListAjax(c *gin.Context) {
	req := new(common_vo.ConfigPageReq)
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
	var config = service.GetConfigServiceInstance()
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
	util.SuccessData(c, "")
}

func (w *ConfigApi) EditSave(c *gin.Context) {
	req := new(model.SysConfig)
	err := c.ShouldBind(req)
	lv_err.HasErrAndPanic(err)
	w.FillInUpdate(c, &req.BaseModel)
	service.GetConfigServiceInstance().EditSave(req)
	util.SuccessData(c, "")
}

func (w *ConfigApi) Remove(c *gin.Context) {
	var configIds = c.Param("configIds")
	var configService service.ConfigService
	configService.DeleteByIds(configIds)
	util.SuccessData(c, "")
}

func (w *ConfigApi) Export(c *gin.Context) {
	var req common_vo.ConfigPageReq
	err := c.ShouldBind(&req)
	lv_err.HasErrAndPanic(err)
	listMap, _, err := service.GetConfigServiceInstance().FindPage(&req)
	lv_err.HasErrAndPanic(err)
	headerMap := []map[string]string{
		{"key": "ConfigName", "title": "参数名称", "width": "10"},
		{"key": "configKey", "title": "参数键名", "width": "15"},
		{"key": "configValue", "title": "参数键值", "width": "10"},
		{"key": "configType", "title": "系统内置", "width": "10"},
		{"key": "remark", "title": "备注信息", "width": "10"},
		{"key": "updateTime", "title": "更新时间", "width": "10"},
		{"key": "createTime", "title": "创建时间", "width": "10"},
	}
	ex := util.NewMyExcel()
	ex.ExportToWeb(c, headerMap, *listMap)
}

func (w *ConfigApi) RefreshCacheConfig(c *gin.Context) {
	redisCache := lv_redis.GetInstance(0)
	keys, _, err := redisCache.Scan(0, global.SysConfigCacheKey+"*", global.ScanCountMax)
	if err != nil {
		util.Fail(c, err.Error())
	}
	for _, key := range keys {
		err = redisCache.Del(key)
		lv_err.HasErrAndPanic(err)
	}
	util.SuccessData(c, nil)
}
