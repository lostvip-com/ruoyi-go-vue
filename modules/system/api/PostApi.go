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

type PostApi struct {
}

func (w *PostApi) GetPostInfo(c *gin.Context) {
	var postId = c.Param("postId")
	post := new(model.SysPost)
	post, err := post.FindById(cast.ToInt64(postId))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, post)
}

// GetPostOptionSelect 列表分页数据
func (w *PostApi) GetPostOptionSelect(c *gin.Context) {
	var req *vo.SelectPostPageReq
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("岗位管理", req).WriteJsonExit()
		return
	}
	result, err := dao.GetSysPostDaoInstance().ListAll(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, result)
}

// ListAjax 列表分页数据
func (w *PostApi) ListAjax(c *gin.Context) {
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

// 新增页面保存
func (w *PostApi) AddSave(c *gin.Context) {
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
func (w *PostApi) EditSave(c *gin.Context) {
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
func (w *PostApi) Remove(c *gin.Context) {
	var postIds = c.Param("postIds")
	err := service.GetSysPostServiceInstance().DeleteRecordByIds(postIds)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

func (w *PostApi) Export(c *gin.Context) {
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
