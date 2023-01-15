package convert

import "math"

// Float32ToInt64 float32类型转int64类型
func Float32ToInt64(v float32) int64 {
	return int64(math.Float32bits(v))
}

// Float32FromInt64 int64类型转float32类型
func Float32FromInt64(v int64) float32 {
	return math.Float32frombits(uint32(v))
}

// Float64ToInt64 float64类型转int64类型
func Float64ToInt64(v float64) int64 {
	return int64(math.Float64bits(v))
}

// Float64FromInt64 int64类型转float64类型
func Float64FromInt64(v int64) float64 {
	return math.Float64frombits(uint64(v))
}
