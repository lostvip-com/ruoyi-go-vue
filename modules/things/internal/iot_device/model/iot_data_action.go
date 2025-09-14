package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/namedsql"
	"time"
)

type IotDataAction struct {
	Id       int    `gorm:"type:bigint;size:20;primary_key;auto_increment;comment:主键;" json:"dictId"`
	DeviceId int    `gorm:"type:bigint;comment:设备ID"`
	Code     string `gorm:"type:string;size:32;comment:属性名"`
	Tags     string `gorm:"type:string;size:64;comment:TAG"`
	Name     string `gorm:"type:string;size:32;comment:属性名"`

	Value            string    `gorm:"type:string;size:32;comment:采集值"`
	DevTime          time.Time `gorm:"type:datetime;comment:采集时间" time_format:"2006-01-02 15:04:05"`
	models.BaseModel `gorm:"embedded"`
}

func (d *IotDataAction) TableName() string {
	return "iot_data_action"
}

// 增
func (e *IotDataAction) Save() error {
	return lv_db.GetOrmDefault().Save(e).Error
}

// 查
func (e *IotDataAction) FindById(id int) (*IotDataAction, error) {
	err := lv_db.GetOrmDefault().Take(e, id).Error
	return e, err
}

// 查第一条
func (e *IotDataAction) FindOne() (*IotDataAction, error) {
	tb := lv_db.GetOrmDefault().Table(e.TableName())

	err := tb.First(e).Error
	return e, err
}

// 改
func (e *IotDataAction) Updates() error {
	return lv_db.GetOrmDefault().Table(e.TableName()).Updates(e).Error
}

// 删
func (e *IotDataAction) Delete() error {
	return lv_db.GetOrmDefault().Delete(e).Error
}

func (e *IotDataAction) Count() (int64, error) {
	sql := " select count(*) from iot_data_his where del_flag = 0 "

	return namedsql.Count(lv_db.GetOrmDefault(), sql, e)
}
