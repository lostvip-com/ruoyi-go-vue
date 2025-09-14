package system

import (
	cm_model "common/models"
	"common/myconf"
	"common/schedule"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	_ "system/api"
	"system/model"
)

//自动建表

func init() {
	//自动同步表结构
	cfg := myconf.GetConfigInstance()
	migrate := cfg.GetAutoMigrate()
	if migrate == "create" || migrate == "update" || migrate == "true" {
		lv_log.Warn("######### 开始同步表结构: ############## migrate" + migrate)
		err := lv_db.GetOrmDefault().AutoMigrate(
			schedule.SysJobLog{}, schedule.SysJob{},
			cm_model.SysDept{}, model.SysPost{}, model.SysUser{}, model.SysDictType{}, cm_model.SysDictData{},
			model.SysMenu{}, model.SysRole{}, model.SysConfig{}, model.SysOperLog{}, model.SysLoginInfo{},
		)
		lv_err.HasErrAndPanic(err)
	} else {
		lv_log.Warn("========== 已关闭表结构同步功能========== migrate" + migrate)
	}
}
