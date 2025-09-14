package vo

import (
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"system/model"
)

type GenTableVO struct {
	model.GenTable
	TreeCode       string                 `gorm:"-"` // 树编码字段
	TreeParentCode string                 `gorm:"-"` // 树父编码字段
	TreeName       string                 `gorm:"-"` // 树名称字段
	Columns        []model.GenTableColumn `gorm:"-"` // 表列信息
	PkColumn       model.GenTableColumn   `gorm:"-"` // 表列信息
}

type EditGenTableVO struct {
	model.GenTable
	Tree    bool                   `json:"tree"`
	Crud    bool                   `json:"crud"`
	Sub     bool                   `json:"sub"`
	Columns []model.GenTableColumn `json:"columns"`
}

type GenTableParams struct {
	TreeCode       string `form:"treeCode"`
	TreeParentCode string `form:"treeParentCode"`
	TreeName       string `form:"treeName"`
}

// 分页请求参数
type GenTablePageReq struct {
	TableName    string `form:"tableName"`    //表名称
	TableComment string `form:"tableComment"` //表描述
	BeginTime    string `form:"beginTime"`    //开始时间
	EndTime      string `form:"endTime"`      //结束时间
	lv_dto.Paging
}
