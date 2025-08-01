package service

import (
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	dao2 "system/dao"
	"system/model"
)

type TableColumnService struct {
}

// Insert 新增业务字段
func (svc TableColumnService) Insert(entity *model.GenTableColumn) (int64, error) {
	err := entity.Save()
	if err != nil {
		return 0, err
	}
	return entity.ColumnId, err
}

// Update  修改业务字段
func (svc TableColumnService) Update(entity *model.GenTableColumn) error {
	return entity.Updates()
}

// FindById 根据主键查询数据
func (svc TableColumnService) FindById(id int64) (*model.GenTableColumn, error) {
	entity := &model.GenTableColumn{ColumnId: id}
	_, err := entity.FindOne()
	return entity, err
}

// DeleteById 根据主键删除数据
func (svc TableColumnService) DeleteById(id int64) bool {
	err := (&model.GenTableColumn{ColumnId: id}).Delete()
	if err == nil {
		return true
	}
	return false
}

// DeleteByIds 批量删除数据记录
func (svc TableColumnService) DeleteByIds(ids string) error {
	idarr := lv_conv.ToInt64Array(ids, ",")
	err := lv_db.GetOrmDefault().Exec("delete from gen_table_column where column_id in ? ", idarr).Error
	return err
}

// 查询业务字段列表
func (svc TableColumnService) SelectGenTableColumnListByTableId(tableId int64) ([]model.GenTableColumn, error) {
	var tool dao2.GenTableColumnDao
	return tool.SelectGenTableColumnListByTableId(tableId)
}
