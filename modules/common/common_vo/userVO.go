package common_vo

import (
	"github.com/lostvip-com/lv_framework/web/lv_dto"
)

// Fill with you ideas below.

// 修改用户资料请求参数
type ProfileReq struct {
	NickName    string `json:"nickName"  binding:"required,min=1,max=30"`
	Phonenumber string `json:"phonenumber"  binding:"required,len=11"`
	Email       string `json:"email"  binding:"required,email"`
	Sex         string `json:"sex"  binding:"required"`
}

// 修改密码请求参数
type PasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required,min=5,max=30"`
	NewPassword string `json:"newPassword" binding:"required,min=5,max=30"`
}

// 重置密码请求参数
type ResetPwdReq struct {
	UserId   int    `form:"userId"  binding:"required"`
	Password string `form:"password" binding:"required,min=5,max=30"`
}

// 检查email请求参数
type CheckEmailReq struct {
	UserId int    `form:"userId"  binding:"required,min=5,max=30"`
	Email  string `form:"email"  binding:"required,email"`
}

// 检查email请求参数
type CheckEmailAllReq struct {
	Email string `form:"email"  binding:"required,email"`
}

// 检查phone请求参数
type CheckUserNameReq struct {
	UserName string `form:"UserName"  binding:"required"`
}

// 检查phone请求参数
type CheckPhoneReq struct {
	UserId      int    `form:"userId"  binding:"required`
	Phonenumber string `form:"phonenumber"  binding:"required,len=11"`
}

// 检查phone请求参数
type CheckPhoneAllReq struct {
	Phonenumber string `form:"phonenumber"  binding:"required,len=11"`
}

// 检查密码请求参数
type CheckPasswordReq struct {
	Password string `form:"password"  binding:"required"`
}

// 查询用户列表请求参数
type UserPageReq struct {
	UserName    string `form:"UserName"`    //登录名
	Status      string `form:"status"`      //状态
	Phonenumber string `form:"phonenumber"` //手机号码
	BeginTime   string `form:"beginTime"`   //数据范围
	EndTime     string `form:"endTime"`     //开始时间
	DeptId      int    `form:"deptId"`      //结束时间
	SortName    string `form:"sortName"`    //排序字段
	SortOrder   string `form:"sortOrder"`   //排序方式
	Ancestors   string `form:"ancestors"`   //排序方式
	//
	TenantId int `form:"tenantId"`
	lv_dto.Paging
}

// 新增用户资料请求参数
type AddUserReq struct {
	UserName    string `json:"userName" binding:"required"`
	Phonenumber string `json:"phonenumber"`
	Email       string `json:"email"`
	NickName    string `json:"nickName"  binding:"required"`
	Password    string `json:"password"  binding:"required,min=6,max=30"`
	DeptId      int    `json:"deptId" binding:"required`
	Sex         string `json:"sex"`
	Status      string `json:"status"`
	RoleIds     []int  `json:"roleIds"`
	PostIds     []int  `json:"postIds"`
	Remark      string `json:"remark"`
}

// 新增用户资料请求参数
type EditUserReq struct {
	UserId      int    `json:"userId" binding:"required`
	UserName    string `json:"userName" binding:"required"`
	Phonenumber string `json:"phonenumber"`
	Email       string `json:"email"`
	NickName    string `json:"nickName"  binding:"required"`
	DeptId      int    `json:"deptId" binding:"required`
	Sex         string `json:"sex"  binding:"required"`
	Status      string `json:"status"`
	RoleIds     []int  `json:"roleIds"`
	PostIds     []int  `json:"postIds"`
	Remark      string `json:"remark"`
}
