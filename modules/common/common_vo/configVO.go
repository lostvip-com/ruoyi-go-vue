package common_vo

import (
	"github.com/lostvip-com/lv_framework/web/lv_dto"
)

type ConfigPageReq struct {
	ConfigName string `form:"configName"` //参数名称
	ConfigKey  string `form:"configKey"`  //参数键名
	ConfigType string `form:"configType"` //状态
	BeginTime  string `form:"beginTime"`  //开始时间
	EndTime    string `form:"endTime"`    //结束时间
	lv_dto.Paging
}

// 检查参数键名请求参数
type CheckConfigKeyReq struct {
	ConfigId  int    `form:"configId"  binding:"required"`
	ConfigKey string `form:"configKey"  binding:"required"`
}

//
//// 检查参数键名请求参数
//type CheckPostCodeALLReq struct {
//	ConfigKey string `form:"configKey"  binding:"required"`
//}
