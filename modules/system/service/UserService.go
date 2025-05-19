package service

import (
	"common/common_vo"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_cache"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_gen"
	"github.com/lostvip-com/lv_framework/utils/lv_net"
	"github.com/lostvip-com/lv_framework/utils/lv_secret"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
	"system/dao"
	"system/model"
	"time"
)

type UserService struct{}

var userService *UserService

func GetUserServiceInstance() *UserService {
	if userService == nil {
		userService = &UserService{}
	}
	return userService
}

func (svc *UserService) FindById(id int64) (*model.SysUser, error) {
	entity := &model.SysUser{UserId: id}
	err := entity.FindOne()
	return entity, err
}

func (svc *UserService) FindList(param *common_vo.UserPageReq) (*[]map[string]any, int64, error) {
	var deptService DeptService
	var dept, _ = deptService.FindById(param.DeptId)
	if dept != nil { //数据权限
		param.Ancestors = dept.Ancestors
	}
	var d = dao.GetUserDaoInstance()
	return d.FindPage(param)
}

// 新增用户
func (svc *UserService) AddSave(req *common_vo.AddUserReq, c *gin.Context) (int64, error) {
	var u model.SysUser
	u.UserName = req.UserName
	u.UserName = req.UserName
	u.Email = req.Email
	u.Phonenumber = req.Phonenumber
	u.Status = req.Status
	u.Sex = req.Sex
	u.DeptId = req.DeptId
	u.Remark = req.Remark
	t := time.Now()
	u.LoginDate = &t
	//生成密码
	newSalt := lv_gen.GenerateSubId(6)
	u.Password, _ = lv_secret.PasswordHash(newSalt)
	u.CreateTime = time.Now()
	createUser := svc.GetCurrUser(c)

	if createUser != nil {
		u.CreateBy = createUser.UserName
	}
	u.DelFlag = "0"

	err := lv_db.GetMasterGorm().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&u).Error; err != nil {
			return err
		}
		//增加岗位数据
		if req.PostIds != "" {
			postIds := lv_conv.ToInt64Array(req.PostIds, ",")
			userPosts := make([]model.SysUserPost, 0)
			for i := range postIds {
				if postIds[i] > 0 {
					var userPost model.SysUserPost
					userPost.UserId = u.UserId
					userPost.PostId = postIds[i]
					userPosts = append(userPosts, userPost)
				}
			} //end for
			if len(userPosts) > 0 {
				if err := tx.CreateInBatches(userPosts, 1).Error; err != nil {
					return err
				}
			}
		}
		//增加角色数据
		if req.RoleIds != "" {
			roleIds := lv_conv.ToInt64Array(req.RoleIds, ",")
			userRoles := make([]model.SysUserRole, 0)
			for i := range roleIds {
				if roleIds[i] > 0 {
					var userRole model.SysUserRole
					userRole.UserId = u.UserId
					userRole.RoleId = roleIds[i]
					userRoles = append(userRoles, userRole)
				}
			}
			if len(userRoles) > 0 {
				if err := tx.CreateInBatches(userRoles, 1).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	return u.UserId, err
}

// 新增用户
func (svc *UserService) EditSave(req *common_vo.EditUserReq, c *gin.Context) error {
	userPtr := &model.SysUser{UserId: req.UserId}
	err := userPtr.FindOne()
	if err != nil {
		return err
	}
	userPtr.UserName = req.UserName
	userPtr.Email = req.Email
	userPtr.Phonenumber = req.Phonenumber
	userPtr.Status = req.Status
	userPtr.Sex = req.Sex
	userPtr.DeptId = req.DeptId
	userPtr.Remark = req.Remark
	userPtr.UpdateTime = time.Now()
	updateUser := svc.GetCurrUser(c)

	if updateUser != nil {
		userPtr.UpdateBy = updateUser.UserName
	}
	err = lv_db.GetMasterGorm().Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(userPtr).Error; err != nil {
			return err
		}
		//增加岗位数据
		if req.PostIds != "" {
			postIds := lv_conv.ToInt64Array(req.PostIds, ",")
			userPosts := make([]model.SysUserPost, 0)
			for i := range postIds {
				if postIds[i] > 0 {
					var userPost model.SysUserPost
					userPost.UserId = userPtr.UserId
					userPost.PostId = postIds[i]
					userPosts = append(userPosts, userPost)
				}
			} //end for
			if len(userPosts) > 0 {
				tx.Exec("delete from sys_user_post where user_id=?", userPtr.UserId)
				if err := tx.Save(userPosts).Error; err != nil {
					return err
				}
			}
		}
		//增加角色数据
		if req.RoleIds != "" {
			roleIds := lv_conv.ToInt64Array(req.RoleIds, ",")
			userRoles := make([]model.SysUserRole, 0)
			for i := range roleIds {
				if roleIds[i] > 0 {
					var userRole model.SysUserRole
					userRole.UserId = userPtr.UserId
					userRole.RoleId = roleIds[i]
					userRoles = append(userRoles, userRole)
				}
			} //end for
			if len(userRoles) > 0 {
				tx.Exec("delete from sys_user_role where user_id=?", userPtr.UserId)
				if err := tx.Save(userRoles).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	return err
}

// 根据主键删除用户信息
func (svc *UserService) DeleteById(id int64) error {
	entity := &model.SysUser{UserId: id}
	err := entity.Delete()
	return err
}

// 批量删除用户记录
func (svc *UserService) DeleteByIds(ids string) error {
	idarr := lv_conv.ToInt64Array(ids, ",")
	idarr = lv_conv.RemoveOne(idarr, 1) //去掉admin的id
	if len(idarr) == 0 {
		return errors.New("ids can not be empty ")
	}
	err := lv_db.GetMasterGorm().Transaction(func(tx *gorm.DB) error {
		err := tx.Table("sys_user").Where("user_id in ? and user_id!=1 ", idarr).Update("del_flag", 1).Error
		if err != nil {
			return err
		}
		err = tx.Table("sys_user_post").Where("user_id in ? and user_id!=1 ", idarr).Update("del_flag", 1).Error
		if err != nil {
			return err
		}
		err = tx.Table("sys_user_role").Where("user_id in ? and user_id!=1 ", idarr).Update("del_flag", 1).Error
		if err != nil {
			return err
		}
		return err
	})
	return err
}

// 判断是否是系统管理员
func (svc *UserService) IsAdmin(userId int64) bool {
	if userId == 1 {
		return true
	} else {
		return false
	}
}

// 检查账号是否符合规范,存在返回false,否则true
func (svc *UserService) CheckPassport(UserName string) bool {
	entity := model.SysUser{UserName: UserName}
	if err := entity.FindOne(); err != nil {
		return false
	} else {
		return true
	}
}

// 获得用户信息详情
func (svc *UserService) GetCurrUser(c *gin.Context) *model.SysUser {
	token := lv_net.GetParam(c, "token")
	key := "login:" + token
	userId, _ := lv_cache.GetCacheClient().HGet(key, "userId")
	u := new(model.SysUser)
	u.UserId = cast.ToInt64(userId)
	err := u.FindOne()
	if err != nil {
		panic(err)
	}
	// 部门
	u.Dept, _ = GetDeptServiceInstance().FindById(u.DeptId)
	u.Roles, _ = dao.GetRoleDaoInstance().FindRoles(u.UserId)
	return u
}

// 更新用户信息详情
func (svc *UserService) UpdateProfile(profile *common_vo.ProfileReq, c *gin.Context) error {
	user := svc.GetCurrUser(c)

	if profile.UserName != "" {
		user.UserName = profile.UserName
	}

	if profile.Email != "" {
		user.Email = profile.Email
	}

	if profile.Phonenumber != "" {
		user.Phonenumber = profile.Phonenumber
	}

	if profile.Sex != "" {
		user.Sex = profile.Sex
	}

	err := user.Updates()
	if err != nil {
		return errors.New("保存数据失败")
	}

	//SaveUserToSession(user, c)
	return nil
}

// 更新用户头像
func (svc *UserService) UpdateAvatar(avatar string, c *gin.Context) error {
	user := svc.GetCurrUser(c)

	if avatar != "" {
		user.Avatar = avatar
	}

	err := user.Updates()
	if err != nil {
		return errors.New("保存数据失败")
	}

	//SaveUserToSession(user, c)
	return nil
}

// 修改用户密码
func (svc *UserService) UpdatePassword(profile *common_vo.PasswordReq, c *gin.Context) error {
	user := svc.GetCurrUser(c)

	if profile.OldPassword == "" {
		return errors.New("旧密码不能为空")
	}

	if profile.NewPassword == "" {
		return errors.New("新密码不能为空")
	}

	if profile.Confirm == "" {
		return errors.New("确认密码不能为空")
	}

	if profile.NewPassword == profile.OldPassword {
		return errors.New("新旧密码不能相同")
	}

	if profile.Confirm != profile.NewPassword {
		return errors.New("确认密码不一致")
	}

	//校验密码
	oldPwd, _ := lv_secret.PasswordHash(profile.OldPassword)
	if oldPwd != user.Password {
		return errors.New("原密码不正确")
	}

	//新校验密码
	newPwd := lv_gen.GenerateSubId(6)
	newPwd, _ = lv_secret.PasswordHash(newPwd)
	user.Password = newPwd

	err := user.Updates()
	if err != nil {
		return errors.New("保存数据失败")
	}

	//SaveUserToSession(user, c)
	return nil
}

func (svc *UserService) ResetPassword(params *common_vo.ResetPwdReq) error {
	user := model.SysUser{UserId: params.UserId}
	if err := user.FindOne(); err != nil {
		return errors.New("用户不存在")
	}
	//新校验密码
	newPwd := lv_gen.GenerateSubId(6)
	newPwd, _ = lv_secret.PasswordHash(newPwd)
	user.Password = newPwd
	err := user.Updates()
	return err
}

func (svc *UserService) CheckPassword(user *model.SysUser, password string) bool {
	if user == nil || user.UserId <= 0 {
		return false
	}
	//校验密码
	pwd, _ := lv_secret.PasswordHash(password)

	if strings.Compare(pwd, user.Password) == 0 {
		return true
	} else {
		return false
	}
}

// 根据登录名查询用户信息
func (svc *UserService) SelectUserByUserName(UserName string) (*model.SysUser, error) {
	var vo dao.SysUserDao
	return vo.SelectUserByUserName(UserName)
}

// 根据手机号查询用户信息
func (svc *UserService) SelectUserByPhoneNumber(phonenumber string) (*model.SysUser, error) {
	var vo dao.SysUserDao
	return vo.SelectUserByPhoneNumber(phonenumber)
}

// 查询已分配用户角色列表
func (svc *UserService) SelectAllocatedList(roleId int64, UserName, phonenumber string) (*[]map[string]any, error) {
	var vo dao.SysUserDao
	return vo.SelectAllocatedList(roleId, UserName, phonenumber)
}

// 查询未分配用户角色列表
func (svc *UserService) SelectUnallocatedList(roleId int64, UserName, phonenumber string) (*[]map[string]any, error) {
	var vo dao.SysUserDao
	return vo.SelectUnallocatedList(roleId, UserName, phonenumber)
}

// 查询未分配用户角色列表
func (svc *UserService) GetRoleKeys(userId int64) (string, error) {
	if userId == 1 {
		return "admin", nil
	}
	var sql = " SELECT GROUP_CONCAT(r.role_key) roles from sys_user_role ur,sys_role r where ur.user_id=? and ur.role_id = r.role_id "
	var roles string
	err := lv_db.GetMasterGorm().Raw(sql, userId).Scan(&roles).Error
	return roles, err
}

func (svc *UserService) GetRoles(userId int64) ([]model.SysRole, error) {
	sql := " select r.* from sys_user_role ur,sys_role r where ur.user_id=? and ur.role_id = r.role_id "
	roles := make([]model.SysRole, 0)
	err := lv_db.GetMasterGorm().Raw(sql, userId).Scan(&roles).Error
	return roles, err
}

func (svc *UserService) CountCol(column, value string) (int64, error) {
	var total int64
	err := lv_db.GetMasterGorm().Table("sys_user").Where("del_flag=0 and "+column+"=?", value).Count(&total).Error
	return total, err
}
