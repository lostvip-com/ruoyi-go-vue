package api

import (
	"common/common_vo"
	"common/models"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"system/dao"
	"system/service"
)

type DictDataApi struct {
}

func (w *DictDataApi) GetDictDataByDictType(c *gin.Context) {
	dictType := cast.ToString(c.Param("dictType"))
	list, err := dao.GetDictDataDaoInstance().FindAll("", dictType)
	if err != nil {
		util.Fail(c, err.Error())
		return
	} else {
		util.SuccessData(c, list)
	}
}

// 列表分页数据
func (w *DictDataApi) ListAjax(c *gin.Context) {
	var req *common_vo.SelectDictDataPageReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("字典数据管理", req).WriteJsonExit()
		return
	}
	rows := make([]models.SysDictData, 0)
	var dictService service.DictDataService
	result, total, err := dictService.FindPage(req)

	if err == nil && len(*result) > 0 {
		rows = *result
	}

	util.BuildTable(c, total, rows).WriteJsonExit()
}

// 新增页面保存
func (w *DictDataApi) AddSave(c *gin.Context) {
	var req *common_vo.AddDictDataReq

	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	var dictService service.DictDataService
	rid, err := dictService.AddSave(req, c)

	if err != nil || rid <= 0 {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessData(c, nil)
}

// 修改页面保存
func (w *DictDataApi) EditSave(c *gin.Context) {
	var req *common_vo.EditDictDataReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	var dictService service.DictDataService
	err := dictService.EditSave(req, c)
	if err == nil {
		util.SuccessData(c, nil)
	} else {
		util.Fail(c, err.Error())
	}
}

// Remove 删除数据
func (w *DictDataApi) Remove(c *gin.Context) {
	var dictCodes = c.Param("dictCodes")
	var dictService service.DictDataService
	err := dictService.DeleteByIds(dictCodes)
	if err == nil {
		util.SuccessData(c, nil)
	} else {
		util.Fail(c, err.Error())
	}
}

// Export 导出
func (w *DictDataApi) Export(c *gin.Context) {
	var req *common_vo.SelectDictDataPageReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("字典数据导出", req).WriteJsonExit()
		return
	}
	var dictService service.DictDataService
	url, err := dictService.Export(req)
	if err == nil {
		util.SuccessData(c, url)
	} else {
		util.Fail(c, err.Error())
	}
}
