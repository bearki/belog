package field

//------------------------------ 值类型转换 ------------------------------//

// Error 格式化error类型字段信息
func Error(name string, value error) Field {
	if value == nil {
		return nullField(name)
	}
	return Field{Key: name, ValType: TypeError, String: value.Error()}
}

//------------------------------ 指针类型转换 ------------------------------//

// Errorp 格式化*error类型字段信息
func Errorp(name string, valuep *error) Field {
	if valuep == nil {
		return nullField(name)
	}
	return Error(name, *valuep)
}

//------------------------------ 切片类型转换 ------------------------------//

// Errors 格式化[]error类型字段信息
func Errors(name string, values []error) Field {
	return Field{Key: name, ValType: TypeErrors, Interface: values}
}
