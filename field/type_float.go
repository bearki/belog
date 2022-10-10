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
)

//------------------------------ 值类型转换 ------------------------------//

// Float32 格式化float32类型字段信息
func Float32(key string, val float32) Field {
	return Field{Key: key, Type: TypeFloat32, Integer: int64(math.Float32bits(val))}
}

// Float64 格式化float64类型字段信息
func Float64(key string, val float64) Field {
	return Field{Key: key, Type: TypeFloat64, Integer: int64(math.Float64bits(val))}
}

//------------------------------ 指针类型转换 ------------------------------//

// Float32p 格式化*float32类型字段信息
func Float32p(key string, valp *float32) Field {
	if valp == nil {
		return nullField(key)
	}
	return Float32(key, *valp)
}

// Float64p 格式化*float64类型字段信息
func Float64p(key string, valp *float64) Field {
	if valp == nil {
		return nullField(key)
	}
	return Float64(key, *valp)
}

//------------------------------ 切片类型转换 ------------------------------//

// Float32s 格式化[]float32类型字段信息
func Float32s(key string, vals []float32) Field {
	return Field{Key: key, Type: TypeFloat32s, Interface: vals}
}

// Float64s 格式化[]float64类型字段信息
func Float64s(key string, vals []float64) Field {
	return Field{Key: key, Type: TypeFloat64s, Interface: vals}
}
