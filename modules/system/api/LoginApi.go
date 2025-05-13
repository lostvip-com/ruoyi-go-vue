package api

import (
	global2 "common/global"
	"common/middleware/auth"
	util2 "common/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_net"
	"github.com/lostvip-com/lv_framework/utils/lv_try"
	"github.com/mssola/user_agent"
	"system/model"
	"system/service"

	"net/http"
	"time"
)

type LoginApi struct {
}
type RegisterReq struct {
	UserName string `form:"username"  binding:"required,min=4,max=30"`
	Password string `form:"password" binding:"required,min=6,max=30"`
	//
	ValidateCode string `form:"validateCode" `
	IdKey        string `form:"idkey" `
}

/**
 * 调用公共服务提供的验证方法
 */
func (w *LoginApi) VerifyCaptcha(IdKey string, ValidateCode string) (string, error) {
	return "", nil
}

// 验证登录
func (w *LoginApi) Login(c *gin.Context) {
	var req = RegisterReq{}

	if err := c.ShouldBind(&req); err != nil {
		util2.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	clientIp := lv_net.GetRemoteClientIp(c.Request)
	var loginSvc = service.GetLoginServiceInstance()
	errTimes4Ip := loginSvc.GetPasswordCounts(clientIp)
	if errTimes4Ip > 5 { //超过5次错误开始校验验证码
		//比对验证码
		_, err := w.VerifyCaptcha(req.IdKey, req.ValidateCode)
		if err != nil {
			util2.ErrorResp(c).SetData(errTimes4Ip).SetMsg("验证码不正确").WriteJsonExit()
			return
		}
	}
	isLock := loginSvc.CheckLock(req.UserName)
	if isLock {
		util2.ErrorResp(c).SetMsg("账号已锁定，请30分钟后再试").WriteJsonExit()
		return
	}
	var userService = service.GetUserServiceInstance()
	//验证账号密码
	user, err := service.GetSessionServiceInstance().SignIn(req.UserName, req.Password)
	if err != nil {
		loginSvc.SetPasswordCounts(clientIp)
		errTimes4UserName := loginSvc.SetPasswordCounts(req.UserName)
		having := global2.ErrTimes2Lock - errTimes4UserName
		w.SaveLogs(c, &req, "账号或密码不正确") //记录日志
		if having <= 5 {
			util2.Fail(c, "账号或密码不正确,还有"+lv_conv.String(having)+"次之后账号将锁定")
		} else {
			util2.Fail(c, "账号或密码不正确")
		}
		return
	}
	//保存在线状态
	newUUID, _ := uuid.NewUUID()
	tokenId := newUUID.String()
	token := auth.CreateToken(user.UserName, user.UserId, user.DeptId, tokenId)
	// 生成token
	ua := c.Request.Header.Get("User-Agent")
	roles, err := userService.GetRoleKeys(user.UserId)
	lv_err.HasErrAndPanic(err)
	var svc service.SessionService
	ip := lv_net.GetRemoteClientIp(c.Request)
	err = svc.SaveUserToSession(tokenId, user, roles)
	go func() {
		err := lv_try.Catch(func() {
			w.SaveLogs(c, &req, "login success") //记录日志
			err = svc.SaveLoginLog2DB(token, user, ua, ip)
		})
		if err != nil {
			lv_log.Error(err.Error())
		}
	}()
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"code": 200, "msg": "success", "token": token})
}

func (w *LoginApi) SaveLogs(c *gin.Context, req *RegisterReq, msg string) {
	var loginInfo model.SysLoginInfo
	loginInfo.UserName = req.UserName
	loginInfo.Ipaddr = c.ClientIP()
	userAgent := c.Request.Header.Get("User-Agent")
	ua := user_agent.New(userAgent)
	loginInfo.Os = ua.OS()
	loginInfo.Browser, _ = ua.Browser()
	loginInfo.LoginTime = time.Now()
	loginInfo.LoginLocation = util2.GetCityByIp(loginInfo.Ipaddr)
	loginInfo.Msg = msg
	loginInfo.Status = "0"
	loginInfo.Insert()
}

// 注销
func (w *LoginApi) Logout(c *gin.Context) {
	var user service.SessionService
	tokenStr := lv_net.GetParam(c, "token")
	err := user.SignOut(tokenStr)
	if err != nil {
		return
	}
	util2.Success(c, nil)
}
