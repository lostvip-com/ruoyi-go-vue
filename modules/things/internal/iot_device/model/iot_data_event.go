package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/namedsql"
	"time"
)

type IotDataEvent struct {
	Id        int    `gorm:"type:bigint;size:20;primary_key;auto_increment;comment:主键;" json:"id"`
	DeviceId  int    `gorm:"type:bigint;comment:设备ID"`
	ProductId int    `gorm:"type:bigint;comment:产品ID"`
	DeptId    int    `gorm:"type:bigint;comment:组织ID"`
	Code      string `gorm:"type:string;size:32;comment:唯一识别号"`
	Type      string `gorm:"type:string;size:1;comment:事件类型，event:0,alert 1,fault 2"`
	Level     string `gorm:"type:string;size:1;comment:级别"`

	Content     string    `gorm:"type:string;size:64;comment:内容"`
	Message     string    `gorm:"type:string;size:64;comment:处理意见"`
	TriggerTime time.Time `gorm:"type:datetime;comment:触发时间" time_format:"2006-01-02 15:04:05" json:"updateTime"`
	RecoverTime time.Time `gorm:"type:datetime;comment:恢复时间" time_format:"2006-01-02 15:04:05" json:"updateTime"`
	TreatedTime time.Time `gorm:"type:datetime;comment:处理时间" time_format:"2006-01-02 15:04:05" json:"updateTime"`

	Status           string `gorm:"type:string;size:1;comment:状态"`
	models.BaseModel `gorm:"embedded"`
}

func (d *IotDataEvent) TableName() string {
	return "iot_data_event"
}

func (e *IotDataEvent) Save() error {
	return lv_db.GetOrmDefault().Save(e).Error
}

func (e *IotDataEvent) FindById(id int) (*IotDataEvent, error) {
	err := lv_db.GetOrmDefault().Table(e.TableName()).Take(e, id).Error
	return e, err
}

func (e *IotDataEvent) FindOne() (*IotDataEvent, error) {
	tb := lv_db.GetOrmDefault().Table(e.TableName())

	err := tb.First(e).Error
	return e, err
}

func (e *IotDataEvent) Updates() error {
	return lv_db.GetOrmDefault().Table(e.TableName()).Updates(e).Error
}

func (e *IotDataEvent) Delete() error {
	return lv_db.GetOrmDefault().Delete(e).Error
}

func (e *IotDataEvent) Count() (int64, error) {
	sql := " select count(*) from iot_alert_ai_di where del_flag = 0 "
	return namedsql.Count(lv_db.GetOrmDefault(), sql, e)
}
