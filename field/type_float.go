/**
 * @Title 浮点型及浮点型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc float32|float64
 * @Desc *float32|*float64
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

import (
	"math"
	"unsafe"
)

//------------------------------ 值类型转换 ------------------------------//

// Float32 格式化float32类型字段信息
func Float32(name string, value float32) Field {
	return Field{
		Key:     name,
		ValType: TypeFloat32,
		Integer: int64(math.Float32bits(value)),
	}
}

// Float64 格式化float64类型字段信息
func Float64(name string, value float64) Field {
	return Field{
		Key:     name,
		ValType: TypeFloat64,
		Integer: int64(math.Float64bits(value)),
	}
}

//------------------------------ 指针类型转换 ------------------------------//

// Float32p 格式化*float32类型字段信息
func Float32p(name string, valuep *float32) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Float32(name, *valuep)
}

// Float64p 格式化*float64类型字段信息
func Float64p(name string, valuep *float64) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Float64(name, *valuep)
}
