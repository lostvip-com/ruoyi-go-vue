// ==========================================================================
// LV自动生成业务逻辑层相关代码: 只生成一次,按需修改,再次生成不会覆盖.
// date  : 2025-07-24 02:41:35 &#43;0000 UTC
// author: lv
// ==========================================================================
package service

import (
    "demo/internal/product/dao"
    "demo/internal/product/model"
    "demo/internal/product/vo"
    "github.com/lostvip-com/lv_framework/utils/lv_conv"
    "github.com/lostvip-com/lv_framework/utils/lv_err"
    "github.com/lostvip-com/lv_framework/utils/lv_reflect"
)

type IotProductService struct{}

var productService *IotProductService

func GetIotProductServiceInstance() *IotProductService {
	if productService == nil {
		productService = &IotProductService{}
	}
	return productService
}

// FindById 根据主键查询数据
func (svc IotProductService) FindById(id int64) (*model.IotProduct, error) {
	var po = new(model.IotProduct)
	po, err := po.FindById(id)
	return po, err
}

// DeleteById 根据主键删除数据
func (svc IotProductService) DeleteById(id int64) error {
	err := (&model.IotProduct{Id: id}).Delete()
	return err
}

// DeleteByIds 批量删除数据记录
func (svc IotProductService) DeleteByIds(ids string) (int64, error) {
	ida := lv_conv.ToInt64Array(ids, ",")
	var productDao = dao.GetIotProductDaoInstance()
	rows, err := productDao.DeleteByIds(ida)
	return rows, err
}

// AddSave 添加数据
func (svc IotProductService) AddSave(form *model.IotProduct) (*model.IotProduct, error) {
	err := form.Save()
	lv_err.HasErrAndPanic(err)
	return form, err
}

// EditSave 修改数据
func (svc IotProductService) EditSave(form *model.IotProduct) (*model.IotProduct, error) {
	var po = new(model.IotProduct)
	po, err := po.FindById(form.Id)
	if err != nil {
		return nil, err
	}
	_ = lv_reflect.CopyProperties(form, po)
	err = po.Updates()
	return po, err
}

// ListByPage 根据条件分页查询数据
func (svc IotProductService) ListByPage(params *vo.IotProductReq) (*[]vo.IotProductResp, int64, error) {
	var productDao = dao.GetIotProductDaoInstance()
	return productDao.ListByPage(params)
}

// ExportAll 导出excel
func (svc IotProductService) ExportAll(param *vo.IotProductReq) (*[]map[string]string, *[]map[string]any, error) {
	var productDao = dao.GetIotProductDaoInstance()
	listMap, _, err := productDao.ListMapByPage(param)
	headerMap := []map[string]string{
		map[string]string{"key": "id", "title": "主键", "width": "15"},
		map[string]string{"key": "key", "title": "产品编码,对应可监控类型ID", "width": "15"},
		map[string]string{"key": "name", "title": "名字", "width": "15"},
		map[string]string{"key": "cloudProductId", "title": "云产品ID", "width": "15"},
		map[string]string{"key": "cloudInstanceId", "title": "云实例ID", "width": "15"},
		map[string]string{"key": "platform", "title": "平台", "width": "15"},
		map[string]string{"key": "protocol", "title": "协议", "width": "15"},
		map[string]string{"key": "nodeType", "title": "节点类型", "width": "15"},
		map[string]string{"key": "netType", "title": "网络类型", "width": "15"},
		map[string]string{"key": "dataFormat", "title": "数据类型", "width": "15"},
		map[string]string{"key": "lastSyncTime", "title": "最后一次同步时间", "width": "15"},
		map[string]string{"key": "factory", "title": "工厂名称", "width": "15"},
		map[string]string{"key": "description", "title": "描述", "width": "15"},
		map[string]string{"key": "status", "title": "产品状态", "width": "15"},
		map[string]string{"key": "extra", "title": "扩展字段", "width": "15"},
		map[string]string{"key": "delFlag", "title": "删除标记", "width": "15"},
		map[string]string{"key": "createTime", "title": "创建日期", "width": "15"},
		map[string]string{"key": "updateTime", "title": "更新日期", "width": "15"},
		map[string]string{"key": "updateBy", "title": "更新者", "width": "15"},
		map[string]string{"key": "createBy", "title": "创建者", "width": "15"},
		map[string]string{"key": "manufacturer", "title": "生产厂商", "width": "15"},
		map[string]string{"key": "tenantId", "title": "租户id", "width": "15"},
	}
	return &headerMap, listMap, err
}
