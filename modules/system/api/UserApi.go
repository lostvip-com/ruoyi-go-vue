package api

import (
	"common/common_vo"
	util "common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
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
	var userService = service2.GetUserServiceInstance()
	err := userService.DeleteByIds(userIds)
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
	var req *common_vo.UserPageReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	if req.DeptId == 0 {
		user := w.GetCurrUser(c)
		req.DeptId = user.DeptId
	}
	req.BeginTime = c.DefaultQuery("params[beginTime]", "")
	req.EndTime = c.DefaultQuery("params[endTime]", "")

	var userService = service2.GetUserServiceInstance()
	result, total, err := userService.FindList(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessPage(c, result, total)
}

func (w *UserApi) AddSave(c *gin.Context) {
	var req *common_vo.AddUserReq
	//获取参数
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	var userService = service2.GetUserServiceInstance()
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
	var req *common_vo.ResetPwdReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	var userService = service2.GetUserServiceInstance()
	err := userService.ResetPassword(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, nil)
}

func (w *UserApi) EditSave(c *gin.Context) {
	var req *common_vo.EditUserReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	var userService = service2.GetUserServiceInstance()
	err := userService.EditSave(req, c)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, req.UserId)
}

// 导出
func (w *UserApi) Export(c *gin.Context) {
	var req *common_vo.UserPageReq
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("导出Excel", req).WriteJsonExit()
	}
	req.BeginTime = c.DefaultQuery("params[beginTime]", "")
	req.EndTime = c.DefaultQuery("params[endTime]", "")
	listMap, err := dao.GetUserDaoInstance().SelectExportList(req)
	lv_err.HasErrAndPanic(err)
	headerMap := []map[string]string{
		map[string]string{"key": "userId", "title": "用户序号", "width": "10"},
		map[string]string{"key": "deptId", "title": "部门编号", "width": "15"},
		map[string]string{"key": "deptId", "title": "部门编号", "width": "15"},
		map[string]string{"key": "nickName", "title": "用户名称", "width": "20"},
		map[string]string{"key": "phonenumber", "title": "手机号码", "width": "20"},
		map[string]string{"key": "status", "title": "帐号状态", "width": "10"},
		map[string]string{"key": "loginIp", "title": "最后登录IP", "width": "30"},
		map[string]string{"key": "loginDate", "title": "最后登录时间", "width": "60"},
		map[string]string{"key": "deptName", "title": "部门名称", "width": "50"},
		map[string]string{"key": "leader", "title": "部门负责人", "width": "30"},
	}
	ex := util.NewMyExcel()
	ex.ExportToWeb(c, headerMap, *listMap)
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

func (w *UserApi) GetAuthUserRole(c *gin.Context) {
	userIdStr := c.Param("userId")
	userId := cast.ToInt64(userIdStr)
	// 登录者的权限
	roleSvc := service2.GetRoleServiceInstance()
	var roles []model.SysRole
	roles = roleSvc.FindRolePermissionsById(userId)
	user := w.GetCurrUser(c)
	c.JSON(http.StatusOK, gin.H{
		"msg":   "操作成功",
		"code":  http.StatusOK,
		"user":  user,
		"roles": roles,
	})
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
