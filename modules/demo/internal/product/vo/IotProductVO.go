// ==========================================================================
// LV自动生成model扩展代码列表 按需修改
// 生成日期: 2025-07-25 07:15:20 &#43;0000 UTC
// 生成人: lv
// ==========================================================================
package vo

import (
  "github.com/lostvip-com/lv_framework/web/lv_dto"
  "time"
)

//分页请求参数 
type IotProductReq struct {
  Key string `form:"key"` //产品编码
  Name string `form:"name"` //名字
  CloudProductId string `form:"cloudProductId"` //云产品ID
  CloudInstanceId string `form:"cloudInstanceId"` //云实例ID
  BeginTime string `form:"beginTime" json:"beginTime"` //数据范围
  EndTime   string `form:"endTime" json:"endTime"`     //开始时间
  lv_dto.Paging
}


//分页请求结果
type IotProductResp struct {
  Id    int64  `json:"id"`
      Key string `json:"key"`
      Name string `json:"name"`
      CloudProductId string `json:"cloudProductId"`
      CloudInstanceId string `json:"cloudInstanceId"`
      Platform string `json:"platform"`
      Protocol string `json:"protocol"`
      NodeType string `json:"nodeType"`
      NetType string `json:"netType"`
      DataFormat string `json:"dataFormat"`
      LastSyncTime int64 `json:"lastSyncTime"`
      Factory string `json:"factory"`
      Description string `json:"description"`
      Status string `json:"status"`
      Extra string `json:"extra"`
      Manufacturer string `json:"manufacturer"`
  UpdateBy string      `json:"updateBy"`
  UpdateTime time.Time `json:"updateTime" time_format:"2006-01-02 15:04:05"`
  CreateTime time.Time `json:"createTime" time_format:"2006-01-02 15:04:05"`
  CreateBy   string    `json:"createBy"`
}
