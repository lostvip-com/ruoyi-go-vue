package test

import (
	"fmt"
	"github.com/lostvip-com/lv_framework/lv_log"
	"github.com/lostvip-com/lv_framework/utils/lv_json"
	"github.com/lostvip-com/lv_framework/utils/lv_reflect"
	"github.com/spf13/cast"
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
	user2 := &User{
		ID:   2,
		Name: "李四",
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
	// 通过json标签获取
	if val, ok := lv_reflect.GetFieldValue(user, "name"); ok {
		fmt.Println("Name by tag:", val) // 输出: Name by tag: 张三
	}

	// 获取嵌套字段
	if val, ok := lv_reflect.GetFieldValue(user, "Address.City"); ok {
		fmt.Println("City:", val) // 输出: City: 北京
	}
	// 使用简化版本（仅支持简单字段）
	if val, ok := lv_reflect.GetFieldValueSimple(user, "Name"); ok {
		fmt.Println("Name (simple):", val) // 输出: Name (simple): 张三
	}
	list := make([]User, 0)
	list = append(list, *user)
	list = append(list, *user2)
	mp := make(map[string]User)
	mapKey := "ID"
	for i := range list {
		it := list[i]
		value, ok := lv_reflect.GetFieldValueSimple(it, mapKey)
		if ok {
			mp[cast.ToString(value)] = it
		} else {
			lv_log.Warn("mapKey not found", mapKey)
		}
	}
	fmt.Println("Name (simple):", lv_json.ToJsonStr(list)) // 输出: Name (simple): 张三
	fmt.Println("Name (simple):", lv_json.ToJsonStr(mp))   // 输出: Name (simple): 张三
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
