package field

// Type 字段的值类型
type Type uint8

const (
	// 时间类型
	TypeTime Type = iota

	// 有符号整型
	TypeInt8
	TypeInt16
	TypeInt
	TypeInt32
	TypeDuration
	TypeInt64

	// 无符号整型
	TypeUint8
	TypeByte
	TypeUint16
	TypeUint
	TypeUintptr
	TypeUint32
	TypeUint64

	// 浮点型
	TypeFloat32
	TypeFloat64

	// 复数型
	TypeComplex64
	TypeComplex128

	// 空指针型
	TypeNull

	// 布尔型
	TypeBool

	// 字符串型
	TypeString
	// 错误类型
	TypeError
)

// Field 键值对序列化结构体
type Field struct {
	Key       string      // 键的字节流
	ValType   Type        // 值类型
	Integer   int64       // 可转为整型的值
	String    string      // 可转为字节流的值
	Interface interface{} // 无法转换的值
}
