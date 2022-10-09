package convert

// BoolToInt 布尔型转整型
func BoolToInt(value bool) int {
	if !value {
		return 0
	}
	return 1
}

// IntToBool 整型转布尔型
func IntToBool(value int) bool {
	return value != 0
}
