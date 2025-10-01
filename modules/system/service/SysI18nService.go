// ==========================================================================
// LV自动生成业务逻辑层相关代码: 只生成一次,按需修改,再次生成不会覆盖.
// date  : 2025-08-11 07:41:35 &#43;0000 UTC
// author: lostvip
// ==========================================================================
package service

import (
	"errors"
	"github.com/lostvip-com/lv_framework/lv_db"
	"github.com/lostvip-com/lv_framework/lv_db/lv_dao"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_reflect"
	"gorm.io/gorm"
	"system/dao"
	"system/model"
	"system/vo"
	"time"
)

type SysI18nService struct{}

var i18nService *SysI18nService

func GetSysI18nServiceInstance() *SysI18nService {
	if i18nService == nil {
		i18nService = &SysI18nService{}
	}
	return i18nService
}

// FindById 根据主键查询数据
func (svc SysI18nService) FindById(id int) (*model.SysI18n, error) {
	var po = new(model.SysI18n)
	po, err := po.FindById(id)
	return po, err
}

// DeleteById 根据主键删除数据
func (svc SysI18nService) DeleteById(id int) error {
	err := (&model.SysI18n{Id: id}).Delete()
	return err
}

// DeleteByIds 批量删除数据记录
func (svc SysI18nService) DeleteByIds(ids string) (int64, error) {
	ida := lv_conv.ToIntArray(ids, ",")
	var i18nDao = dao.GetSysI18nDaoInstance()
	rows, err := i18nDao.DeleteByIds(ida)
	return rows, err
}

func (d SysI18nService) DeleteByKeys(locale string, localeKeys []string) error {
	var total int64
	err := lv_db.GetOrmDefault().Table("sys_i18n").Select("count(*)").Where("locale=? and locale_key in ? ", locale, localeKeys).Find(&total).Error
	if total == 0 {
		return errors.New("no data found ")
	}
	err = lv_db.GetOrmDefault().Table("sys_i18n").Where("locale=? and locale_key in ? ", locale, localeKeys).Delete(&model.SysI18n{}).Error
	return err
}

// AddSave 添加数据
func (svc SysI18nService) AddSave(form *model.SysI18n) (*model.SysI18n, error) {
	err := form.Save()
	lv_err.HasErrAndPanic(err)
	return form, err
}

// EditSave 修改数据
func (svc SysI18nService) EditSave(form *model.SysI18n) (*model.SysI18n, error) {
	var po = new(model.SysI18n)
	po, err := po.FindById(form.Id)
	if err != nil {
		return nil, err
	}
	_ = lv_reflect.CopyProperties(form, po)
	err = po.Updates()
	return po, err
}
func (svc SysI18nService) SaveOrUpdateBatch(listPtr *[]model.SysI18n) (*[]model.SysI18n, error) {
	list := *listPtr
	listUpdate := make([]model.SysI18n, 0)
	listCreate := make([]model.SysI18n, 0)
	for _, it := range list {
		po, err := new(model.SysI18n).FindOne(it.Locale, it.LocaleKey)
		if err != nil { //不存在
			//it.CreateTime = date
			listCreate = append(listCreate, it)
		} else { //已经存在
			//it.UpdateTime = date
			it.Id = po.Id
			listUpdate = append(listUpdate, it)
		}
	}
	err := lv_dao.Transaction(lv_db.GetOrmDefault(), 18*time.Second, func(tx *gorm.DB) error {
		if len(listCreate) > 0 {
			err := lv_db.GetOrmDefault().Save(&listCreate).Error
			if err != nil {
				return err
			}
		}
		if len(listUpdate) > 0 {
			for _, it := range listUpdate {
				id := it.Id
				it.Id = 0         //不更新ID
				it.LocaleKey = "" //不更新LocaleKey
				err := lv_db.GetOrmDefault().Table(it.TableName()).Where("id=?", id).Updates(&it).Error
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	return listPtr, err
}

// ListByPage 根据条件分页查询数据
func (svc SysI18nService) ListByPage(params *vo.SysI18nReq) (*[]vo.SysI18nResp, int64, error) {
	var i18nDao = dao.GetSysI18nDaoInstance()
	return i18nDao.ListByPage(params)
}

// ExportAll 导出excel
func (svc SysI18nService) ExportAll(param *vo.SysI18nReq) (*[]map[string]string, *[]map[string]any, error) {
	var i18nDao = dao.GetSysI18nDaoInstance()
	listMap, _, err := i18nDao.ListMapByPage(param)
	headerMap := []map[string]string{
		map[string]string{"key": "id", "title": "", "width": "15"},
		map[string]string{"key": "locale", "title": "本地标识", "width": "15"},
		map[string]string{"key": "localeKey", "title": "国际化key", "width": "15"},
		map[string]string{"key": "localeName", "title": "国际化名称", "width": "15"},
		map[string]string{"key": "sort", "title": "字典排序", "width": "15"},
		map[string]string{"key": "remark", "title": "备注", "width": "15"},
		map[string]string{"key": "createTime", "title": "创建日期", "width": "15"},
		map[string]string{"key": "updateTime", "title": "更新日期", "width": "15"},
		map[string]string{"key": "updateBy", "title": "更新者", "width": "15"},
		map[string]string{"key": "createBy", "title": "创建者", "width": "15"},
	}
	return &headerMap, listMap, err
}
