// ==========================================================================
// 国际化表
// 生成日期：2024-08-04 10:51:57 +0800 CST
// 生成人：lv
// ==========================================================================
package models

import (
	"github.com/lostvip-com/lv_framework/lv_db"
)

// SysI18n 通用业务数据国际化表表
type SysI18n struct {
	Id         int    `gorm:"type:bigint;size:20;primary_key;auto_increment;id;" json:"id"`
	Locale     string `gorm:"type:varchar(64);index:idx_local;defalut:zh,comment:本地标识" json:"locale"`
	LocaleKey  string `gorm:"type:varchar(64);index:idx_locale_key;comment:国际化key" json:"localeKey"`
	LocaleName string `gorm:"type:varchar(64);comment:国际化名称;" json:"localeName"`
	Sort       int    `gorm:"type:int(11);defalut:100;not null;comment:字典排序;" json:"sort"`
	Remark     string `gorm:"type:varchar(100);comment:备注;" json:"remark"`
	BaseModel
}

func (e *SysI18n) TableName() string {
	return "sys_i18n"
}

// 增
func (e *SysI18n) Save() error {
	return lv_db.GetOrmDefault().Save(e).Error
}

// 查
func (e *SysI18n) FindById(id int) (*SysI18n, error) {
	err := lv_db.GetOrmDefault().Take(e, id).Error
	return e, err
}

// 查第一条
func (e *SysI18n) FindOne(locale, localeKey string) (*SysI18n, error) {
	tb := lv_db.GetOrmDefault().Table(e.TableName())
	tb = tb.Where("locale=? and locale_key=? ", locale, localeKey)
	err := tb.First(e).Error
	return e, err
}

// 改
func (e *SysI18n) Updates() error {
	return lv_db.GetOrmDefault().Table(e.TableName()).Updates(e).Error
}

// 删
func (e *SysI18n) Delete() error {
	return lv_db.GetOrmDefault().Delete(e).Error
}
