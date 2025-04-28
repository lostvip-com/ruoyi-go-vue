package api

import (
	"common/common_vo"
	util2 "common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"net/http"
	"system/model"
	"system/service"
)

type DictTypeController struct {
}

// ListAjax 列表分页数据
func (w *DictTypeController) ListAjax(c *gin.Context) {
	var req *common_vo.DictTypePageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetMsg(err.Error()).Log("字典类型管理", req).WriteJsonExit()
		return
	}
	rows := make([]model.SysDictType, 0)
	var dictTypeService service.DictTypeService
	result, total, err := dictTypeService.SelectListByPage(req)

	if err == nil && len(result) > 0 {
		rows = result
	}

	util2.BuildTable(c, total, rows).WriteJsonExit()
}

// 新增页面保存
func (w *DictTypeController) AddSave(c *gin.Context) {
	var req common_vo.AddDictTypeReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg(err.Error()).Log("字典管理", req).WriteJsonExit()
		return
	}
	var dictTypeService service.DictTypeService
	exist, err := dictTypeService.CheckDictTypeUniqueAll(req.DictType)
	lv_err.HasErrAndPanic(err)
	if exist {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg("字典类型已存在").Log("字典管理", req).WriteJsonExit()
		return
	}

	rid, err := dictTypeService.AddSave(&req, c)

	if err != nil || rid <= 0 {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).Log("字典管理", req).WriteJsonExit()
		return
	}
	util2.SucessResp(c).SetData(rid).Log("字典管理", req).WriteJsonExit()
}

// EditSave 修改页面保存
func (w *DictTypeController) EditSave(c *gin.Context) {
	var req *common_vo.EditDictTypeReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).Log("字典类型管理", req).WriteJsonExit()
		return
	}
	var dictTypeService service.DictTypeService
	if req.DictId == 0 && dictTypeService.IsDictTypeExist(req.DictType) {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg("字典类型已存在").Log("字典类型管理", req).WriteJsonExit()
		return
	}

	rs, err := dictTypeService.EditSave(req, c)

	if err != nil || rs <= 0 {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).Log("字典类型管理", req).WriteJsonExit()
		return
	}
	util2.SucessResp(c).Log("字典类型管理", req).WriteJsonExit()
}

// Remove 删除数据
func (w *DictTypeController) Remove(c *gin.Context) {
	var req *lv_dto.IdsReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Del).SetMsg(err.Error()).Log("字典管理", req).WriteJsonExit()
		return
	}
	var dictTypeService service.DictTypeService
	rs := dictTypeService.DeleteRecordByIds(req.Ids)

	if rs == nil {
		util2.SucessResp(c).SetBtype(lv_dto.Buniss_Del).Log("字典管理", req).WriteJsonExit()
	} else {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Del).Log("字典管理", req).WriteJsonExit()
	}
}

// 导出
func (w *DictTypeController) Export(c *gin.Context) {
	var req *common_vo.DictTypePageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetMsg(err.Error()).Log("字典管理", req).WriteJsonExit()
		return
	}
	var dictTypeService service.DictTypeService
	url, err := dictTypeService.Export(req)
	if err != nil {
		util2.ErrorResp(c).SetMsg(err.Error()).Log("字典管理", req).WriteJsonExit()
		return
	}

	util2.SucessResp(c).SetMsg(url).Log("导出Excel", req).WriteJsonExit()
}

// 检查字典类型是否唯一不包括本参数
func (w *DictTypeController) CheckDictTypeUnique(c *gin.Context) {
	var req *common_vo.CheckDictTypeALLReq
	if err := c.ShouldBind(&req); err != nil {
		c.Writer.WriteString("1")
		return
	}
	var dictTypeService service.DictTypeService
	yes := dictTypeService.IsDictTypeExist(req.DictType)
	var msg string
	if yes {
		msg = "已经存在！"
	} else {
		msg = "不存在！"
	}
	util2.Success(c, yes, msg)
}

// 检查字典类型是否唯一
func (w *DictTypeController) CheckDictTypeUniqueAll(c *gin.Context) {
	var req *common_vo.CheckDictTypeALLReq
	if err := c.ShouldBind(&req); err != nil {
		c.Writer.WriteString("1")
		return
	}
	var dictTypeService service.DictTypeService
	exist, err := dictTypeService.CheckDictTypeUniqueAll(req.DictType)
	lv_err.HasErrAndPanic(err)
	util2.Success(c, exist, "exist or not ")
}

// 加载部门列表树结构的数据
func (w *DictTypeController) TreeData(c *gin.Context) {
	var dictTypeService service.DictTypeService
	result := dictTypeService.SelectDictTree(nil)
	c.JSON(http.StatusOK, result)
}
