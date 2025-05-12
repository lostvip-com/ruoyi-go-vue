package dao

import (
	"common/common_vo"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/namedsql"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/spf13/cast"
	"system/model"
)

// Fill with you ideas below.

// 修改用户资料请求参数
type SysUserDao struct {
}

var userDao *SysUserDao

func GetUserDaoInstance() *SysUserDao {
	if userDao == nil {
		userDao = new(SysUserDao)
	}
	return userDao
}
func (e SysUserDao) DeleteByIds(ida []int64) int64 {
	db := lv_db.GetMasterGorm().Table("sys_user").Where("user_id in ? and user_id!=1 ", ida).Update("del_flag", 1)
	if db.Error != nil {
		panic(db.Error)
	}
	return db.RowsAffected
}

// 根据条件分页查询用户列表
func (d SysUserDao) FindPage(param *common_vo.UserPageReq) (*[]map[string]string, int64, error) {
	db := lv_db.GetMasterGorm()
	sqlParams, sql := d.GetSql(param)
	lv_log.Info("============sqlParams:", sqlParams)
	limitSql := sql + " order by u.user_id desc "
	limitSql += "  limit " + cast.ToString(param.GetStartNum()) + "," + cast.ToString(param.GetPageSize())
	result, err := namedsql.ListMapStr(db, limitSql, sqlParams, true)
	countSql := "select count(*) from (" + sql + ") t "
	total, err := namedsql.Count(db, countSql, sqlParams)
	return result, total, err
}

func (d SysUserDao) GetSql(param *common_vo.UserPageReq) (map[string]interface{}, string) {
	sqlParams := make(map[string]interface{})
	sql := `
            select u.user_id, u.dept_id, u.user_name, u.nick_name, u.email, u.avatar, u.phonenumber, u.password,u.sex, u.status, u.del_flag, 
            u.login_ip, u.login_date, u.create_by, u.create_time, u.remark,d.dept_name, d.leader
            from sys_user u left join sys_dept d on  u.dept_id = d.dept_id where u.del_flag =0 
           `
	if param != nil {
		if param.UserName != "" {
			sql += " and  u.user_name like @UserName "
			sqlParams["UserName"] = "%" + param.UserName + "%"
		}

		if param.Phonenumber != "" {
			sql += " and  u.phonenumber like @phonenumber "
			sqlParams["phonenumber"] = "%" + param.Phonenumber + "%"
		}

		if param.Status != "" {
			sql += " and  u.status = @status "
			sqlParams["status"] = param.Status
		}
		if param.TenantId != 0 {
			if param.TenantId != 0 {
				sql += " and  u.tenant_id = @TenantId "
				sqlParams["TenantId"] = param.TenantId
			}
		}
		if param.BeginTime != "" {
			sql += " and  u.create_time >= @BeginTime "
			sqlParams["BeginTime"] = param.BeginTime
		}
		if param.EndTime != "" {
			sql += " and  u.create_time <= @EndTime "
			sqlParams["EndTime"] = param.EndTime
		}
		if param.Ancestors != "" {
			sql += " and u.dept_id IN ( SELECT t.dept_id FROM sys_dept t WHERE t.ancestors like @Ancestors )"
			sqlParams["Ancestors"] = param.Ancestors + "%"
			lv_log.Info("============sqlParams:", sqlParams)
		}
		lv_log.Info("============sqlParams:", sqlParams)
	}
	lv_log.Info("============sqlParams:", sqlParams)
	return sqlParams, sql
}

// 导出excel
func (d SysUserDao) SelectExportList(param *common_vo.UserPageReq) (*[]map[string]any, error) {
	db := lv_db.GetMasterGorm()
	sqlParams, sql := d.GetSql(param)
	limitSql := sql + " order by u.user_id desc "
	result, err := namedsql.ListMap(db, limitSql, &sqlParams, true)
	return result, err
}

// 根据条件分页查询已分配用户角色列表
func (d SysUserDao) SelectAllocatedList(roleId int64, UserName, phonenumber string) (*[]map[string]string, error) {
	db := lv_db.GetMasterGorm()
	sqlParams := make(map[string]interface{})
	sql := `
            select distinct u.user_id, u.dept_id, u.user_name, u.nick_name, u.email, u.avatar, u.phonenumber,u.status, u.create_time
            from sys_user u 
            left join sys_dept d on  u.dept_id = d.dept_id 
            left join sys_user_role ur on  u.user_id = ur.user_id
             left join sys_role r on  r.role_id = ur.role_id
            where u.del_flag =0 and  r.role_id = ` + cast.ToString(roleId)

	if UserName != "" {
		sql += " and u.user_name like @UserName "
		sqlParams["UserName"] = "%" + UserName + "%"
	}

	if phonenumber != "" {
		sql += " and  u.phonenumber like @phonenumber "
		sqlParams["phonenumber"] = "%" + phonenumber + "%"
	}

	result, err := namedsql.ListMapStr(db, sql, &sqlParams, true)
	return result, err
}

// 根据条件分页查询未分配用户角色列表
func (d SysUserDao) SelectUnallocatedList(roleId int64, UserName, phonenumber string) (*[]map[string]string, error) {
	db := lv_db.GetMasterGorm()
	sqlParams := make(map[string]interface{})
	sql := `
            select distinct u.user_id, u.dept_id, u.user_name, u.nick_name, u.email, u.avatar, u.phonenumber,u.status, u.create_time
            from sys_user u 
            where u.del_flag =0  `
	sql += " and u.user_id not in (select u.user_id from sys_user u inner join sys_user_role ur on u.user_id = ur.user_id and ur.role_id = " + cast.ToString(roleId) + ") "
	if UserName != "" {
		sql += " and u.user_name like @UserName "
		sqlParams["UserName"] = "%" + UserName + "%"
	}

	if phonenumber != "" {
		sql += " and  u.phonenumber like @phonenumber "
		sqlParams["phonenumber"] = "%" + phonenumber + "%"
	}

	result, err := namedsql.ListMapStr(db, sql, &sqlParams, true)
	return result, err
}

// CountPhone 检查手机号是否已使用 ,存在返回true,否则false
func (d SysUserDao) CountPhone(phone string) (int64, error) {
	db := lv_db.GetMasterGorm()
	var total int64
	err := db.Table("sys_user").Where("phonenumber = ?", phone).Count(&total).Error
	return total, err
}

// SelectUserByUserName 根据登录名查询用户信息
func (d SysUserDao) SelectUserByUserName(UserName string) (*model.SysUser, error) {
	var entity = new(model.SysUser)
	entity.UserName = UserName
	err := entity.FindOne()
	return entity, err
}

// 根据手机号查询用户信息
func (d SysUserDao) SelectUserByPhoneNumber(phonenumber string) (*model.SysUser, error) {
	var entity = new(model.SysUser)
	entity.Phonenumber = phonenumber
	err := entity.FindOne()
	return entity, err
}
