package dao

import (
	"common/common_vo"
	"common/models"
	"errors"
	"github.com/lostvip-com/lv_framework/lv_db"
)

type DictDataDao struct {
}

var dictDataDao *DictDataDao

func GetDictDataDaoInstance() *DictDataDao {
	if dictDataDao == nil {
		dictDataDao = &DictDataDao{}
	}
	return dictDataDao
}

// 根据条件分页查询数据
func (dao *DictDataDao) FindPage(param *common_vo.SelectDictDataPageReq) (*[]models.SysDictData, int, error) {
	db := lv_db.GetOrmDefault()
	tb := db.Table("sys_dict_data t")
	if param != nil {
		if param.DictLabel != "" {
			tb.Where("t.dict_label like ?", "%"+param.DictLabel+"%")
		}

		if param.Status != "" {
			tb.Where("t.status = ?", param.Status)
		}

		if param.DictType != "" {
			tb.Where("t.dict_type like ?", "%"+param.DictType+"%")
		}
	}
	var total int64
	var result []models.SysDictData
	tb.Count(&total).Offset(param.GetStartNum()).Limit(param.GetPageSize()).Order("dict_sort asc").Find(&result)
	return &result, int(total), nil
}

// FindAll 获取所有数据
func (dao *DictDataDao) FindAll(dictLabel, dictType string) ([]models.SysDictData, error) {
	db := lv_db.GetOrmDefault()
	if db == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	tb := db.Table("sys_dict_data t ")

	if dictLabel != "" {
		tb.Where("t.dict_label like ?", "%"+dictLabel+"%")
	}
	tb.Where("t.status =? ", 0)

	if dictType != "" {
		tb.Where("t.dict_type =?", dictType)
	}
	tb.Order("dict_sort asc")
	var result []models.SysDictData
	err := tb.Find(&result).Error
	return result, err
}

// DeleteBatch 批量删除
func (d *DictDataDao) DeleteBatch(codes ...int) error {
	err := lv_db.GetOrmDefault().Where("dict_code in ?").Delete(codes).Error
	return err
}

// 批量删除
func (d *DictDataDao) DeleteByType(dictType string) error {
	err := lv_db.GetOrmDefault().Delete(&models.SysDictData{}, "dict_type=?", dictType).Error
	return err
}
