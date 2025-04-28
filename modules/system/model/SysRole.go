package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
	"time"
)

type SysRole struct {
	RoleId            int64     `json:"roleId" gorm:"column:role_id;primaryKey"` //表示主键
	RoleName          string    `json:"roleName" gorm:"role_name"`
	RoleKey           string    `json:"roleKey" gorm:"role_key"`
	RoleSort          int       `json:"roleSort" gorm:"role_sort"`
	DataScope         string    `json:"dataScope" gorm:"data_scope"`
	Status            string    `json:"status" gorm:"status"`
	MenuCheckStrictly bool      `json:"menuCheckStrictly" gorm:"menu_check_strictly"`
	DeptCheckStrictly bool      `json:"deptCheckStrictly" gorm:"dept_check_strictly"`
	DelFlag           string    `json:"delFlag" gorm:"del_flag"`
	CreateBy          string    `json:"createBy" gorm:"create_by"`
	CreateTime        time.Time `json:"createTime" gorm:"column:create_time;type:datetime;autoCreateTime"`
	UpdateBy          string    `json:"updateBy" gorm:"update_by"`
	UpdateTime        time.Time `json:"updateTime" gorm:"column:update_time;type:datetime;autoCreateTime"`
	Remark            string    `json:"remark" gorm:"remark"`
	models.BaseModel
}

// 映射数据表
func (r *SysRole) TableName() string {
	return "sys_role"
}

// 插入数据
func (r *SysRole) Insert() error {
	return lv_db.GetMasterGorm().Save(r).Error
}

// 更新数据
func (r *SysRole) Update() error {
	return lv_db.GetMasterGorm().Updates(r).Error
}

// 删除
func (r *SysRole) Delete() error {
	return lv_db.GetMasterGorm().Delete(r).Error
}

// 根据结构体中已有的非空数据来获得单条数据
func (e *SysRole) FindOne() error {
	tb := lv_db.GetMasterGorm()
	if e.RoleId != 0 {
		tb = tb.Where("role_id=? and del_flag=0", e.RoleId)
	}
	if e.RoleKey != "" {
		tb = tb.Where("role_key=? and del_flag=0", e.RoleKey)
	}
	err := tb.First(e).Error
	return err
}
