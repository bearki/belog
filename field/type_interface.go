package field

import "time"

// Interface 格式化任意类型的字段信息
//	注意：未在枚举列表的类型将使用反射进行格式化，性能会有所降低
func Interface(key string, val interface{}) Field {
	switch v := val.(type) {
	case time.Time:
		return Time(key, v)
	case *time.Time:
		return Timep(key, v)
	case []time.Time:
		return Times(key, v)
	case int8:
		return Int8(key, v)
	case *int8:
		return Int8p(key, v)
	case []int8:
		return Int8s(key, v)
	case int16:
		return Int16(key, v)
	case *int16:
		return Int16p(key, v)
	case []int16:
		return Int16s(key, v)
	case int:
		return Int(key, v)
	case *int:
		return Intp(key, v)
	case []int:
		return Ints(key, v)
	case int32:
		return Int32(key, v)
	case *int32:
		return Int32p(key, v)
	case []int32:
		return Int32s(key, v)
	case int64:
		return Int64(key, v)
	case *int64:
		return Int64p(key, v)
	case []int64:
		return Int64s(key, v)
	case time.Duration:
		return Duration(key, v)
	case *time.Duration:
		return Durationp(key, v)
	case []time.Duration:
		return Durations(key, v)
	case uint8: // or byte
		return Uint8(key, v)
	case *uint8:
		return Uint8p(key, v)
	case []uint8:
		return Uint8s(key, v)
	case uint16:
		return Uint16(key, v)
	case *uint16:
		return Uint16p(key, v)
	case []uint16:
		return Uint16s(key, v)
	case uint:
		return Uint(key, v)
	case *uint:
		return Uintp(key, v)
	case []uint:
		return Uints(key, v)
	case uint32:
		return Uint32(key, v)
	case *uint32:
		return Uint32p(key, v)
	case []uint32:
		return Uint32s(key, v)
	case uint64:
		return Uint64(key, v)
	case *uint64:
		return Uint64p(key, v)
	case []uint64:
		return Uint64s(key, v)
	case uintptr:
		return Uintptr(key, v)
	case *uintptr:
		return Uintptrp(key, v)
	case []uintptr:
		return Uintptrs(key, v)
	case float32:
		return Float32(key, v)
	case *float32:
		return Float32p(key, v)
	case []float32:
		return Float32s(key, v)
	case float64:
		return Float64(key, v)
	case *float64:
		return Float64p(key, v)
	case []float64:
		return Float64s(key, v)
	case complex64:
		return Complex64(key, v)
	case *complex64:
		return Complex64p(key, v)
	case []complex64:
		return Complex64s(key, v)
	case complex128:
		return Complex128(key, v)
	case *complex128:
		return Complex128p(key, v)
	case []complex128:
		return Complex128s(key, v)
	case nil:
		return nullField(key)
	case bool:
		return Bool(key, v)
	case *bool:
		return Boolp(key, v)
	case []bool:
		return Bools(key, v)
	case string:
		return String(key, v)
	case *string:
		return Stringp(key, v)
	case []string:
		return Strings(key, v)
	case error:
		return Error(key, v)
	case *error:
		return Errorp(key, v)
	case []error:
		return Errors(key, v)
	case Objecter:
		return Object(key, v)
	case []Objecter:
		return Objects(key, v)

	default:
		return Field{
			Key:       key,
			Type:      TypeUnknown,
			Interface: val,
		}
	}
}

// Interface 格式化任意类型的字段信息
// @Desc 注意：未在枚举列表的类型将使用反射进行格式化，性能会有所降低
func Any(key string, val interface{}) Field {
	return Interface(key, val)
}
