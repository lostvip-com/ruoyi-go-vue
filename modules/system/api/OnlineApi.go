package api

import (
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"strings"
	"system/service"
	"system/vo"
)

type OnlineApi struct {
}

// ListAjax 列表分页数据
func (w *OnlineApi) ListAjax(c *gin.Context) {
	var param vo.OnlinePageReq

	err := c.ShouldBind(&param)
	lv_err.HasErrAndPanic(err)
	rows := make([]vo.OnlineVo, 0)
	db := lv_db.GetMasterGorm()
	tb := db.Table("sys_user_online t")
	if param.SessionId != "" {
		tb.Where("t.session_id = ?", param.SessionId)
	}

	if param.UserName != "" {
		tb.Where("t.user_name like ?", "%"+param.UserName+"%")
	}

	if param.DeptName != "" {
		tb.Where("t.dept_name like ?", "%"+param.DeptName+"%")
	}

	if param.Ipaddr != "" {
		tb.Where("t.ipaddr = ?", param.Ipaddr)
	}

	if param.LoginLocation != "" {
		tb.Where("t.login_location = ?", param.LoginLocation)
	}

	if param.Browser != "" {
		tb.Where("t.browser = ?", param.Browser)
	}

	if param.Os != "" {
		tb.Where("t.os = ?", param.Os)
	}

	if param.Status != "" {
		tb.Where("t.status = ?", param.Status)
	}

	if param.BeginTime != "" {
		tb.Where("t.create_time >= ? ", param.BeginTime)
	}

	if param.EndTime != "" {
		tb.Where("date_format(t.create_time,'%y%m%d') <= date_format(?,'%y%m%d') ", param.EndTime)
	}
	var total int64
	tb = tb.Count(&total)
	tb.Order("t.last_access_time desc").Offset((param.PageNum - 1) * param.PageSize).Limit(param.PageSize).Find(&rows)
	err = tb.Error
	lv_err.HasErrAndPanic(err)
	util.SuccessPage(c, rows, total)
}

// 用户强退
func (w *OnlineApi) ForceLogout(c *gin.Context) {
	sessionId := c.Query("sessionId")
	var userService service.SessionService
	err := userService.ForceLogout(sessionId)
	lv_err.HasErrAndPanic(err)
	util.Success(c, gin.H{"sessionId": sessionId})
}

// 批量强退
func (w *OnlineApi) BatchForceLogout(c *gin.Context) {
	var userService service.SessionService
	ids := c.Query("ids")
	if ids == "" {
		util.ErrorResp(c).SetMsg("参数错误").Log("批量强退", gin.H{"ids": ids}).WriteJsonExit()
		return
	}
	ids = strings.ReplaceAll(ids, "[", "")
	ids = strings.ReplaceAll(ids, "]", "")
	ids = strings.ReplaceAll(ids, `"`, "")
	idarr := strings.Split(ids, ",")
	if len(idarr) > 0 {
		for _, sessionId := range idarr {
			if sessionId != "" {
				userService.ForceLogout(sessionId)
			}
		}
	}
	util.SucessResp(c).Log("批量强退", gin.H{"ids": ids}).WriteJsonExit()
}
