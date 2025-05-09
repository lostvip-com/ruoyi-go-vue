package api

import (
	userModel "common/common_vo"
	util2 "common/util"
	"github.com/gin-gonic/gin"
	db2 "github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	service2 "system/service"
)

type UserApi struct {
}

// ListAjax 用户列表分页数据
func (w *UserApi) ListAjax(c *gin.Context) {
	var req *userModel.SelectUserPageReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetMsg(err.Error()).Log("用户管理", req).WriteJsonExit()
		return
	}
	var userService service2.UserService
	result, total, err := userService.SelectRecordList(req)
	lv_err.HasErrAndPanic(err)
	util2.BuildTable(c, total, result).WriteJsonExit()
}

// 保存新增用户数据
func (w *UserApi) AddSave(c *gin.Context) {
	var req *userModel.AddUserReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Add).SetMsg(err.Error()).Log("新增用户", req).WriteJsonExit()
		return
	}
	var userService service2.UserService
	//判断登录名是否已注册
	count, err := userService.CountCol("username", req.UserName)
	if count > 0 {
		util2.Fail(c, "登录名已经存在")
		return
	}
	//判断手机号码是否已注册
	count, _ = userService.CountCol("phonenumber", req.Phonenumber)
	if count > 0 {
		util2.Fail(c, "手机号码已经存在")
		return
	}
	uid, err := userService.AddSave(req, c)
	lv_err.HasErrAndPanic(err)
	util2.Success(c, uid)
}

// 修改页面保存
func (w *UserApi) ChangeStatus(c *gin.Context) {
	userId := c.Query("userId")
	status := c.Query("status")
	sql := " update sys_user set status=? where user_id = ? "
	rows := db2.GetMasterGorm().Exec(sql, status, userId).RowsAffected
	util2.Success(c, rows)
}

// 重置密码保存
func (w *UserApi) ResetPwdSave(c *gin.Context) {
	var req *userModel.ResetPwdReq
	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).Log("重置密码", req).WriteJsonExit()
	}
	var userService service2.UserService
	result, err := userService.ResetPassword(req)

	if err != nil || !result {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).Log("重置密码", req).WriteJsonExit()
	} else {
		util2.SucessResp(c).SetBtype(lv_dto.Buniss_Edit).Log("重置密码", req).WriteJsonExit()
	}
}

// 保存修改用户数据
func (w *UserApi) EditSave(c *gin.Context) {
	var req *userModel.EditUserReq
	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	var userService service2.UserService
	err := userService.EditSave(req, c)
	if err != nil {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).WriteJsonExit()
		return
	}
	util2.SucessResp(c).SetData(req.UserId).SetBtype(lv_dto.Buniss_Edit).WriteJsonExit()
}

// 删除数据
func (w *UserApi) Remove(c *gin.Context) {
	var req *lv_dto.IdsReq

	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Del).SetMsg(err.Error()).WriteJsonExit()
	}
	var userService service2.UserService
	err := userService.DeleteRecordByIds(req.Ids)
	if err == nil {
		util2.SucessResp(c).SetBtype(lv_dto.Buniss_Del).WriteJsonExit()
	} else {
		util2.ErrorResp(c).SetBtype(lv_dto.Buniss_Del).WriteJsonExit()
	}
}

// 导出
func (w *UserApi) Export(c *gin.Context) {
	var req *userModel.SelectUserPageReq

	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetMsg(err.Error()).Log("导出Excel", req).WriteJsonExit()
	}
	var userService service2.UserService
	url, err := userService.Export(req)

	if err != nil {
		util2.ErrorResp(c).SetMsg(err.Error()).Log("导出Excel", req).WriteJsonExit()
		return
	}
	util2.SucessResp(c).SetMsg(url).Log("导出Excel", req).WriteJsonExit()
}
