package model

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/namedsql"
	"github.com/lostvip-com/lv_framework/utils/lv_time"
	"gorm.io/gorm"
)

type GenTable struct {
	TableId int64 `gorm:"type:int(20);primary_key;auto_increment;参数主键;" json:"tableId" form:"tableId"`

	ParentMenuId   int64  `json:"parentMenuId" gorm:"size:20;comment:上级菜单ID;"`
	Table_Name     string `json:"tableName,omitempty" gorm:"size:32"`
	TableComment   string `json:"tableComment" gorm:"size:32"`
	SubTableName   string `json:"subTableName" gorm:"size:32"`
	SubTableFkName string `json:"subTableFkName" gorm:"size:32"`
	ClassName      string `json:"className" gorm:"size:32"`
	TplCategory    string `json:"tplCategory" gorm:"size:32"`
	PackageName    string `json:"packageName" gorm:"size:32"`
	ModuleName     string `json:"moduleName" gorm:"size:32"`
	BusinessName   string `json:"businessName" gorm:"size:32"`
	FunctionName   string `json:"functionName" gorm:"size:32"`
	FunctionAuthor string `json:"functionAuthor" gorm:"size:32"`
	GenType        string `json:"genType" gorm:"size:32"`
	GenPath        string `json:"genPath" gorm:"size:64"`
	Options        string `json:"options" gorm:"size:128"`
	models.BaseModel
	HasEditTime string `gorm:"-"` //1需要导入time.Time 0 不需要
}

// TableName 映射数据表
func (r *GenTable) TableName() string {
	return "gen_table"
}

// BeforeCreate 实现钩子
func (u *GenTable) BeforeCreate(db *gorm.DB) error {
	u.CreateTime = lv_time.GetCurrentTime() // 设置创建时的更新时间
	u.UpdateTime = u.CreateTime             // 设置创建时的更新时间
	return nil
}

// BeforeUpdate 实现 BeforeUpdate 钩子
func (u *GenTable) BeforeUpdate(db *gorm.DB) error {
	u.CreateTime = lv_time.GetCurrentTime() // 设置更新时的更新时间
	return nil
}

// Save 增
func (e *GenTable) Save() error {
	return lv_db.GetMasterGorm().Save(e).Error
}

// FindById 查
func (e *GenTable) FindById(id int64) (*GenTable, error) {
	err := lv_db.GetMasterGorm().Take(e, id).Error
	return e, err
}

// FindOne 查第一条
func (e *GenTable) FindOne() (*GenTable, error) {
	tb := lv_db.GetMasterGorm().Table(e.TableName())

	if e.Table_Name != "" {
		tb = tb.Where("table_name=?", e.Table_Name)
	}
	if e.TableComment != "" {
		tb = tb.Where("table_comment=?", e.TableComment)
	}
	if e.ClassName != "" {
		tb = tb.Where("class_name=?", e.ClassName)
	}
	if e.TplCategory != "" {
		tb = tb.Where("tpl_category=?", e.TplCategory)
	}
	err := tb.First(e).Error
	return e, err
}

// Updates 改
func (e *GenTable) Updates() error {
	return lv_db.GetMasterGorm().Table(e.TableName()).Updates(e).Error
}

// Delete 删
func (e *GenTable) Delete() error {
	return lv_db.GetMasterGorm().Delete(e).Error
}

func (e *GenTable) Count() (int64, error) {
	sql := " select count(*) from gen_table where del_flag = 0 "

	if e.Table_Name != "" {
		sql += " and table_name = @Table_Name "
	}
	if e.TableComment != "" {
		sql += " and table_comment = @TableComment "
	}
	if e.ClassName != "" {
		sql += " and class_name = @ClassName "
	}
	if e.TplCategory != "" {
		sql += " and tpl_category = @TplCategory "
	}

	return namedsql.Count(lv_db.GetMasterGorm(), sql, e)
}
