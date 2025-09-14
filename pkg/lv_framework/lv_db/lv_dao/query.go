package lv_dao

import (
	"context"
	"errors"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_batis"
	"github.com/lostvip-com/lv_framework/lv_db/namedsql"
	"gorm.io/gorm"
	"time"
)

func CountColumnDelFlag0(db *gorm.DB, table, column, value string) (int64, error) {
	var total int64
	err := db.Table(table).Where("del_flag=0 and "+column+"=?", value).Count(&total).Error
	return total, err
}
func CountColumnAll(db *gorm.DB, table, column, value string) (int64, error) {
	var total int64
	err := db.Table(table).Where(column+"=?", value).Count(&total).Error
	return total, err
}
func ListMapByNamedSql(db *gorm.DB, sql string, req any, isCamel bool) (*[]map[string]any, error) {
	return namedsql.ListMap(db, sql, req, isCamel)
}
func ListMap2Map(db *gorm.DB, limitSql string, req any, mapKey string, isCamel bool) (*map[string]map[string]any, error) {
	return namedsql.ListMap2Map(db, limitSql, req, mapKey, isCamel)
}

func ListData2Map[T any](db *gorm.DB, limitSql string, req any, mapKey string) (*map[string]T, error) {
	return namedsql.ListData2Map[T](db, limitSql, req, mapKey)
}
func GetOneMapByNamedSql(db *gorm.DB, sql string, req any, isCamel bool) (*map[string]any, error) {
	d := lv_db.GetOrmDefault()
	mp, err := namedsql.GetOneMapByNamedSql(d, sql, req, isCamel)
	if err != nil {
		return nil, err
	}
	return mp, err
}

// ListByNamedSql 通用泛型查询
func ListByNamedSql[T any](db *gorm.DB, sql string, req any) (*[]T, error) {
	return namedsql.ListData[T](lv_db.GetOrmDefault(), sql, req)
}

func CountByNamedSql(db *gorm.DB, sql string, params any) (int64, error) {
	return namedsql.Count(db, sql, params)
}

/**
 * 通用泛型查询
 */
func GetPageByNamedSql[T any](db *gorm.DB, sqlfile string, sqlTag string, req any) (*[]T, int64, error) {
	//解析sql
	ibatis := lv_batis.NewInstance(sqlfile)
	limitSql, countSql, err := ibatis.GetPageSql(sqlTag, req)
	if err != nil {
		return nil, 0, err
	}
	rows, err := namedsql.ListData[T](db, limitSql, req)
	if err != nil {
		return nil, 0, err
	}
	count, err := namedsql.Count(db, countSql, req)
	if err != nil {
		return nil, 0, err
	}
	return rows, count, err
}

func DeleteByIds(db *gorm.DB, tableName, column string, ids []int64) error {
	var total int64
	err := db.Table(tableName).Select("count(*)").Where("id in ? ", ids).Find(&total).Error
	if total == 0 {
		return errors.New("no data found ")
	}
	delSql := "delete from " + tableName + " where id in ?"
	err = db.Exec(delSql, ids).Error
	return err
}

func Transaction(db *gorm.DB, timeout time.Duration, fn func(tx *gorm.DB) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := db.WithContext(ctx).Transaction(fn)
	return err
}
