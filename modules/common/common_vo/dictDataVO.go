package common_vo

import (
	"github.com/lostvip-com/lv_framework/web/lv_dto"
)

// AddDictDataReq 新增页面请求参数
type AddDictDataReq struct {
	DictLabel string `json:"dictLabel" form:"dictLabel"  binding:"required"`
	DictValue string `json:"dictValue" form:"dictValue"  binding:"required"`
	DictType  string `json:"dictType" form:"dictType"  binding:"required"`
	DictSort  int    `json:"dictSort" form:"dictSort" binding:"required"`
	CssClass  string `json:"cssClass" form:"cssClass"`
	ListClass string `json:"listClass" form:"listClass" binding:"required"`
	IsDefault string `json:"isDefault" form:"isDefault"`
	Status    string `json:"status" form:"status"    binding:"required"`
	Remark    string `json:"remark" form:"remark"`
}

// EditDictDataReq 修改页面请求参数
type EditDictDataReq struct {
	DictCode  int    `form:"dictCode" binding:"required"`
	DictLabel string `form:"dictLabel"  binding:"required"`
	DictValue string `form:"dictValue"  binding:"required"`
	DictType  string `form:"dictType"`
	DictSort  int    `form:"dictSort"  `
	CssClass  string `form:"cssClass"`
	ListClass string `form:"listClass" `
	IsDefault string `form:"isDefault" `
	Status    string `form:"status"    `
	Remark    string `form:"remark"`
}

// SelectDictDataPageReq 分页请求参数
type SelectDictDataPageReq struct {
	DictType  string `form:"dictType"`  //字典名称
	DictLabel string `form:"dictLabel"` //字典标签
	Status    string `form:"status"`    //状态
	BeginTime string `form:"beginTime"` //开始时间
	EndTime   string `form:"endTime"`   //结束时间
	lv_dto.Paging
}
