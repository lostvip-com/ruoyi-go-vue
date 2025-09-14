package util

import "reflect"

// IsStructPtrCanAddr 判断 v 是否为“结构体指针”
// 返回 true 的条件：
//  1. v 本身是指针类型；
//  2. 指针指向的元素是结构体。
func IsStructPtrCanAddr(v any) bool {
	if v == nil {
		return false
	}
	t := reflect.TypeOf(v)
	// 先判断是否为指针类型
	if t.Kind() != reflect.Ptr {
		return false
	}
	// 再判断指针指向的是否为结构体
	if t.Elem().Kind() != reflect.Struct {
		return false
	}
	// 最后检查 Elem() 是否可取地址（避免 panic）
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return false
	}
	return val.Elem().CanAddr()
}

// IsStruct 判断 v 是否为结构体（非指针）
func IsStruct(v any) bool {
	if v == nil {
		return false
	}
	t := reflect.TypeOf(v)
	// 先去掉指针层
	return t.Kind() == reflect.Struct
}
