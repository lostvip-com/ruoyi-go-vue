package lv_conv

import (
	"github.com/spf13/cast"
	"strings"
)

func String(str any) string {
	return cast.ToString(str)
}
func Int64(str any) int64 {
	return cast.ToInt64(str)
}

func SubStr(str string, startIndex, endIndex int) string {
	rs := []rune(str)
	return string(rs[startIndex:endIndex])
}

// 将带分隔符的字符串切成int64数组
func ToIntArray(str, split string) []int {
	result := make([]int, 0)
	if str == "" {
		return result
	}
	arr := strings.Split(str, split)
	if len(arr) > 0 {
		for i := range arr {
			if arr[i] != "" {
				result = append(result, cast.ToInt(arr[i]))
			}
		}
	}
	return result
}
