package service

import (
	"common/global"
	"errors"
	"github.com/lostvip-com/lv_framework/lv_cache"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_json"
	"github.com/lostvip-com/lv_framework/utils/lv_secret"
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
	loginKey := global.LoginCacheKey + uuid
	num, err := lv_cache.GetCacheClient().Exists(loginKey)
	return err == nil && num > 0
}

func (svc *SessionService) SignIn(loginnName, password string) (*model.SysUser, error) {
	//查询用户信息
	user := model.SysUser{UserName: loginnName}
	err := user.FindOne()

	if err != nil {
		return nil, err
	}
	pwd, _ := lv_secret.PasswordHash(password)
	lv_log.Error("------------" + pwd)
	//校验密码
	if !lv_secret.PasswordVerify(password, user.Password) {
		return nil, errors.New("密码错误")
	}
	return &user, nil
}

func (svc *SessionService) SignOut(uuid string) error {
	return lv_cache.GetCacheClient().Del(global.LoginCacheKey + uuid)
}

func (svc *SessionService) ForceLogout(token string) error {
	return svc.SignOut(token)
}

func (svc *SessionService) SaveSessionToRedis(loginInfo *model.SysLoginInfo) error {
	fieldMap := lv_json.StructToMap(loginInfo)
	fieldMap["username"] = loginInfo.UserName
	key := global.LoginCacheKey + loginInfo.TokenId
	err := lv_cache.GetCacheClient().HMSet(key, fieldMap, 12*time.Hour)
	lv_err.HasErrAndPanic(err)
	err = lv_cache.GetCacheClient().Expire(key, 12*time.Hour)
	return err
}

func (svc *SessionService) Refresh(token string) {
	token = "login:" + token
	lv_cache.GetCacheClient().Expire(token, 8*time.Hour)
}

func (svc *SessionService) SaveLogs(login *model.SysLoginInfo, msg string) error {
	login.Msg = msg
	return login.Insert()
}
