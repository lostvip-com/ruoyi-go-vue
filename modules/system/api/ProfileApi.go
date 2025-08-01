package api

import (
	"common/api"
	"common/common_vo"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_global"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"net/http"
	"strconv"
	"system/dao"
	"system/service"
	"time"
)

type ProfileApi struct {
	api.BaseApi
}

// 用户资料页面
func (w *ProfileApi) GetProfile(c *gin.Context) {
	var userService = service.GetUserServiceInstance()
	user := w.GetCurrUser(c)
	roles, err := userService.GetRoles(user.UserId)
	roleNames := ""
	if err == nil {
		for _, role := range roles {
			roleNames += role.RoleName + ","
		}
	}
	posts, err := dao.GetSysPostDaoInstance().FindPostsByUserId(user.UserId)
	postGroup := ""
	if err == nil {
		for _, post := range *posts {
			postGroup += post.PostName + ","
		}
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"msg":       "success",
		"code":      http.StatusOK,
		"data":      user,
		"roleGroup": roleNames,
		"postGroup": postGroup,
	})
}

// 修改用户信息
func (w *ProfileApi) Update(c *gin.Context) {
	var req common_vo.ProfileReq

	if err := c.ShouldBindJSON(&req); err != nil {
		lv_log.Error(err)
		util.Fail(c, "参数错误！")
		return
	}
	var userService service.UserService
	err := userService.UpdateProfile(&req, c)
	if err != nil {
		util.Fail(c, err.Error())
	} else {
		util.Success(c, "")
	}
}

// 修改用户密码
func (w *ProfileApi) UpdatePassword(c *gin.Context) {
	var req common_vo.PasswordReq
	err := c.ShouldBind(&req)
	lv_err.HasErrAndPanic(err)
	var userService service.UserService
	err = userService.UpdatePassword(&req, c)
	lv_err.HasErrAndPanic(err)
	util.Success(c, nil)
}

// 保存头像
func (w *ProfileApi) UpdateAvatar(c *gin.Context) {
	var userService service.UserService
	user := w.GetCurrUser(c)
	saveDir := lv_global.Config().GetUploadPath() + "/profile"
	fileHead, err := c.FormFile("avatarfile")
	if err != nil {
		util.Fail(c, err.Error())
	}
	filename := user.UserName + strconv.FormatInt(time.Now().UnixNano(), 10) + ".png"
	dts := saveDir + "/" + filename
	if err := c.SaveUploadedFile(fileHead, dts); err != nil {
		util.Fail(c, err.Error())
	}
	httpUrl := "static/upload/profile/" + filename //http访问路径，静态根目录是static ->访问： /upload/profile/filename
	err = userService.UpdateAvatar(httpUrl, c)
	if err != nil {
		util.Fail(c, err.Error())
	} else {
		util.Success(c, httpUrl)
	}
}
