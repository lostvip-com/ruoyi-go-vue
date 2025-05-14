// ==========================================================================
// LV自动生成model扩展代码列表、增、删，改、查、导出，只生成一次，按需修改,再次生成不会覆盖.
// 生成日期：2020-02-17 14:03:51
// 生成路径: app/model_cmn/online/online.go
// 生成人：yunjie
// ==========================================================================
package vo

import (
	"github.com/lostvip-com/lv_framework/web/lv_dto"
)

type JobReq struct {
	lv_dto.Paging
	JobName      string `form:"jobName"`
	JobGroup     string `form:"jobGroup"`
	Status       string `form:"status"`
	InvokeTarget string `form:"invokeTarget"`
	BeginTime    string `form:"beginTime"` //开始时间
	EndTime      string `form:"endTime"`   //结束时间
}

type JobStatus struct {
	JobId    int64  `form:"jobId" form:"jobId"`
	Status   string `form:"status" form:"status"`
	JobGroup string `form:"jobGroup" form:"jobGroup"`
}
