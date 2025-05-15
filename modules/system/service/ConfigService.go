package service

import (
	"common/common_vo"
	"common/global"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_cache"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	dao2 "system/dao"
	"system/model"
	"time"
)

type ConfigService struct {
}

var configService *ConfigService

func GetConfigServiceInstance() *ConfigService {
	if configService == nil {
		configService = &ConfigService{}
	}
	return configService
}

func (svc *ConfigService) GetValue(configKey string) string {
	unitedKey := global.SysConfigCacheKey + configKey
	result, err := lv_cache.GetCacheClient().Get(unitedKey)
	if err != nil {
		po := &model.SysConfig{ConfigKey: configKey}
		po, err := po.FindOne()
		lv_err.HasErrAndPanic(err)
		err = svc.SetCache(po)
		lv_err.HasErrAndPanic(err)
	}
	return result
}
func (svc *ConfigService) SetCache(po *model.SysConfig) error {
	unitedKey := global.SysConfigCacheKey + po.ConfigKey
	err := lv_cache.GetCacheClient().Set(unitedKey, po, 10*time.Minute)
	return err
}
func (svc *ConfigService) RemoveCache(po *model.SysConfig) error {
	unitedKey := global.SysConfigCacheKey + po.ConfigKey
	err := lv_cache.GetCacheClient().Del(unitedKey)
	return err
}

// FindConfigById 根据主键查询数据
func (svc *ConfigService) FindConfigById(id int64) (*model.SysConfig, error) {
	po := &model.SysConfig{ConfigId: id}
	po, err := po.FindOne()
	return po, err
}

// DeleteConfigById 根据主键删除数据
func (svc *ConfigService) DeleteConfigById(id int64) error {
	entity := &model.SysConfig{ConfigId: id}
	po, err := entity.FindOne()
	lv_err.HasErrAndPanic(err)
	err = po.Delete()
	return err
}

// DeleteByIds 批量删除数据记录
func (svc *ConfigService) DeleteByIds(ids string) {
	idArr := lv_conv.ToInt64Array(ids, ",")
	cfg := new(model.SysConfig)
	for _, id := range idArr {
		cfg, err := cfg.FindById(cast.ToInt64(id))
		lv_err.HasErrAndPanic(err)
		cfg.Delete()
	}
}

// 修改数据
func (svc *ConfigService) EditSave(req *common_vo.EditConfigReq, c *gin.Context) {
	po := &model.SysConfig{ConfigId: req.ConfigId}
	po, err := po.FindOne()
	lv_err.HasErrAndPanic(err)
	po.ConfigName = req.ConfigName
	po.ConfigKey = req.ConfigKey
	po.ConfigValue = req.ConfigValue
	po.Remark = req.Remark
	po.ConfigType = req.ConfigType
	po.UpdateTime = time.Now()
	po.UpdateBy = ""
	var userService UserService
	user := userService.GetProfile(c)

	if user == nil {
		po.UpdateBy = user.UserName
	}

	err = po.Update()
	lv_err.HasErrAndPanic(err)
	//保存到缓存
	_ = svc.SetCache(po)
}

func (svc *ConfigService) FindAll(params *common_vo.SelectConfigPageReq) ([]model.SysConfig, error) {
	var config dao2.ConfigDao
	return config.FindAll(params)
}

// 根据条件分页查询角色数据
func (svc *ConfigService) FindPage(params *common_vo.SelectConfigPageReq) (*[]map[string]string, int64, error) {
	var config dao2.ConfigDao
	return config.FindPage(params)
}

// 导出excel
func (svc *ConfigService) Export(param *common_vo.SelectConfigPageReq) (string, error) {
	//head := []string{"参数主键", "参数名称", "参数键名", "参数键值", "系统内置（Y是 N否）", "状态"}
	//col := []string{"config_id", "config_name", "config_key", "config_value", "config_type"}
	//var d dao2.ConfigDao
	//listMap, err := d.SelectExportList(param)
	//lv_err.HasErrAndPanic(err)
	return "", nil
}

func (svc *ConfigService) CountKey(key string) (int64, error) {
	var config dao2.ConfigDao
	return config.Count(key)
}
