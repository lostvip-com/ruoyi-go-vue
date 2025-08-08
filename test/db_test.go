package test

import (
	"common/myconf"
	"fmt"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_json"
	"system/service"
	"system/vo"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	cfg := myconf.GetConfigInstance()
	lv_db.GetDB("db-log")
	fmt.Println(cfg)
	req := &vo.OperLogPageReq{}
	list, totlal, err := service.GetOperLogServiceInstance().FindPage(req)

	lv_log.Info("============", lv_json.ToJsonStr(list), totlal, err)
	fmt.Println("============", lv_json.ToJsonStr(list), totlal, err)
	if err != nil {
		t.Log("=======err=====", lv_json.ToJsonStr(list), totlal, err)
	} else {
		t.Log("======success======", lv_json.ToJsonStr(list), totlal, err)
	}
	time.Sleep(time.Second * 1)
}
