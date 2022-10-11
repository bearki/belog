package field

// Type 字段的值类型
type Type uint8

const (
	// 未知类型
	TypeUnknown Type = iota

	/*---------------------- 普通类型开始 ----------------------*/
	NormalTypeStart

	// 时间类型
	TypeTime
	// 有符号整型
	TypeInt8
	TypeInt16
	TypeInt
	TypeInt32
	TypeInt64
	TypeDuration
	// 无符号整型
	TypeUint8
	TypeUint16
	TypeUint
	TypeUint32
	TypeUint64
	TypeByte
	TypeUintptr
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
	// 自定义类型
	TypeObjecter

	/*---------------------- 普通类型结束 ----------------------*/
	NormalTypeEnd

	/*---------------------- 切片类型开始 ----------------------*/
	SliceTypeStart

	// 时间类型切片
	TypeTimes
	// 有符号整型切片
	TypeInt8s
	TypeInt16s
	TypeInts
	TypeInt32s
	TypeDurations
	TypeInt64s
	// 无符号整型切片
	TypeUint8s
	TypeBytes
	TypeUint16s
	TypeUints
	TypeUintptrs
	TypeUint32s
	TypeUint64s
	// 浮点型切片
	TypeFloat32s
	TypeFloat64s
	// 复数型切片
	TypeComplex64s
	TypeComplex128s
	// 布尔型切片
	TypeBools
	// 字符串型切片
	TypeStrings
	// 错误类型切片
	TypeErrors
	// 自定义类型切片
	TypeObjecters

	/*---------------------- 切片类型结束 ----------------------*/
	SliceTypeEnd
)
