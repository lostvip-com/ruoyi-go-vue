package model

import (
	"common/models"
	"gorm.io/gorm"
)

type IotAuth struct {
	Id               int    `gorm:"type:bigint;size:20;primary_key;auto_increment;comment:主键;" json:"dictId"`
	ResourceId       string `gorm:"type:string;size:32;comment:资源ID"`
	ResourceType     string `gorm:"type:string;size:32;comment:资源类型"`
	ClientId         string `gorm:"uniqueIndex;size:32;comment:客户端ID"`
	UserName         string `gorm:"type:string;size:32;comment:用户名"`
	Password         string `gorm:"type:string;size:32;comment:密码"`
	models.BaseModel `gorm:"embedded"`
}

func (d *IotAuth) TableName() string {
	return "iot_auth"
}

func (d *IotAuth) Get() interface{} {
	return *d
}

func (d *IotAuth) BeforeCreate(tx *gorm.DB) (err error) {

	return nil
}
