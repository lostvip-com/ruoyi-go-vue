package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/spf13/cast"
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

// 清空记录
func (w *OperateLogApi) Clean(c *gin.Context) {
	var operlogService service.OperLogService
	err := operlogService.DeleteRecordAll()
	lv_err.HasErrAndPanic(err)
	util.Success(c, nil)
}

// 删除数据
func (w *OperateLogApi) Remove(c *gin.Context) {
	var req *lv_dto.IdsReq

	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	err := service.GetOperLogServiceInstance().DeleteRecordByIds(req.Ids)
	lv_err.HasErrAndPanic(err)
	util.Success(c, nil)
}

func (w *OperateLogApi) DelectOperlog(context *gin.Context) {
	var operId = context.Param("operId")
	service.GetOperLogServiceInstance().DeleteRecordById(cast.ToInt64(operId))
}
