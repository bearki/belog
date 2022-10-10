package field

// Objecter 自定义类型编码（构造）器
type Objecter interface {
	// ToString 将自定义类型转换为字符串
	ToString() string
	// ToJSON 将自定义类型转换为JSON
	ToJSON() []byte
}

//------------------------------ 值类型转换 ------------------------------//

// Object 格式化自定义类型字段信息
func Object(key string, val Objecter) Field {
	if val == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeObjecter, Interface: val}
}

//------------------------------ 切片类型转换 ------------------------------//

// Objects 格式化[]Objecter字段信息
func Objects(key string, vals []Objecter) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeObjecters, Interface: vals}
}
