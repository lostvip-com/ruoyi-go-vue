package util

import (
	"github.com/spf13/cast"
	"strconv"
	"strings"
)

func SplitToInt(data string, sep string) []int {
	var sa = strings.Split(data, sep)
	var sarr []int
	for i := 0; i < len(sa); i++ {
		var v = sa[i]
		var v1, _ = strconv.Atoi(v)
		sarr = append(sarr, v1)
	}
	return sarr
}

func SplitToInt64(data string, sep string) []int {
	var sa = strings.Split(data, sep)
	var sarr []int
	for i := 0; i < len(sa); i++ {
		var v = sa[i]
		sarr = append(sarr, cast.ToInt(v))
	}
	return sarr
}

func RemoveOne(nums []int, val int) []int {
	newNums := make([]int, 0)
	for _, num := range nums {
		if num != val {
			newNums = append(newNums, num)
		}
	}
	return newNums
}
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
