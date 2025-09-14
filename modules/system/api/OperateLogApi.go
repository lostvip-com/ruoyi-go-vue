package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_json"
	"github.com/spf13/cast"
	"strings"
	"system/service"
	"system/vo"
)

type OperateLogApi struct {
}

// ListAjax 用户列表分页数据
func (w *OperateLogApi) ListAjax(c *gin.Context) {
	var req vo.OperLogPageReq

	err := c.ShouldBind(&req)
	lv_err.HasErrAndPanic(err)
	var operlogService service.OperLogService
	rows, total, err := operlogService.FindPage(&req)
	lv_err.HasErrAndPanic(err)
	util.SuccessPage(c, rows, total)
}

func (w *OperateLogApi) Export(c *gin.Context) {
	var req vo.OperLogPageReq
	err := c.ShouldBind(&req)
	lv_err.HasErrAndPanic(err)
	var operlogService service.OperLogService
	rows, _, err := operlogService.FindPage(&req)
	lv_err.HasErrAndPanic(err)

	headerMap := []map[string]string{
		{"key": "title", "title": "模块标题", "width": "15"},
		{"key": "businessType", "title": "业务类型", "width": "10"},
		{"key": "method", "title": "方法名称", "width": "10"},
		{"key": "RequestMethod", "title": "请求方式", "width": "10"},
		{"key": "operatorType", "title": "操作类别", "width": "10"},
		{"key": "operUrl", "title": "请求URL", "width": "10"},
		{"key": "operLocation", "title": "操作地点", "width": "10"},
		{"key": "status", "title": "操作状态", "width": "10"},
		{"key": "operTime", "title": "操作时间", "width": "10"},
	}
	ex := util.NewMyExcel()
	listMap := make([]map[string]any, 0)
	for i := range *rows {
		it := lv_json.StructToMap((*rows)[i])
		listMap = append(listMap, it)
	}
	ex.ExportToWeb(c, headerMap, listMap)
}

// 清空记录
func (w *OperateLogApi) Clean(c *gin.Context) {
	err := service.GetOperLogServiceInstance().TruncateLogTable()
	lv_err.HasErrAndPanic(err)
	util.SuccessData(c, nil)
}

func (w *OperateLogApi) DelLogs(c *gin.Context) {
	var operIds = c.Param("operIds")
	ids := strings.Split(operIds, ",")
	db := lv_db.GetOrmDefault().Table("sys_oper_log").Delete("oper_id in ? ", ids)
	lv_err.HasErrAndPanic(db.Error)
	util.SuccessMsg(c, "rows:"+cast.ToString(db.RowsAffected))
}
