package model

// ==========================================================================
// LV自动生成数据库操作代码，无需手动修改，重新生成会自动覆盖.
// 生成日期：2024-08-12 21:03:00 +0800 CST
// 生成人：lv
// ==========================================================================

import (
	"github.com/lostvip-com/lv_framework/lv_db"
	"time"
)

// SysLogininfor 系统访问记录
type SysLoginInfo struct {
	InfoId        int       `gorm:"type:bigint;size:20;primary_key;auto_increment;访问ID;" json:"infoId"`
	UserName      string    `gorm:"type:varchar(50);comment:登录账号;" json:"userName"`
	Ipaddr        string    `gorm:"type:varchar(50);comment:登录IP地址;" json:"ipaddr"`
	LoginLocation string    `gorm:"type:varchar(255);comment:登录地点;" json:"loginLocation"`
	Browser       string    `gorm:"type:varchar(50);comment:浏览器类型;" json:"browser"`
	Os            string    `gorm:"type:varchar(50);comment:操作系统;" json:"os"`
	Status        string    `gorm:"type:char(1);comment:登录状态（0成功 1失败）;" json:"status"`
	Msg           string    `gorm:"type:varchar(255);comment:提示消息;" json:"msg"`
	LoginTime     time.Time `gorm:"type:datetime;comment:访问时间;" json:"loginTime" time_format:"2006-01-02 15:04:05"`

	UserId   int    `gorm:"-" json:"userId"`
	DeptId   int    `gorm:"-" json:"deptId"`
	TokenId  string `gorm:"-" json:"tokenId"`
	RoleKeys string `gorm:"-" json:"roleKeys"`
	DeptName string `gorm:"-" json:"deptName"`
	Avatar   string `gorm:"-" json:"avatar"` //登录名称对应的头像
	TenantId int    `gorm:"-" json:"avatar"`
}

func (*SysLoginInfo) TableName() string {
	return "sys_logininfor"
}

// Insert 插入数据
func (r *SysLoginInfo) Insert() error {
	return lv_db.GetOrmDefault().Save(r).Error
}

// Update 更新数据
func (r *SysLoginInfo) Update() error {
	return lv_db.GetOrmDefault().Updates(r).Error
}

// Delete 删除
func (r *SysLoginInfo) Delete() error {
	return lv_db.GetOrmDefault().Delete(r).Error
}

// 查
func (e *SysLoginInfo) FindById() error {
	err := lv_db.GetOrmDefault().Take(e, e.InfoId).Error
	return err
}
