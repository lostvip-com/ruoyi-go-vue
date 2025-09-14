package api

import (
	api2 "common/api"
	"common/common_vo"
	"common/util"
	"github.com/gin-gonic/gin"
	db2 "github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/spf13/cast"
	"net/http"
	"strings"
	"system/dao"
	"system/model"
	"system/service"
	"system/vo"
)

type RoleApi struct {
	api2.BaseApi
}

func (w *RoleApi) GetRoleOptionSelect(c *gin.Context) {
	arr, err := dao.GetRoleDaoInstance().FindAll(nil)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessData(c, arr)
}
func (w *RoleApi) GetRoleInfo(c *gin.Context) {
	roleId := c.Param("roleId")
	role := new(model.SysRole)
	role, err := role.FindById(cast.ToInt(roleId))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessData(c, role)
}

// 列表分页数据
func (w *RoleApi) ListAjax(c *gin.Context) {
	var req *common_vo.RolePageReq

	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).Log("角色管理", req).WriteJsonExit()
		return
	}
	rows := make([]model.SysRole, 0)
	roleService := service.RoleService{}
	result, total, err := roleService.FindPage(req)

	if err == nil && len(result) > 0 {
		rows = result
	}
	util.BuildTable(c, total, rows).WriteJsonExit()
}

// 新增页面保存
func (w *RoleApi) AddSave(c *gin.Context) {
	var req *model.SysRole
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	roleDao := dao.SysRoleDao{}
	count, err := roleDao.FindCount("", req.RoleName)
	if count >= 1 {
		util.Fail(c, err.Error())
		return
	}
	count, err = roleDao.FindCount(req.RoleKey, "")
	if count >= 1 {
		util.Fail(c, err.Error())
		return
	}
	roleService := service.RoleService{}
	rid, err := roleService.AddSave(req, c)

	if err != nil || rid <= 0 {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessData(c, nil)
}

func (w *RoleApi) EditSave(c *gin.Context) {
	var req *model.SysRole
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	roleService := service.GetRoleServiceInstance()
	user := w.GetCurrUser(c)
	err := roleService.EditSave(req, user)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessData(c, nil)
}

func (w *RoleApi) ChangeStatus(c *gin.Context) {
	var req *vo.RoleStatusReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	sql := " update sys_role set status=? where role_id = ? "
	rows := db2.GetOrmDefault().Exec(sql, req.Status, req.RoleId).RowsAffected
	util.SuccessData(c, rows)
}

func (w *RoleApi) GetUnAllocatedList(c *gin.Context) {
	roleId := cast.ToInt(c.Query("roleId"))
	UserName := c.Query("userName")
	phonenumber := c.Query("phonenumber")
	var rows []map[string]any
	var userService service.UserService
	userList, err := userService.SelectUnallocatedList(roleId, UserName, phonenumber)

	if err == nil && userList != nil {
		rows = *userList
	}

	c.JSON(http.StatusOK, lv_dto.TableDataInfo{
		Code:  200,
		Msg:   "操作成功",
		Total: len(rows),
		Rows:  rows,
	})
}

func (w *RoleApi) Remove(c *gin.Context) {
	var roleIds = c.Param("roleIds")
	roleService := service.GetRoleServiceInstance()
	err := roleService.DeleteByIds(roleIds)

	if err != nil {
		util.Fail(c, err.Error())
	} else {
		util.SuccessData(c, nil)
	}
}

func (w *RoleApi) PutDataScope(c *gin.Context) {
	var req *common_vo.DataScopeReq
	if err := c.ShouldBind(&req); err != nil {
		util.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	roleService := service.RoleService{}
	if !roleService.CheckRoleAllowed(req.RoleId) {
		util.ErrorResp(c).SetMsg("不允许操作超级管理员角色").WriteJsonExit()
		return
	}
	_, err := roleService.AuthDataScope(req, c)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessData(c, nil)
}

// 查询已分配用户角色列表
func (w *RoleApi) AllocatedList(c *gin.Context) {
	roleId := cast.ToInt(c.Query("roleId"))
	UserName := c.Query("UserName")
	phonenumber := c.Query("phonenumber")
	var rows []map[string]any

	var userService service.UserService
	userList, err := userService.SelectAllocatedList(roleId, UserName, phonenumber)

	if err == nil && userList != nil {
		rows = *userList
	}

	c.JSON(http.StatusOK, lv_dto.TableDataInfo{
		Code:  200,
		Msg:   "操作成功",
		Total: len(rows),
		Rows:  rows,
	})
}

func (w *RoleApi) AuthRoleToUsers(c *gin.Context) {
	roleId := cast.ToInt(c.Query("roleId"))
	userIds := c.Query("userIds")
	if roleId <= 0 {
		util.Fail(c, "roleId can not be empty")
		return
	}
	if userIds == "" {
		util.Fail(c, "userIds can not be empty")
		return
	}
	roleService := service.RoleService{}
	err := roleService.InsertRoleUserIds(roleId, strings.Split(userIds, ","))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessData(c, nil)
}
func (w *RoleApi) CancelAll(c *gin.Context) {
	roleId := cast.ToInt(c.Query("roleId"))
	userIds := c.Query("userIds")
	roleService := service.GetRoleServiceInstance()
	err := roleService.DeleteUserRoleInfos(roleId, userIds)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessData(c, nil)
}

func (w *RoleApi) Cancel(c *gin.Context) {
	var req vo.UserRoleReq
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	roleService := service.GetRoleServiceInstance()
	err := roleService.DeleteUserRoleInfo(cast.ToInt(req.UserId), cast.ToInt(req.RoleId))
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.SuccessData(c, nil)
}

func (w *RoleApi) GetDeptTreeRole(c *gin.Context) {
	roleId := c.Param("roleId")
	roleSvc := service.GetRoleServiceInstance()
	checkedKeys := roleSvc.GetDeptTreeRole(cast.ToInt(roleId))
	depts := service.GetDeptServiceInstance().SelectDeptTreeList()
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"msg":         "操作成功",
		"code":        http.StatusOK,
		"checkedKeys": checkedKeys,
		"depts":       depts,
	})
}
