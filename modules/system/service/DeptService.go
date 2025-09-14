package service

import (
	"common/common_vo"
	"common/models"
	"errors"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_dao"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_reflect"
	"github.com/spf13/cast"
	"system/dao"
	"system/vo"
)

type DeptService struct{}

var deptService *DeptService

func GetDeptServiceInstance() *DeptService {
	if userService == nil {
		deptService = &DeptService{}
	}
	return deptService
}

// AddSave 新增保存信息
func (svc *DeptService) AddSave(req *models.SysDept) (*models.SysDept, error) {
	parent := new(models.SysDept)
	parent, err := parent.FindById(req.ParentId)
	if err == nil {
		if parent.Status != "0" {
			return nil, errors.New("部门停用，不允许新增")
		}
	} else {
		return nil, errors.New("父部门不能为空")
	}
	err = req.Save()
	if err != nil {
		return nil, err
	}
	//这里跟原版不一样了，多加了一级自己的ID，以方便数据权限控制
	req.Ancestors = parent.Ancestors + "," + cast.ToString(req.DeptId)
	err = req.Update()
	lv_err.HasErrAndPanic(err)
	return req, err
}

// EditSave 修改保存信息
func (svc *DeptService) EditSave(req *models.SysDept) (*models.SysDept, error) {
	po := &models.SysDept{DeptId: req.DeptId}
	po, err := po.FindById(req.DeptId)
	lv_err.HasErrAndPanic(err)
	parent := &models.SysDept{}
	parent, err = parent.FindById(req.ParentId)
	lv_err.HasErrAndPanic(err)
	if parent.Status != "0" {
		return nil, errors.New("部门停用，不允许新增")
	} else {
		_ = lv_reflect.CopyProp(req, po, true)
		err := po.Update()
		lv_err.HasErrAndPanic(err)
		//递归修改所有子项目
		svc.UpdateChildrenAncestors(po, parent.Ancestors)
		return po, err
	}
}

// UpdateChildrenAncestors 修改子元素关系（替换前半部分）
func (svc *DeptService) UpdateChildrenAncestors(dept *models.SysDept, parentCodes string) {
	dept.Ancestors = parentCodes + "," + cast.ToString(dept.DeptId)
	lv_db.GetOrmDefault().Table("sys_dept").Where("dept_id=", dept.DeptId).Update("ancestors", dept.Ancestors)
	// ancestors 上级ancestors发生变化，修改下级
	deptList, _ := svc.FindChildren(dept.DeptId)
	if deptList == nil || len(deptList) <= 0 {
		return
	}
	for _, child := range deptList {
		svc.UpdateChildrenAncestors(child, dept.Ancestors)
	}
}

// FindAll 根据分页查询部门管理数据
func (svc *DeptService) FindAll(param *common_vo.DeptPageReq) (*[]models.SysDept, error) {
	if param == nil {
		return svc.SelectDeptList(0, "", "", param.TenantId)
	} else {
		return svc.SelectDeptList(param.ParentId, param.DeptName, param.Status, param.TenantId)
	}
}

// 根据角色ID查询部门
func (svc *DeptService) SelectRoleDeptTree(roleId int) ([]any, error) {
	sql := ` select concat(d.dept_id, d.dept_name) as DeptName 
             from sys_dept d 
             left join sys_role_dept rd  on d.dept_id = rd.dept_id 
             where d.del_flag =0 and rd.role_id = @roleId
             order by d.parent_id, d.order_num
             `
	param := map[string]any{}
	param["roleId"] = roleId
	listMap, err := lv_dao.ListMapByNamedSql(lv_db.GetOrmDefault(), sql, param, false)
	lv_err.HasErrAndPanic(err)
	var result []any
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

// 删除部门管理信息
func (svc *DeptService) DeleteDeptById(deptId int) error {
	var entity models.SysDept
	err := entity.UpdateDelFlag(deptId)
	return err
}

// FindById 根据部门ID查询信息
func (svc *DeptService) FindById(deptId int) (*models.SysDept, error) {
	var dept = new(models.SysDept)
	dept, err := dept.FindById(deptId)
	lv_err.HasErrAndPanic(err)
	if dept.ParentId > 0 {
		var parent = new(models.SysDept)
		parent, err := parent.FindById(dept.ParentId)
		if err == nil {
			dept.ParentName = parent.DeptName
		}
	}
	return dept, err
}

func (svc *DeptService) FindChildren(parentId int) ([]*models.SysDept, error) {
	db := lv_db.GetOrmDefault()
	var rs []*models.SysDept
	err := db.Table("sys_dept").Where("parent_id=?", parentId).Find(&rs).Error
	return rs, err
}

// 查询部门管理数据
func (svc *DeptService) SelectDeptList(parentId int, deptName, status string, tenantId int) (*[]models.SysDept, error) {
	var dao dao.SysDeptDao
	return dao.SelectDeptList(parentId, deptName, status, tenantId)
}

func (svc *DeptService) SelectDeptTreeList() []vo.TreeSelect {
	var deptResults []vo.TreeSelect
	var depts []models.SysDept
	err := lv_db.GetOrmDefault().Where("del_flag = '0'").Order("parent_id, order_num").Find(&depts).Error
	lv_err.HasErrAndPanic(err)
	for i := 0; i < len(depts); i++ {
		dept := depts[i]
		deptId := dept.DeptId
		name := dept.DeptName
		pId := dept.ParentId
		if pId == 0 {
			tChild := svc.getChildList(depts, deptId)
			treeSelect := vo.TreeSelect{
				Id:       deptId,
				Label:    name,
				Children: tChild,
			}
			deptResults = append(deptResults, treeSelect)
		}
	}

	return deptResults
}

func (svc *DeptService) getChildList(depts []models.SysDept, deptId int) []vo.TreeSelect {
	var tlist []vo.TreeSelect
	for i := 0; i < len(depts); i++ {
		dept1 := depts[i]
		id := dept1.DeptId
		pId := dept1.ParentId
		name := dept1.DeptName

		if pId == deptId {
			tChild := svc.getChildList(depts, id)
			tree := vo.TreeSelect{
				Id:       id,
				Label:    name,
				Children: tChild,
			}
			tlist = append(tlist, tree)
		}

	}
	return tlist
}

// 查询部门是否存在用户
func (svc *DeptService) CheckDeptExistUser(deptId int) bool {
	sql := " select count(*) from sys_user where del_flag = 0 "
	param := map[string]any{}
	param["deptId"] = deptId
	sql += " and dept_id= @deptId "
	count, err := lv_dao.CountByNamedSql(lv_db.GetOrmDefault(), sql, param)
	if err != nil {
		panic(err)
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}

// 查询部门人数
func (svc *DeptService) SelectDeptCount(deptId, parentId int) int {
	sql := " select count(*) from sys_dept where del_flag = 0 "
	param := map[string]any{}
	if deptId > 0 {
		param["deptId"] = deptId
		sql += " and dept_id= @deptId "
	}
	if parentId > 0 {
		param["parentId"] = parentId
		sql += " and parent_id= @parentId "
	}
	count, err := lv_dao.CountByNamedSql(lv_db.GetOrmDefault(), sql, param)
	if err != nil {
		panic(err)
	}
	return int(count)
}
