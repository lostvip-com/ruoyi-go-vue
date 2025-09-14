// ==========================================================================
// LV自动生成数据库操作代码，无需手动修改，重新生成会自动覆盖.
// 生成日期：2024-07-19 17:09:35 +0800 CST
// 生成人：lv
// ==========================================================================
package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
	"time"
)

// IotDevice IotDevice
type IotDevice struct {
	Id        int `gorm:"type:bigint;size:20;primary_key;auto_increment;comment:主键;" json:"dictId"`
	DeptId    int `gorm:"type:bigint(20);comment:上级单位ID;" json:"deptId"`
	ProductId int `gorm:"type:bigint(20);comment:产品ID;" json:"productId"`

	GatewayId       int        `gorm:"type:bigint;comment:所属网关ID，只对网关子设备有效;" json:"gatewayId"`
	DriverId        int        `gorm:"type:bigint;comment:驱动实例ID;" json:"driverId"`
	NodeType        string     `gorm:"type:string;size:1;comment:节点类型3设备2网关;" json:"nodeType"`
	Name            string     `gorm:"type:varchar(32);comment:名称;" json:"name"`
	Sn              string     `gorm:"type:varchar(32);comment:序列号;" json:"sn"`
	DevNo           string     `gorm:"type:varchar(2);comment:从机地址;" json:"devNo"`
	Tags            string     `gorm:"type:string;size:32;comment:标签，多个;" json:"tags"`
	Status          string     `gorm:"type:int;default:0;comment:设备状态;0工程态，1 已接开通" json:"status"`
	Secret          string     `gorm:"type:varchar(32);comment:密钥;" json:"secret"`
	Description     string     `gorm:"type:text;comment:描述;" json:"description"`
	InstallLocation string     `gorm:"type:varchar(32);comment:安装地址;" json:"installLocation"`
	LastSyncTime    *time.Time `gorm:"type:datetime;comment:最后一次同步时间;" json:"lastSyncTime" time_format:"2006-01-02 15:04:05"`
	LastOnlineTime  *time.Time `gorm:"type:datetime;comment:最后一次在线时间;" json:"lastOnlineTime" time_format:"2006-01-02 15:04:05"`
	models.BaseModel
	// 非映射字段
	ProductName  string `gorm:"-" json:"productName"`
	NodeTypeName string `gorm:"-" json:"nodeTypeName"`
}

func (e *IotDevice) TableName() string {
	return "iot_device"
}

// 增
func (e *IotDevice) Save() error {
	return lv_db.GetOrmDefault().Save(e).Error
}

// 查
func (e *IotDevice) FindById(deviceId int) (*IotDevice, error) {
	err := lv_db.GetOrmDefault().Take(e, deviceId).Error
	return e, err
}

// 查第一条
func (e *IotDevice) FindOne() error {
	tb := lv_db.GetOrmDefault().Table(e.TableName())

	if e.Name != "" {
		tb = tb.Where("name like ?", "%"+e.Name+"%")
	}
	if e.ProductId != 0 {
		tb = tb.Where("product_id=?", e.ProductId)
	}

	err := tb.First(e).Error
	return err
}

// 改
func (e *IotDevice) Updates() error {
	return lv_db.GetOrmDefault().Table(e.TableName()).Updates(e).Error
}

// 删
func (e *IotDevice) Delete() error {
	return lv_db.GetOrmDefault().Table(e.TableName()).Delete(e).Error
}
