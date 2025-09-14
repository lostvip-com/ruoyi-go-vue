package lv_if

/**
 * 自定义三元表达式
 */
func IfTrue(condition bool, trueResult, falseResult interface{}) interface{} {
	if condition {
		return trueResult
	}
	return falseResult
}

// IfEmpty 模仿sql中的ifnull函数
func IfEmpty(strToCheck, replacedStr string) string {
	if strToCheck == "" {
		return replacedStr
	}
	return strToCheck
}

// If0 模仿sql中的ifnull函数
func If0(intToCheck, replacedInt int) int {
	if intToCheck == 0 {
		return replacedInt
	}
	return intToCheck
}
