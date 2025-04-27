// ==========================================================================
// LV自动生成数据库操作代码，无需手动修改，重新生成会自动覆盖.
// 生成日期：2024-01-06 17:33:26 +0800 CST
// 生成人：lv
// ==========================================================================

package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
	"time"
)

// SysUser 用户信息
type SysUser struct {
	UserId      int64      `gorm:"size:20;primary_key;auto_increment;用户ID;"     json:"userId"  form:"userId"`
	DeptId      int64      `gorm:"size:20;comment:部门ID;" json:"deptId" form:"deptId"`
	UserName    string     `gorm:"size:32;comment:登录账号;" json:"UserName" form:"UserName"`
	NickName    string     `gorm:"size:32;comment:用户昵称;" json:"userName" form:"userName"`
	UserType    string     `gorm:"size:2;comment:用户类型（00系统用户）;" json:"userType" form:"userType"`
	Email       string     `gorm:"size:50;comment:用户邮箱;" json:"email" form:"email"`
	Phonenumber string     `gorm:"size:11;comment:手机号码;" json:"phonenumber" form:"phonenumber"`
	Sex         string     `gorm:"size:1;comment:用户性别（0男 1女 2未知）;" json:"sex" form:"sex"`
	Avatar      string     `gorm:"size:100;comment:头像路径;" json:"avatar" form:"avatar"`
	Password    string     `gorm:"size:50;comment:密码;" json:"password" form:"password"`
	Salt        string     `gorm:"size:20;comment:盐加密;" json:"salt" form:"salt"`
	Status      string     `gorm:"size:char(1;comment:帐号状态（0正常 1停用）;" json:"status" form:"status"`
	LoginIp     string     `gorm:"size:50;comment:最后登陆IP;" json:"loginIp" form:"loginIp"`
	LoginDate   *time.Time `gorm:"size:datetime;comment:最后登陆时间;" json:"loginDate" form:"loginDate" time_format:"2006-01-02 15:04:05"`
	UpdateBy    string     `gorm:"size:64;comment:更新者;" json:"updateBy" form:"updateBy"`
	Remark      string     `gorm:"size:500;comment:备注;" json:"remark"   form:"remark"`
	CreateBy    string     `gorm:"size:32;comment:创建人;column:create_by;"  json:"createBy"`
	models.BaseModel
	//临时属性
	Dept  *models.SysDept `gorm:"-"  json:"dept"`
	Roles []SysRole       `gorm:"-"  json:"roles"`
}

// 映射数据表
func (e *SysUser) TableName() string {
	return "sys_user"
}

// 增
func (e *SysUser) Insert() error {
	return lv_db.GetMasterGorm().Save(e).Error
}

// 查
func (e *SysUser) GetById() error {
	err := lv_db.GetMasterGorm().Take(e).Error
	return err
}

// 查
func (e *SysUser) FindById() error {
	tb := lv_db.GetMasterGorm()
	err := tb.Take(e, e.UserId).Error
	return err
}

// 查
func (e *SysUser) FindOne() error {
	tb := lv_db.GetMasterGorm()
	if e.UserId != 0 {
		tb = tb.Where("user_id=?", e.UserId)
	}

	if e.UserName != "" && e.UserName != "" {
		tb = tb.Where("user_name=?", e.UserName)
	}
	err := tb.First(e).Error
	return err
}

// 改
func (e *SysUser) Updates() error {
	return lv_db.GetMasterGorm().Table(e.TableName()).Updates(e).Error
}

// 删
func (e *SysUser) Delete() error {
	return lv_db.GetMasterGorm().Delete(e).Error
}

func (e *SysUser) GetRoleKeys() []string {
	keys := make([]string, 0)
	if e.Roles != nil {
		for _, role := range e.Roles {
			keys = append(keys, role.RoleKey)
		}
	}
	return keys
}
