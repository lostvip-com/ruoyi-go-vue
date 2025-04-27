package service

import (
	"github.com/lostvip-com/lv_framework/lv_db/lv_dao"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/spf13/cast"
	"html/template"
	"strings"
	dao2 "system/dao"
	"system/model"
)

type PermissionService struct{}

var permissionService *PermissionService

func GetPermissionServiceInstance() *PermissionService {
	if permissionService == nil {
		permissionService = &PermissionService{}
	}
	return permissionService
}

// PermButton 根据用户id和权限字符串判断是否输出控制按钮
func (svc *PermissionService) PermButton(u interface{}, permission, funcName, text, aclassName, iclassName string) template.HTML {
	result := svc.HasPermi(u, permission)

	htmlstr := ""
	if result == "" {
		htmlstr = `<a class="` + aclassName + `" onclick="` + funcName + `" hasPermission="` + permission + `">
                    <i class="` + iclassName + `"></i> ` + text + `
                </a>`
	}

	return template.HTML(htmlstr)
}

// HasPerm 根据用户id和权限字符串判断是否有此权限
func (svc *PermissionService) HasPermi(u interface{}, permission string) string {
	if u == nil {
		return "disabled"
	}

	uid := cast.ToInt64(u)
	var userService UserService
	if uid <= 0 {
		return "disabled"
	}
	//获取权限信息
	var menus []model.SysMenu
	var dao dao2.MenuDao
	if userService.IsAdmin(uid) {
		menus, _ = dao.SelectMenuNormalAll("")
	} else {
		menus, _ = dao.SelectMenusByUserId(uid, "")
	}

	if menus != nil && len(menus) > 0 {
		for i := range menus {
			if strings.EqualFold((menus)[i].Perms, permission) {
				return ""
			}
		}
	}

	return "disabled"
}

func (svc *PermissionService) FindPerms(roles []string) []string {
	sql := `
	 SELECT distinct m.perms from sys_menu m 
	 left join sys_role_menu rm on m.menu_id = rm.menu_id
     left join sys_role r on rm.role_id = r.role_id		
	 where r.role_key in (@roles)
  `
	params := map[string]any{"roles": roles}
	listPerms, err := lv_dao.ListMapByNamedSql(sql, params, false)
	if err != nil {
		lv_log.Error(err)
		return []string{}
	}
	arr := make([]string, len(*listPerms))
	for i, mp := range *listPerms {
		arr[i] = cast.ToString(mp["perms"])
	}
	return arr
}
