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

func SplitToInt64(data string, sep string) []int64 {
	var sa = strings.Split(data, sep)
	var sarr []int64
	for i := 0; i < len(sa); i++ {
		var v = sa[i]
		sarr = append(sarr, cast.ToInt64(v))
	}
	return sarr
}
