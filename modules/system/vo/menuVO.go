package vo

import (
	"github.com/lostvip-com/lv_framework/web/lv_dto"
)

// 检查菜单名称请求参数
type CheckMenuNameReq struct {
	MenuId   int    `form:"menuId"  binding:"required"`
	ParentId int    `form:"parentId"  binding:"required"`
	MenuName string `form:"menuName"  binding:"required"`
}

// 检查菜单名称请求参数
type CheckMenuNameALLReq struct {
	ParentId int    `form:"parentId"  binding:"required"`
	MenuName string `form:"menuName"  binding:"required"`
}

// 分页请求参数
type SelectMenuPageReq struct {
	MenuName  string `form:"menuName"`  //菜单名称
	Visible   string `form:"visible"`   //可见
	BeginTime string `form:"beginTime"` //开始时间
	EndTime   string `form:"endTime"`   //结束时间
	PageNum   int    `form:"pageNum"`   //当前页码
	PageSize  int    `form:"pageSize"`  //每页数
	SortName  string `form:"sortName"`  //排序字段
	SortOrder string `form:"sortOrder"` //排序方式
	Status    string `form:"status"`    //状态
	lv_dto.Paging
}

type MenuVo struct {
	Name       string   `json:"name"`
	Path       string   `json:"path,omitempty"`
	Hidden     bool     `json:"hidden" `
	Redirect   string   `json:"redirect,omitempty"`
	Component  string   `json:"component,omitempty" `
	Query      string   `json:"query,omitempty"`
	AlwaysShow bool     `json:"alwaysShow,omitempty" `
	MetaVo     MetaVo   `json:"meta" `
	Children   []MenuVo `json:"children,omitempty"`
}

type MetaVo struct {
	Title   string `json:"title"`
	Icon    string `json:"icon" `
	NoCache bool   `json:"noCache" `
	Link    string `json:"link,omitempty" `
}

type MenuTreeSelect struct {
	Id       int              `json:"id"`
	Label    string           `json:"label"`
	Children []MenuTreeSelect `json:"children,omitempty"`
}
