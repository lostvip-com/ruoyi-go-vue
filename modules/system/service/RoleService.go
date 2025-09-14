package service

import (
	"common/common_vo"
	"common/util"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"system/dao"
	"system/model"
	"time"
)

type RoleService struct {
	BaseService
}

var roleService *RoleService

func GetRoleServiceInstance() *RoleService {
	if roleService == nil {
		roleService = &RoleService{}
	}
	return roleService
}

// 根据主键查询数据
func (svc *RoleService) FindById(id int) (*model.SysRole, error) {
	entity := &model.SysRole{RoleId: id}
	err := entity.FindOne()
	return entity, err
}

// 根据主键查询数据
func (svc *RoleService) FindPage(params *common_vo.RolePageReq) ([]model.SysRole, int, error) {
	var d dao.SysRoleDao
	return d.SelectListPage(params)
}

// 添加数据
func (svc *RoleService) AddSave(req *model.SysRole, c *gin.Context) (int, error) {
	role := new(model.SysRole)
	role.RoleName = req.RoleName
	role.RoleKey = req.RoleKey
	role.Status = req.Status
	role.Remark = req.Remark
	role.CreateTime = time.Now()
	//数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）
	role.DataScope = "2"
	role.DelFlag = "0"
	var userService UserService
	user := userService.GetCurrUser(c)
	if user != nil {
		role.CreateBy = user.UserName
	}
	session := lv_db.GetOrmDefault().Begin()
	err := session.Save(role).Error
	lv_err.HasErrAndPanic(err)
	if err != nil {
		session.Rollback()
		return 0, err
	}
	if req.MenuIds != nil {
		err = svc.saveRoleMenu(session, req, role)
	}
	err = session.Commit().Error
	return role.RoleId, err
}

// 修改数据
func (svc *RoleService) EditSave(req *model.SysRole, user *model.SysUser) error {
	r := &model.SysRole{RoleId: req.RoleId}
	err := r.FindOne()
	lv_err.HasErrAndPanic(err)
	r.RoleName = req.RoleName
	r.RoleKey = req.RoleKey
	r.Status = req.Status
	r.Remark = req.Remark
	r.UpdateTime = time.Now()
	r.UpdateBy = user.UserName
	r.RoleSort = req.RoleSort
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = lv_db.GetOrmDefault().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//更新role表
		if err := tx.Updates(&r).Error; err != nil {
			return err
		}
		//删除旧的功能权限授权，中间表
		err = svc.saveRoleMenu(tx, req, r)
		return err
	})

	return err
}

// saveRoleMenu 角色菜单关联
func (svc *RoleService) saveRoleMenu(tx *gorm.DB, req *model.SysRole, r *model.SysRole) (err error) {
	if len(req.MenuIds) == 0 {
		return nil
	}
	roleMenus := make([]model.SysRoleMenu, 0)
	for i := range req.MenuIds {
		if req.MenuIds[i] > 0 {
			var tmp model.SysRoleMenu
			tmp.RoleId = r.RoleId
			tmp.MenuId = req.MenuIds[i]
			roleMenus = append(roleMenus, tmp)
		}
	}
	err = tx.Exec(" delete from sys_role_menu where role_id=? ", r.RoleId).Error
	if err != nil {
		return err
	}
	//插入新的功能权限授权，中间表
	err = tx.CreateInBatches(roleMenus, len(roleMenus)).Error
	return err
}

// 保存数据权限
func (svc *RoleService) AuthDataScope(req *common_vo.DataScopeReq, c *gin.Context) (int, error) {
	entity := &model.SysRole{RoleId: req.RoleId}
	err := entity.FindOne()
	lv_err.HasErrAndPanic(err)
	if req.DataScope != "" {
		entity.DataScope = req.DataScope
	}
	var userService UserService
	user := userService.GetCurrUser(c)
	if user != nil {
		entity.UpdateBy = user.UserName
	}
	entity.UpdateTime = time.Now()

	db := lv_db.GetOrmDefault()
	err = db.Transaction(func(tx *gorm.DB) error {
		//更新role表
		if err := tx.Updates(&entity).Error; err != nil {
			return err
		}
		//删除旧的功能权限授权，中间表
		if req.DeptIds != nil {
			if len(req.DeptIds) > 0 {
				roleDepts := make([]model.SysRoleDept, 0)
				for i := range req.DeptIds {
					if req.DeptIds[i] > 0 {
						var tmp model.SysRoleDept
						tmp.RoleId = entity.RoleId
						tmp.DeptId = req.DeptIds[i]
						roleDepts = append(roleDepts, tmp)
					}
				}
				if len(roleDepts) > 0 {
					tx.Exec("delete from  sys_role_dept where role_id=?", entity.RoleId)
					err := tx.CreateInBatches(roleDepts, len(roleDepts)).Error
					return err
				}
			}
		}
		return err
	})

	return 1, err

}

// DeleteByIds 批量删除数据记录
func (svc *RoleService) DeleteByIds(ids string) error {
	idArr := util.ToIntArray(ids, ",")
	idsDel := make([]int, 0)
	for _, id := range idArr {
		if id != 1 { //忽略admin
			idsDel = append(idsDel, id)
		}
	}
	err := lv_db.GetOrmDefault().Exec("delete from sys_role where role_id in ?  ", idsDel).Error
	return err
}

func (svc *RoleService) DeleteRolesByUserIds(userIds []int) {
	err := svc.GetDb().Exec("delete from sys_user_role where user_id in (?)", userIds).Error
	lv_err.HasErrAndPanic(err)
}

// SelectRoleContactVo 根据用户ID查询角色
func (svc *RoleService) SelectRoleContactVo(userId int) ([]common_vo.SysRoleFlag, error) {
	var paramsPost *common_vo.RolePageReq
	var roleDao = dao.GetRoleDaoInstance()
	roleAll, err := roleDao.FindAll(paramsPost)
	if err != nil || roleAll == nil {
		return nil, errors.New("未查询到角色数据")
	}
	userRole, err := roleDao.FindRoles(userId)
	if userRole != nil {
		for i := range userRole {
			for j := range roleAll {
				if userRole[i].RoleId == roleAll[j].RoleId {
					roleAll[j].Flag = true
					break
				}
			}
		}
	}
	return roleAll, nil
}

func (svc *RoleService) InsertRoleUserIds(roleId int, userIds []string) error {
	var roleUserList []model.SysUserRole
	for _, userId := range userIds {
		var tmp model.SysUserRole
		tmp.RoleId = roleId
		tmp.UserId = cast.ToInt(userId)
		roleUserList = append(roleUserList, tmp)
	}
	err := lv_db.GetOrmDefault().CreateInBatches(roleUserList, len(roleUserList)).Error
	return err
}
func (svc *RoleService) InsertUserRoleIds(userId int, arrRoleIds []string) error {
	var roleUserList []model.SysUserRole
	for _, roleId := range arrRoleIds {
		var tmp model.SysUserRole
		tmp.UserId = userId
		tmp.RoleId = cast.ToInt(roleId)
		roleUserList = append(roleUserList, tmp)
	}
	err := lv_db.GetOrmDefault().CreateInBatches(roleUserList, len(roleUserList)).Error
	return err
}

// DeleteUserRoleInfo 取消授权用户角色
func (svc *RoleService) DeleteUserRoleInfo(userId, roleId int) error {
	entity := &model.SysUserRole{UserId: userId, RoleId: roleId}
	if entity.RoleId == 1 {
		return errors.New("Not Allowed!")
	}
	entity, err := entity.FindOne()
	if err != nil {
		return err
	}
	return entity.Delete()
}

// DeleteUserRoleInfos 批量取消授权用户角色
func (svc *RoleService) DeleteUserRoleInfos(roleId int, ids string) error {
	idarr := util.ToIntArray(ids, ",")

	idStr := ""
	for _, item := range idarr {
		if item == 1 {
			continue
		}
		if item > 0 {
			if idStr != "" {
				idStr += "," + lv_conv.String(item)
			} else {
				idStr = lv_conv.String(item)
			}
		}
	}
	err := lv_db.GetOrmDefault().Exec("delete from sys_user_role where role_id=? and user_id in (?)", roleId, idStr).Error
	return err
}

// IsRoleNameExist 检查角色名是否唯一
func (svc *RoleService) IsRoleNameExist(roleName string) (bool, error) {
	var roleDao = dao.GetRoleDaoInstance()
	_, err := roleDao.FindRoleByName(roleName)
	if err == nil {
		return true, err
	}
	return false, err
}

// 检查角色键是否唯一
func (svc *RoleService) IsRoleKeyExist(roleKey string) (bool, error) {
	var roleDao = dao.GetRoleDaoInstance()
	_, err := roleDao.FindRoleByRoleKey(roleKey)
	if err == nil {
		return true, err
	}
	return false, err
}

func (svc *RoleService) IsAdminRoleId(id int) bool {
	if id == 1 {
		return true
	} else {
		return false
	}
}

// 校验角色是否允许操作
func (svc *RoleService) CheckRoleAllowed(roleId int) bool {
	if svc.IsAdminRoleId(roleId) {
		return false
	} else {
		return true
	}
}

func (svc *RoleService) CheckRoleDataScope(roleId int) bool {
	const baseSql = `
         select distinct r.role_id, r.role_name, r.role_key, r.role_sort, r.data_scope, r.menu_check_strictly, r.dept_check_strictly,
		 r.status, r.del_flag, r.create_time, r.remark 
		 from sys_role r 
		 left join sys_user_role ur on ur.role_id = r.role_id 
		 left join sys_user u on u.user_id = ur.user_id 
		 left join sys_dept d on u.dept_id = d.dept_id 
        `
	var sql = baseSql + " where r.del_flag = '0' AND r.role_id = " + cast.ToString(roleId)
	var count int64
	err := lv_db.GetOrmDefault().Raw(sql).Count(&count).Error
	lv_err.HasErrAndPanic(err)
	return count < 1
}

func (svc *RoleService) GetDeptTreeRole(roleId int) []int {
	role := new(model.SysRole)
	role, err := role.FindById(roleId)
	deptCheckStrictly := role.DeptCheckStrictly
	sql := ` select d.dept_id from sys_dept d
             left join sys_role_dept rd on d.dept_id = rd.dept_id 
             where rd.role_id = ` + cast.ToString(roleId)
	if deptCheckStrictly {
		sql += " and d.dept_id not in (select d.parent_id from sys_dept d inner join sys_role_dept rd on d.dept_id = rd.dept_id and rd.role_id = " + cast.ToString(roleId) + " ) "
	}
	sql += " order by d.parent_id, d.order_num "
	var count []int
	err = lv_db.GetOrmDefault().Raw(sql).Find(&count).Error
	lv_err.HasErrAndPanic(err)
	return count
}

func (svc *RoleService) FindRolePermissionsById(userId int) []model.SysRole {
	var roles []model.SysRole
	var sql = "select distinct r.* " +
		"from sys_role r " +
		"left join sys_user_role ur on ur.role_id = r.role_id " +
		"left join sys_user u on u.user_id = ur.user_id " +
		"left join sys_dept d on u.dept_id = d.dept_id " +
		"where r.del_flag = '0'"
	if !svc.IsAdmin(userId) {
		sql += " and ur.user_id = " + cast.ToString(userId)
	}
	err := lv_db.GetOrmDefault().Raw(sql).Scan(&roles).Error
	lv_err.HasErrAndPanic(err)
	return roles
}

func (svc *RoleService) buildRolePermissionByUserId(roles []model.SysRole) []string {

	var roleNames []string
	for i := 0; i < len(roles); i++ {
		var role = roles[i]
		var roleName = role.RoleName
		roleNames = append(roleNames, roleName)
	}
	return roleNames
}

func (svc *RoleService) GetRolePermission(user model.SysUser) []string {
	return svc.GetRolePermissionById(user.UserId)
}

func (svc *RoleService) GetRolePermissionById(userId int) []string {
	if svc.IsAdmin(userId) {
		return []string{"admin"}
	} else {
		var roles = svc.FindRolePermissionsById(userId)
		return svc.buildRolePermissionByUserId(roles)
	}
}

func (svc *RoleService) GetMenuPermission(user *model.SysUser) []string {
	if svc.IsAdmin(user.UserId) {
		return []string{"*:*:*"}
	} else {
		var permissions []string
		err := lv_db.GetOrmDefault().Raw("select distinct m.perms "+
			"from sys_menu m "+
			"left join sys_role_menu rm on m.menu_id = rm.menu_id "+
			"left join sys_user_role ur on rm.role_id = ur.role_id "+
			"left join sys_role r on r.role_id = ur.role_id "+
			"where m.status = '0' and r.status = '0' and ur.user_id = ?", user.UserId).
			Scan(&permissions).Error
		lv_err.HasErrAndPanic(err)
		return svc.RemoveEmpty(permissions)
	}
}

func (svc *RoleService) RemoveEmpty(arr []string) (result []string) {
	for _, value := range arr {
		if len(value) > 0 { // 判断字符串长度大于零则添加到结果数组中
			result = append(result, value)
		}
	}
	return
}

func (svc *RoleService) SelectRolesByUserName(userName string) string {
	var roles []model.SysRole
	var result = ""
	err := svc.GetDb().Raw("select distinct r.* "+
		"from sys_role r "+
		"left join sys_user_role ur on ur.role_id = r.role_id "+
		"left join sys_user u on u.user_id = ur.user_id "+
		"left join sys_dept d on u.dept_id = d.dept_id "+
		"WHERE r.del_flag = '0' and u.user_name = ?", userName).
		Scan(&roles).Error
	lv_err.HasErrAndPanic(err)
	if roles != nil {
		for i := range roles {
			sysRoles := roles[i]
			if i == 0 {
				result = sysRoles.RoleName
			} else {
				result += "," + sysRoles.RoleName
			}
		}
	}
	return result
}
