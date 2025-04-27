package model

import (
	"github.com/lostvip-com/lv_framework/lv_db"
	"time"
)

type SysMenu struct {
	MenuId     int64     `gorm:"size:20;primary_key;auto_increment;comment:菜单ID;" json:"menu_id"` // 菜单ID
	MenuName   string    `gorm:"size:32;" json:"menu_name"`                                       // 菜单名称
	ParentId   int64     `gorm:"size:32;" json:"parent_id"`                                       // 父菜单ID
	OrderNum   int32     `gorm:"size:32;" json:"order_num"`                                       // 显示顺序
	Path       string    `gorm:"size:32;" json:"path"`                                            // 路由地址
	Component  string    `gorm:"size:32;" json:"component"`                                       // 组件路径
	Query      string    `gorm:"size:32;" json:"query"`                                           // 路由参数
	RouteName  string    `gorm:"size:32;" json:"route_name"`                                      // 路由名称
	IsFrame    string    `gorm:"size:1;" json:"is_frame"`                                         // 是否为外链（0是 1否）
	IsCache    string    `gorm:"size:1;" json:"is_cache"`                                         // 是否缓存（0缓存 1不缓存）
	MenuType   string    `gorm:"size:1;" json:"menu_type"`                                        // 菜单类型（M目录 C菜单 F按钮）
	Visible    string    `gorm:"size:1;" json:"visible"`                                          // 菜单状态（0显示 1隐藏）
	Status     string    `gorm:"size:1;" json:"status"`                                           // 菜单状态（0正常 1停用）
	Perms      string    `gorm:"size:32;" json:"perms"`                                           // 权限标识
	Icon       string    `gorm:"size:32;" json:"icon"`                                            // 菜单图标
	CreateBy   string    `gorm:"size:32;" json:"create_by"`                                       // 创建者
	CreateTime time.Time `gorm:"size:32;" json:"create_time"`                                     // 创建时间
	UpdateBy   string    `gorm:"size:32;" json:"update_by"`                                       // 更新者
	UpdateTime time.Time `gorm:"size:32;" json:"update_time"`                                     // 更新时间
	Remark     string    `gorm:"size:32;" json:"remark"`                                          // 备注
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
