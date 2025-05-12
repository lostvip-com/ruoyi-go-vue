// ///////////////////////////////////////////////////////////////////////////
// 业务逻辑处理类的基类，简单的直接在model中处理即可，不需要service
//
// //////////////////////////////////////////////////////////////////////////
package namedsql

import (
	"database/sql"
	"github.com/lostvip-com/lv_framework/lv_global"
	"github.com/lostvip-com/lv_framework/lv_log"
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

// ListMap sql查询返回map isCamel key是否按驼峰式命名,有些数据会出现2进制输出
func ListMap(db *gorm.DB, sqlQuery string, params any, isCamel bool) (*[]map[string]any, error) {
	// 1. 执行查询
	var rows *sql.Rows
	var err error
	if lv_global.IsDebug {
		db = db.Debug()
	}
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
	defer rows.Close()

	// 2. 获取列信息
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	// 3. 初始化扫描缓冲区
	values := make([]interface{}, len(cols))
	scanArgs := make([]interface{}, len(cols))
	for i, colType := range colTypes {
		switch colType.DatabaseTypeName() {
		case "DATE", "DATETIME", "TIMESTAMP", "TIME":
			var t time.Time
			values[i] = &t
		case "INT", "BIGINT", "SMALLINT", "TINYINT":
			var n int64
			values[i] = &n
		case "FLOAT", "DOUBLE", "DECIMAL":
			var f float64
			values[i] = &f
		case "BOOLEAN", "BOOL":
			var b bool
			values[i] = &b
		case "BLOB", "BINARY", "VARBINARY":
			var blob []byte
			values[i] = &blob
		default:
			// 默认处理文本和其他类型
			var s string
			values[i] = &s
		}
		scanArgs[i] = values[i]
	}

	// 4. 遍历结果集
	result := make([]map[string]any, 0)
	for rows.Next() {
		if err = rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		rowData := make(map[string]any)
		for i, colName := range cols {
			val := reflect.Indirect(reflect.ValueOf(values[i])).Interface()

			// 特殊处理时间类型
			if t, ok := val.(time.Time); ok {
				val = t.Format("2006-01-02 15:04:05")
			}

			// 处理列名驼峰转换
			key := colName
			if isCamel {
				key = lv_sql.ToCamel(key)
			}
			rowData[key] = val
		}
		result = append(result, rowData)
	}

	// 5. 检查遍历错误
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListMapStr 所有数据转为字符串格式返回，params可是是strtuc指针或map,isCamel key是否按驼峰式命名
func ListMapStr(db *gorm.DB, sqlQuery string, params any, isCamel bool) (*[]map[string]string, error) {
	var rows *sql.Rows
	var err error
	if lv_global.IsDebug {
		db = db.Debug()
	}
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
		lv_log.Info(err)
		return nil, err
	}
	cols, err := rows.Columns()
	if err != nil {
		lv_log.Info(err)
		return nil, err
	}
	result := make([]map[string]string, 0)
	values := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			lv_log.Info(err)
			return nil, err
		}
		var value string
		resultC := map[string]string{}
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			colKey := cols[i]
			if isCamel {
				colKey = lv_sql.ToCamel(colKey)
			}
			resultC[colKey] = value
		}
		result = append(result, resultC)
	}
	return &result, err
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
