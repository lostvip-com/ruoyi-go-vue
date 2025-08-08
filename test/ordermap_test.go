package test

import (
	"fmt"
	"github.com/iancoleman/orderedmap"
	"github.com/lostvip-com/lv_framework/utils/lv_json"
	"testing"
)

func TestOrderMap(t *testing.T) {

	user1 := &User{
		ID:   1,
		Name: "张三",
		Age:  25,
	}

	user2 := &User{
		ID:   2,
		Name: "张三",
		Age:  25,
	}

	user3 := &User{
		ID:   1,
		Name: "张三",
		Age:  25,
	}
	o := orderedmap.New()

	o.Set("3", user3)
	o.Set("1", user1)
	o.Set("2", user2)

	str := lv_json.ToJsonStr(o)
	fmt.Println("map :", str) // 输出: Name (simple): 张三

	fmt.Println("map :", str) // 输出: Name (simple): 张三

	fmt.Println("map :", str) // 输出: Name (simple): 张三
	//fmt.Println("Name (simple):", lv_json.ToJsonStr(mp))   // 输出: Name (simple): 张三
}
