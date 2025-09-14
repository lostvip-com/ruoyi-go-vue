package util

import (
	"bytes"
	"common/service"
	"errors"
	"github.com/lostvip-com/lv_framework/lv_log"
	"reflect"
	"text/template"
)

// ProcessTpl 后面优化：使用缓存
func ProcessTpl(locale, key string, dataPtr any) string {
	if locale == "" {
		return ""
	}
	tpl, err := template.New("i18n").Parse(key)
	if err != nil {
		lv_log.Error(" ProcessTpl error !!!!", locale, key, err.Error())
		return ""
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, dataPtr)
	if err != nil {
		lv_log.Error("XXXXXXXXXXXXXX" + err.Error())
		return ""
	}
	localeKey := buf.String()
	localName := service.GetI18nServiceInstance().GetI18nText(locale, localeKey)
	return localName
}

// TranslateI18nTag 解析结构体中所有带有i18n标签的字段
func TranslateI18nTag(local string, dataPtr any) error {
	if local == "" {
		return nil
	}
	v := reflect.ValueOf(dataPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("ptr must be pointer to struct")
	}
	v = v.Elem() // 拿到结构体值
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i) // 字段元信息
		tag := sf.Tag.Get("i18n")
		if tag == "" {
			continue
		}

		localVal := ProcessTpl(local, tag, dataPtr)
		if localVal == "" {
			continue
		}

		fv := v.Field(i)  // 字段值
		if !fv.CanSet() { // 1. 检查是否可写
			continue
		}
		if fv.Kind() != reflect.String { // 2. 检查类型是否匹配
			continue
		}
		fv.SetString(localVal)
	}
	return nil
}
