package lv_dao

import (
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_batis"
	"github.com/lostvip-com/lv_framework/lv_db/namedsql"
	"github.com/lostvip-com/lv_framework/lv_global"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
)

func CountColumnDelFlag0(table, column, value string) (int64, error) {
	var total int64
	err := lv_db.GetOrmDefault().Table(table).Where("del_flag=0 and "+column+"=?", value).Count(&total).Error
	return total, err
}
func CountColumnAll(table, column, value string) (int64, error) {
	var total int64
	err := lv_db.GetOrmDefault().Table(table).Where(column+"=?", value).Count(&total).Error
	return total, err
}
func ListMapByNamedSql(sql string, req any, isCamel bool) (*[]map[string]any, error) {
	d := lv_db.GetOrmDefault()
	return namedsql.ListMap(d, sql, req, isCamel)
}

func GetOneMapByNamedSql(sql string, req any, isCamel bool) (*map[string]any, error) {
	d := lv_db.GetOrmDefault()
	mp, err := namedsql.GetOneMapByNamedSql(d, sql, req, isCamel)
	if err != nil {
		return nil, err
	}
	return mp, err
}

func ListByNamed(sql string, req any, isCamel bool) (*[]map[string]any, error) {
	d := lv_db.GetOrmDefault()
	return namedsql.ListMap(d, sql, req, isCamel)
}

/**
 * 通用泛型查询
 */
func ListByNamedSql[T any](sql string, req any) (*[]T, error) {
	return namedsql.ListData[T](lv_db.GetOrmDefault(), sql, req)
}

func CountByNamedSql(sql string, params any) (int64, error) {
	return namedsql.Count(lv_db.GetOrmDefault(), sql, params)
}

/**
 * 通用泛型查询
 */
func GetPageByNamedSql[T any](sqlfile string, sqlTag string, req any) (*[]T, int64, error) {
	//解析sql
	ibatis := lv_batis.NewInstance(sqlfile)
	sql, err := ibatis.GetLimitSql(sqlTag, req)
	lv_err.HasErrAndPanic(err)
	//查询数据
	rows, err := namedsql.ListData[T](lv_db.GetOrmDefault(), sql, req)
	lv_err.HasErrAndPanic(err)
	count, err := namedsql.Count(lv_db.GetOrmDefault(), sql, req)
	return rows, count, err
}

func DeleteByIds(tableName string, ids []int64) error {
	delSql := "delete from " + tableName + " where id in (?)"
	err := lv_db.GetOrmDefault().Exec(delSql, ids).Error
	return err
}

func DeleteSoftByIds(tableName string, ids []int64) error {
	delSql := "update " + tableName + "set del_flag=? where id in (?)"
	err := lv_db.GetOrmDefault().Exec(delSql, lv_global.FLAG_DEL_YES, ids).Error
	return err
}
