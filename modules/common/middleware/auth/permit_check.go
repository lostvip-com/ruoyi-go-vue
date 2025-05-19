package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/lostvip-com/lv_framework/web/router"
	"net/http"
	"strings"
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
	//获取用户信息
	userSvc := service.GetUserServiceInstance()
	userPtr := userSvc.GetCurrUser(c)
	c.Set("userId", userPtr.UserId) //供api使用
	c.Set("user", userPtr)          //供api使用
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
	//获取用户菜单列表
	menuSvc := service.GetMenuServiceInstance()
	hasPermission, err := menuSvc.IsRolePermited(userPtr.GetRoleKeys(), permission)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, lv_dto.CommonRes{
			Code: http.StatusForbidden,
			Msg:  "Operate Forbidden",
		})
		return
	}
	//fmt.Println("url:" + url + "---permission:" + permission)
	if !hasPermission {
		c.AbortWithStatusJSON(http.StatusForbidden, lv_dto.CommonRes{
			Code: http.StatusForbidden,
			Msg:  "Operate Forbidden",
		})
	}

	c.Next()
}
