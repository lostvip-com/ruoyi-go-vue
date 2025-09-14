package things

import (
	"common/myconf"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	model3 "things/internal/iot_device/model"
	_ "things/internal/iot_product"
	model2 "things/internal/iot_product/model"
)

//自动建表

func init() {
	//自动同步表结构
	cfg := myconf.GetConfigInstance()
	migrate := cfg.GetAutoMigrate()
	if migrate == "create" || migrate == "update" || migrate == "true" {
		lv_log.Warn("######### 开始同步表结构: ############## migrate" + migrate)
		err := lv_db.GetOrmDefault().AutoMigrate(
			model2.IotProduct{},
			model2.IotPrdProperty{},
			model2.IotPrdAction{},
			model2.IotPrdEvent{},
			model3.IotDevice{},
			model3.IotDataEvent{},
			model3.IotDataAction{},
			model3.IotDataEvent{},
		)
		lv_err.HasErrAndPanic(err)
	} else {
		lv_log.Warn("========== 已关闭表结构同步功能========== migrate" + migrate)
	}
}
