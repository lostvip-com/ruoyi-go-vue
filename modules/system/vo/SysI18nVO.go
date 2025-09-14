// ==========================================================================
// LV自动生成model扩展代码列表 按需修改
// 生成日期: 2025-08-11 07:41:35 &#43;0000 UTC
// 生成人: dpc
// ==========================================================================
package vo

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
)

// 分页请求参数
type SysI18nReq struct {
	Locale     string `form:"locale" json:"locale"`         //本地标识
	LocaleKey  string `form:"localeKey" json:"localeKey"`   //国际化key
	LocaleName string `form:"localeName" json:"localeName"` //国际化名称
	Sort       int    `form:"sort" json:"sort"`             //字典排序
	BeginTime  string `form:"beginTime" json:"beginTime"`   //数据范围
	EndTime    string `form:"endTime" json:"endTime"`       //开始时间
	lv_dto.Paging
}

// 分页请求结果
type SysI18nResp struct {
	Id         int              `json:"id"`
	Locale     string           `json:"locale"`
	LocaleKey  string           `json:"localeKey"`
	LocaleName string           `json:"localeName"`
	Sort       int              `json:"sort"`
	Remark     string           `json:"remark"`
	UpdateBy   string           `json:"updateBy"`
	UpdateTime models.LocalTime `json:"updateTime"`
	CreateTime models.LocalTime `json:"createTime"`
	CreateBy   string           `json:"createBy"`
}
