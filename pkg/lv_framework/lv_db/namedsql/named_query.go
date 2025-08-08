// ///////////////////////////////////////////////////////////////////////////
// 业务逻辑处理类的基类，简单的直接在model中处理即可，不需要service
//
// //////////////////////////////////////////////////////////////////////////
package namedsql

import (
	"database/sql"
	"fmt"
	"github.com/lostvip-com/lv_framework/lv_global"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_reflect"
	"github.com/lostvip-com/lv_framework/utils/lv_sql"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"time"
)

func Exec(db *gorm.DB, dmlSql string, req map[string]any) (int64, error) {
	if lv_global.IsDebug {
		db = db.Debug()
	}
	if strings.Contains(dmlSql, "@") {
		kvMap, isMap := checkAndExtractMap(req)
		if isMap {
			req = kvMap
		}
		tx := db.Exec(dmlSql, req)
		return tx.RowsAffected, tx.Error
	} else {
		tx := db.Exec(dmlSql)
		return tx.RowsAffected, tx.Error
	}
}

func GetOneMapByNamedSql(db *gorm.DB, limitSql string, req any, isCamel bool) (result *map[string]any, err error) {
	list, err := ListMap(db, limitSql, req, isCamel)
	if err == nil {
		var mpList = *list
		if mpList == nil || len(mpList) == 0 {
			err = gorm.ErrRecordNotFound
		} else {
			result = &mpList[0]
		}
	}
	return result, err
}

func toCamelMap(result *map[string]any) *map[string]any {
	mp := make(map[string]any)
	for k, v := range *result {
		mp[cast.ToString(lv_sql.ToCamel(k))] = v
	}
	return &mp
}

/**
 * 通用泛型查询
 */
func ListData[T any](db *gorm.DB, limitSql string, req any) (*[]T, error) {
	var list = make([]T, 0)
	var err error
	if lv_global.IsDebug {
		db = db.Debug()
	}
	if strings.Contains(limitSql, "@") {
		kvMap, isMap := checkAndExtractMap(req)
		if isMap {
			req = kvMap
		}
		err = db.Raw(limitSql, req).Scan(&list).Error
	} else {
		err = db.Raw(limitSql).Scan(&list).Error
	}

	return &list, err
}

// ListData2Map 通用泛型查询
func ListData2Map[T any](db *gorm.DB, limitSql string, req any, mapKey string) (*map[string]T, error) {
	listPtr, err := ListData[T](db, limitSql, req)
	if err != nil {
		return nil, err
	}
	mp := make(map[string]T)
	list := *listPtr
	for i := range list {
		it := list[i]
		value, ok := lv_reflect.GetFieldValueSimple(it, mapKey)
		if ok {
			mp[cast.ToString(value)] = it
		} else {
			lv_log.Warn("mapKey not found", mapKey)
		}
	}
	return &mp, err
}
func Count(db *gorm.DB, countSql string, params any) (int64, error) {
	if lv_global.IsDebug {
		db = db.Debug()
	}

	if !strings.Contains(countSql, "count") {
		countSql = " select count(*) from (" + countSql + ") t where 1=1  "
	}
	if !strings.Contains(countSql, "limit") {
		countSql = countSql + "   limit 1  "
	}

	var rows *sql.Rows
	var err error
	if strings.Contains(countSql, "@") {
		kvMap, isMap := checkAndExtractMap(params)
		if isMap {
			params = kvMap
		}
		rows, err = db.Raw(countSql, params).Rows()
	} else {
		rows, err = db.Raw(countSql).Rows()
	}
	if err != nil {
		lv_log.Info(err)
		return 0, err
	}
	//查总数
	var count int64
	if rows != nil {
		for rows.Next() {
			rows.Scan(&count)
		}
	}
	return count, err
}

/**
 * gorm中参数为map指针时，无法正常传参数！！
 * 处理方式：把map的指针转为值类型。
 */
func checkAndExtractMap(value interface{}) (map[string]any, bool) {
	// 判断是否是指针类型
	if ptr, ok := value.(*map[string]any); ok {
		// 指针指向Map类型
		return *ptr, true
	}
	return nil, false
}

func ListMapAny(db *gorm.DB, sqlQuery string, params any, isCamel bool) (*[]map[string]any, error) {
	return ListMap(db, sqlQuery, params, isCamel)
}

// ListMap sql查询返回map isCamel key是否按驼峰式命名,有些数据会出现2进制输出
func ListMap(db *gorm.DB, sqlQuery string, params any, isCamel bool) (*[]map[string]any, error) {
	// Validate inputs
	if db == nil {
		return nil, fmt.Errorf("db cannot be nil")
	}
	if sqlQuery == "" {
		return nil, fmt.Errorf("sqlQuery cannot be empty")
	}
	// Debug mode
	if lv_global.IsDebug {
		db = db.Debug()
	}
	// 1. Execute query
	var rows *sql.Rows
	var err error

	if strings.Contains(sqlQuery, "@") {
		kvMap, isMap := checkAndExtractMap(params)
		if isMap {
			rows, err = db.Raw(sqlQuery, kvMap).Rows()
		} else {
			// Try with original params if not a map
			rows, err = db.Raw(sqlQuery, params).Rows()
		}
	} else {
		rows, err = db.Raw(sqlQuery).Rows()
	}
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	if rows == nil {
		return nil, fmt.Errorf("query returned nil rows")
	}
	defer rows.Close()

	// 2. Get column information
	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to get column types: %w", err)
	}
	// 3. Initialize scan buffers
	values := make([]interface{}, len(cols))
	scanArgs := make([]interface{}, len(cols))

	for i, colType := range colTypes {
		// Handle nullable types
		nullable, nullableOk := colType.Nullable()
		dbType := strings.ToUpper(colType.DatabaseTypeName())

		switch {
		case strings.Contains(dbType, "INT"):
			if nullable && nullableOk {
				var n sql.NullInt64
				values[i] = &n
			} else {
				var n int64
				values[i] = &n
			}
		case strings.Contains(dbType, "FLOAT") || strings.Contains(dbType, "DOUBLE") || strings.Contains(dbType, "DECIMAL"):
			if nullable && nullableOk {
				var f sql.NullFloat64
				values[i] = &f
			} else {
				var f float64
				values[i] = &f
			}
		case strings.Contains(dbType, "BOOL"):
			if nullable && nullableOk {
				var b sql.NullBool
				values[i] = &b
			} else {
				var b bool
				values[i] = &b
			}
		case strings.Contains(dbType, "DATE") || strings.Contains(dbType, "TIME"):
			if nullable && nullableOk {
				var t sql.NullTime
				values[i] = &t
			} else {
				var t time.Time
				values[i] = &t
			}
		case strings.Contains(dbType, "BLOB") || strings.Contains(dbType, "BINARY"):
			var blob []byte
			values[i] = &blob
		default:
			// Handle string types and unknowns
			if nullable && nullableOk {
				var s sql.NullString
				values[i] = &s
			} else {
				var s string
				values[i] = &s
			}
		}
		scanArgs[i] = values[i]
	}
	// 4. Process result set
	result := make([]map[string]any, 0)
	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}

		rowData := make(map[string]any)
		for i, colName := range cols {
			val := reflect.Indirect(reflect.ValueOf(values[i])).Interface()
			// Handle NULL values
			switch v := val.(type) {
			case sql.NullString:
				if v.Valid {
					val = v.String
				} else {
					val = nil
				}
			case sql.NullInt64:
				if v.Valid {
					val = v.Int64
				} else {
					val = nil
				}
			case sql.NullFloat64:
				if v.Valid {
					val = v.Float64
				} else {
					val = nil
				}
			case sql.NullBool:
				if v.Valid {
					val = v.Bool
				} else {
					val = nil
				}
			case sql.NullTime:
				if v.Valid {
					val = v.Time.Format("2006-01-02 15:04:05")
				} else {
					val = nil
				}
			case time.Time:
				val = v.Format("2006-01-02 15:04:05")
			}

			// Handle column name case
			key := colName
			if isCamel {
				key = lv_sql.ToCamel(key)
			}
			rowData[key] = val
		}
		result = append(result, rowData)
	}
	// 5. Check for iteration errors
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration failed: %w", err)
	}

	return &result, nil
}

func ListArrStr(db *gorm.DB, sqlQuery string, params any) (*[][]string, error) {
	if lv_global.IsDebug {
		db = db.Debug()
	}
	var rows *sql.Rows
	var err error
	if strings.Contains(sqlQuery, "@") {
		kvMap, isMap := checkAndExtractMap(params)
		if isMap {
			params = kvMap
		}
		rows, err = db.Raw(sqlQuery, params).Rows()
	} else {
		rows, err = db.Raw(sqlQuery).Rows()
	}
	if err != nil {
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	listRows := make([][]string, 0)
	for rows.Next() {
		row := make([]string, 0)
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		var value string
		for _, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			row = append(row, value)
		}
		listRows = append(listRows, row)
	}
	return &listRows, err
}

// ListOneColStr 查询某一列，放到数组中
func ListOneColStr(db *gorm.DB, sqlQuery string, params any) ([]string, error) {
	if lv_global.IsDebug {
		db = db.Debug()
	}
	var rows *sql.Rows
	var err error
	if strings.Contains(sqlQuery, "@") {
		kvMap, isMap := checkAndExtractMap(params)
		if isMap {
			params = kvMap
		}
		rows, err = db.Raw(sqlQuery, params).Rows()
	} else {
		rows, err = db.Raw(sqlQuery).Rows()
	}
	if err != nil {
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	arr := make([]string, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		for _, col := range values {
			if col != nil {
				arr = append(arr, string(col))
			}
		}
	}
	return arr, err
}
