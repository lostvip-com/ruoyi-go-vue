package dao

import (
	"common/common_vo"
	"common/models"
	"errors"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/namedsql"
	"github.com/lostvip-com/lv_framework/utils/lv_office"
	"xorm.io/builder"
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
func (dao *DictDataDao) FindPage(param *common_vo.SelectDictDataPageReq) (*[]models.SysDictData, int64, error) {
	db := lv_db.GetMasterGorm()
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
	return &result, total, nil
}

// 导出excel
func (dao *DictDataDao) FindListExport(param *common_vo.SelectDictDataPageReq, head, col []string) (string, error) {
	db := lv_db.GetMasterGorm()
	build := builder.Select(col...).From("sys_dict_data", "t")
	if param != nil {
		if param.DictLabel != "" {
			build.Where(builder.Like{"t.dict_label", param.DictLabel})
		}
		if param.Status != "" {
			build.Where(builder.Eq{"t.status": param.Status})
		}
		if param.DictType != "" {
			build.Where(builder.Like{"t.dict_type", param.DictType})
		}
	}
	sqlStr, _ := build.ToBoundSQL()
	arr, err := namedsql.ListArrStr(db, sqlStr, nil)
	path, err := lv_office.DownlaodExcel(head, *arr)
	return path, err
}

// FindAll 获取所有数据
func (dao *DictDataDao) FindAll(dictLabel, dictType string) ([]models.SysDictData, error) {
	db := lv_db.GetMasterGorm()
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
func (d *DictDataDao) DeleteBatch(codes ...int64) error {
	err := lv_db.GetMasterGorm().Where("dict_code in ?").Delete(codes).Error
	return err
}

// 批量删除
func (d *DictDataDao) DeleteByType(dictType string) error {
	err := lv_db.GetMasterGorm().Delete(&models.SysDictData{}, "dict_type=?", dictType).Error
	return err
}
