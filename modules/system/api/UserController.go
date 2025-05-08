package api

import (
	userModel "common/common_vo"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	service2 "system/service"
)

// ListAjax 用户列表分页数据
func (w *UserApi) ListAjax(c *gin.Context) {
	var req *userModel.SelectUserPageReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("用户管理", req).WriteJsonExit()
		return
	}
	var userService service2.UserService
	result, total, err := userService.SelectRecordList(req)
	lv_err.HasErrAndPanic(err)
	util.BuildTable(c, total, result).WriteJsonExit()
}

// 保存新增用户数据
func (w *UserApi) AddSave(c *gin.Context) {
	var req *userModel.AddUserReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg(err.Error()).Log("新增用户", req).WriteJsonExit()
		return
	}
	var userService service2.UserService
	//判断登录名是否已注册
	count, err := userService.CountCol("username", req.UserName)
	if count > 0 {
		util.Fail(c, "登录名已经存在")
		return
	}
	//判断手机号码是否已注册
	count, _ = userService.CountCol("phonenumber", req.Phonenumber)
	if count > 0 {
		util.Fail(c, "手机号码已经存在")
		return
	}
	uid, err := userService.AddSave(req, c)
	lv_err.HasErrAndPanic(err)
	util.Success(c, uid)
}
