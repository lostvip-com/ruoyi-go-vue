// ==========================================================================
// LV自动生成model扩展代码列表、增、删，改、查、导出，只生成一次，按需修改,再次生成不会覆盖.
// 生成日期：2025-07-24 02:41:35 &#43;0000 UTC
// 生成人：lv
// ==========================================================================
package dao

import (
    "github.com/lostvip-com/lv_framework/lv_db"
    "github.com/lostvip-com/lv_framework/lv_db/lv_batis"
    "github.com/lostvip-com/lv_framework/lv_db/lv_dao"
    "github.com/lostvip-com/lv_framework/utils/lv_err"
    "demo/internal/product/vo"
    "demo/internal/product/model"
)

type IotProductDao struct { }
var productDao *IotProductDao

func GetIotProductDaoInstance() *IotProductDao {
    if productDao == nil {
        productDao = &IotProductDao{}
    }
    return productDao
}
// ListMapByPage 根据条件分页查询数据
func (d IotProductDao) ListMapByPage(req *vo.IotProductReq) (*[]map[string]any, int64, error) {
	ibatis := lv_batis.NewInstance("product/iot_product_mapper.sql") //under the mapper directory
	// 约定用方法名ListByPage对应sql文件中的同名tagName
	limitSQL, err := ibatis.GetLimitSql("ListIotProduct", req)
	//查询数据
	rows, err := lv_dao.ListMapByNamedSql(limitSQL, req,true)
	lv_err.HasErrAndPanic(err)
	count, err := lv_dao.CountByNamedSql(ibatis.GetCountSql(), req)
	lv_err.HasErrAndPanic(err)
	return rows, count, nil
}

// ListByPage 根据条件分页查询数据
func (d IotProductDao) ListByPage(req *vo.IotProductReq) (*[]vo.IotProductResp, int64, error) {
	ibatis := lv_batis.NewInstance("product/iot_product_mapper.sql") //under the mapper directory
	// 对应sql文件中的同名tagName
	limitSQL, err := ibatis.GetLimitSql("ListIotProduct", req)
	//查询数据
	rows, err := lv_dao.ListByNamedSql[vo.IotProductResp](limitSQL, req)
	lv_err.HasErrAndPanic(err)
	count, err := lv_dao.CountByNamedSql(ibatis.GetCountSql(), req)
	lv_err.HasErrAndPanic(err)
	return rows, count, nil
}

// ListAll 导出excel使用
func (d IotProductDao) ListAll(req *vo.IotProductReq, isCamel bool) (*[]map[string]any, error) {
	ibatis := lv_batis.NewInstance("product/iot_product_mapper.sql")
	// 约定用方法名ListByPage对应sql文件中的同名tagName
	sql, err := ibatis.GetSql("ListIotProduct", req)
	lv_err.HasErrAndPanic(err)

	arr, err := lv_dao.ListMapByNamedSql(sql, req, isCamel)
	return arr, err
}

// Find 根据条件查询
func (d IotProductDao) Find(where, order string) (*[]model.IotProduct, error) {
	var list []model.IotProduct
	err := lv_db.GetMasterGorm().Table("iot_product").Where(where).Order(order).Find(&list).Error
	return &list, err
}

// Find 通过主键批量删除
func (d IotProductDao) DeleteByIds(ida []int64) (int64,error) {
	db := lv_db.GetMasterGorm().Table("iot_product").Where("id in ? ", ida).Update("del_flag", 1)
    return db.RowsAffected, db.Error
}
