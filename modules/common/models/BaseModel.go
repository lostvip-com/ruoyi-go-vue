package models

import (
	"github.com/lostvip-com/lv_framework/utils/lv_time"
	"gorm.io/gorm"
	"time"
)

// BaseModel 共享属性
type BaseModel struct {
	CreateTime time.Time `gorm:"type:datetime;comment:创建日期;autoCreateTime" time_format:"2006-01-02 15:04:05" json:"createTime,omitzero"`
	UpdateTime time.Time `gorm:"type:datetime;comment:更新日期;autoUpdateTime" time_format:"2006-01-02 15:04:05" json:"updateTime,omitzero"`
	UpdateBy   string    `gorm:"type:string;size:32;comment:更新者;" json:"updateBy"`
	CreateBy   string    `gorm:"type:string;size:32;comment:创建者;" json:"createBy"`
	TenantId   int64     `gorm:"type:string;size:32;comment:租户id;" json:"tenantId" form:"tenantId"`
}

// BeforeCreate 实现钩子
func (u *BaseModel) BeforeCreate(db *gorm.DB) error {
	u.CreateTime = lv_time.GetCurrentTime() // 设置创建时的更新时间
	u.UpdateTime = u.CreateTime             // 设置创建时的更新时间
	return nil
}

// BeforeUpdate 实现 BeforeUpdate 钩子
func (u *BaseModel) BeforeUpdate(db *gorm.DB) error {
	u.CreateTime = lv_time.GetCurrentTime() // 设置更新时的更新时间
	return nil
}
