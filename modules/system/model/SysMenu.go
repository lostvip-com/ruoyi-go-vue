package model

import (
	"github.com/lostvip-com/lv_framework/lv_db"
	"time"
)

type SysMenu struct {
	MenuId     int64     `json:"menuId" gorm:"column:menu_id;primaryKey"` //表示主键
	MenuName   string    `json:"menuName" gorm:"menu_name"`
	ParentId   int64     `json:"parentId" gorm:"parent_id"`
	OrderNum   int64     `json:"orderNum" gorm:"order_num"`
	MenuType   string    `json:"menuType" gorm:"menu_type"`
	Visible    string    `json:"visible" gorm:"visible"`
	Perms      string    `json:"perms" gorm:"perms"`
	Query      string    `json:"query" gorm:"query"`
	IsFrame    string    `json:"isFrame" gorm:"is_frame"`
	Icon       string    `json:"icon" gorm:"icon"`
	Path       string    `json:"path" gorm:"path"`
	Status     string    `json:"status" gorm:"status"`
	IsCache    string    `json:"isCache" gorm:"is_cache"`
	Component  string    `json:"component" gorm:"component"`
	CreateBy   string    `json:"createBy" gorm:"create_by"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time;type:datetime;autoCreateTime"`
	UpdateBy   string    `json:"updateBy" gorm:"update_by"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time;type:datetime;autoCreateTime"`
	Remark     string    `json:"remark" gorm:"remark"` // 备注
	//
	Children   []SysMenu `gorm:"-" json:"children"`
	ParentName string    `gorm:"-" json:"parentName"`
}

func (e *SysMenu) TableName() string {
	return "sys_menu"
}

// 增
func (e *SysMenu) Save() error {
	return lv_db.GetMasterGorm().Save(e).Error
}

// 查
func (e *SysMenu) FindById() error {
	err := lv_db.GetMasterGorm().Take(e, e.MenuId).Error
	return err
}

// 查第一条
func (e *SysMenu) FindOne() error {
	tb := lv_db.GetMasterGorm().Table("sys_menu")
	if e.MenuId != 0 {
		tb = tb.Where("menu_id=?", e.MenuId)
	}
	if e.ParentId != 0 {
		tb = tb.Where("parent_id=?", e.ParentId)
	}
	if e.MenuName != "" {
		tb = tb.Where("menu_name=?", e.MenuName)
	}

	if e.Perms != "" {
		tb = tb.Where("perms=?", e.Perms)
	}

	err := tb.First(e).Error
	return err
}

// 查第一条
func (e *SysMenu) FindLastOne() error {
	tb := lv_db.GetMasterGorm().Table("sys_menu")
	if e.MenuId != 0 {
		tb = tb.Where("menu_id=?", e.MenuId)
	}
	if e.ParentId != 0 {
		tb = tb.Where("parent_id=?", e.ParentId)
	}
	if e.MenuName != "" {
		tb = tb.Where("menu_name=?", e.MenuName)
	}
	if e.Perms != "" {
		tb = tb.Where("perms=?", e.Perms)
	}
	tb.Order("menu_id desc")
	tb.Limit(1)
	err := tb.First(e).Error
	return err
}

// 改
func (e *SysMenu) Update() error {
	return lv_db.GetMasterGorm().Table(e.TableName()).Updates(e).Error
}

// 删
func (e *SysMenu) Delete() error {
	return lv_db.GetMasterGorm().Table(e.TableName()).Delete(e).Error
}
