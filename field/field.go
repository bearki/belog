package field

// Field 键值对序列化结构体
type Field struct {
	Key       string      // 键的字节流
	Type      Type        // 值类型
	Integer   int64       // 可转为整型的值
	String    string      // 可转为字符串的值
	Interface interface{} // 无法转换的值
}
