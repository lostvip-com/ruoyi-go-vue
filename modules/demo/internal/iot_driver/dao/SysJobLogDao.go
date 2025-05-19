// ==========================================================================
// LV自动生成model扩展代码列表、增、删，改、查、导出，只生成一次，按需修改,再次生成不会覆盖.
// 生成日期：2025-05-17 13:18:47 &#43;0000 UTC
// 生成人：lv
// ==========================================================================
package dao

import (
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_batis"
	"github.com/lostvip-com/lv_framework/lv_db/lv_dao"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
    "demo/internal/iot_driver/vo"
    "demo/internal/iot_driver/model"
)

type SysJobLogDao struct { }
var logDao *SysJobLogDao

func GetSysJobLogDaoInstance() *SysJobLogDao {
	if logDao == nil {
		logDao = &SysJobLogDao{}
	}
	return logDao
}
// ListMapByPage 根据条件分页查询数据
func (d SysJobLogDao) ListMapByPage(req *vo.SysJobLogReq) (*[]map[string]any, int64, error) {
	ibatis := lv_batis.NewInstance("iot_driver/sys_job_log_mapper.sql") //under the mapper directory
	// 约定用方法名ListByPage对应sql文件中的同名tagName
	limitSQL, err := ibatis.GetLimitSql("ListSysJobLog", req)
	//查询数据
	rows, err := lv_dao.ListMapByNamedSql(limitSQL, req,true)
	lv_err.HasErrAndPanic(err)
	count, err := lv_dao.CountByNamedSql(ibatis.GetCountSql(), req)
	lv_err.HasErrAndPanic(err)
	return rows, count, nil
}

// ListByPage 根据条件分页查询数据
func (d SysJobLogDao) ListByPage(req *vo.SysJobLogReq) (*[]vo.SysJobLogResp, int64, error) {
	ibatis := lv_batis.NewInstance("iot_driver/sys_job_log_mapper.sql") //under the mapper directory
	// 对应sql文件中的同名tagName
	limitSQL, err := ibatis.GetLimitSql("ListSysJobLog", req)
	//查询数据
	rows, err := lv_dao.ListByNamedSql[vo.SysJobLogResp](limitSQL, req)
	lv_err.HasErrAndPanic(err)
	count, err := lv_dao.CountByNamedSql(ibatis.GetCountSql(), req)
	lv_err.HasErrAndPanic(err)
	return rows, count, nil
}

// ListAll 导出excel使用
func (d SysJobLogDao) ListAll(req *vo.SysJobLogReq, isCamel bool) (*[]map[string]any, error) {
	ibatis := lv_batis.NewInstance("iot_driver/sys_job_log_mapper.sql")
	// 约定用方法名ListByPage对应sql文件中的同名tagName
	sql, err := ibatis.GetSql("ListSysJobLog", req)
	lv_err.HasErrAndPanic(err)

	arr, err := lv_dao.ListMapByNamedSql(sql, req, isCamel)
	return arr, err
}

// Find 根据条件查询
func (d SysJobLogDao) Find(where, order string) (*[]model.SysJobLog, error) {
	var list []model.SysJobLog
	err := lv_db.GetMasterGorm().Table("sys_job_log").Where(where).Order(order).Find(&list).Error
	return &list, err
}

// Find 通过主键批量删除
func (d SysJobLogDao) DeleteByIds(ida []int64) int64 {
	db := lv_db.GetMasterGorm().Table("sys_job_log").Where("job_log_id in ? ", ida).Update("del_flag", 1)
    if db.Error != nil {
        panic(db.Error)
    }
    return db.RowsAffected
}
