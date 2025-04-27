package service

import (
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_dao"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/web/lv_dto"
	"github.com/spf13/cast"
	"strings"
	"system/dao"
	"system/model"
	"system/vo"
	"time"
)

type MenuService struct {
}

var menuService *MenuService

func GetMenuServiceInstance() *MenuService {
	if menuService == nil {
		menuService = &MenuService{}
	}
	return menuService
}

// 根据主键查询数据
func (svc *MenuService) SelectRecordById(id int64) (*model.SysMenu, error) {
	return svc.SelectRecordById(id)
}

// 根据条件查询数据
func (svc *MenuService) SelectListAll(params *vo.SelectMenuPageReq) ([]model.SysMenu, error) {
	return svc.SelectListAll(params)
}

// 根据条件分页查询数据
func (svc *MenuService) SelectListPage(params *vo.SelectMenuPageReq) (*[]model.SysMenu, *lv_dto.Paging, error) {
	return svc.SelectListPage(params)
}

// 根据主键删除数据
func (svc *MenuService) DeleteRecordById(id int64) error {
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
func (svc *MenuService) Edit(req *model.SysMenu) (int64, error) {
	entity := &model.SysMenu{MenuId: req.MenuId}
	err := entity.FindOne()
	if err != nil {
		return 0, err
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
	return entity.MenuId, err
}

// 批量删除数据记录
func (svc *MenuService) DeleteRecordByIds(ids string) int64 {
	idarr := lv_conv.ToInt64Array(ids, ",")
	var dao dao.MenuDao
	result, err := dao.DeleteBatch(idarr...)
	if err != nil {
		return 0
	}
	return result
}

// MenuTreeData 加载所有菜单列表树
func (svc *MenuService) MenuTreeData(userId int64, menuType string) (*[]lv_dto.Ztree, error) {
	var result *[]lv_dto.Ztree
	menuList, err := svc.SelectMenuNormalByUser(userId, menuType)
	if err != nil {
		return nil, err
	}
	result, err = svc.InitZtree(menuList, nil, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 获取用户的菜单数据
func (svc *MenuService) SelectMenuNormalByUser(userId int64, menuType string) (*[]model.SysMenu, error) {
	var userService UserService
	//从数据库中读取
	var dao dao.MenuDao
	if userService.IsAdmin(userId) {
		menus, err := dao.SelectMenuNormalAll(menuType)
		return &menus, err
	} else {
		menus, err := dao.SelectMenusByUserId(userId, menuType)
		return &menus, err
	}
}

// 获取用户的菜单数据
func (svc *MenuService) ListMenuNormalByUser(userId int64, menuType string) (*[]model.SysMenu, error) {
	var dao dao.MenuDao
	var userService UserService
	if userService.IsAdmin(userId) {
		menus, err := dao.SelectMenuNormalAll(menuType)
		return &menus, err
	} else {
		menus, err := dao.SelectMenusByUserId(userId, menuType)
		return &menus, err
	}
}

// SelectMenuNormalAll 获取管理员菜单数据,不区分资源类型传空即可
func (svc *MenuService) SelectMenuNormalAll(userId int64, menuType string) ([]vo.RouterVO, error) {
	var menus []model.SysMenu
	menuDao := dao.GetMenuDaoInstance()
	var err error
	if userId == 0 {
		menus, err = menuDao.SelectMenuNormalAll(menuType)
	} else {
		menus, err = menuDao.SelectMenusByUserId(userId, menuType)
	}
	lv_err.HasErrAndPanic(err)
	arr := make([]vo.RouterVO, 0)
	pcMap := svc.InitParentChildMap(menus)
	list0 := pcMap[0]
	for i := 1; i < len(list0); i++ {
		it := &list0[i]
		lv_log.Infof("%v --- %v", i, it)
		fillChildrenTree(it, pcMap)
	}
	//存入缓存
	return arr, nil
}

func fillChildrenTree(parent *vo.RouterVO, pcFlatMap map[int64][]vo.RouterVO) {
	children := pcFlatMap[parent.MenuId] // 获取本节点的子节点
	parent.Children = children
	lv_log.Infof("当前节点 %d 的子节点数：%d", parent.MenuId, len(children)) // 调试日志
	if parent.Children == nil || len(parent.Children) == 0 {
		return
	}
	for i := 0; i < len(parent.Children); i++ {
		child := &parent.Children[i]                                       // 明确引用子节点
		lv_log.Infof("子节点 %d 的 ID：%d，内容：%v", i, child.MenuId, child) // 打印子节点详细信息
		fillChildrenTree(child, pcFlatMap)                                 // 递归处理子节点
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
	len := len(menus)
	for i := 0; i < len; i++ {
		if menus[i].MenuType == "F" { //忽略按钮
			continue
		}
		childrenMap[menus[i].MenuId] = make([]vo.RouterVO, 0) //每个menu都预设子菜单项
	}

	for i := 0; i < len; i++ {
		menu := menus[i]
		if menu.MenuType == "F" { //忽略按钮
			continue
		}
		router := svc.menu2RouteVo(&menu) //格式变化
		pid := menu.ParentId
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
		Name:      menu.RouteName,
		Path:      menu.Path,
		Hidden:    menu.Visible == "1",
		Component: menu.Component,
		Meta:      meta}
	return router
}

// 根据角色ID查询菜单
func (svc *MenuService) RoleMenuTreeData(roleId, userId int64, menuType string) (*[]lv_dto.Ztree, error) {
	var result *[]lv_dto.Ztree
	menuList, err := svc.ListMenuNormalByUser(userId, menuType)
	if err != nil {
		return nil, err
	}
	var dao dao.MenuDao
	if roleId > 0 {
		roleMenuList, err := dao.SelectMenuTree(roleId)
		if err != nil || roleMenuList == nil {
			result, err = svc.InitZtree(menuList, nil, true)
		} else {
			result, err = svc.InitZtree(menuList, &roleMenuList, true)
		}
	} else {
		result, err = svc.InitZtree(menuList, nil, true)
	}

	return result, nil
}

// 对象转菜单树
func (svc *MenuService) InitZtree(menuList *[]model.SysMenu, roleMenuList *[]string, permsFlag bool) (*[]lv_dto.Ztree, error) {
	var result []lv_dto.Ztree
	isCheck := false
	if roleMenuList != nil && len(*roleMenuList) > 0 {
		isCheck = true
	}

	for _, obj := range *menuList {
		var ztree lv_dto.Ztree
		ztree.Title = obj.MenuName
		ztree.Id = obj.MenuId
		ztree.Name = svc.transMenuName(obj.MenuName, permsFlag)
		ztree.Pid = obj.ParentId
		if isCheck {
			tmp := cast.ToString(obj.MenuId) + obj.Perms
			tmpcheck := false
			for j := range *roleMenuList {
				if strings.Compare((*roleMenuList)[j], tmp) == 0 {
					tmpcheck = true
					break
				}
			}
			ztree.Checked = tmpcheck
		}
		result = append(result, ztree)
	}

	return &result, nil
}

func (svc *MenuService) transMenuName(menuName string, permsFlag bool) string {
	if permsFlag {
		return "<font color=\"#888\">&nbsp;&nbsp;&nbsp;" + menuName + "</font>"
	} else {
		return menuName
	}
}

func (svc *MenuService) IsRolePermited(roles []string, perm string) (bool, interface{}) {
	sql := "SELECT count(*) from sys_menu m,sys_role_menu rm,sys_role r where m.menu_id=rm.menu_id and rm.role_id = r.role_id and r.role_key in @roles and m.perms=@perm"
	count, err := lv_dao.CountByNamedSql(sql, map[string]interface{}{"roles": roles, "perm": perm})
	return count > 0, err
}
