package service

import (
	"errors"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_dao"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"system/dao"
	"system/model"
	"system/vo"
	"time"
)

type MenuService struct {
	BaseService
}

var menuService *MenuService

func GetMenuServiceInstance() *MenuService {
	if menuService == nil {
		menuService = &MenuService{}
	}
	return menuService
}

// 根据主键查询数据
func (svc *MenuService) FindById(id int64) (*model.SysMenu, error) {
	return svc.FindById(id)
}

// 根据条件查询数据
func (svc *MenuService) FindAll(params *vo.SelectMenuPageReq) ([]model.SysMenu, error) {
	return svc.FindAll(params)
}

// 根据条件分页查询数据
func (svc *MenuService) SelectListPage(params *vo.SelectMenuPageReq) (*[]model.SysMenu, *lv_dto.Paging, error) {
	return svc.SelectListPage(params)
}

// 根据主键删除数据
func (svc *MenuService) DeleteById(id int64) error {
	err := (&model.SysMenu{MenuId: id}).Delete()
	if err == nil {
		lv_db.GetMasterGorm().Exec("delete from sys_menu where parent_id=?", id)
	}
	return err
}

// 添加数据
func (svc *MenuService) AddSave(req *model.SysMenu) (int64, error) {
	req.CreateTime = time.Now()
	err := req.Save()
	return req.MenuId, err
}

// 修改数据
func (svc *MenuService) Edit(req *model.SysMenu) error {
	entity := &model.SysMenu{MenuId: req.MenuId}
	err := entity.FindOne()
	if err != nil {
		return err
	}
	entity.MenuName = req.MenuName
	entity.Visible = req.Visible
	entity.ParentId = req.ParentId
	entity.Remark = req.Remark
	entity.MenuType = req.MenuType
	entity.Component = req.Component
	entity.Perms = req.Perms
	entity.IsFrame = req.IsFrame
	entity.Visible = req.Visible
	entity.IsCache = req.IsCache
	entity.OrderNum = req.OrderNum
	entity.UpdateTime = time.Now()
	err = entity.Update()
	return err
}

// 批量删除数据记录
func (svc *MenuService) DeleteByIds(ids string) int64 {
	idarr := lv_conv.ToInt64Array(ids, ",")
	var dao dao.MenuDao
	result, err := dao.DeleteBatch(idarr...)
	if err != nil {
		return 0
	}
	return result
}

// SelectMenuTree 加载所有菜单列表树
func (svc *MenuService) SelectMenuTree(userId int64, menu *model.SysMenu) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	sql := "select distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible, m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time "
	sql += "from sys_menu m "
	if svc.IsAdmin(userId) {
		sql += "where 1=1 "
	} else {
		sql += "left join sys_role_menu rm on m.menu_id = rm.menu_id "
		sql += "left join sys_user_role ur on rm.role_id = ur.role_id "
		sql += "left join sys_role ro on ur.role_id = ro.role_id "
		sql += "where ur.user_id = " + cast.ToString(userId)
	}
	if menu != nil {
		var menuName = menu.MenuName
		if menuName != "" {
			sql += "AND m.menu_name like concat('%', " + menuName + ", '%')"
		}
		var visible = menu.Visible
		if visible != "" {
			sql += "AND m.visible = " + visible
		}
		var status = menu.Status
		if status != "" {
			sql += "AND m.status = " + status
		}
	}

	err := lv_db.GetMasterGorm().Raw(sql).Scan(&menus).Error
	return menus, err
}

func fillChildrenTree(parent *vo.RouterVO, pcFlatMap map[int64][]vo.RouterVO) {
	children := pcFlatMap[parent.MenuId] // 获取本节点的子节点
	parent.Children = children
	//lv_log.Info("当前节点:", parent.MenuId, "children: ", len(children)) // 调试日志
	if parent.Children == nil || len(parent.Children) == 0 {
		return
	}
	for i := 0; i < len(parent.Children); i++ {
		child := &parent.Children[i] // 明确引用子节点
		//lv_log.Info("子节点:", i, " menuId:", child.MenuId, "child:", child) // 打印子节点详细信息
		fillChildrenTree(child, pcFlatMap) // 递归处理子节点
	}
}

func (svc *MenuService) GenMenus(parent *model.SysMenu, allMenus []model.SysMenu) {
	if parent.MenuType == "F" {
		return
	}
	if parent.Children == nil {
		parent.Children = make([]model.SysMenu, 0)
	}
	for i := range allMenus {
		if allMenus[i].ParentId == parent.MenuId { //发现子菜单
			parent.Children = append(parent.Children, allMenus[i])
			svc.GenMenus(&allMenus[i], allMenus) //接着找下一级的子菜单
		}
	}
	if len(parent.Children) == 0 {
		return
	}
}
func (svc *MenuService) InitParentChildMap(menus []model.SysMenu) map[int64][]vo.RouterVO {
	childrenMap := make(map[int64][]vo.RouterVO)
	num := len(menus)
	for i := 0; i < num; i++ {
		menu := menus[i]
		router := svc.menu2RouteVo(&menu) //格式变化
		pid := menu.ParentId

		if childrenMap[pid] == nil {
			list := make([]vo.RouterVO, 0)
			childrenMap[pid] = list
		}
		childrenMap[pid] = append(childrenMap[pid], router) //组织父子关系
	}
	return childrenMap
}

func (svc *MenuService) menu2RouteVo(menu *model.SysMenu) vo.RouterVO {
	meta := vo.Meta{Title: menu.MenuName, Icon: menu.Icon, NoCache: menu.IsCache == "1"}
	if menu.IsFrame == "0" {
		meta.Link = menu.Path
	}
	router := vo.RouterVO{
		MenuId:    menu.MenuId,
		ParentId:  menu.ParentId,
		Name:      getRouteName(menu),
		Path:      getRouterPath(menu),
		Component: getComponent(menu),
		Hidden:    menu.Visible == "1",
		Meta:      meta}
	if "M" == menu.MenuType {
		if !isInnerLink(menu.Path) {
			router.AlwaysShow = true
			router.Redirect = "noRedirect"
		}
	}
	return router
}

func (svc *MenuService) IsRolePermited(roles []string, perm string) (bool, interface{}) {
	sql := "SELECT count(*) from sys_menu m,sys_role_menu rm,sys_role r where m.menu_id=rm.menu_id and rm.role_id = r.role_id and r.role_key in @roles and m.perms=@perm"
	count, err := lv_dao.CountByNamedSql(sql, map[string]interface{}{"roles": roles, "perm": perm})
	return count > 0, err
}

func (svc *MenuService) BuildMenuTreeSelect(lists []model.SysMenu) []vo.MenuTreeSelect {
	var menuTreeSelect []vo.MenuTreeSelect
	for i := 0; i < len(lists); i++ {
		var menu = lists[i]
		menuId := menu.MenuId
		parentId := menu.ParentId
		if 0 == parentId {
			var menuVo = vo.MenuTreeSelect{
				Id:    menuId,
				Label: menu.MenuName,
			}
			menuVo.Children = svc.BuildChildMenusTreeSelect(menuId, lists)
			menuTreeSelect = append(menuTreeSelect, menuVo)
		}
	}
	return menuTreeSelect
}

func (svc *MenuService) BuildChildMenusTreeSelect(parentId int64, menus []model.SysMenu) []vo.MenuTreeSelect {
	var List []vo.MenuTreeSelect
	for i := 0; i < len(menus); i++ {
		var menu = menus[i]
		var menuId = menu.MenuId
		var pId = menu.ParentId
		if pId == parentId {
			var menuVo = vo.MenuTreeSelect{
				Id:    menuId,
				Label: menu.MenuName,
			}
			menuVo.Children = svc.BuildChildMenusTreeSelect(menuId, menus)
			List = append(List, menuVo)
		}
	}
	return List
}

func (svc *MenuService) SelectMenuListByRoleId(roleId string, menuCheckStrictly bool) ([]int, error) {
	var sql = "select m.menu_id from sys_menu m " +
		"left join sys_role_menu rm on m.menu_id = rm.menu_id " +
		"where rm.role_id = " + roleId + " "
	if menuCheckStrictly {
		sql += "and m.menu_id not in (select m.parent_id from sys_menu m inner join sys_role_menu rm on m.menu_id = rm.menu_id and rm.role_id = " + roleId + ")"
	}
	sql += "order by m.parent_id, m.order_num"
	var menuIds []int
	err := lv_db.GetMasterGorm().Raw(sql).Scan(&menuIds).Error
	return menuIds, err
}

// FindRouterTreeAll 获取所有菜单（管理员不授权直接加载所有）
func (svc *MenuService) FindRouterTreeAll() ([]model.SysMenu, error) {
	var menus []model.SysMenu
	sql := ` select distinct m.* from sys_menu m 
             where m.menu_type in ('M', 'C') and m.status = 0 
             order by m.parent_id, m.order_num
           `
	err := lv_db.GetMasterGorm().Raw(sql).Scan(&menus).Error
	return menus, err
}

// FindRouterTreeAllByUserId 获取所有菜单（非管理员获取已经授权的菜单）
func (svc *MenuService) FindRouterTreeAllByUserId(id int64) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	sql := ` select distinct m.* from sys_menu m 
             left join sys_role_menu rm on m.menu_id = rm.menu_id 
             left join sys_user_role ur on rm.role_id = ur.role_id 
             left join sys_role ro on ur.role_id = ro.role_id 
             left join sys_user u on ur.user_id = u.user_id 
             where u.user_id = ? and m.menu_type in ('M', 'C') and m.status = 0  AND ro.status = 0
             order by m.parent_id, m.order_num
           `
	err := lv_db.GetMasterGorm().Raw(sql, id).Scan(&menus).Error
	return menus, err
}

func (svc *MenuService) BuildMenus(lists []model.SysMenu) []vo.MenuVo {
	var menuVos []vo.MenuVo
	for i := 0; i < len(lists); i++ {
		var menu = &lists[i]
		MenuId := menu.MenuId
		parentId := menu.ParentId
		if 0 == parentId {
			var path = ""
			if isInnerLink(menu.Path) {
				path = menu.Path
			}
			var menuVo = vo.MenuVo{
				Hidden: "1" == menu.Visible,
				Query:  menu.Query,
				MetaVo: vo.MetaVo{
					Title:   menu.MenuName,
					Icon:    menu.Icon,
					NoCache: "1" == menu.IsCache,
					Link:    path,
				},
				Name:      getRouteName(menu),
				Path:      getRouterPath(menu),
				Component: getComponent(menu),
			}
			if "M" == menu.MenuType {
				if !isInnerLink(menu.Path) {
					menuVo.AlwaysShow = true
					menuVo.Redirect = "noRedirect"
					menuVo.Children = svc.BuildChildMenus(MenuId, lists)
				}
			}
			menuVos = append(menuVos, menuVo)
		}

	}
	return menuVos
}

func (svc *MenuService) BuildChildMenus(ParentId int64, lists []model.SysMenu) []vo.MenuVo {
	var List []vo.MenuVo
	for i := 0; i < len(lists); i++ {
		var menu = &lists[i]
		var menuId = menu.MenuId
		var pId = menu.ParentId
		if pId == ParentId {
			var path = ""
			if isInnerLink(menu.Path) {
				path = menu.Path
			}
			var menuVo = vo.MenuVo{
				Hidden: "1" == menu.Visible,
				Query:  menu.Query,
				MetaVo: vo.MetaVo{
					Title:   menu.MenuName,
					Icon:    menu.Icon,
					NoCache: "1" == menu.IsCache,
					Link:    path,
				},
				Name:      getRouteName(menu),
				Path:      getRouterPath(menu),
				Component: getComponent(menu),
			}
			if "M" == menu.MenuType {
				if !isInnerLink(menu.Path) {
					menuVo.AlwaysShow = true
					menuVo.Redirect = "noRedirect"
					menuVo.Children = svc.BuildChildMenus(menuId, lists)
				}
			}
			List = append(List, menuVo)
		}
	}
	return List
}

func getRouteName(menu *model.SysMenu) string {
	var name = FirstUpper(menu.Path)
	if isMenuFrame(menu) {
		return ""
	}
	return name
}

func getRouterPath(menu *model.SysMenu) string {
	var routerPath = menu.Path
	if isInnerLink(routerPath) {
		return routerPath
	}
	// 非外链并且是一级目录（类型为目录）
	if 0 == menu.ParentId && "M" == menu.MenuType && "1" == menu.IsFrame {
		routerPath = "/" + menu.Path
	} else if isMenuFrame(menu) {
		routerPath = "/"
	}
	return routerPath
}

func getComponent(menu *model.SysMenu) string {
	var component = "Layout"
	if "" != menu.Component && !isMenuFrame(menu) {
		component = menu.Component
	} else if "" == menu.Component && isInnerLink(menu.Path) {
		component = "InnerLink"
	} else if "" == menu.Component && isParentView(menu) {
		component = "ParentView"
	}
	return component
}

func isParentView(menu *model.SysMenu) bool {
	return menu.ParentId != 0 && "M" == menu.MenuType
}

// 是否为外链
func isInnerLink(path string) bool {
	return strings.Contains(path, "http://") || strings.Contains(path, "https://")
}

func isMenuFrame(menu *model.SysMenu) bool {
	return menu.ParentId == 0 && "C" == menu.MenuType && menu.IsFrame == "1"
}

/*首字母大写*/
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func SelectMenuListByRoleId(roleId string, menuCheckStrictly bool) []int {
	var sql = "select m.menu_id from sys_menu m " +
		"left join sys_role_menu rm on m.menu_id = rm.menu_id " +
		"where rm.role_id = " + roleId + " "
	if menuCheckStrictly {
		sql += "and m.menu_id not in (select m.parent_id from sys_menu m inner join sys_role_menu rm on m.menu_id = rm.menu_id and rm.role_id = " + roleId + ")"
	}
	sql += "order by m.parent_id, m.order_num"
	var menuIds []int
	err := lv_db.GetMasterGorm().Raw(sql).Scan(&menuIds).Error
	if err != nil {
		panic(errors.New(err.Error()))
	}
	return menuIds
}
func FindMenuInfoById(menuId string) model.SysMenu {
	sql := "select menu_id, menu_name, parent_id, order_num, path, " +
		"component, `query`, is_frame, is_cache, menu_type, visible, status, ifnull(perms,'') as perms, icon, create_time " +
		"from sys_menu "
	sql = sql + "where menu_id = " + menuId
	var list model.SysMenu
	err := lv_db.GetMasterGorm().Raw(sql).First(&list).Error
	if err != nil {
		panic(errors.New(err.Error()))
	}
	return list
}

func hasChildByMenuId(menuId string) model.SysMenu {
	var menu model.SysMenu
	lv_db.GetMasterGorm().Where("parent_id = ? ", menuId).First(&menu)
	return menu
}

func hasChildCountByMenuId(menuId string) int {
	var menuCount int
	err := lv_db.GetMasterGorm().Where("parent_id = ? ", menuId).Scan(&menuCount).Error
	if err != nil {
		panic(errors.New(err.Error()))
	}
	return menuCount
}

func checkMenuExistRole(menuId string) int {
	var menuCount int
	err := lv_db.GetMasterGorm().Raw("select count(1) from sys_role_menu where menu_id = " + menuId).Scan(&menuCount).Error
	if err != nil {
		panic(errors.New(err.Error()))
	}
	return menuCount
}

func DeleteMenu(menuIds string) error {
	menu := hasChildByMenuId(menuIds)
	if menu.MenuId != 0 {
		return errors.New("存在子菜单,不允许删除")
	}

	menuCount := checkMenuExistRole(menuIds)
	if menuCount > 0 {
		return errors.New("菜单已分配,不允许删除")
	}

	err := lv_db.GetMasterGorm().Exec("delete from sys_menu where menu_id in (?) ", menuIds).Error
	if err == nil {
		errors.New("删除部门关联用户失败")
	}
	return err
}

func checkMenuNameUnique(parentId int, menuName string) int {
	var menuCount int
	err := lv_db.GetMasterGorm().Raw("select count(1) from sys_menu where menu_name = "+menuName+" and parent_id = "+strconv.Itoa(parentId), &menuCount).Error
	if err != nil {
		panic(errors.New(err.Error()))
	}
	return menuCount
}
