package test

import (
	"fmt"
	"github.com/lostvip-com/lv_framework/utils/lv_sql"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"testing"
)

func TestSQL(t *testing.T) {
	sql := "select * from test Order By name asc "
	sqlLimit, _ := lv_sql.GetLimitSql(sql, lv_dto.Paging{PageNum: 1, PageSize: 10})
	fmt.Println("============" + sqlLimit)
	sql = lv_sql.GetCountSql(sql)
	fmt.Println("============" + sql)
}
