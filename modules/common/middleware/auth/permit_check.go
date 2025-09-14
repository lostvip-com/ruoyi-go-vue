package auth

import (
	"common/global"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/lostvip-com/lv_framework/web/router"
	"net/http"
	"strings"
	"system/model"
	"system/service"
)

// 鉴权中间件，只有登录成功之后才能通过
func PermitCheck(c *gin.Context) {
	//判断是否登录
	//根据url判断是否有权限
	url := c.Request.URL.Path
	strEnd := url[len(url)-1 : len(url)]
	if strings.EqualFold(strEnd, "/") {
		url = strings.TrimRight(url, "/")
	}
	//获取用户菜单列表
	u, ok := c.Get(global.KEY_GIN_USER_PTR)
	if !ok {
		util.Fail(c, "获取用户信息失败")
	}
	userPtr := u.(*model.SysUser)
	//获取用户信息
	userSvc := service.GetUserServiceInstance()
	lv_log.Debug("=====================token check ", userPtr, ok)
	if userSvc.IsAdmin(userPtr.UserId) {
		c.Next()
		return
	}
	//获取权限标识
	permission := router.FindPermission(url)
	if permission == "" {
		c.Next()
		return
	}

	hasPermission, err := service.GetMenuServiceInstance().IsRolePermited(userPtr.GetRoleKeys(), permission)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, lv_dto.CommonRes{Code: http.StatusForbidden, Msg: "Operate Forbidden"})
		return
	}
	//fmt.Println("url:" + url + "---permission:" + permission)
	if !hasPermission {
		c.AbortWithStatusJSON(http.StatusForbidden, lv_dto.CommonRes{Code: http.StatusForbidden, Msg: "Operate Forbidden"})
	}
	c.Next()
}
