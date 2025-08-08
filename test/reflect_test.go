package test

import (
	"fmt"
	"github.com/lostvip-com/lv_framework/utils/lv_reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	user := &User{
		ID:   1,
		Name: "张三",
		Age:  25,
		Address: Address{
			City:   "北京",
			Street: "长安街",
		},
	}

	// 获取简单字段
	if val, ok := lv_reflect.GetFieldValue(user, "Name"); ok {
		fmt.Println("Name:", val) // 输出: Name: 张三
	}

	// 获取嵌套字段
	if val, ok := lv_reflect.GetFieldValue(user, "Address.City"); ok {
		fmt.Println("City:", val) // 输出: City: 北京
	}

	// 通过json标签获取
	if val, ok := lv_reflect.GetFieldValue(user, "name"); ok {
		fmt.Println("Name by tag:", val) // 输出: Name by tag: 张三
	}

	// 使用简化版本（仅支持简单字段）
	if val, ok := lv_reflect.GetFieldValueSimple(user, "Name"); ok {
		fmt.Println("Name (simple):", val) // 输出: Name (simple): 张三
	}
}

// 示例结构体
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address Address
}

type Address struct {
	City   string `json:"city"`
	Street string `json:"street"`
}
