package vo

import (
	"github.com/lostvip-com/lv_framework/web/lv_dto"
)

// 查询列表请求参数
type OperLogPageReq struct {
	Title         string `json:"title" form:"title"`                 //系统模块
	OperName      string `json:"operName" form:"operName"`           //操作人员
	BusinessTypes string `json:"businessTypes" form:"businessTypes"` //操作类型
	Status        string `json:"status" form:"status"`               //操作状态
	BeginTime     string `json:"beginTime" form:"beginTime"`         //数据范围
	EndTime       string `json:"endTime" form:"endTime"`             //开始时间
	lv_dto.Paging
}

func (req *OperLogPageReq) GetBeginTime() string {
	if req.BeginTime != "" {
		return req.BeginTime + " 00:00:00"
	}
	return req.BeginTime
}
func (req *OperLogPageReq) GetEndTime() string {
	if req.EndTime != "" {
		return req.EndTime + " 23:59:59"
	}
	return req.EndTime
}
