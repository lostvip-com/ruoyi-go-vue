// ==========================================================================
// LV自动生成数据库操作代码，无需手动修改，重新生成会自动覆盖.
// 生成日期：2024-07-19 17:16:13 +0800 CST
// 生成人：lv
// ==========================================================================
package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
)

type IotProduct struct {
	Id         int    `gorm:"type:bigint;size:20;primary_key;auto_increment;comment:主键;" json:"dictId"`
	Name       string `json:"name" gorm:"type:string;size:255;comment:名字"`
	Key        string `json:"key" gorm:"type:string;size:255;comment:产品标识"`
	Platform   string `json:"platform" gorm:"type:string;size:255;comment:平台"`
	Protocol   string `json:"protocol" gorm:"type:string;size:255;comment:协议"`
	NodeType   string `json:"nodeType" gorm:"type:string;size:255;comment:节点类型"`
	NetType    string `json:"netType" gorm:"type:string;size:255;comment:网络类型"`
	DataFormat string `json:"dataFormat" gorm:"type:string;size:255;comment:数据类型"`
	Factory    string `json:"factory" gorm:"type:string;size:255;comment:工厂名称"`
	Remark     string `json:"remark" gorm:"size:64;comment:描述"`
	Status     string `json:"status" gorm:"type:string;size:255;comment:产品状态"`
	//Extra           MapStringString           `gorm:"type:string;size:255;comment:扩展字段"`
	//Properties      []Properties              `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 物模型的属性列表
	//Events          []Events                  `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 物模型的事件列表
	//Actions         []Actions                 `gorm:"foreignKey:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // 物模型的动作列表
	Manufacturer string `gorm:"type:varchar(32);comment:厂商;" json:"manufacturer"`
	//
	NodeTypeName string `gorm:"-" json:"nodeTypeName"`
	ProtocolName string `gorm:"-" json:"protocolName"`
	models.BaseModel
}

func (e *IotProduct) TableName() string {
	return "iot_product"
}

// 增
func (e *IotProduct) Save() error {
	return lv_db.GetOrmDefault().Save(e).Error
}

// 查
func (e *IotProduct) FindById(productId int) (*IotProduct, error) {
	err := lv_db.GetOrmDefault().Take(e, productId).Error
	return e, err
}

// 查第一条
func (e *IotProduct) FindOne() error {
	tb := lv_db.GetOrmDefault().Table(e.TableName())

	if e.Name != "" {
		tb = tb.Where("name like ?", "%"+e.Name+"%")
	}

	err := tb.First(e).Error
	return err
}

// 改
func (e *IotProduct) Updates() error {
	return lv_db.GetOrmDefault().Table(e.TableName()).Updates(e).Error
}

// 删
func (e *IotProduct) Delete() error {
	return lv_db.GetOrmDefault().Delete(e).Error
}
