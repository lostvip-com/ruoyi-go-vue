package api

import (
	userModel "common/common_vo"
	util "common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/spf13/cast"
	"net/http"
	"strings"
	"system/dao"
	"system/model"
	service2 "system/service"
)

type UserApi struct {
	BaseApi
}

// 删除数据
func (w *UserApi) Remove(c *gin.Context) {
	userIds := c.Param("userIds")
	var userService = service2.GetUserService()
	err := userService.DeleteRecordByIds(userIds)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

func (w *UserApi) GetUserInfo(c *gin.Context) {
	userId := cast.ToInt64(c.Param("userId"))
	user := new(model.SysUser)
	user, err := user.FindById(userId)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	var roleIds []int64
	var postIds []int64
	// 登录者的权限
	roles, err := dao.GetRoleDaoInstance().FindRoles(userId)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	for _, it := range roles {
		roleIds = append(roleIds, it.RoleId)
	}
	posts, err := dao.GetSysPostDaoInstance().FindPostsByUserId(userId)
	if posts != nil {
		for _, it := range *posts {
			postIds = append(postIds, it.PostId)
		}
	}
	var result = gin.H{
		"msg":     "success",
		"code":    http.StatusOK,
		"data":    user,
		"roles":   roles,
		"posts":   posts,
		"postIds": postIds,
		"roleIds": roleIds,
	}
	c.JSON(http.StatusOK, result)
}

// ListAjax 用户列表分页数据
func (w *UserApi) ListAjax(c *gin.Context) {
	var req *userModel.SelectUserPageReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	var userService = service2.GetUserService()
	result, total, err := userService.SelectRecordList(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.BuildTable(c, total, result).WriteJsonExit()
}

func (w *UserApi) AddSave(c *gin.Context) {
	var req *userModel.AddUserReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	var userService = service2.GetUserService()
	//判断登录名是否已注册
	count, err := userService.CountCol("username", req.UserName)
	if count > 0 {
		util.Fail(c, "用户名已经存在")
		return
	}
	count, _ = userService.CountCol("phonenumber", req.Phonenumber)
	if count > 0 { //判断手机号码是否已注册
		util.Fail(c, "手机号码已经存在")
		return
	}
	uid, err := userService.AddSave(req, c)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, uid)
}

func (w *UserApi) ChangeStatus(c *gin.Context) {
	userId := c.Query("userId")
	status := c.Query("status")
	sql := " update sys_user set status=? where user_id = ? "
	rows := lv_db.GetMasterGorm().Exec(sql, status, userId).RowsAffected
	util.Success(c, rows)
}

func (w *UserApi) ResetPwdSave(c *gin.Context) {
	var req *userModel.ResetPwdReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	var userService = service2.GetUserService()
	err := userService.ResetPassword(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

func (w *UserApi) EditSave(c *gin.Context) {
	var req *userModel.EditUserReq
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	var userService = service2.GetUserService()
	err := userService.EditSave(req, c)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, req.UserId)
}

// 导出
func (w *UserApi) Export(c *gin.Context) {
	var req *userModel.SelectUserPageReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("导出Excel", req).WriteJsonExit()
	}
	var userService = service2.GetUserService()
	url, err := userService.Export(req)

	if err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("导出Excel", req).WriteJsonExit()
		return
	}
	util.SucessResp(c).SetMsg(url).Log("导出Excel", req).WriteJsonExit()
}

func (w *UserApi) PutAuthUserRoleIds(c *gin.Context) {
	uIdStr, _ := c.GetQuery("userId")
	roleIds, _ := c.GetQuery("roleIds")
	userId := cast.ToInt64(uIdStr)
	uIds := []int64{userId}
	arr := strings.Split(roleIds, ",")
	if len(arr) == 0 {
		util.Fail(c, "请选择角色")
		return
	}
	service2.GetRoleServiceInstance().DeleteRolesByUserIds(uIds)
	err := service2.GetRoleServiceInstance().InsertUserRoleIds(userId, arr)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

func (w *UserApi) GetUserDeptTree(c *gin.Context) {
	var deptSvc = service2.GetDeptServiceInstance()
	list := deptSvc.SelectDeptTreeList()
	util.Success(c, list)
}

func (w *UserApi) ExportExport(context *gin.Context) {

}

func (w *UserApi) ImportUserData(context *gin.Context) {

}

func (w *UserApi) ImportTemplate(context *gin.Context) {

}
