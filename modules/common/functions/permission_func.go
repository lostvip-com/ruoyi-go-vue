package functions

import (
	"github.com/spf13/cast"
	"html/template"
	"strings"
	"system/dao"
	"system/model"
)

// 根据用户id和权限字符串判断是否输出控制按钮
func PermButton(u interface{}, permission, funcName, text, aclassName, iclassName string) template.HTML {

	result := HasPermi(u, permission)

	htmlstr := ""
	if result == "" {
		htmlstr = `<a class="` + aclassName + `" onclick="` + funcName + `" hasPermission="` + permission + `">
                    <i class="` + iclassName + `"></i> ` + text + `
                </a>`
	}

	return template.HTML(htmlstr)
}
func IsAdmin(userId int) bool {
	if userId == 1 {
		return true
	} else {
		return false
	}
}

// 根据用户id和权限字符串判断是否有此权限
func HasPermi(u interface{}, permission string) string {
	if u == nil {
		return "disabled"
	}
	uid := cast.ToInt(u)
	if uid <= 0 {
		return "disabled"
	}
	//获取权限信息
	var menuDao dao.MenuDao
	var menus []model.SysMenu
	if IsAdmin(uid) {
		menus, _ = menuDao.FindMenuNormalAll(false)
	} else {
		menus, _ = menuDao.FindMenusByUserId(uid, nil)
	}

	if menus != nil && len(menus) > 0 {
		for i := range menus {
			if strings.EqualFold(menus[i].Perms, permission) {
				return ""
			}
		}
	}
	return "disabled"
}

func UpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
func Substr(s string, start, length int) string {
	if len(s) == 0 {
		return s
	}
	if start < 0 {
		start = 0
	}
	if start > len(s) {
		start = len(s)
	}
	end := start + length
	if end > len(s) {
		end = len(s)
	}
	return s[start:end]
}
