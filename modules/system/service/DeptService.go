package service

import (
	"common/common_vo"
	"common/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_dao"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_reflect"
	"github.com/spf13/cast"
	"system/dao"
	"system/vo"
	"time"
)

type DeptService struct{}

var deptService *DeptService

func GetDeptServiceInstance() *DeptService {
	if userService == nil {
		deptService = &DeptService{}
	}
	return deptService
}

// 新增保存信息
func (svc *DeptService) AddSave(req *common_vo.AddDeptReq, c *gin.Context) (*models.SysDept, error) {
	if req.OrderNum == 0 {
		req.OrderNum = 100
	}
	parent := new(models.SysDept)
	parent, err := parent.FindById(req.ParentId)
	if err == nil {
		if parent.Status != "0" {
			return nil, errors.New("部门停用，不允许新增")
		}
	} else {
		return nil, errors.New("父部门不能为空")
	}
	dept0 := new(models.SysDept)
	var userService UserService
	lv_reflect.CopyProperties(req, dept0)
	user := userService.GetProfile(c)
	if user != nil && user.UserId > 0 {
		dept0.CreateBy = user.UserName
	}
	dept0.CreateTime = time.Now()
	//这里跟原版不一样了，多加了一级自己的ID，以方便数据权限控制
	err = dept0.Save()
	if err != nil {
		return nil, err
	}
	dept0.Ancestors = parent.Ancestors + "," + cast.ToString(dept0.DeptId)
	dept0.Update()
	return dept0, err
}

// 修改保存信息
func (svc *DeptService) EditSave(req *common_vo.EditDeptReq, c *gin.Context) (int64, error) {
	dept0 := &models.SysDept{DeptId: req.DeptId}
	err := dept0.FindOne()
	if err != nil {
		return 0, errors.New("数据不存在")
	}
	pdept := &models.SysDept{DeptId: req.ParentId}
	var dao dao.SysDeptDao
	err = pdept.FindOne()
	if pdept.Status != "0" {
		return 0, errors.New("部门停用，不允许新增")
	} else {
		dept0.DeptName = req.DeptName
		dept0.Status = req.Status
		dept0.ParentId = req.ParentId
		dept0.DeptType = req.NodeType
		dept0.Leader = req.Leader
		dept0.Phone = req.Phone
		dept0.OrderNum = req.OrderNum
		var userService UserService
		user := userService.GetProfile(c)
		if user != nil && user.UserId > 0 {
			dept0.UpdateBy = user.UserName
		}
		dept0.UpdateTime = time.Now()
		dept0.Update()
		//递归修改所有子项目
		dao.UpdateDeptChildrenAncestors(dept0, pdept.Ancestors)
		return 1, nil
	}
}

// 根据分页查询部门管理数据
func (svc *DeptService) SelectListAll(param *common_vo.DeptPageReq) (*[]models.SysDept, error) {
	if param == nil {
		return svc.SelectDeptList(0, "", "", param.TenantId)
	} else {
		return svc.SelectDeptList(param.ParentId, param.DeptName, param.Status, param.TenantId)
	}
}

// 删除部门管理信息
func (svc *DeptService) DeleteDeptById(deptId int64) error {
	var dao dao.SysDeptDao
	return dao.DeleteDeptById(deptId)
}

// FindById 根据部门ID查询信息
func (svc *DeptService) FindById(deptId int64) (*models.SysDept, error) {
	var dao dao.SysDeptDao
	dept, err := dao.SelectDeptById(deptId)
	if dept.ParentId > 0 {
		parentDept, err := dao.SelectDeptById(dept.ParentId)
		if err == nil {
			dept.ParentName = parentDept.DeptName
		}
	}
	return dept, err
}

// 根据ID查询所有子部门
func (svc *DeptService) SelectChildrenDeptById(deptId int64) []*models.SysDept {
	var dao dao.SysDeptDao
	return dao.SelectChildrenDeptById(deptId)
}

//
//// 加载部门列表树
//func (svc *DeptService) SelectDeptTree(parentId int64, deptName, status string, tenantId int64) (*[]lv_dto.Ztree, error) {
//	var dao dao.SysDeptDao
//	list, err := dao.SelectDeptList(parentId, deptName, status, tenantId)
//	if err != nil {
//		return nil, err
//	}
//
//	return svc.InitZtree(list, nil), nil
//
//}

// 查询部门管理数据
func (svc *DeptService) SelectDeptList(parentId int64, deptName, status string, tenantId int64) (*[]models.SysDept, error) {
	var dao dao.SysDeptDao
	return dao.SelectDeptList(parentId, deptName, status, tenantId)
}

func (svc *DeptService) SelectDeptTreeList() []vo.TreeSelect {
	var deptResults []vo.TreeSelect
	var depts []models.SysDept
	err := lv_db.GetMasterGorm().Where("del_flag = '0'").Order("parent_id, order_num").Find(&depts).Error
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

func (svc *DeptService) getChildList(depts []models.SysDept, deptId int64) []vo.TreeSelect {
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

//
//// 根据角色ID查询部门（数据权限）
//func (svc *DeptService) RoleDeptTreeData(roleId int64, tenantId int64) (*[]lv_dto.Ztree, error) {
//	var result *[]lv_dto.Ztree
//	var dao dao.SysDeptDao
//	deptList, err := dao.SelectDeptList(0, "", "", tenantId)
//	if err != nil {
//		return nil, err
//	}
//
//	if roleId > 0 {
//		roleDeptList, err := dao.SelectRoleDeptTree(roleId)
//		if err != nil || roleDeptList == nil {
//			result = svc.InitZtree(deptList, nil)
//		} else {
//			result = svc.InitZtree(deptList, &roleDeptList)
//		}
//	} else {
//		result = svc.InitZtree(deptList, nil)
//	}
//	return result, nil
//}

//
//// 对象转部门树
//func (svc *DeptService) InitZtree(deptList *[]models.SysDept, roleDeptList *[]string) *[]lv_dto.Ztree {
//	var result []lv_dto.Ztree
//	isCheck := false
//	if roleDeptList != nil && len(*roleDeptList) > 0 {
//		isCheck = true
//	}
//
//	for i := range *deptList {
//		if (*deptList)[i].Status == "0" {
//			var ztree lv_dto.Ztree
//			ztree.Id = (*deptList)[i].DeptId
//			ztree.Pid = (*deptList)[i].ParentId
//			ztree.Name = (*deptList)[i].DeptName
//			ztree.Title = (*deptList)[i].DeptName
//			ztree.NodeType = (*deptList)[i].DeptType
//			if isCheck {
//				tmp := lv_conv.String((*deptList)[i].DeptId) + (*deptList)[i].DeptName
//				tmpcheck := false
//				for j := range *roleDeptList {
//					if strings.EqualFold((*roleDeptList)[j], tmp) {
//						tmpcheck = true
//						break
//					}
//				}
//				ztree.Checked = tmpcheck
//			}
//			result = append(result, ztree)
//		}
//	}
//	return &result
//}

// 查询部门是否存在用户
func (svc *DeptService) CheckDeptExistUser(deptId int64) bool {
	sql := " select count(*) from sys_user where del_flag = 0 "
	param := map[string]any{}
	param["deptId"] = deptId
	sql += " and dept_id= @deptId "
	count, err := lv_dao.CountByNamedSql(sql, param)
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
func (svc *DeptService) SelectDeptCount(deptId, parentId int64) int64 {
	var dao dao.SysDeptDao
	return dao.SelectDeptCount(deptId, parentId)
}

//
///**
// * 0
// */
//func (svc *DeptService) IsExistDeptName(deptName string, parentId int64) (bool, error) {
//	var dao dao.SysDeptDao
//	dept, err := dao.FindOne(deptName, parentId)
//	if err != nil {
//		return false, err
//	}
//	if err == nil && dept.DeptId > 0 {
//		return true, err
//	} else {
//		return false, err
//	}
//}
