package convert

import "math"

func Float32ToInt64(v float32) int64 {
	return int64(math.Float32bits(v))
}

func Float32FromInt64(v int64) float32 {
	return math.Float32frombits(uint32(v))
}

func Float64ToInt64(v float64) int64 {
	return int64(math.Float64bits(v))
}

func Float64FromInt64(v int64) float64 {
	return math.Float64frombits(uint64(v))
}
