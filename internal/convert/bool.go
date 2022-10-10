package convert

// BoolToInt64 布尔型转整型
func BoolToInt64(value bool) int64 {
	if !value {
		return 0
	}
	return 1
}

// BoolFromInt64 整型转布尔型
func BoolFromInt64(value int64) bool {
	return value != 0
}
