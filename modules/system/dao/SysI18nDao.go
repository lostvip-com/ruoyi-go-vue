// ==========================================================================
// LV自动生成model扩展代码列表、增、删，改、查、导出，只生成一次，按需修改,再次生成不会覆盖.
// 生成日期：2025-08-11 07:41:35 &#43;0000 UTC
// 生成人：dpc
// ==========================================================================
package dao

import (
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_batis"
	"github.com/lostvip-com/lv_framework/lv_db/lv_dao"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"system/model"
	"system/vo"
)

type SysI18nDao struct{}

var i18nDao *SysI18nDao

func GetSysI18nDaoInstance() *SysI18nDao {
	if i18nDao == nil {
		i18nDao = &SysI18nDao{}
	}
	return i18nDao
}

// ListMapByPage 根据条件分页查询数据
func (d SysI18nDao) ListMapByPage(req *vo.SysI18nReq) (*[]map[string]any, int64, error) {
	ibatis := lv_batis.NewInstance("system/sys_i18n_mapper.sql") //under the mapper directory
	// 约定用方法名ListByPage对应sql文件中的同名tagName
	db := lv_db.GetOrmDefault()
	limitSQL, countSql, err := ibatis.GetPageSql("ListSysI18n", req)
	//查询数据
	rows, err := lv_dao.ListMapByNamedSql(db, limitSQL, req, true)
	lv_err.HasErrAndPanic(err)
	count, err := lv_dao.CountByNamedSql(db, countSql, req)
	lv_err.HasErrAndPanic(err)
	return rows, count, nil
}

// ListByPage 根据条件分页查询数据
func (d SysI18nDao) ListByPage(req *vo.SysI18nReq) (*[]vo.SysI18nResp, int64, error) {
	ibatis := lv_batis.NewInstance("system/sys_i18n_mapper.sql") //under the mapper directory
	// 对应sql文件中的同名tagName
	db := lv_db.GetOrmDefault()
	limitSQL, countSql, err := ibatis.GetPageSql("ListSysI18n", req)
	//查询数据
	rows, err := lv_dao.ListByNamedSql[vo.SysI18nResp](db, limitSQL, req)
	lv_err.HasErrAndPanic(err)
	count, err := lv_dao.CountByNamedSql(db, countSql, req)
	lv_err.HasErrAndPanic(err)
	return rows, count, nil
}

// ListAll 导出excel使用
func (d SysI18nDao) ListAll(req *vo.SysI18nReq, isCamel bool) (*[]map[string]any, error) {
	ibatis := lv_batis.NewInstance("system/sys_i18n_mapper.sql")
	// 约定用方法名ListByPage对应sql文件中的同名tagName
	sql, err := ibatis.GetSql("ListSysI18n", req)
	lv_err.HasErrAndPanic(err)

	arr, err := lv_dao.ListMapByNamedSql(lv_db.GetOrmDefault(), sql, req, isCamel)
	return arr, err
}

// Find 根据条件查询
func (d SysI18nDao) Find(where, order string) (*[]model.SysI18n, error) {
	var list []model.SysI18n
	err := lv_db.GetOrmDefault().Table("sys_i18n").Where(where).Order(order).Find(&list).Error
	return &list, err
}

// Find 通过主键批量删除
func (d SysI18nDao) DeleteByIds(ida []int) (int64, error) {
	db := lv_db.GetOrmDefault().Table("sys_i18n").Where("id in ? ", ida).Update("del_flag", 1)
	return db.RowsAffected, db.Error
}
