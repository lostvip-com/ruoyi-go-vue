// ==========================================================================
// LV自动生成业务逻辑层相关代码: 只生成一次,按需修改,再次生成不会覆盖.
// date  : {{.CreateTime}}
// author: {{.FunctionAuthor}}
// ==========================================================================
package service

import (
	"time"
	"github.com/lostvip-com/lv_framework/utils/lv_conv"
	"github.com/lostvip-com/lv_framework/utils/lv_err"
	"github.com/lostvip-com/lv_framework/utils/lv_office"
	"github.com/lostvip-com/lv_framework/utils/lv_reflect"
    "{{.ModuleName}}/internal/{{.PackageName}}/model"
    "{{.ModuleName}}/internal/{{.PackageName}}/dao"
    "{{.ModuleName}}/internal/{{.PackageName}}/vo"
)
type {{.ClassName}}Service struct{}

var {{.BusinessName}}Service *{{.ClassName}}Service

func Get{{.ClassName}}ServiceInstance() *{{.ClassName}}Service {
	if {{.BusinessName}}Service == nil {
		{{.BusinessName}}Service = &{{.ClassName}}Service{}
	}
	return {{.BusinessName}}Service
}

// FindById 根据主键查询数据
func (svc {{.ClassName}}Service) FindById(id {{.PkColumn.GoType}}) (*model.{{.ClassName}}, error) {
	var po = new(model.{{.ClassName}})
    po,err := po.FindById(id)
	return po, err
}

// DeleteById 根据主键删除数据
func (svc {{.ClassName}}Service) DeleteById(id {{.PkColumn.GoType}}) error {
	err := (&model.{{.ClassName}}{ {{.PkColumn.GoField}}: id}).Delete()
	return err
}

// DeleteByIds 批量删除数据记录
func (svc {{.ClassName}}Service) DeleteByIds(ids string) int64 {
	ida := lv_conv.ToInt64Array(ids, ",")
    var d dao.{{.ClassName}}Dao
	rowsAffected := d.DeleteByIds(ida)
	return rowsAffected
}

// AddSave 添加数据
func (svc {{.ClassName}}Service) AddSave(form *model.{{.ClassName}}) (*model.{{.ClassName}}, error) {
	err := form.Save()
	lv_err.HasErrAndPanic(err)
	return form, err
}

// EditSave 修改数据
func (svc {{.ClassName}}Service) EditSave(form *model.{{.ClassName}})  (*model.{{.ClassName}}, error) {
	var po = new(model.{{.ClassName}})
	po,err := po.FindById(form.{{.PkColumn.GoField}})
    if err!=nil{
        return nil,err
    }
	_ = lv_reflect.CopyProperties(req, po)
	err = po.Updates()
	return po,err
}

// ListByPage 根据条件分页查询数据
func (svc {{.ClassName}}Service) ListByPage(params *vo.{{.ClassName}}Req) (*[]vo.{{.ClassName}}Resp,int64, error) {
	var {{.BusinessName}}Dao = dao.Get{{.ClassName}}DaoInstance()
	return {{.BusinessName}}Dao.ListByPage(params)
}

// ExportAll 导出excel
func (svc {{.ClassName}}Service) ExportAll(param *vo.{{.ClassName}}Req) (string, error) {
    var {{.BusinessName}}Dao = dao.Get{{.ClassName}}DaoInstance()
    var err error
    var listMap *[]map[string]any
    if param.PageNum > 0 { //分页导出
        listMap, _, err = {{.BusinessName}}Dao.ListMapByPage(param)
    } else { //全部导出
        listMap, err = {{.BusinessName}}Dao.ListAll(param, true)
    }
    lv_err.HasErrAndPanic(err)
	heads := []string{ {{range $index, $column := .Columns}} {{if eq $index 0}}"{{$column.ColumnComment}}"{{else}},"{{$column.ColumnComment}}"{{end}}{{end}}}
	keys  := []string{ {{range $index, $column := .Columns}} {{if eq $index 0}}"{{$column.HtmlField}}"{{else}},"{{$column.HtmlField}}"{{end}}{{end}}}
	url, err := lv_office.DownlaodExcelByListMap(&heads, &keys, listMap)
	return url, err
}