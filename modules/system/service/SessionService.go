package service

import (
	"common/global"
	"errors"
	"github.com/lostvip-com/lv_framework/lv_cache"
	"github.com/lostvip-com/lv_framework/utils/lv_secret"
	"strings"
	"system/model"
	"time"
)

type SessionService struct{}

var sessionService *SessionService

func GetSessionServiceInstance() *SessionService {
	if sessionService == nil {
		sessionService = &SessionService{}
	}
	return sessionService
}
func (svc *SessionService) IsSignedIn(uuid string) bool {
	loginKey := global.LOGIN_TOKEN_KEY + uuid
	num, err := lv_cache.GetCacheClient().Exists(loginKey)
	return err == nil && num > 0
}

// 用户登录，成功返回用户信息，否则返回nil; passport应当会md5值字符串
func (svc *SessionService) SignIn(loginnName, password string) (*model.SysUser, error) {
	//查询用户信息
	user := model.SysUser{UserName: loginnName}
	err := user.FindOne()

	if err != nil {
		return nil, err
	}

	//校验密码
	pwdNew := user.UserName + password + user.Salt

	pwdNew = lv_secret.MustEncryptString(pwdNew)

	if strings.Compare(user.Password, pwdNew) == -1 {
		return nil, errors.New("密码错误")
	}
	return &user, nil
}

// SignOut 用户注销
func (svc *SessionService) SignOut(tokenStr string) error {
	return lv_cache.GetCacheClient().Del(global.LOGIN_TOKEN_KEY + tokenStr)
}

// ForceLogout 强退用户
func (svc *SessionService) ForceLogout(token string) error {
	return svc.SignOut(token)
}

func (svc *SessionService) SaveUserToSession(tokenId string, user *model.SysUser, roleKeys string) error {
	//记录到redis
	fieldMap := make(map[string]interface{})
	fieldMap["userName"] = user.UserName
	fieldMap["userId"] = user.UserId
	fieldMap["UserName"] = user.UserName
	fieldMap["avatar"] = user.Avatar
	fieldMap["roleKeys"] = roleKeys
	fieldMap["deptId"] = user.DeptId
	//fieldMap["tenantId"] = user.TenantId //租户ID
	//其它
	key := global.LOGIN_TOKEN_KEY + tokenId
	err := lv_cache.GetCacheClient().HSet(key, fieldMap)
	if err != nil {
		panic("redis 故障！" + err.Error())
	}
	err = lv_cache.GetCacheClient().Expire(key, time.Hour)
	if err != nil {
		panic("redis 故障！" + err.Error())
	}
	return err
}

func (svc *SessionService) SaveLoginLog2DB(token string, user *model.SysUser, userAgent string, ip string) error {
	//save to lv_db
	//ua := user_agent.New(userAgent)
	//os := ua.OS()
	//browser, _ := ua.Browser()
	////移除登录次数记录
	//var logininforService LoginService
	//logininforService.RemovePasswordCounts(user.UserName)
	//
	//var userOnline vo.OnlineVo
	//// 保存用户信息到session
	//loginLocation := util.GetCityByIp(ip)
	//userOnline.UserName = user.UserName
	//userOnline.Browser = browser
	//userOnline.Os = os
	//userOnline.DeptName = ""
	//userOnline.Ipaddr = ip
	//userOnline.StartTimestamp = lv_time.GetCurrentTimeStr()
	//userOnline.LastAccessTime = lv_time.GetCurrentTimeStr()
	//userOnline.CreateTime = userOnline.StartTimestamp
	//userOnline.Status = "on_line"
	//userOnline.LoginLocation = loginLocation
	//userOnline.SessionId = token
	//err := userOnline.Delete()
	//if err != nil {
	//	return err
	//}
	//err = userOnline.Save()
	//if err != nil {
	//	return err
	//}
	return nil
}

func (svc *SessionService) Refresh(token string) {
	token = "login:" + token
	lv_cache.GetCacheClient().Expire(token, 8*time.Hour)
}
