package service

import (
	"bytes"
	"common/global"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_file"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"system/vo"
)

// ////////////////////////////////////////////////////////////////
// 存入本地的sqlite
//
// ////////////////////////////////////////////////////////////////
type CodeGenService struct {
}

type TplInfo struct {
	PathSrc string
	NameSrc string

	PathDist string
	NameDist string
}

var funcMap = template.FuncMap{
	"contains":   Contains,
	"upperFirst": UpperFirst,
	"substr":     Substr,
	"replace":    strings.Replace,
	"index":      strings.Index,
}

func (e *CodeGenService) ListTpl() []TplInfo {
	var list []TplInfo
	for _, dir := range global.BaseFilePathArr {
		dir = dir + string(os.PathSeparator) + "tpl_gen"
		if !lv_file.IsFileExist(dir) {
			continue
		}
		err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".tpl") {
				tpl := TplInfo{PathSrc: filepath.Dir(path), NameSrc: filepath.Base(path)}
				list = append(list, tpl)
			} else {
				lv_log.Warn("。。。。。。skip! not end with .tpl " + path)
			}
			return nil
		})
		if err != nil {
			return nil
		}
	}
	return list
}

func (e *CodeGenService) PreviewCode(tab *vo.GenTableVO) map[string]map[string]string {
	mapAll := make(map[string]map[string]string)
	listTpl := e.ListTpl()
	for _, tpl := range listTpl {
		b1, err := e.GenCodeByTpl(tab, &tpl)
		if err != nil {
			lv_log.Error(err)
			continue
		}
		groupKey := filepath.Ext(tpl.NameDist) //获取后缀做为分组key 如 java,go,sql
		if mapAll[groupKey] == nil {           //是否存在此组
			mapAll[groupKey] = make(map[string]string)
		}
		mapAll[groupKey][tpl.NameDist] = b1.String()
	}
	return mapAll
}

func (e *CodeGenService) GenCodeByTpl(tab *vo.GenTableVO, tpl *TplInfo) (*bytes.Buffer, error) {
	e.replaceTplVar(tpl, tab)
	file := filepath.Join(tpl.PathSrc, tpl.NameSrc)
	t1, err := template.New(tpl.NameSrc).Funcs(funcMap).ParseFiles(file)
	if err != nil {
		return nil, err
	}
	var b1 bytes.Buffer
	err = t1.Execute(&b1, tab)
	return &b1, err
}

func (e *CodeGenService) GenCode(tab *vo.GenTableVO, overwrite bool) {
	//内部模板
	srcTpl := e.ListTpl()
	for _, tpl := range srcTpl {
		buff, err := e.GenCodeByTpl(tab, &tpl)
		targetPath := lv_file.GetCurrentPath() + tpl.PathDist + "/" + tpl.NameDist
		if overwrite {
			targetPath, err = lv_file.FileCreate(buff, targetPath)
		} else if !lv_file.IsFileExist(targetPath) {
			targetPath, err = lv_file.FileCreate(buff, targetPath)
		}
		lv_log.Info("生成文件:", err, targetPath)
	}
}

// 读取模板
//func (e *CodeGenService) LoadTemplate(templateName string, data interface{}) (string, error) {
//	cur, err := os.Getwd()
//	if err != nil {
//		return "", err
//	}
//	b, err := os.ReadFile(filepath.Join(cur, "resources", "template", templateName))
//	if err != nil {
//		return "", err
//	}
//	templateStr := string(b)
//	tmpl, err := template.New(templateName).Funcs(funcMap).Parse(templateStr) //建立一个模板，内容是"hello, {{OssUrl}}"
//	if err != nil {
//		return "", err
//	}
//	buffer := bytes.NewBufferString("")
//	err = tmpl.Execute(buffer, data) //将string与模板合成，变量name的内容会替换掉{{OssUrl}}
//	return buffer.String(), err
//}

/**
 * 替换模板变量中的路径变量
 */
func (e *CodeGenService) replaceTplVar(tpl *TplInfo, tab *vo.GenTableVO) {
	//替换路径中的占位符
	if tpl.PathDist == "" {
		tpl.PathDist = strings.ReplaceAll(tpl.PathSrc, "resources"+string(os.PathSeparator)+"tpl_gen"+string(os.PathSeparator), "")
	}
	tpl.PathDist = strings.ReplaceAll(tpl.PathDist, "{{.PackageName}}", tab.PackageName)
	tpl.PathDist = strings.ReplaceAll(tpl.PathDist, "{{.ModuleName}}", tab.ModuleName)
	tpl.PathDist = strings.ReplaceAll(tpl.PathDist, "{{.BusinessName}}", tab.BusinessName)
	tpl.PathDist = strings.ReplaceAll(tpl.PathDist, "{{.ClassName}}", tab.ClassName)
	tpl.PathDist = strings.ReplaceAll(tpl.PathDist, "{{.FuncName}}", tab.FunctionName)
	tpl.PathDist = strings.ReplaceAll(tpl.PathDist, "{{.FunctionName}}", tab.FunctionName)
	tpl.PathDist = strings.ReplaceAll(tpl.PathDist, "{{.Table_Name}}", tab.Table_Name)
	tpl.PathDist = strings.ReplaceAll(tpl.PathDist, "{{.Tmp}}", "tmp")
	_ = lv_file.PathCreateIfNotExist(tpl.PathDist)
	//替换文件名中的占位符
	tpl.NameDist = tpl.NameSrc[0 : len(tpl.NameSrc)-4]
	tpl.NameDist = strings.ReplaceAll(tpl.NameDist, "{{.ModuleName}}", tab.ModuleName)
	tpl.NameDist = strings.ReplaceAll(tpl.NameDist, "{{.BusinessName}}", tab.PackageName)
	tpl.NameDist = strings.ReplaceAll(tpl.NameDist, "{{.BusinessName}}", tab.BusinessName)
	tpl.NameDist = strings.ReplaceAll(tpl.NameDist, "{{.ClassName}}", tab.ClassName)
	tpl.NameDist = strings.ReplaceAll(tpl.NameDist, "{{.FuncName}}", tab.FunctionName)
	tpl.NameDist = strings.ReplaceAll(tpl.NameDist, "{{.FunctionName}}", tab.FunctionName)
	tpl.NameDist = strings.ReplaceAll(tpl.NameDist, "{{.Table_Name}}", tab.Table_Name)
}
