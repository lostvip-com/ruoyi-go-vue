// ==========================================================================
// LV自动生成数据库操作代码，无需手动修改，重新生成会自动覆盖.
// 生成日期: 2025-07-24 02:41:35 &#43;0000 UTC
// 生成人: lv
// ==========================================================================
package model

import (
    "common/models"
    "github.com/lostvip-com/lv_framework/lv_db"
    "github.com/lostvip-com/lv_framework/lv_db/namedsql"
    "time"
)

// IotProduct IotProduct
type IotProduct struct {
    Id  int64  `gorm:"type:bigint(20);primary_key;auto_increment;主键;" json:"id"`
    Key  string  `gorm:"type:varchar(32);comment:产品编码,对应可监控类型ID;" json:"key"`
    Name  string  `gorm:"type:varchar(32);comment:名字;" json:"name"`
    CloudProductId  string  `gorm:"type:varchar(32);comment:云产品ID;" json:"cloudProductId"`
    CloudInstanceId  string  `gorm:"type:varchar(32);comment:云实例ID;" json:"cloudInstanceId"`
    Platform  string  `gorm:"type:varchar(32);comment:平台;" json:"platform"`
    Protocol  string  `gorm:"type:varchar(32);comment:协议;" json:"protocol"`
    NodeType  string  `gorm:"type:varchar(32);comment:节点类型;" json:"nodeType"`
    NetType  string  `gorm:"type:varchar(32);comment:网络类型;" json:"netType"`
    DataFormat  string  `gorm:"type:varchar(32);comment:数据类型;" json:"dataFormat"`
    LastSyncTime  int64  `gorm:"type:bigint(20);comment:最后一次同步时间;" json:"lastSyncTime"`
    Factory  string  `gorm:"type:varchar(32);comment:工厂名称;" json:"factory"`
    Description  string  `gorm:"type:text;comment:描述;" json:"description"`
    Status  string  `gorm:"type:varchar(32);comment:产品状态;" json:"status"`
    Extra  string  `gorm:"type:varchar(32);comment:扩展字段;" json:"extra"`
    DelFlag  int  `gorm:"type:tinyint(1);comment:删除标记;" json:"delFlag"`
    CreateTime time.Time `gorm:"type:datetime;comment:创建日期;" json:"createTime" time_format:"2006-01-02 15:04:05"`
    UpdateTime time.Time `gorm:"type:datetime;comment:更新日期;" json:"updateTime" time_format:"2006-01-02 15:04:05"`
    UpdateBy  string  `gorm:"type:varchar(64);comment:更新者;" json:"updateBy"`
    CreateBy  string  `gorm:"type:varchar(64);comment:创建者;" json:"createBy"`
    Manufacturer  string  `gorm:"type:varchar(32);comment:生产厂商;" json:"manufacturer"`
    TenantId  int64  `gorm:"type:bigint(20);comment:租户id;" json:"tenantId"`
    models.BaseModel
}

func (e *IotProduct) TableName() string {
    return "iot_product"
}

func (e *IotProduct) Save() error {
    return lv_db.GetMasterGorm().Save(e).Error
}

func (e *IotProduct) FindById(id int64) (*IotProduct,error) {
    err := lv_db.GetMasterGorm().Take(e,id).Error
    return e,err
}

func (e *IotProduct) FindOne() (*IotProduct,error) {
    tb := lv_db.GetMasterGorm().Table(e.TableName())

    if e.Key != "" {
         tb = tb.Where("key=?", e.Key)
    }
    if e.Name != "" {
         tb = tb.Where("name=?", e.Name)
    }
    if e.CloudProductId != "" {
         tb = tb.Where("cloud_product_id=?", e.CloudProductId)
    }
    if e.CloudInstanceId != "" {
         tb = tb.Where("cloud_instance_id=?", e.CloudInstanceId)
    }
    err := tb.First(e).Error
    return e,err
}

func (e *IotProduct) Updates() error {
    return lv_db.GetMasterGorm().Table(e.TableName()).Updates(e).Error
}

func (e *IotProduct) Delete() error {
    return lv_db.GetMasterGorm().Delete(e).Error
}

func (e *IotProduct) Count() (int64, error) {
    sql := " select count(*) from iot_product where del_flag = 0 "
  
     if e.Key != "" {
        sql += " and key = @Key "
     }
     if e.Name != "" {
        sql += " and name = @Name "
     }
     if e.CloudProductId != "" {
        sql += " and cloud_product_id = @CloudProductId "
     }
     if e.CloudInstanceId != "" {
        sql += " and cloud_instance_id = @CloudInstanceId "
     }

    return namedsql.Count(lv_db.GetMasterGorm(), sql, e)
}