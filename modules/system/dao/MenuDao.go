package dao

import (
	"errors"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	"system/model"
	"system/vo"
)

var menuDao *MenuDao

// 修改页面请求参数
type MenuDao struct {
}

func GetMenuDaoInstance() *MenuDao {
	if menuDao == nil {
		menuDao = &MenuDao{}
	}
	return menuDao
}

// 批量删除
func (r *MenuDao) DeleteBatch(ids ...int) (int, error) {
	db := lv_db.GetOrmDefault().Table("sys_menu").Where("menu_id in ? ", ids).Update("del_flag", 1)
	return int(db.RowsAffected), db.Error
}

func (r *MenuDao) DeleteChildren(parentId int) (int, error) {
	tb := lv_db.GetOrmDefault().Table("sys_menu").Where("parent_id=?", parentId).Update("del_flag", 1)
	return int(tb.RowsAffected), tb.Error
}

// 根据主键查询数据
func (dao *MenuDao) FindById(id int) (*model.SysMenu, error) {
	tb := lv_db.GetOrmDefault()
	if tb == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	var result model.SysMenu
	tb = tb.Table("sys_menu t")
	tb.Select("t.menu_id, t.parent_id, t.menu_name, t.order_num, t.url, t.target, t.menu_type, t.visible, t.perms, t.icon, t.remark,(SELECT menu_name FROM sys_menu WHERE menu_id = t.parent_id) parent_name")
	tb.Where("t.menu_id=?", id)
	err := tb.First(&result).Error
	if err != nil {
		return nil, errors.New("获取数据失败")
	}
	return &result, nil
}

// 根据条件分页查询数据
func (dao *MenuDao) SelectListPage(param *vo.SelectMenuPageReq) (*[]model.SysMenu, int, error) {
	tb := lv_db.GetOrmDefault()
	tb = tb.Table("sys_menu")
	if param != nil {
		if param.MenuName != "" {
			tb.Where("menu_name like ?", "%"+param.MenuName+"%")
		}

		if param.Visible != "" {
			tb.Where("visible = ", param.Visible)
		}

		if param.BeginTime != "" {
			tb.Where("date_format(create_time,'%y%m%d') >= date_format(?,'%y%m%d') ", param.BeginTime)
		}

		if param.EndTime != "" {
			tb.Where("date_format(create_time,'%y%m%d') <= date_format(?,'%y%m%d') ", param.EndTime)
		}
	}

	var result []model.SysMenu
	err := tb.Offset(param.GetStartNum()).Limit(param.GetPageSize()).Find(&result).Error
	lv_err.HasErrAndPanic(err)
	err = tb.Offset(0).Limit(-1).Count(&param.Total).Error
	lv_err.HasErrAndPanic(err)
	return &result, int(param.Total), nil
}

// 获取所有数据
func (dao *MenuDao) FindAll(sysMenu *vo.SelectMenuPageReq) ([]model.SysMenu, error) {
	tb := lv_db.GetOrmDefault()
	var rows []model.SysMenu
	var sql = `select menu_id, menu_name, parent_id, order_num, path, component, query, is_frame, is_cache, menu_type, 
                      visible,status, ifnull(perms,'') as perms, icon, create_time 
               from sys_menu where 1 = 1 `
	var name = sysMenu.MenuName
	if name != "" {
		sql += " AND menu_name LIKE CONCAT('%', '" + name + "', '%')"
	}
	var visible = sysMenu.Visible
	if visible != "" {
		sql += " AND visible = " + visible
	}
	var status = sysMenu.Status
	if status != "" {
		sql += " AND status = " + status
	}
	err := tb.Raw(sql).Find(&rows).Error
	return rows, err
}

// FindMenuNormalAll 获取管理员菜单数据
func (dao *MenuDao) FindMenuNormalAll(noF bool) ([]model.SysMenu, error) {
	var result []model.SysMenu

	tb := lv_db.GetOrmDefault()
	tb = tb.Table("sys_menu as m")
	tb.Where(" m.visible = 0")
	if noF == true {
		tb.Where(" m.menu_type!='F' ")
	}
	tb.Order("m.parent_id, m.order_num")
	err := tb.Find(&result).Error

	if err != nil {
		return nil, err
	} else {
		return result, err
	}
}
func (dao *MenuDao) FindMenusByUserId(userId int, sysMenu *vo.SelectMenuPageReq) ([]model.SysMenu, error) {
	var sql = "select distinct m.menu_id, m.parent_id, m.menu_name, m.path, m.component, m.`query`, m.visible," +
		" m.status, ifnull(m.perms,'') as perms, m.is_frame, m.is_cache, m.menu_type, m.icon, m.order_num, m.create_time " +
		"from sys_menu m left join sys_role_menu rm on m.menu_id = rm.menu_id " +
		"left join sys_user_role ur on rm.role_id = ur.role_id " +
		"left join sys_role ro on ur.role_id = ro.role_id "
	sql += "where ur.user_id = " + cast.ToString(userId) + " "
	var menuName = sysMenu.MenuName
	if menuName != "" {
		sql += "AND m.menu_name like concat(%" + menuName + "%) "
	}
	if sysMenu.Visible != "" {
		sql += "AND m.visible = " + sysMenu.Visible + " "
	}
	var status = sysMenu.Status
	if status != "" {
		sql += "AND m.status = " + status + " "
	}
	sql += "order by m.parent_id, m.order_num"
	var list []model.SysMenu
	err := lv_db.GetOrmDefault().Raw(sql).Scan(&list).Error
	return list, err
}

// FindMenus 根据用户ID读取菜单数据， noF 非按钮
func (dao *MenuDao) FindMenus(userId int, noF bool, params *model.SysMenu) ([]model.SysMenu, error) {
	var result []model.SysMenu
	db := lv_db.GetOrmDefault()
	tb := db.Table("sys_menu as m")
	tb.Joins("LEFT join sys_role_menu as rm on m.menu_id = rm.menu_id")
	tb.Joins("LEFT join sys_user_role as ur on rm.role_id = ur.role_id")
	tb.Joins("LEFT join sys_role as ro on ur.role_id = ro.role_id")
	tb.Select("m.*")
	tb.Where("ur.user_id = ? and  m.visible = 0  AND ro.status = 0", userId)
	if params != nil {
		tb.Where("m.menu_id = ?", params.MenuId)
		if params.MenuName != "" {
			tb.Where("AND m.menu_name like concat('%', " + params.MenuName + ", '%')")
		}
		var visible = params.Visible
		if visible != "" {
			tb.Where("AND m.visible =? ", visible)
		}
		var status = params.Status
		if status != "" {
			tb.Where("AND m.status = ", status)
		}
	}
	if noF == true {
		tb.Where(" m.menu_type!='F' ")
	}
	tb.Order("m.parent_id, m.order_num")
	err := tb.Find(&result).Error

	if err != nil {
		return nil, err
	} else {
		return result, err
	}
}

// 根据角色ID查询菜单
func (dao *MenuDao) SelectMenuTree(roleId int) ([]string, error) {
	tb := lv_db.GetOrmDefault()
	if tb == nil {
		return nil, errors.New("获取数据库连接失败")
	}

	var result []string

	if tb == nil {
		return nil, errors.New("获取数据库连接失败")
	}
	tb = tb.Table("sys_menu as m ")
	tb.Joins(" LEFT join sys_role_menu as rm on m.menu_id = rm.menu_id")
	tb.Where(" rm.role_id = ?", roleId)
	tb.Order(" m.parent_id, m.order_num desc ")
	tb.Select("concat(m.menu_id, ifnull(m.perms,'')) as perms")
	var list []model.SysMenu
	err := tb.Find(&list).Error
	if err != nil {
		return nil, err
	}
	for _, record := range list {
		if record.Perms != "" {
			result = append(result, record.Perms)
		}
	}

	return result, err
}
