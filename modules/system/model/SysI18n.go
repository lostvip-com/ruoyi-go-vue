// ==========================================================================
// LV自动生成数据库操作代码，无需手动修改，重新生成会自动覆盖.
// 生成日期: 2025-08-11 07:41:35 &#43;0000 UTC
// 生成人: dpc
// ==========================================================================
package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
)

// SysI18n 国际化翻译
type SysI18n struct {
	Id         int    `gorm:"type:bigint(20);primary_key;auto_increment;;" json:"id"`
	Locale     string `gorm:"type:varchar(64);comment:本地标识;" json:"locale"`
	LocaleKey  string `gorm:"type:varchar(64);comment:国际化key;" json:"localeKey"`
	LocaleName string `gorm:"type:varchar(64);comment:国际化名称;" json:"localeName"`
	Sort       int    `gorm:"type:int(11);comment:字典排序;" json:"sort"`
	Remark     string `gorm:"type:varchar(100);comment:备注;" json:"remark"`
	models.BaseModel
}

func (e *SysI18n) TableName() string {
	return "sys_i18n"
}

func (e *SysI18n) Save() error {
	return lv_db.GetOrmDefault().Save(e).Error
}

func (e *SysI18n) FindById(id int) (*SysI18n, error) {
	err := lv_db.GetOrmDefault().Take(e, id).Error
	return e, err
}

func (e *SysI18n) FindOne(locale, localeKey string) (*SysI18n, error) {
	tb := lv_db.GetOrmDefault().Table(e.TableName())

	if locale != "" {
		tb = tb.Where("locale=?", locale)
	}
	if localeKey != "" {
		tb = tb.Where("locale_key=?", localeKey)
	}
	err := tb.First(e).Error
	return e, err
}

func (e *SysI18n) Updates() error {
	return lv_db.GetOrmDefault().Table(e.TableName()).Updates(e).Error
}

func (e *SysI18n) Delete() error {
	return lv_db.GetOrmDefault().Delete(e).Error
}
