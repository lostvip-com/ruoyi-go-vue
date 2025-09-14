package jobs

import (
	"common/schedule"
	"github.com/lostvip-com/lv_framework/lv_log"
)

func init() {
	//在这里注册定时任务的执行函数不知
	go func() {
		schedule.GetSchedulerServiceInstance().RegFunc("SayHello", SayHello)
	}()
}
func SayHello() error {
	lv_log.Info("hello job !!!!!!!!!!!!!!!!!!!!!!!")
	return nil
}
