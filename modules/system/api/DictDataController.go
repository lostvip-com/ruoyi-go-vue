package api

import (
	"common/common_vo"
	"common/models"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"system/service"
)

type DictDataController struct {
}

// 列表分页数据
func (w *DictDataController) ListAjax(c *gin.Context) {
	var req *common_vo.SelectDictDataPageReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("字典数据管理", req).WriteJsonExit()
		return
	}
	rows := make([]models.SysDictData, 0)
	var dictService service.DictDataService
	result, total, err := dictService.SelectListByPage(req)

	if err == nil && len(*result) > 0 {
		rows = *result
	}

	util.BuildTable(c, total, rows).WriteJsonExit()
}

// 新增页面保存
func (w *DictDataController) AddSave(c *gin.Context) {
	var req *common_vo.AddDictDataReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg(err.Error()).Log("字典数据管理", req).WriteJsonExit()
		return
	}
	var dictService service.DictDataService
	rid, err := dictService.AddSave(req, c)

	if err != nil || rid <= 0 {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).Log("字典数据管理", req).WriteJsonExit()
		return
	}
	util.SucessResp(c).SetData(rid).SetBtype(lv_dto.Buniss_Add).Log("字典数据管理", req).WriteJsonExit()
}

// 修改页面保存
func (w *DictDataController) EditSave(c *gin.Context) {
	var req *common_vo.EditDictDataReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	var dictService service.DictDataService
	err := dictService.EditSave(req, c)
	if err == nil {
		util.Success(c, nil)
	} else {
		util.Fail(c, err.Error())
	}
}

// Remove 删除数据
func (w *DictDataController) Remove(c *gin.Context) {
	var dictCodes = c.Param("dictCodes")
	var dictService service.DictDataService
	err := dictService.DeleteRecordByIds(dictCodes)
	if err == nil {
		util.Success(c, nil)
	} else {
		util.Fail(c, err.Error())
	}
}

// Export 导出
func (w *DictDataController) Export(c *gin.Context) {
	var req *common_vo.SelectDictDataPageReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("字典数据导出", req).WriteJsonExit()
		return
	}
	var dictService service.DictDataService
	url, err := dictService.Export(req)
	if err == nil {
		util.Success(c, url)
	} else {
		util.Fail(c, err.Error())
	}
}
