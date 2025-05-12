package api

import (
	"common/common_vo"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/spf13/cast"
	"system/dao"
	"system/model"
	"system/service"
)

type DictTypeApi struct {
}

func (w *DictTypeApi) GetTypeDict(c *gin.Context) {
	dictId := c.Param("dictId")
	dictType := new(model.SysDictType)
	dictType, err := dictType.FindById(cast.ToInt64(dictId))
	lv_err.HasErrAndPanic(err)
	util.Success(c, dictType)
}

// ListAjax 列表分页数据
func (w *DictTypeApi) ListAjax(c *gin.Context) {
	var req *common_vo.DictTypePageReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	beginTime := c.DefaultQuery("params[beginTime]", "")
	endTime := c.DefaultQuery("params[endTime]", "")
	req.BeginTime = beginTime
	req.EndTime = endTime

	rows := make([]model.SysDictType, 0)
	var dictTypeService service.DictTypeService
	result, total, err := dictTypeService.FindPage(req)

	if err == nil && len(result) > 0 {
		rows = result
	}

	util.BuildTable(c, total, rows).WriteJsonExit()
}

// 新增页面保存
func (w *DictTypeApi) AddSave(c *gin.Context) {
	var req common_vo.AddDictTypeReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg(err.Error()).Log("字典管理", req).WriteJsonExit()
		return
	}
	var dictTypeService service.DictTypeService
	exist, err := dictTypeService.CheckDictTypeUniqueAll(req.DictType)
	lv_err.HasErrAndPanic(err)
	if exist {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg("字典类型已存在").Log("字典管理", req).WriteJsonExit()
		return
	}

	rid, err := dictTypeService.AddSave(&req, c)

	if err != nil || rid <= 0 {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).Log("字典管理", req).WriteJsonExit()
		return
	}
	util.SucessResp(c).SetData(rid).Log("字典管理", req).WriteJsonExit()
}

// EditSave 修改页面保存
func (w *DictTypeApi) EditSave(c *gin.Context) {
	var req *common_vo.EditDictTypeReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).Log("字典类型管理", req).WriteJsonExit()
		return
	}
	var dictTypeService service.DictTypeService
	if req.DictId == 0 && dictTypeService.IsDictTypeExist(req.DictType) {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg("字典类型已存在").Log("字典类型管理", req).WriteJsonExit()
		return
	}

	rs, err := dictTypeService.EditSave(req, c)

	if err != nil || rs <= 0 {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).Log("字典类型管理", req).WriteJsonExit()
		return
	}
	util.SucessResp(c).Log("字典类型管理", req).WriteJsonExit()
}

// Remove 删除数据
func (w *DictTypeApi) Remove(c *gin.Context) {
	dictIds := c.Param("dictIds")
	var dictTypeService service.DictTypeService
	err := dictTypeService.DeleteByIds(dictIds)
	if err == nil {
		util.Success(c, nil)
	} else {
		util.Fail(c, err.Error())
	}
}

// GetOptionSelect 加载部门列表树结构的数据
//func (w *DictTypeApi) GetOptionSelect(c *gin.Context) {
//	var dictTypeService service.DictTypeService
//	result := dictTypeService.FindDictTree(nil)
//	c.JSON(http.StatusOK, result)
//}

// 导出
func (w *DictTypeApi) Export(c *gin.Context) {
	var req *common_vo.DictTypePageReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("字典管理", req).WriteJsonExit()
		return
	}
	var dictTypeService service.DictTypeService
	url, err := dictTypeService.Export(req)
	if err == nil {
		util.Success(c, url)
	} else {
		util.Fail(c, err.Error())
	}
}

// GetOptionSelect 加载部门列表树结构的数据
func (w *DictTypeApi) GetOptionSelect(c *gin.Context) {
	result, err := dao.GetSysDictTypeDaoInstance().FindAll(nil)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, result)
}
