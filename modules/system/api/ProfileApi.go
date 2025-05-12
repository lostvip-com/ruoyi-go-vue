package api

import (
	"common/common_vo"
	"common/global"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"net/http"
	"os"
	"strconv"
	"system/model"
	"system/service"
	"time"
)

type ProfileApi struct {
}

// 用户资料页面
func (w *ProfileApi) Profile(c *gin.Context) {
	var userService service.UserService
	user := userService.GetProfile(c)
	util.Success(c, user)
}

// 修改用户信息
func (w *ProfileApi) Update(c *gin.Context) {
	var req common_vo.ProfileReq

	if err := c.ShouldBind(&req); err != nil {
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

// 检查登录名是否存在
func (w *ProfileApi) CheckUserNameOK(c *gin.Context) {
	var req common_vo.CheckUserNameReq
	if err := c.ShouldBind(&req); err != nil {
		c.Writer.WriteString("1")
		return
	}
	entity := model.SysUser{UserName: req.UserName}
	err := entity.FindOne()
	if err == nil { //查到了
		c.Writer.WriteString(global.JQ_BE_NO)
	} else {
		c.Writer.WriteString(global.JQ_BE_OK)
	}
}

// CheckPhoneOK 检查手机号是否存在 1 存在，0不存在
func (w *ProfileApi) CheckPhoneOK(c *gin.Context) {
	var req common_vo.CheckPhoneAllReq
	err := c.ShouldBind(&req)
	lv_err.HasErrAndPanic(err)
	var userService service.UserService
	count, err := userService.CountCol("phonenumber", req.Phonenumber)
	lv_err.HasError1(err)
	if count > 0 {
		c.Writer.WriteString(global.JQ_BE_NO)
	} else {
		c.Writer.WriteString(global.JQ_BE_OK)
	}
}

// 检查手机号是否存在 1 存在，0不存在
func (w *ProfileApi) CheckEmailOK(c *gin.Context) {
	email := c.Query("email")
	var userService service.UserService
	count, err := userService.CountCol("email", email)
	lv_err.HasError1(err)
	if count > 0 {
		c.Writer.WriteString(global.JQ_BE_NO)
	} else {
		c.Writer.WriteString(global.JQ_BE_OK)
	}
}

// 校验密码是否正确
func (w *ProfileApi) CheckPassword(c *gin.Context) {
	var req common_vo.CheckPasswordReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, lv_dto.CommonRes{
			Code: 500,
			Msg:  err.Error(),
		})
	}
	var userService service.UserService
	user := userService.GetProfile(c)

	result := userService.CheckPassword(user, req.Password)

	if result {
		c.Writer.WriteString(global.JQ_BE_NO)
	} else {
		c.Writer.WriteString(global.JQ_BE_OK)
	}
}

// 保存头像
func (w *ProfileApi) UpdateAvatar(c *gin.Context) {
	var userService service.UserService
	user := userService.GetProfile(c)
	curDir, err := os.Getwd()

	if err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).Log("保存头像", gin.H{"userid": user.UserId}).WriteJsonExit()
	}
	saveDir := curDir + "/static/upload/"
	fileHead, err := c.FormFile("avatarfile")

	if err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg("没有获取到上传文件").Log("保存头像", gin.H{"userid": user.UserId}).WriteJsonExit()
	}
	curdate := time.Now().UnixNano()
	filename := user.UserName + strconv.FormatInt(curdate, 10) + ".png"
	dts := saveDir + filename

	if err := c.SaveUploadedFile(fileHead, dts); err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).Log("保存头像", gin.H{"userid": user.UserId}).WriteJsonExit()
	}

	avatar := "/upload/" + filename

	err = userService.UpdateAvatar(avatar, c)

	if err != nil {
		util.ErrorResp(c).SetBtype(lv_dto.Buniss_Edit).SetMsg(err.Error()).Log("保存头像", gin.H{"userid": user.UserId}).WriteJsonExit()
	} else {
		util.SucessResp(c).SetBtype(lv_dto.Buniss_Edit).Log("保存头像", gin.H{"userid": user.UserId}).WriteJsonExit()
	}
}
