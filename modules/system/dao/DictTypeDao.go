package dao

import (
	"common/common_vo"
	"errors"
	"github.com/lostvip-com/lv_framework/lv_db"
	"system/model"
)

type DictTypeDao struct {
}

var dictTypeDao *DictTypeDao

func GetSysDictTypeDaoInstance() *DictTypeDao {
	if dictTypeDao == nil {
		dictTypeDao = new(DictTypeDao)
	}
	return dictTypeDao
}

// 根据条件分页查询数据
func (dao *DictTypeDao) FindPage(param *common_vo.DictTypePageReq) ([]model.SysDictType, int, error) {
	db := lv_db.GetOrmDefault()
	if db == nil {
		return nil, 0, errors.New("获取数据库连接失败")
	}
	tb := db.Table("sys_dict_type t")
	if param != nil {
		if param.DictName != "" {
			tb.Where("t.dict_name like ?", "%"+param.DictName+"%")
		}
		if param.DictType != "" {
			tb.Where("t.dict_type like ?", "%"+param.DictType+"%")
		}
		if param.Status != "" {
			tb.Where("t.status = ", param.Status)
		}
		if param.BeginTime != "" {
			tb.Where("t.create_time >= ?", param.BeginTime)
		}
		if param.EndTime != "" {
			tb.Where("t.create_time <= ?", param.EndTime)
		}
	}
	var total int64
	tb = tb.Count(&total)
	var result []model.SysDictType
	tb = tb.Order("dict_id desc").Offset(param.GetStartNum()).Limit(param.GetPageSize())
	tb = tb.Find(&result)
	return result, int(total), tb.Error
}

// FindAll 获取所有数据
func (dao *DictTypeDao) FindAll(param *common_vo.DictTypePageReq) ([]model.SysDictType, error) {
	gdb := lv_db.GetOrmDefault()
	if gdb == nil {
		return nil, errors.New("获取数据库连接失败")
	}

	tb := gdb.Table("sys_dict_type t")

	if param != nil {
		if param.DictName != "" {
			tb.Where("t.dict_name like ?", "%"+param.DictName+"%")
		}

		if param.DictType != "" {
			tb.Where("t.dict_type like ?", "%"+param.DictType+"%")
		}

		if param.Status != "" {
			tb.Where("t.status = ", param.Status)
		}
	}
	tb.Order(" t.dict_id desc ")
	var result []model.SysDictType
	tb = tb.Find(&result)
	return result, tb.Error
}
