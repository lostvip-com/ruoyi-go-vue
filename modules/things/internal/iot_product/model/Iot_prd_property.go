// ==========================================================================
// LV自动生成数据库操作代码，无需手动修改，重新生成会自动覆盖.
// 生成日期：2024-07-30 21:59:38 +0800 CST
// 生成人：lv
// ==========================================================================
package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
)

// IotPrdProperty 属性
type IotPrdProperty struct {
	Id               int    `gorm:"type:bigint;size:20;primary_key;auto_increment;comment:主键;" json:"dictId"`
	ProductId        int    `gorm:"type:bigint(20);comment:产品ID;" json:"productId"`
	Name             string `gorm:"type:varchar(32);comment:名字;" json:"name"`
	Code             string `gorm:"type:varchar(32);comment:标识符;" json:"code"`
	Tag              string `gorm:"type:varchar(50);comment:标签;" json:"tag"`
	AccessMode       string `gorm:"type:char(2);comment:读写模型RO,RW;" json:"accessMode"`
	DataType         string `gorm:"type:varchar(50);comment:数据类型;" json:"dataType"`
	DataRate         int    `gorm:"type:varchar(50);comment:倍率;" json:"dataRate"`
	Unit             string `gorm:"type:varchar(50);comment:单位;" json:"unit"`
	Remark           string `gorm:"type:longtext;comment:描述;" json:"remark"`
	models.BaseModel `gorm:"embedded"`
}

func (e *IotPrdProperty) TableName() string {
	return "iot_prd_property"
}

// Save 增
func (e *IotPrdProperty) Save() error {
	return lv_db.GetOrmDefault().Save(e).Error
}

// FindById 查
func (e *IotPrdProperty) FindById(deviceId int) (*IotPrdProperty, error) {
	err := lv_db.GetOrmDefault().Take(e, deviceId).Error
	return e, err
}

// FindOne 查第一条
func (e *IotPrdProperty) FindOne() (*IotPrdProperty, error) {
	tb := lv_db.GetOrmDefault().Table(e.TableName())

	if e.ProductId != 0 {
		tb = tb.Where("product_id=?", e.ProductId)
	}
	if e.Name != "" {
		tb = tb.Where("name=?", e.Name)
	}
	if e.Code != "" {
		tb = tb.Where("code=?", e.Code)
	}
	if e.Tag != "" {
		tb = tb.Where("tag=?", e.Tag)
	}
	err := tb.First(e).Error
	return e, err
}

// Updates 改
func (e *IotPrdProperty) Updates() error {
	return lv_db.GetOrmDefault().Table(e.TableName()).Updates(e).Error
}

// Delete 删
func (e *IotPrdProperty) Delete() error {
	return lv_db.GetOrmDefault().Delete(e).Error
}
