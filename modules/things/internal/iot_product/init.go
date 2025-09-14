package iot_product

import (
	"common/myconf"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"things/internal/iot_device/model"
)

func init() {
	//自动同步表结构
	cfg := myconf.GetConfigInstance()
	migrate := cfg.GetAutoMigrate()
	if migrate == "create" || migrate == "update" || migrate == "true" {
		lv_log.Warn("######### 开始同步表结构: ############## migrate" + migrate)
		err := lv_db.GetOrmDefault().AutoMigrate(
			//model.IotProduct{},
			model.IotDevice{},
			//model.IotPrdProperty{},
			//model.IotPrdAction{},
			//model.IotPrdEvent{},
		)
		lv_err.HasErrAndPanic(err)
	} else {
		lv_log.Warn("========== 已关闭表结构同步功能========== migrate" + migrate)
	}
}
