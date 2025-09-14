// ==========================================================================
// LV自动生成数据库操作代码，无需手动修改，重新生成会自动覆盖.
// 生成日期: 2025-03-20 01:34:11 &#43;0000 UTC
// 生成人: lv
// ==========================================================================
package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/namedsql"
	"time"
)

// IotDataProp IotDataProp
type IotDataProp struct {
	Id       int       `gorm:"type:bigint;size:20;primary_key;auto_increment;comment:主键;" json:"id"`
	DeviceId int       `gorm:"type:bigint;comment:设备ID;" json:"deviceId"`
	Code     string    `gorm:"type:varchar(32);comment:属性名;" json:"code"`
	Value    string    `gorm:"type:varchar(32);comment:采集值;" json:"value"`
	Name     string    `gorm:"type:varchar(32);comment:属性名;" json:"name"`
	DevTime  time.Time `gorm:"type:datetime;comment:采集时间;" json:"devTime" time_format:"2006-01-02 15:04:05"`
	DelFlag  string    `gorm:"type:char(1);comment:删除标记;" json:"delFlag"`
	Tags     string    `gorm:"type:varchar(64);comment:TAG;" json:"tags"`
	models.BaseModel
}

func (e *IotDataProp) TableName() string {
	return "iot_data_prop"
}

func (e *IotDataProp) Save() error {
	return lv_db.GetOrmDefault().Save(e).Error
}

func (e *IotDataProp) FindById(id int) (*IotDataProp, error) {
	err := lv_db.GetOrmDefault().Take(e, id).Error
	return e, err
}

func (e *IotDataProp) FindOne() (*IotDataProp, error) {
	tb := lv_db.GetOrmDefault().Table(e.TableName())

	if e.Name != "" {
		tb = tb.Where("name=?", e.Name)
	}
	//if e.TenantId != "" {
	//	tb = tb.Where("tenant_id=?", e.TenantId)
	//}
	if e.Code != "" {
		tb = tb.Where("code=?", e.Code)
	}
	if e.DeviceId != 0 {
		tb = tb.Where("device_id=?", e.DeviceId)
	}
	err := tb.First(e).Error
	return e, err
}

func (e *IotDataProp) Updates() error {
	return lv_db.GetOrmDefault().Table(e.TableName()).Updates(e).Error
}

func (e *IotDataProp) Delete() error {
	return lv_db.GetOrmDefault().Delete(e).Error
}

func (e *IotDataProp) Count() (int64, error) {
	sql := " select count(*) from iot_data_his where del_flag = 0 "

	if e.Name != "" {
		sql += " and name = @Name "
	}
	//if e.TenantId != "" {
	//	sql += " and tenant_id = @TenantId "
	//}
	if e.Code != "" {
		sql += " and code = @Code "
	}
	if e.DeviceId != 0 {
		sql += " and device_id = @DeviceId "
	}

	return namedsql.Count(lv_db.GetOrmDefault(), sql, e)
}
