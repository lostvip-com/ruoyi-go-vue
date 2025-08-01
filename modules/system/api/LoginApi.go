package api

import (
	global2 "common/global"
	"common/middleware/auth"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_net"
	"github.com/mssola/user_agent"
	"github.com/spf13/cast"
	"strings"
	"system/model"
	"system/service"
	"time"

	"net/http"
)

type LoginApi struct {
}

// LoginParam 登录参数
type LoginReq struct {
	Code     string `json:"code"`
	Password string `json:"password"`
	UserName string `json:"username"`
	Uuid     string `json:"uuid"`
}

/**
 * 调用公共服务提供的验证方法
 */
func (w *LoginApi) VerifyCaptcha(IdKey string, ValidateCode string) (string, error) {
	return "", nil
}

// 验证登录
func (w *LoginApi) Login(c *gin.Context) {
	var req = LoginReq{}
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	clientIp := lv_net.GetRemoteClientIp(c.Request)
	var loginSvc = service.GetLoginServiceInstance()
	errTimes4Ip := loginSvc.GetPasswordCounts(clientIp)
	if errTimes4Ip > 5 { //超过5次错误开始校验验证码
		_, err := w.VerifyCaptcha(req.Uuid, req.Code)
		if err != nil {
			util.Fail(c, "验证码不正确")
			return
		}
	}
	isLock := loginSvc.CheckLock(req.UserName)
	if isLock {
		util.Fail(c, "账号已锁定，请30分钟后再试")
		return
	}
	user, err := service.GetSessionServiceInstance().SignIn(req.UserName, req.Password)
	lv_err.HasErrAndPanic(err)
	roles, err := service.GetUserServiceInstance().GetRoleKeys(user.UserId)
	lv_err.HasErrAndPanic(err)
	newUUID, _ := uuid.NewUUID()
	tokenId := newUUID.String()
	token := auth.CreateToken(user.UserName, user.UserId, user.DeptId, tokenId)
	var svc = service.GetSessionServiceInstance()
	loginInfo := new(model.SysLoginInfo)
	loginInfo.UserName = req.UserName
	loginInfo.Ipaddr = lv_net.GetRemoteClientIp(c.Request)
	userAgent := c.Request.Header.Get("User-Agent")
	ua := user_agent.New(userAgent)
	loginInfo.Os = ua.OS()
	loginInfo.Browser, _ = ua.Browser()
	loginInfo.LoginTime = time.Now()
	loginInfo.LoginLocation = util.GetCityByIp(loginInfo.Ipaddr)
	if err != nil {
		loginSvc.SetPasswordCounts(clientIp)
		errTimes4UserName := loginSvc.SetPasswordCounts(req.UserName)
		having := global2.ErrTimes2Lock - errTimes4UserName
		svc.SaveLogs(loginInfo, "账号或密码不正确") //记录日志
		if having <= 5 {
			util.Fail(c, "账号或密码不正确,还有"+lv_conv.String(having)+"次之后账号将锁定")
		} else {
			util.Fail(c, "账号或密码不正确")
		}
		return
	}
	loginInfo.Status = "0" //成功
	dept, err := service.GetDeptServiceInstance().FindById(user.DeptId)
	lv_err.HasErrAndPanic(err)
	svc.SaveSessionToRedis(tokenId, roles, dept.DeptName, user)
	svc.SaveLogs(loginInfo, "login success") //记录日志
	if err != nil {
		lv_log.Error(err.Error())
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "token": token})
}

// 注销
func (w *LoginApi) Logout(c *gin.Context) {
	tokenStr := c.Request.Header.Get(auth.Authorization)
	if len(tokenStr) <= 0 {
		util.Fail(c, "token is null")
		return
	}
	mySigningKey := []byte(auth.Secret)
	tokenStr = strings.ReplaceAll(tokenStr, auth.Bearer, "")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil { //token 过期这里会出现错误，只记录即可
		lv_log.Error(err.Error())
	}
	if token != nil {
		uuidStr := token.Claims.(jwt.MapClaims)["uuid"]
		var user service.SessionService
		err = user.SignOut(cast.ToString(uuidStr))
		util.Success(c, nil)
	} else {
		util.Fail(c, "bad token！")
	}

}
