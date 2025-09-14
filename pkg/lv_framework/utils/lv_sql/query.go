package lv_sql

import (
	"errors"
	"fmt"
	"github.com/morrisxyang/xreflect"
	"github.com/spf13/cast"
	"reflect"
	"strings"
)

const (
	// FromQueryTag tag标记
	FromQueryTag = "lv_sql"
	// Mysql 数据库标识
	Mysql = "mysql"
	// Postgres 数据库标识
	Postgres = "postgres"
)

// ResolveSearchQuery 解析
/**
 * 	exact / iexact 等于
 * 	contains / icontains 包含
 *	gt / gte 大于 / 大于等于
 *	lt / lte 小于 / 小于等于
 *	startswith / istartswith 以…起始
 *	endswith / iendswith 以…结束
 *	in
 *	isnull
 *  order 排序		e.g. order[key]=desc     order[key]=asc
 */
func ResolveSearchQuery(driver string, q interface{}, condition Condition) {
	qType := reflect.TypeOf(q)
	qValue := reflect.ValueOf(q)
	var tag string
	var ok bool
	var t *resolveSearchTag
	for i := 0; i < qType.NumField(); i++ {
		tag, ok = "", false
		tag, ok = qType.Field(i).Tag.Lookup(FromQueryTag)
		if !ok {
			//递归调用
			ResolveSearchQuery(driver, qValue.Field(i).Interface(), condition)
			continue
		}
		switch tag {
		case "-":
			continue
		}
		t = makeTag(tag)
		if qValue.Field(i).IsZero() {
			continue
		}
		//解析
		switch t.Type {
		case "left":
			//左关联
			join := condition.SetJoinOn(t.Type, fmt.Sprintf(
				"left join `%s` on `%s`.`%s` = `%s`.`%s`",
				t.Join,
				t.Join,
				t.On[0],
				t.Table,
				t.On[1],
			))
			ResolveSearchQuery(driver, qValue.Field(i).Interface(), join)
		case "exact", "iexact":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` = ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "contains", "icontains":
			//fixme mysql不支持ilike
			if driver == Postgres && t.Type == "icontains" {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String() + "%"})
			} else {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String() + "%"})
			}
		case "gt":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` > ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "gte":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` >= ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "lt":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` < ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "lte":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` <= ?", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "startswith", "istartswith":
			if driver == Postgres && t.Type == "istartswith" {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{qValue.Field(i).String() + "%"})
			} else {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{qValue.Field(i).String() + "%"})
			}
		case "endswith", "iendswith":
			if driver == Postgres && t.Type == "iendswith" {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` ilike ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String()})
			} else {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` like ?", t.Table, t.Column), []interface{}{"%" + qValue.Field(i).String()})
			}
		case "in":
			condition.SetWhere(fmt.Sprintf("`%s`.`%s` in (?)", t.Table, t.Column), []interface{}{qValue.Field(i).Interface()})
		case "isnull":
			if !(qValue.Field(i).IsZero() && qValue.Field(i).IsNil()) {
				condition.SetWhere(fmt.Sprintf("`%s`.`%s` isnull", t.Table, t.Column), make([]interface{}, 0))
			}
		case "order":
			switch strings.ToLower(qValue.Field(i).String()) {
			case "desc", "asc":
				condition.SetOrder(fmt.Sprintf("`%s`.`%s` %s", t.Table, t.Column, qValue.Field(i).String()))
			}
		}
	}
}

func GetLimitSql(sql string, params interface{}) (string, error) {
	var pageNum, pageSize any
	paramType := reflect.TypeOf(params).Kind()
	if paramType == reflect.Map {
		paramMap := params.(map[string]interface{})
		pageNum = paramMap["pageNum"]
		pageSize = paramMap["pageSize"]
	} else {
		pageNum, _ = xreflect.FieldValue(params, "PageNum")
		pageSize, _ = xreflect.FieldValue(params, "PageSize")
	}
	if pageSize == nil || pageNum == nil {
		return sql, errors.New("PageSize or PageNum nil error ")
	}
	start := cast.ToInt64(pageSize) * (cast.ToInt64(pageNum) - 1)
	//sql = sql + " limit  " + cast.ToString(start) + "," + cast.ToString(pageSize)
	// 改为可以兼容mysql和postgresql的分页方式
	sql = sql + " limit  " + cast.ToString(pageSize) + " offset " + cast.ToString(start)
	return sql, nil
}

func GetCountSql(sql string) string {
	// 移除SQL中的ORDER BY子句，因为计数查询不需要排序
	sqlWithoutOrder := removeOrderBy(sql)
	return " select count(*)  from (" + sqlWithoutOrder + ") t "
}

// removeOrderBy 移除SQL中的ORDER BY子句，保留原始大小写
// removeOrderBy 移除SQL中的ORDER BY子句，保留原始大小写
func removeOrderBy(sql string) string {
	// 先转为大写用于查找位置，但不改变原始字符串
	sqlUpper := strings.ToUpper(sql)

	// 查找ORDER BY的位置（不区分大小写）
	orderbyIndex := strings.Index(sqlUpper, " ORDER BY ")
	if orderbyIndex == -1 {
		return sql // 如果没有找到ORDER BY，直接返回原SQL
	}

	// 查找ORDER BY之后是否有LIMIT子句
	limitIndex := strings.Index(sqlUpper[orderbyIndex+len(" ORDER BY "):], " LIMIT ")
	if limitIndex != -1 {
		// 如果有LIMIT，只移除ORDER BY部分到LIMIT之前
		limitStart := orderbyIndex + len(" ORDER BY ") + limitIndex
		// 找到LIMIT前一个ORDER BY部分的结束位置
		//orderbyPart := sql[orderbyIndex:limitStart]
		// 保留LIMIT及其后的内容
		return sql[:orderbyIndex] + sql[limitStart:]
	} else {
		// 如果没有LIMIT，移除ORDER BY及其后所有内容
		return sql[:orderbyIndex]
	}
}
