package field

// 字符串类型键值对序列化符号
var (
	stringValPrefix = [1]byte{'"'}
	stringValSuffix = [1]byte{'"'}
)

// Type 字段的值类型
type Type uint8

// 类型枚举
const (
	TypeInt8 Type = iota
	TypeInt16
	TypeInt
	TypeInt32
	TypeInt64

	TypeByte
	TypeUint8
	TypeUint16
	TypeUint
	TypeUintptr
	TypeUint32
	TypeUint64

	TypeFloat32
	TypeFloat64

	TypeComplex64
	TypeComplex128

	TypeBool
	TypeNull
	TypeString
)

// IsValidRange 是否在有效范围内
func IsValidRange(minType, valType, maxType Type) bool {
	return minType <= valType && valType <= maxType
}

// Field 键值对序列化结构体
type Field struct {
	Key       []byte      // 键的字节流
	ValType   Type        // 值类型
	Integer   int64       // 可转为整型的值
	Bytes     []byte      // 可转为字节流的值
	Interface interface{} // 无法转换的值
	Prefix    []byte      // 值的前缀字节流，如: "、{、[
	Suffix    []byte      // 值的后缀字节流，如: "、}、]
}
