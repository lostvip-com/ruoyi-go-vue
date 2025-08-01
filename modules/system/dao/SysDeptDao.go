package dao

import (
	"common/models"
	"github.com/lostvip-com/lv_framework/lv_db/lv_dao"
)

// Fill with you ideas below.

// GenTable is the golang structure for table sys_dept.
type SysDeptDao struct {
}

// 查询部门管理数据
func (d SysDeptDao) SelectDeptList(parentId int, deptName, status string, tenantId int) (*[]models.SysDept, error) {
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
