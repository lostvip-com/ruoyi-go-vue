package api

import (
	"common/common_vo"
	"common/models"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/spf13/cast"
	"strings"
	"system/service"
)

type DeptApi struct {
	BaseApi
}

/*排除节点*/
func (w *DeptApi) ExcludeDept(c *gin.Context) {
	deptId := c.Param("deptId")
	svc := service.GetDeptServiceInstance()
	dept, err := svc.FindById(cast.ToInt64(deptId))
	listPtr, err := svc.FindAll(&common_vo.DeptPageReq{})
	if err != nil {
		util.Fail(c, err.Error())
	} else {
		list := *listPtr
		dist := make([]models.SysDept, 0)
		for _, it := range list {
			if strings.Contains(it.Ancestors, dept.Ancestors) {
				continue
			}
			dist = append(dist, it)
		}
		util.Success(c, dist)
	}
}

// ListAjax 列表分页数据
func (w *DeptApi) ListAjax(c *gin.Context) {
	var req = common_vo.DeptPageReq{}
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	result, err := service.GetDeptServiceInstance().FindAll(&req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, result)
}

// AddSave 新增页面保存
func (w *DeptApi) AddSave(c *gin.Context) {
	var req *models.SysDept
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	w.FillInCreate(c, &req.BaseModel)
	rid, err := service.GetDeptServiceInstance().AddSave(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, rid)
}

// EditSave 修改页面保存
func (w *DeptApi) EditSave(c *gin.Context) {
	var req *models.SysDept
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, err.Error())
		return
	}
	w.FillInUpdate(c, &req.BaseModel)
	po, err := service.GetDeptServiceInstance().EditSave(req)
	if err != nil {
		util.Fail(c, err.Error())
		return
	}
	util.Success(c, po)
}

// 删除数据
func (w *DeptApi) Remove(c *gin.Context) {
	id := lv_conv.Int64(c.Param("deptId"))
	err := service.GetDeptServiceInstance().DeleteDeptById(id)
	if err != nil {
		util.Fail(c, err.Error())
	} else {
		util.Success(c, id)
	}
}

// 删除数据
func (w *DeptApi) GetDept(c *gin.Context) {
	id := lv_conv.Int64(c.Param("deptId"))
	service := service.GetDeptServiceInstance()
	dept, err := service.FindById(id)
	if err != nil {
		util.Fail(c, err.Error())
	} else {
		util.Success(c, dept)
	}
}

//
//// 加载部门列表树结构的数据
//func (w *DeptApi) TreeData(c *gin.Context) {
//	var service service.DeptService
//	tenantId := session.GetTenantId(c)
//	result, _ := service.SelectDeptTree(0, "", "", tenantId)
//	c.JSON(http.StatusOK, result)
//}
//
//// 加载角色部门（数据权限）列表树
//func (w *DeptApi) RoleDeptTreeData(c *gin.Context) {
//	var service service.DeptService
//	tenantId := session.GetTenantId(c)
//	roleId := lv_conv.Int64(c.Query("roleId"))
//	result, err := service.RoleDeptTreeData(roleId, tenantId)
//
//	if err != nil {
//		util.ErrorResp(c).SetMsg(err.Error()).Log("菜单树", gin.H{"roleId": roleId})
//		return
//	}
//
//	c.JSON(http.StatusOK, result)
//}
