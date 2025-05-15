package service

import (
	"common/common_vo"
	"common/global"
	"github.com/lostvip-com/lv_framework/lv_cache"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/namedsql"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
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
	if err != nil {
		_ = svc.RemoveCache(po)
	}
	return err
}

// DeleteByIds 批量删除数据记录
func (svc *ConfigService) DeleteByIds(ids string) {
	idArr := lv_conv.ToInt64Array(ids, ",")
	cfg := new(model.SysConfig)
	for _, id := range idArr {
		cfg, err := cfg.FindById(cast.ToInt64(id))
		lv_err.HasErrAndPanic(err)
		err = cfg.Delete()
		if err != nil {
			_ = svc.RemoveCache(cfg)
		}
	}
}

// 修改数据
func (svc *ConfigService) EditSave(param *model.SysConfig) {
	po := new(model.SysConfig)
	po, err := po.FindById(param.ConfigId)
	lv_err.HasErrAndPanic(err)
	err = lv_reflect.CopyProp(param, po, true)
	lv_err.HasErrAndPanic(err)
	//保存到缓存
	err = po.Update()
	lv_err.HasErrAndPanic(err)
	_ = svc.SetCache(po)
}

// 批量删除
func (d *ConfigService) Count(configKey string) (total int64, err error) {
	err = lv_db.GetMasterGorm().Table("sys_config").Where("config_key=?", configKey).Count(&total).Error
	return total, err
}

// 批量删除
func (d *ConfigService) DeleteBatch(ids ...int64) error {
	err := lv_db.GetMasterGorm().Delete(&model.SysConfig{}, ids).Error
	return err
}

// 根据条件分页查询用户列表
func (d ConfigService) FindPage(param *common_vo.SelectConfigPageReq) (*[]map[string]any, int64, error) {
	db := lv_db.GetMasterGorm()
	sqlParams, sql := d.GetSql(param)
	countSql := "select count(*) from (" + sql + ") t "
	total, err := namedsql.Count(db, countSql, sqlParams)
	lv_err.HasErrAndPanic(err)
	limitSql := sql + " order by u.config_id desc "
	limitSql += "  limit " + cast.ToString(param.GetStartNum()) + "," + cast.ToString(param.GetPageSize())
	result, err := namedsql.ListMap(db, limitSql, sqlParams, true)
	lv_err.HasErrAndPanic(err)
	return result, total, err
}

func (d ConfigService) GetSql(param *common_vo.SelectConfigPageReq) (map[string]interface{}, string) {
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

// 导出excel
func (d ConfigService) SelectExportList(param *common_vo.SelectConfigPageReq) (*[]map[string]any, error) {
	db := lv_db.GetMasterGorm()
	sqlParams, sql := d.GetSql(param)
	limitSql := sql + " order by u.user_id desc "
	result, err := namedsql.ListMapAny(db, limitSql, &sqlParams, false)
	return result, err
}

// 获取所有数据
func (d *ConfigService) FindAll(param *common_vo.SelectConfigPageReq) ([]model.SysConfig, error) {
	db := lv_db.GetMasterGorm()
	tb := db.Table("sys_config t")

	if param != nil {
		if param.ConfigName != "" {
			tb.Where("t.config_name like ?", "%"+param.ConfigName+"%")
		}

		if param.ConfigType != "" {
			tb.Where("t.status = ", param.ConfigType)
		}

		if param.ConfigKey != "" {
			tb.Where("t.config_key like ?", "%"+param.ConfigKey+"%")
		}

		if param.BeginTime != "" {
			tb.Where("t.create_time >= ? ", param.BeginTime)
		}

		if param.EndTime != "" {
			tb.Where("t.create_time <= ? ", param.EndTime)
		}
	}
	var result []model.SysConfig
	err := tb.Find(&result).Error
	return result, err
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
