package service

import (
	"common/global"
	"github.com/lostvip-com/lv_framework/lv_cache"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/spf13/cast"
	"system/model"
	"system/vo"
	"time"
)

const USER_NOPASS_TIME string = "user_nopass_"
const USER_LOCK string = "user_lock_"

type LoginService struct {
}

var loginService *LoginService

func GetLoginServiceInstance() *LoginService {
	if loginService == nil {
		loginService = &LoginService{}
	}
	return loginService
}

// FindPage 根据条件分页查询用户列表
func (svc LoginService) FindPage(param *vo.LoginInfoPageReq) (*[]model.SysLoginInfo, int64, error) {
	db := lv_db.GetMasterGorm()
	tb := db.Table("sys_logininfor")
	if param != nil {
		if param.UserName != "" {
			tb.Where("user_name like ?", "%"+param.UserName+"%")
		}
		if param.Ipaddr != "" {
			tb.Where("ipaddr like ?", "%"+param.Ipaddr+"%")
		}
		if param.Status != "" {
			tb.Where("status = ?", param.Status)
		}

		if param.BeginTime != "" {
			tb.Where("login_time >= ?", param.BeginTime)
		}
		if param.EndTime != "" {
			tb.Where("login_time <= ?", param.EndTime)
		}
	}
	var total int64
	tb = tb.Count(&total).Offset(param.GetStartNum()).Limit(param.GetPageSize())
	if param.OrderByColumn != "" {
		tb.Order("info_id " + param.IsAsc)
	}
	var result []model.SysLoginInfo
	err := tb.Find(&result).Error
	return &result, total, err
}

// FindById 根据主键查询用户信息
func (svc LoginService) FindById(id int64) (*model.SysLoginInfo, error) {
	entity := &model.SysLoginInfo{InfoId: id}
	err := entity.FindById()
	return entity, err
}

// DeleteById 根据主键删除用户信息
func (svc LoginService) DeleteById(id int64) bool {
	entity := &model.SysLoginInfo{InfoId: id}
	err := entity.Delete()
	if err == nil {
		return true
	}

	return false
}

// DeleteByIds 批量删除记录
func (svc LoginService) DeleteByIds(ids string) error {
	idarr := lv_conv.ToInt64Array(ids, ",")
	err := lv_db.GetMasterGorm().Exec("delete from sys_logininfor where info_id in ? ", idarr).Error
	return err
}

// 清空记录
func (svc LoginService) DeleteRecordAll() error {
	db := lv_db.GetMasterGorm()
	err := db.Exec(" truncate table sys_logininfor ").Error
	return err
}

// 导出excel
func (svc LoginService) Export(param *vo.LoginInfoPageReq) (string, error) {
	//head := []string{"访问编号", "登录名称", "登录地址", "登录地点", "浏览器", "操作系统", "登录状态", "操作信息", "登录时间"}
	//col := []string{"info_id", "user_name", "ipaddr", "login_location", "browser", "os", "status", "msg", "login_time"}
	//db := lv_db.GetMasterGorm()
	//build := builder.Select(col...).From(" sys_logininfor ", "t")
	//if param != nil {
	//	if param.UserName != "" {
	//		build.Where(builder.Like{"t.user_name", param.UserName})
	//	}
	//
	//	if param.Ipaddr != "" {
	//		build.Where(builder.Like{"t.ipaddr", param.Ipaddr})
	//	}
	//
	//	if param.Status != "" {
	//		build.Where(builder.Eq{"t.status": param.Status})
	//	}
	//
	//	if param.BeginTime != "" {
	//		build.Where(builder.Gte{"date_format(t.create_time,'%y%m%d')": "date_format('" + param.BeginTime + "','%y%m%d')"})
	//	}
	//
	//	if param.EndTime != "" {
	//		build.Where(builder.Lte{"date_format(t.create_time,'%y%m%d')": "date_format('" + param.EndTime + "','%y%m%d')"})
	//	}
	//}
	//sqlStr, err := build.ToBoundSQL()
	//arr, err := namedsql.ListArrStr(db, sqlStr, nil)
	//path, err := lv_office.DownlaodExcel(head, *arr)
	return "", nil
}

// 记录密码尝试次数
func (svc LoginService) SetPasswordCounts(UserName string) int {
	curTimes := 0
	curTimeObj, err := lv_cache.GetCacheClient().Get(USER_NOPASS_TIME + UserName)
	if err == nil {
		curTimes = cast.ToInt(curTimeObj)
	}
	curTimes = curTimes + 1
	lv_cache.GetCacheClient().Set(USER_NOPASS_TIME+UserName, curTimes, 1*time.Minute)

	if curTimes >= global.ErrTimes2Lock {
		svc.Lock(UserName)
	}
	return curTimes
}

// 记录密码尝试次数
func (svc LoginService) GetPasswordCounts(UserName string) int {
	curTimes := 0
	curTimeObj, err := lv_cache.GetCacheClient().Get(USER_NOPASS_TIME + UserName)
	if err != nil {
		curTimes = cast.ToInt(curTimeObj)
	}
	return curTimes
}

// 移除密码错误次数
func (svc LoginService) RemovePasswordCounts(UserName string) {
	lv_cache.GetCacheClient().Del(USER_NOPASS_TIME + UserName)
}

// 锁定账号
func (svc LoginService) Lock(UserName string) {
	lv_cache.GetCacheClient().Set(USER_LOCK+UserName, true, 30*time.Minute)
}

// 解除锁定
func (svc LoginService) Unlock(UserName string) {
	lv_cache.GetCacheClient().Del(USER_LOCK + UserName)
}

// 检查账号是否锁定
func (svc LoginService) CheckLock(UserName string) bool {
	result := false
	rs, _ := lv_cache.GetCacheClient().Get(USER_LOCK + UserName)
	if rs != "" {
		result = true
	}
	return result
}
