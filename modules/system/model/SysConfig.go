package model

import (
	"common/models"
	"errors"
	"github.com/lostvip-com/lv_framework/lv_db"
)

// SysConfig 参数配置
type SysConfig struct {
	ConfigId    int64  `gorm:"type:int(11);primary_key;auto_increment;参数主键;" json:"configId" form:"configId"`
	ConfigName  string `gorm:"type:varchar(100);comment:参数名称;" json:"configName" form:"configName"`
	ConfigKey   string `gorm:"type:varchar(100);comment:参数键名;" json:"configKey" form:"configKey"`
	ConfigValue string `gorm:"type:varchar(500);comment:参数键值;" json:"configValue" form:"configValue"`
	ConfigType  string `gorm:"type:char(1);comment:系统内置（Y是 N否）;" json:"configType" form:"configType"`
	Remark      string `gorm:"type:varchar(500);comment:备注;" json:"remark" form:"remark"`
	models.BaseModel
}

func (e *SysConfig) TableName() string {
	return "sys_config"
}

// 增
func (e *SysConfig) Save() error {
	return lv_db.GetOrmDefault().Save(e).Error
}

// 查
func (e *SysConfig) FindById(id int64) (*SysConfig, error) {
	err := lv_db.GetOrmDefault().Table(e.TableName()).First(e, id).Error
	return e, err
}

// 查第一条
func (e *SysConfig) FindOne() (*SysConfig, error) {
	tb := lv_db.GetOrmDefault().Table(e.TableName())
	if e.ConfigId != 0 {
		tb = tb.Where("config_id=?", e.ConfigId)
	}
	if e.ConfigKey != "" {
		tb = tb.Where("config_key=?", e.ConfigKey)
	}

	err := tb.First(e).Error
	if err != nil {
		return nil, errors.New("未找到系统配置信息Key:" + e.ConfigKey)
	}
	return e, err
}

// 改
func (e *SysConfig) Update() error {
	return lv_db.GetOrmDefault().Table(e.TableName()).Updates(e).Error
}

// 删
func (e *SysConfig) Delete() error {
	return lv_db.GetOrmDefault().Table(e.TableName()).Delete(e).Error
}
