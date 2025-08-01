package common_vo

import (
	"github.com/lostvip-com/lv_framework/web/lv_dto"
)

// DeptPageReq 分页请求参数
type DeptPageReq struct {
	ParentId  int    `form:"parentId"` //父部门ID
	NodeType  string `form:"nodeType"`
	DeptName  string `form:"deptName"`  //部门名称
	Status    string `form:"status"`    //状态
	BeginTime string `form:"beginTime"` //开始时间
	EndTime   string `form:"endTime"`   //结束时间
	TenantId  int    `form:"tenantId"`  //结束时间
	lv_dto.Paging
}

// 检查菜单名称请求参数
type CheckDeptNameReq struct {
	DeptId   int    `form:"deptId"  binding:"required"`
	ParentId int    `form:"parentId"  binding:"required"`
	DeptName string `form:"deptName"  binding:"required"`
}

// 检查菜单名称请求参数
type CheckDeptNameALLReq struct {
	ParentId int    `form:"parentId"  binding:"required"`
	DeptName string `form:"deptName"  binding:"required"`
}
