package service

import (
	"common/common_vo"
	"common/util"
	"github.com/gin-gonic/gin"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/spf13/cast"
	dao2 "system/dao"
	"system/model"
)

type DictTypeService struct {
}

// 根据主键查询数据
func (svc *DictTypeService) FindById(id int) (*model.SysDictType, error) {
	dictType := &model.SysDictType{}
	dictType, err := dictType.FindById(id)
	return dictType, err
}

// 根据主键删除数据
func (svc *DictTypeService) DeleteById(id int) bool {
	err := (&model.SysDictType{DictId: id}).Delete()
	if err == nil {
		return true
	}
	return false
}

func (svc *DictTypeService) DeleteByIds(ids string) error {
	ida := util.ToIntArray(ids, ",")
	data := dao2.GetDictDataDaoInstance()
	//data.DeleteBatch()
	tp := new(model.SysDictType)
	for _, id := range ida {
		tp.DictId = cast.ToInt(id)
		_, err := tp.FindOne()
		lv_err.HasErrAndPanic(err)
		//delete dictType and dictData
		err = data.DeleteByType(tp.DictType)
		err = tp.Delete()
		lv_err.HasErrAndPanic(err)
	}
	return nil
}

// 添加数据
func (svc *DictTypeService) AddSave(req *common_vo.AddDictTypeReq, c *gin.Context) (int, error) {
	var entity model.SysDictType
	entity.Status = req.Status
	entity.DictType = req.DictType
	entity.DictName = req.DictName
	entity.Remark = req.Remark
	//entity.CreateTime = time.Now()
	entity.CreateBy = ""
	var userService UserService
	user := userService.GetCurrUser(c)
	if user != nil {
		entity.CreateBy = user.UserName
	}
	err := entity.Save()
	return entity.DictId, err
}

// 修改数据
func (svc *DictTypeService) EditSave(req *common_vo.EditDictTypeReq, c *gin.Context) (int, error) {
	entity := &model.SysDictType{DictId: req.DictId}
	entity, err := entity.FindOne()
	if err != nil {
		return 0, err
	}
	entity.Status = req.Status
	entity.DictType = req.DictType
	entity.DictName = req.DictName
	entity.Remark = req.Remark
	//entity.UpdateTime = time.Now()
	entity.UpdateBy = ""
	var userService UserService
	user := userService.GetCurrUser(c)
	if user == nil {
		entity.UpdateBy = user.UserName
	}
	err = entity.Updates()
	return entity.DictId, err
}

// 根据条件分页查询角色数据
func (svc *DictTypeService) FindAll(params *common_vo.DictTypePageReq) ([]model.SysDictType, error) {
	var dao = dao2.GetSysDictTypeDaoInstance()
	return dao.FindAll(params)
}

// 根据条件分页查询角色数据
func (svc *DictTypeService) FindPage(params *common_vo.DictTypePageReq) ([]model.SysDictType, int, error) {
	var dao = dao2.GetSysDictTypeDaoInstance()
	return dao.FindPage(params)
}

// 根据字典类型查询信息
func (svc *DictTypeService) FindByType(dictType string) *model.SysDictType {
	entity := &model.SysDictType{DictType: dictType}
	entity, err := entity.FindOne()
	if err != nil {
		return nil
	}
	return entity
}

// 导出excel
//func (svc *DictTypeService) Export(param *common_vo.DictTypePageReq) (string, error) {
//
//}

// 检查字典类型是否唯一
func (svc *DictTypeService) CheckDictTypeUniqueAll(configKey string) (bool, error) {
	var entity model.SysDictType
	entity.DictType = configKey
	count, err := entity.Count()

	return count > 0, err
}

// 检查字典类型是否唯一
func (svc *DictTypeService) IsDictTypeExist(configKey string) bool {
	var entity model.SysDictType
	entity.DictType = configKey
	count, err := entity.Count()
	lv_err.HasErrAndPanic(err)
	return count > 0
}

func (svc *DictTypeService) transDictName(entity model.SysDictType) string {
	return `(` + entity.DictName + `)&nbsp;&nbsp;&nbsp;` + entity.DictType
}
