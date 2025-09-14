package lv_reflect

import (
	"reflect"
	"strings"
)

// GetFieldValue 高效获取结构体字段值（无缓存版本）
func GetFieldValue(obj interface{}, fieldName string) (interface{}, bool) {
	if obj == nil || fieldName == "" {
		return nil, false
	}

	// 获取对象的反射值
	val := reflect.ValueOf(obj)

	// 解引用指针
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, false
		}
		val = val.Elem()
	}

	// 确保是结构体
	if val.Kind() != reflect.Struct {
		return nil, false
	}

	// 处理嵌套字段名，如 "User.Name"
	fieldNames := strings.Split(fieldName, ".")

	// 逐级访问字段
	for i, name := range fieldNames {
		name = strings.TrimSpace(name)

		// 查找当前级别的字段
		field := findField(val, name)
		if !field.IsValid() {
			return nil, false
		}

		// 更新val为当前字段值
		val = field

		// 如果不是最后一个字段，需要确保是结构体或结构体指针
		if i < len(fieldNames)-1 {
			// 解引用指针
			for val.Kind() == reflect.Ptr {
				if val.IsNil() {
					return nil, false
				}
				val = val.Elem()
			}

			// 确保是结构体
			if val.Kind() != reflect.Struct {
				return nil, false
			}
		}
	}

	// 返回最终字段值
	if val.IsValid() && val.CanInterface() {
		return val.Interface(), true
	}

	return nil, false
}

// findField 在结构体中查找字段（支持字段名和标签）
func findField(val reflect.Value, name string) reflect.Value {
	typ := val.Type()

	// 遍历所有字段
	for i := 0; i < val.NumField(); i++ {
		fieldType := typ.Field(i)
		fieldValue := val.Field(i)

		// 直接匹配字段名（不区分大小写）
		if strings.EqualFold(fieldType.Name, name) {
			return fieldValue
		}

		// 匹配json标签
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			tagName := strings.Split(jsonTag, ",")[0]
			if tagName == name {
				return fieldValue
			}
		}

		// 可以添加其他常用标签的匹配
		// 例如bson标签
		if bsonTag := fieldType.Tag.Get("bson"); bsonTag != "" {
			tagName := strings.Split(bsonTag, ",")[0]
			if tagName == name {
				return fieldValue
			}
		}
	}

	return reflect.Value{}
}

// GetFieldValueSimple 简化版本，只支持简单字段名（不支持嵌套）
func GetFieldValueSimple(obj interface{}, fieldName string) (interface{}, bool) {
	if obj == nil || fieldName == "" {
		return nil, false
	}

	val := reflect.ValueOf(obj)

	// 解引用指针
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, false
		}
		val = val.Elem()
	}

	// 确保是结构体
	if val.Kind() != reflect.Struct {
		return nil, false
	}

	// 直接查找字段
	field := val.FieldByNameFunc(func(n string) bool {
		return strings.EqualFold(n, fieldName)
	})

	// 如果没找到，尝试通过标签查找
	if !field.IsValid() {
		field = findFieldByJsonTag(val, fieldName)
	}

	if field.IsValid() && field.CanInterface() {
		return field.Interface(), true
	}

	return nil, false
}

// findFieldByTag 通过标签查找字段
func findFieldByJsonTag(val reflect.Value, tagName string) reflect.Value {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		// 检查json标签
		if jsonTag := field.Tag.Get("json"); jsonTag != "" {
			if strings.Split(jsonTag, ",")[0] == tagName {
				return val.Field(i)
			}
		}
	}
	return reflect.Value{}
}
