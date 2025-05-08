package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/spf13/cast"
	"system/dao"
	"system/model"
	"system/service"
	"system/vo"
)

type PostController struct {
}

func (w *PostController) GetPostInfo(c *gin.Context) {
	var postId = c.Param("postId")
	post := new(model.SysPost)
	post, err := post.FindById(cast.ToInt64(postId))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, post)
}

// ListAjax 列表分页数据
func (w *PostController) GetPostOptionSelect(c *gin.Context) {
	var req *vo.SelectPostPageReq
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("岗位管理", req).WriteJsonExit()
		return
	}
	result, total, err := dao.GetSysPostDaoInstance().SelectPageList(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessPage(c, result, total)
}

// ListAjax 列表分页数据
func (w *PostController) ListAjax(c *gin.Context) {
	var req *vo.SelectPostPageReq
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("岗位管理", req).WriteJsonExit()
		return
	}
	result, total, err := dao.GetSysPostDaoInstance().SelectPageList(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessPage(c, result, total)
}

// AddSave 新增页面保存
func (w *PostController) AddSave(c *gin.Context) {
	var req *vo.AddPostReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg(err.Error()).Log("岗位管理", req).WriteJsonExit()
		return
	}
	var postService = service.GetSysPostServiceInstance()
	if postService.IsPostCodeExist(req.PostCode) {
		util.Fail(c, "岗位编码已存在")
		return
	}
	pid, err := postService.AddSave(req, c)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, pid)
}

// EditSave 修改页面保存
func (w *PostController) EditSave(c *gin.Context) {
	var req *vo.EditSysPostReq
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).Log("岗位管理", req).WriteJsonExit()
		return
	}
	var postService = service.GetSysPostServiceInstance()
	err := postService.EditSave(req, c)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

// Remove 删除数据
func (w *PostController) Remove(c *gin.Context) {
	var req *lv_dto.IdsReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	var postService = service.GetSysPostServiceInstance()
	err := postService.DeleteRecordByIds(req.Ids)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

// Export 导出
func (w *PostController) Export(c *gin.Context) {
	var req *vo.SelectPostPageReq
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("岗位管理", req).WriteJsonExit()
		return
	}
	var postService = service.GetSysPostServiceInstance()
	url, err := postService.Export(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, url)
}
