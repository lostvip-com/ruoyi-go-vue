package dao

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_dao"
)

// Fill with you ideas below.

// GenTable is the golang structure for table sys_dept.
type SysDeptDao struct {
}

// 根据ID查询所有子部门
func (dao SysDeptDao) SelectChildrenDeptById(deptId int64) []*models.SysDept {
	db := lv_db.GetMasterGorm()
	var rs []*models.SysDept
	db.Table("sys_dept").Where("parent_id=?", deptId).Find(&rs)
	return rs
}

// 删除部门管理信息
func (dao SysDeptDao) DeleteDeptById(deptId int64) error {
	var entity models.SysDept
	err := entity.UpdateDelFlag(deptId)
	return err
}

// 查询部门管理数据
func (d SysDeptDao) SelectDeptList(parentId int64, deptName, status string, tenantId int64) (*[]models.SysDept, error) {
	sql := ` select *  from sys_dept d  where d.del_flag =0  `
	param := map[string]any{}
	if parentId > 0 {
		param["parentId"] = parentId
		sql += " and d.parent_id = @parentId "
	}
	if deptName != "" {
		param["deptName"] = "%" + deptName + "%"
		sql += " and d.dept_name like @deptName "
	}
	if status != "" {
		param["status"] = status
		sql += " and d.status = @status "
	}
	if tenantId != 0 {
		param["tenantId"] = tenantId
		sql += " and d.tenant_id = @tenantId "
	}
	sql += " order by d.parent_id, d.order_num desc "

	return lv_dao.ListByNamedSql[models.SysDept](sql, param)
}

// 根据角色ID查询部门
func (dao SysDeptDao) SelectRoleDeptTree(roleId int64) ([]string, error) {
	sql := ` select concat(d.dept_id, d.dept_name) as DeptName 
             from sys_dept d 
             left join sys_role_dept rd  on d.dept_id = rd.dept_id 
             where d.del_flag =0 and rd.role_id = @roleId
             order by d.parent_id, d.order_num
             `
	param := map[string]any{}
	param["roleId"] = roleId
	listMap, err := lv_dao.ListMapStrByNamedSql(sql, param, false)
	var result []string
	var rs = *listMap
	if err == nil && rs != nil && len(rs) > 0 {
		for _, record := range rs {
			if record["DeptName"] != "" {
				result = append(result, record["DeptName"])
			}
		}
	}
	return result, nil
}
