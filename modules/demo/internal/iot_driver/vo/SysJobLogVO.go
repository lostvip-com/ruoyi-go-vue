// ==========================================================================
// LV自动生成model扩展代码列表 按需修改
// 生成日期: 2025-05-17 13:18:47 &#43;0000 UTC
// 生成人: lv
// ==========================================================================
package vo

import (
  "github.com/lostvip-com/lv_framework/web/lv_dto"
  "time"
)

//分页请求参数 
type SysJobLogReq struct {
  ExceptionInfo string `form:"exceptionInfo"` //异常信息
  Status string `form:"status"` //执行状态（0正常 1失败）
  JobMessage string `form:"jobMessage"` //日志信息
  InvokeTarget string `form:"invokeTarget"` //调用目标字符串
  BeginTime string `form:"beginTime" json:"beginTime"` //数据范围
  EndTime   string `form:"endTime" json:"endTime"`     //开始时间
  lv_dto.Paging
}


//分页请求结果
type SysJobLogResp struct {
  JobLogId    int64  `json:"jobLogId"`
      ExceptionInfo string `json:"exceptionInfo"`
      Status string `json:"status"`
      JobMessage string `json:"jobMessage"`
      InvokeTarget string `json:"invokeTarget"`
      JobGroup string `json:"jobGroup"`
      JobName string `json:"jobName"`
  UpdateBy string      `json:"updateBy"`
  UpdateTime time.Time `json:"updateTime" time_format:"2006-01-02 15:04:05"`
  CreateTime time.Time `json:"createTime" time_format:"2006-01-02 15:04:05"`
  CreateBy   string    `json:"createBy"`
}
