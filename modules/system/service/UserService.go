package service

import (
	"common/common_vo"
	"common/global"
	"common/util"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_cache"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
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

func (svc *UserService) FindById(id int) (*model.SysUser, error) {
	entity := &model.SysUser{UserId: id}
	err := entity.FindOne()
	return entity, err
}

func (svc *UserService) FindList(param *common_vo.UserPageReq) (*[]map[string]any, int, error) {
	var deptService DeptService
	var dept, _ = deptService.FindById(param.DeptId)
	if dept != nil { //数据权限
		param.Ancestors = dept.Ancestors
	}
	var d = dao.GetUserDaoInstance()
	return d.FindPage(param)
}

// 新增用户
func (svc *UserService) AddSave(req *common_vo.AddUserReq, c *gin.Context) (int, error) {
	var u model.SysUser
	u.UserName = req.UserName
	u.UserName = req.UserName
	u.Email = req.Email
	u.Phonenumber = req.Phonenumber
	u.Status = req.Status
	u.Sex = req.Sex
	u.DeptId = req.DeptId
	u.Remark = req.Remark
	//u.LoginDate = time.Now()
	//生成密码
	//newSalt := lv_gen.GenerateSubId(6)
	u.Password, _ = lv_secret.PasswordHash(u.Password)
	//u.CreateTime = u.LoginDate
	createUser := svc.GetCurrUser(c)

	if createUser != nil {
		u.CreateBy = createUser.UserName
	}
	u.DelFlag = "0"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := lv_db.GetOrmDefault().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&u).Error; err != nil {
			return err
		}
		//增加岗位数据
		if len(req.PostIds) > 0 {
			userPosts := make([]model.SysUserPost, 0)
			for i := range req.PostIds {
				if req.PostIds[i] > 0 {
					var userPost model.SysUserPost
					userPost.UserId = u.UserId
					userPost.PostId = req.PostIds[i]
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
		if len(req.RoleIds) > 0 {
			userRoles := make([]model.SysUserRole, 0)
			for i := range req.RoleIds {
				if req.RoleIds[i] > 0 {
					var userRole model.SysUserRole
					userRole.UserId = u.UserId
					userRole.RoleId = req.RoleIds[i]
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
	//userPtr.UpdateTime = time.Now()
	updateUser := svc.GetCurrUser(c)

	if updateUser != nil {
		userPtr.UpdateBy = updateUser.UserName
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = lv_db.GetOrmDefault().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(userPtr).Error; err != nil {
			return err
		}
		//增加岗位数据
		if len(req.PostIds) > 0 {
			userPosts := make([]model.SysUserPost, 0)
			for i := range req.PostIds {
				if req.PostIds[i] > 0 {
					var userPost model.SysUserPost
					userPost.UserId = userPtr.UserId
					userPost.PostId = req.PostIds[i]
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
		if len(req.RoleIds) > 0 {
			userRoles := make([]model.SysUserRole, 0)
			for i := range req.RoleIds {
				if req.RoleIds[i] > 0 {
					var userRole model.SysUserRole
					userRole.UserId = userPtr.UserId
					userRole.RoleId = req.RoleIds[i]
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
func (svc *UserService) DeleteById(id int) error {
	entity := &model.SysUser{UserId: id}
	err := entity.Delete()
	return err
}

// 批量删除用户记录
func (svc *UserService) DeleteByIds(ids string) error {
	idarr := util.ToIntArray(ids, ",")
	idarr = util.RemoveOne(idarr, 1) //去掉admin的id
	if len(idarr) == 0 {
		return errors.New("ids can not be empty ")
	}
	err := lv_db.GetOrmDefault().Transaction(func(tx *gorm.DB) error {
		err := tx.Table("sys_user_post").Delete("user_id in ? and user_id!=1 ", idarr).Error
		if err != nil {
			return err
		}
		err = tx.Table("sys_user_role").Delete("user_id in ? ? and user_id!=1 ", idarr).Error
		if err != nil {
			return err
		}
		err = tx.Table("sys_user").Where("user_id in ? and user_id!=1 ", idarr).Update("del_flag", 1).Error
		return err
	})
	return err
}

// 判断是否是系统管理员
func (svc *UserService) IsAdmin(userId int) bool {
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
	tokenId := c.GetString("tokenId")
	user, err := svc.GetProfile(tokenId)
	lv_err.HasErrAndPanic(err)
	return user
}

// 获得用户信息详情
func (svc *UserService) GetProfile(tokenId string) (*model.SysUser, error) {
	key := global.LoginCacheKey + tokenId
	userId, err := lv_cache.GetCacheClient().HGet(key, "userId")
	lv_err.HasErrAndPanic(err)
	u := new(model.SysUser)
	u.UserId = cast.ToInt(userId)
	err = u.FindOne()
	if err != nil {
		return nil, err
	}
	// 部门
	u.Dept, err = GetDeptServiceInstance().FindById(u.DeptId)
	u.Roles, err = dao.GetRoleDaoInstance().FindRoles(u.UserId)
	return u, err
}

// 更新用户信息详情
func (svc *UserService) UpdateProfile(profile *common_vo.ProfileReq, c *gin.Context) error {
	user := svc.GetCurrUser(c)
	if profile.NickName != "" {
		user.NickName = profile.NickName
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

	if profile.NewPassword == profile.OldPassword {
		return errors.New("新旧密码不能相同")
	}
	//校验密码
	oldPwd, _ := lv_secret.PasswordHash(profile.OldPassword)
	//校验密码
	if !lv_secret.PasswordVerify(oldPwd, user.Password) {
		return errors.New("原密码错误")
	}
	//新校验密码
	newPwd, _ := lv_secret.PasswordHash(profile.NewPassword)
	user.Password = newPwd
	err := user.Updates()
	return err
}

func (svc *UserService) ResetPassword(params *common_vo.ResetPwdReq) error {
	user := model.SysUser{UserId: params.UserId}
	if err := user.FindOne(); err != nil {
		return errors.New("用户不存在")
	}
	//新校验密码
	//newPwd := lv_gen.GenerateSubId(6)
	newPwd, _ := lv_secret.PasswordHash(params.Password)
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
func (svc *UserService) SelectAllocatedList(roleId int, UserName, phonenumber string) (*[]map[string]any, error) {
	var vo dao.SysUserDao
	return vo.SelectAllocatedList(roleId, UserName, phonenumber)
}

// 查询未分配用户角色列表
func (svc *UserService) SelectUnallocatedList(roleId int, UserName, phonenumber string) (*[]map[string]any, error) {
	var vo dao.SysUserDao
	return vo.SelectUnallocatedList(roleId, UserName, phonenumber)
}

// 查询未分配用户角色列表
func (svc *UserService) GetRoleKeys(userId int) (string, error) {
	if userId == 1 {
		return "admin", nil
	}
	var sql = " SELECT GROUP_CONCAT(r.role_key) roles from sys_user_role ur,sys_role r where ur.user_id=? and ur.role_id = r.role_id "
	var roles string
	err := lv_db.GetOrmDefault().Raw(sql, userId).Scan(&roles).Error
	return roles, err
}

func (svc *UserService) GetRoles(userId int) ([]model.SysRole, error) {
	sql := " select r.* from sys_user_role ur,sys_role r where ur.user_id=? and ur.role_id = r.role_id "
	roles := make([]model.SysRole, 0)
	err := lv_db.GetOrmDefault().Raw(sql, userId).Scan(&roles).Error
	return roles, err
}

func (svc *UserService) CountCol(column, value string) (int, error) {
	var total int64
	err := lv_db.GetOrmDefault().Table("sys_user").Where("del_flag=0 and "+column+"=?", value).Count(&total).Error
	return int(total), err
}
