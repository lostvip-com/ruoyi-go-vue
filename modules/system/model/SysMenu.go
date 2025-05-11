package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
)

type SysMenu struct {
	MenuId   int64  `json:"menuId" gorm:"column:menu_id;primaryKey"`
	MenuName string `json:"menuName" gorm:"menu_name"`
	ParentId int64  `json:"parentId" gorm:"parent_id"`
	OrderNum int64  `json:"orderNum" gorm:"order_num"`
	Path     string `json:"path" gorm:"path"`

	Component string `json:"component" gorm:"component"`
	Query     string `json:"query" gorm:"query"`
	RouteName string `json:"query" gorm:"route_name"`
	IsFrame   string `json:"isFrame" gorm:"is_frame"`
	IsCache   string `json:"isCache" gorm:"is_cache"`

	MenuType string `json:"menuType" gorm:"menu_type"`
	Visible  string `json:"visible" gorm:"visible"`
	Status   string `json:"status" gorm:"status"`
	Perms    string `json:"perms" gorm:"perms"`
	Icon     string `json:"icon" gorm:"icon"`

	Remark string `json:"remark" gorm:"remark"` // 备注
	models.BaseModel
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
func (e *SysMenu) FindById(id int64) (*SysMenu, error) {
	err := lv_db.GetMasterGorm().Take(e, id).Error
	return e, err
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
