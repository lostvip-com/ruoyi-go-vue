package models

import (
	"time"
)

type LocalTime time.Time

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(`"2006-01-02 15:04:05"`))
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, "2006-01-02 15:04:05")
	b = append(b, '"')
	return b, nil
}

// BaseModel 共享属性
type BaseModel struct {
	CreateTime LocalTime `gorm:"type:datetime;comment:创建日期;autoCreateTime;<-:create;" json:"createTime,omitzero"`
	UpdateTime LocalTime `gorm:"type:datetime;comment:更新日期;autoUpdateTime"  json:"updateTime,omitzero"`
	UpdateBy   string    `gorm:"type:string;size:32;size:32;comment:更新者;" json:"updateBy"`
	CreateBy   string    `gorm:"type:string;size:32;size:32;comment:创建者;<-:create;" json:"createBy"`
	//TenantId   int     `gorm:"type:string;size:32;index:idx_tenant_id;not null;comment:租户id;" json:"tenantId" form:"tenantId"`
}
