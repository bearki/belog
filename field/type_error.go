package field

//------------------------------ 值类型转换 ------------------------------//

// Error 格式化error类型字段信息
func Error(key string, val error) Field {
	if val == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeError, String: val.Error()}
}

//------------------------------ 指针类型转换 ------------------------------//

// Errorp 格式化*error类型字段信息
func Errorp(key string, valp *error) Field {
	if valp == nil {
		return nullField(key)
	}
	return Error(key, *valp)
}

//------------------------------ 切片类型转换 ------------------------------//

// Errors 格式化[]error类型字段信息
func Errors(key string, vals []error) Field {
	return Field{Key: key, Type: TypeErrors, Interface: vals}
}
