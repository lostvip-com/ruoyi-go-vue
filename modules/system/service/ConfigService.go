package service

import (
	"common/common_vo"
	"common/global"
	"common/util"
	"github.com/lostvip-com/lv_framework/lv_cache"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/namedsql"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_reflect"
	"github.com/spf13/cast"
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
func (svc *ConfigService) FindConfigById(id int) (*model.SysConfig, error) {
	po := &model.SysConfig{ConfigId: id}
	po, err := po.FindOne()
	return po, err
}

// DeleteConfigById 根据主键删除数据
func (svc *ConfigService) DeleteConfigById(id int) error {
	entity := &model.SysConfig{ConfigId: id}
	po, err := entity.FindOne()
	lv_err.HasErrAndPanic(err)
	err = po.Delete()
	if err != nil {
		_ = svc.RemoveCache(po)
	}
	return err
}

// DeleteByIds 批量删除数据记录
func (svc *ConfigService) DeleteByIds(ids string) {
	idArr := util.ToIntArray(ids, ",")
	cfg := new(model.SysConfig)
	for _, id := range idArr {
		cfg, err := cfg.FindById(cast.ToInt(id))
		lv_err.HasErrAndPanic(err)
		err = cfg.Delete()
		if err != nil {
			_ = svc.RemoveCache(cfg)
		}
	}
}

// 修改数据
func (svc *ConfigService) EditSave(param *model.SysConfig) error {
	po := new(model.SysConfig)
	po, err := po.FindById(param.ConfigId)
	lv_err.HasErrAndPanic(err)
	err = lv_reflect.CopyProp(param, po, true)
	lv_err.HasErrAndPanic(err)
	//保存到缓存
	err = po.Update()
	lv_err.HasErrAndPanic(err)
	err = svc.SetCache(po)
	return err
}

// 批量删除
func (svc *ConfigService) Count(configKey string) (int, error) {
	var total int64
	err := lv_db.GetOrmDefault().Table("sys_config").Where("config_key=?", configKey).Count(&total).Error
	return int(total), err
}

// 批量删除
func (svc *ConfigService) DeleteBatch(ids ...int) error {
	err := lv_db.GetOrmDefault().Delete(&model.SysConfig{}, ids).Error
	return err
}

// FindPage 根据条件分页查询用户列表
func (svc *ConfigService) FindPage(param *common_vo.ConfigPageReq) (*[]map[string]any, int, error) {
	db := lv_db.GetOrmDefault()
	sqlParams, sql := svc.GetSql(param)
	countSql := "select count(*) from (" + sql + ") t "
	total, err := namedsql.Count(db, countSql, sqlParams)
	lv_err.HasErrAndPanic(err)
	limitSql := sql + " order by u.config_id desc "
	limitSql += "  limit " + cast.ToString(param.GetStartNum()) + "," + cast.ToString(param.GetPageSize())
	result, err := namedsql.ListMap(db, limitSql, sqlParams, true)
	lv_err.HasErrAndPanic(err)
	return result, int(total), err
}

// 导出excel
func (svc *ConfigService) SelectAllList(param *common_vo.ConfigPageReq) (*[]map[string]any, error) {
	db := lv_db.GetOrmDefault()
	sqlParams, sql := svc.GetSql(param)
	limitSql := sql + " order by u.user_id desc "
	result, err := namedsql.ListMapAny(db, limitSql, &sqlParams, false)
	return result, err
}
func (svc *ConfigService) GetSql(param *common_vo.ConfigPageReq) (map[string]interface{}, string) {
	sqlParams := make(map[string]interface{})
	sql := `
           select * from sys_config u where 1=1
          `
	if param != nil {
		if param.ConfigName != "" {
			sql += " and  u.config_name like @ConfigName "
			sqlParams["ConfigName"] = "%" + param.ConfigName + "%"
		}
		if param.ConfigKey != "" {
			sql += " and  u.config_key like @ConfigKey "
			sqlParams["ConfigKey"] = "%" + param.ConfigKey + "%"
		}
		if param.ConfigType != "" {
			sql += " and  u.config_type like @ConfigType "
			sqlParams["ConfigType"] = "%" + param.ConfigType + "%"
		}

		if param.BeginTime != "" {
			sql += " and  u.create_time >= @BeginTime "
			sqlParams["BeginTime"] = param.BeginTime
		}
		if param.EndTime != "" {
			sql += " and  u.create_time <= @EndTime "
			sqlParams["EndTime"] = param.EndTime
		}

	}
	return sqlParams, sql
}
