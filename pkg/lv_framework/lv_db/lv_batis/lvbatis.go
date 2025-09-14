// //////////////////////////////////////////////////////////////////////////////////////
// from Dotsql
// It is not an ORM, it is not a query builder.
// Dotsql is a library that helps you keep sql files in one place and use it with ease.
// For more usage examples see https://github.com/qustavo/dotsql
// /////////////////////////////////////////////////////////////////////////////////////
package lv_batis

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_file"
	"github.com/lostvip-com/lv_framework/utils/lv_reflect"
	"github.com/lostvip-com/lv_framework/utils/lv_tpl"
	"github.com/morrisxyang/xreflect"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"io"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// Execer is an interface used by Exec.
type Execer interface {
	Exec(query string, args ...interface{}) *gorm.DB
}

// ExecerContext is an interface used by ExecContext.
type ExecerContext interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type LvBatis struct {
	Queries     map[string]string
	Vars        map[string]map[string]any
	TplFile     string
	CurrBaseSql string
}

/**
 * 从mapper目录解析sql文件
 */
func NewInstance(relativePath string) *LvBatis {
	basePath, _ := os.Getwd()
	absolutePath := basePath + "/resources/mapper" //为了方便管理，必须把映射文件放到mapper目录
	if strings.HasPrefix(relativePath, "/") {
		absolutePath = absolutePath + relativePath
	} else {
		absolutePath = absolutePath + "/" + relativePath
	}
	dot, err := LoadFromFile(absolutePath)
	dot.TplFile = relativePath
	if err != nil {
		panic(err)
	}
	return dot
}

func (d *LvBatis) GetSql(tagName string, params interface{}) (string, error) {
	query, err := d.LookupQuery(tagName)
	if err != nil || query == "" {
		panic("tpl文件格式错误!")
	}
	//动态解析
	sql, err := lv_tpl.ParseTemplateStr(query, params)
	if sql == "" || err != nil {
		lv_log.Error(err)
		panic(d.getTplFile() + " 可能存在错误：<p/>1.使用了参数对象中不存在的属性<p/>2.template语法错误！")
	}
	d.CurrBaseSql = sql //缓存当前正在执行的分页sql
	return sql, err
}

/**
 * 从mapper目录解析sql文件
 */

/**
 * 从mapper目录解析sql文件
 */
func (d *LvBatis) GetLimitSqlParams(tagName string, params interface{}) (string, map[string]any, error) {
	var pageNum, pageSize any
	paramType := reflect.TypeOf(params).Kind()
	sqlParams := d.Vars[tagName]
	if paramType == reflect.Map {
		paramMap := params.(map[string]interface{})
		pageNum = paramMap["pageNum"]
		pageSize = paramMap["pageSize"]
		for key, value := range paramMap { //合并参数
			sqlParams[key] = value // 覆盖或新增键值对
		}
	} else {
		pageNum, _ = xreflect.FieldValue(params, "PageNum")
		pageSize, _ = xreflect.FieldValue(params, "PageSize")
		lv_reflect.CopyProperties2Map(params, sqlParams) //合并参数
	}
	if pageSize == nil || pageNum == nil {
		return "", nil, errors.New("pageSize and pageNum can not be empty! ")
	}
	sql, err := d.GetSql(tagName, sqlParams)
	start := cast.ToInt64(pageSize) * (cast.ToInt64(pageNum) - 1)
	sql = sql + " limit  " + cast.ToString(start) + "," + cast.ToString(pageSize)
	return sql, sqlParams, err
}

func (d *LvBatis) GetLimitSql(tagName string, params interface{}) (string, error) {
	var pageNum, pageSize int
	paramType := reflect.TypeOf(params).Kind()
	sqlParams := d.Vars[tagName]
	if paramType == reflect.Map {
		paramMap := params.(map[string]interface{})
		pNum := paramMap["pageNum"]
		pSize := paramMap["pageSize"]
		pageNum  = cast.ToInt(pNum)
		pageSize = cast.ToInt(pSize)
		for key, value := range paramMap { //合并参数
			sqlParams[key] = value // 覆盖或新增键值对
		}
	} else {
		pNum, _ := xreflect.FieldValue(params, "PageNum")
		pSize, _ := xreflect.FieldValue(params, "PageSize")
		pageNum  = cast.ToInt(pNum)
		pageSize = cast.ToInt(pSize)
		lv_reflect.CopyProperties2Map(params, sqlParams) //合并参数
	}

	if pageSize==0{
		pageSize = 1000;
	}
	if pageNum == 0 {
		pageNum = 1
	}
	if d.CurrBaseSql == "" {
		sql,err := d.GetSql(tagName, sqlParams)
		if err!=nil{
			return "",err
		}
		d.CurrBaseSql = sql
	}
	start := cast.ToInt64(pageSize) * (cast.ToInt64(pageNum) - 1)
	//sql = sql + " limit  " + cast.ToString(start) + "," + cast.ToString(pageSize)
	// 改为可以兼容mysql和postgresql的分页方式
	sql := d.CurrBaseSql + " limit  " + cast.ToString(pageSize) + " offset " + cast.ToString(start)
	return sql, nil
}

func (d *LvBatis) GetCountSql(tagName string, sqlParams interface{}) (string, error) {
	if d.CurrBaseSql == "" {
		sql,err := d.GetSql(tagName, sqlParams)
		if err!=nil{
			return "",err
		}
		d.CurrBaseSql = sql
	}
	index := strings.Index(d.CurrBaseSql, " order ")
	noOrderSql := d.CurrBaseSql
	if index >20  { // select * from t where order by
		noOrderSql = d.CurrBaseSql[:index]
	}else{
		noOrderSql = d.CurrBaseSql
	}
	sql := " select count(*)  from (" + noOrderSql + ") t "
	return sql, nil
}

func (d *LvBatis) GetPageSql(tagName string, sqlParams any) (string, string, error) {
	countSql, err := d.GetCountSql(tagName,sqlParams)
	if err != nil {
		return "", "", err
	}
	limitSql, err := d.GetLimitSql(tagName, sqlParams)
	return limitSql, countSql, err
}

func (d *LvBatis) LookupQuery(name string) (query string, err error) {
	query, ok := d.Queries[name]
	if !ok {
		err = fmt.Errorf("sql: '%s' could not be found", name)
	}

	return
}

// Exec 默认只能执行单条sql，除非mysql配置为allowMultiQueries=true
func (d *LvBatis) Exec(db Execer, name string, args ...interface{}) (*gorm.DB, error) {
	query, err := d.LookupQuery(name)
	if err != nil {
		return nil, err
	}
	return db.Exec(query, args...), err
}

// ExecMultiSqlInTransaction 在事务中执行多个SQL语句
func (d *LvBatis) ExecMultiSqlInTransaction(db *gorm.DB, name string, args ...interface{}) (*gorm.DB, error) {
	query, err := d.LookupQuery(name)
	if err != nil {
		return nil, err
	}
	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var result *gorm.DB
	// 按分号拆分SQL语句
	statements := strings.Split(query, ";")
	for _, stmt := range statements {
		trimmedStmt := strings.TrimSpace(stmt)
		if trimmedStmt != "" {
			result = tx.Exec(trimmedStmt, args...)
			if result.Error != nil {
				tx.Rollback()
				return nil, result.Error
			}
		}
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return result, nil
}

// ExecContext is a wrapper for database/sql's ExecContext(), using dotsql named query.
func (d *LvBatis) ExecContext(ctx context.Context, db ExecerContext, name string, args ...interface{}) (sql.Result, error) {
	query, err := d.LookupQuery(name)
	if err != nil {
		return nil, err
	}

	return db.ExecContext(ctx, query, args...)
}

// GetRawSql returns the query, everything after the --name tag
func (d *LvBatis) GetRawSql(name string) (string, error) {
	return d.LookupQuery(name)
}

// GetQueryMap returns a map[string]string of loaded Queries
func (d *LvBatis) GetQueryMap() map[string]string {
	return d.Queries
}

func (d *LvBatis) getTplFile() string {
	return d.TplFile
}

// Load imports sql Queries from any io.Reader.
func Load(r io.Reader) (*LvBatis, error) {
	scanner := &Scanner{}
	queries := scanner.Run(bufio.NewScanner(r))
	varMap := parseVarName(queries)
	dotSql := &LvBatis{
		Queries: queries,
		Vars:    varMap,
	}

	return dotSql, nil
}

func parseVarName(funSql map[string]string) map[string]map[string]any {
	//re := regexp.MustCompile(`\.\w+`)
	// 使用正则匹配模板中的变量
	re := regexp.MustCompile(`{{[^{}]*\.\s*([a-zA-Z_][a-zA-Z0-9_]*)[^{}]*}}`)
	mp := make(map[string]map[string]any)
	for funcKey, sql := range funSql {
		matches := re.FindAllStringSubmatch(sql, -1)
		varMap := make(map[string]any) //变量去重
		mp[funcKey] = varMap
		for _, match := range matches {
			if len(match) > 1 {
				key := match[1]
				varMap[key] = nil
			}
		} //end for

	} //end for
	return mp
}

// LoadFromFile imports SQL Queries from the file.
func LoadFromFile(sqlFile string) (*LvBatis, error) {
	if !lv_file.IsFileExist(sqlFile) {
		return nil, errors.New("file not be found! " + sqlFile)
	}
	f, err := os.Open(sqlFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Load(f)
}

// LoadFromString imports SQL Queries from the string.
func LoadFromString(sql string) (*LvBatis, error) {
	buf := bytes.NewBufferString(sql)
	return Load(buf)
}

// Merge takes one or mord *LvBatis and merge its Queries
// It's in-order, so the last source will override Queries with the same name
// in the previous arguments if any.
func Merge(dots ...*LvBatis) *LvBatis {
	queries := make(map[string]string)

	for _, dot := range dots {
		for k, v := range dot.GetQueryMap() {
			queries[k] = v
		}
	}

	return &LvBatis{
		Queries: queries,
	}
}
