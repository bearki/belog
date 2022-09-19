package field

func Interface(key string, value interface{}) Field {
	switch val := value.(type) {
	case bool:
		return Bool(key, val)
	case *bool:
		return Boolp(key, val)
	case string:
		return String(key, val)
	case *string:
		return Stringp(key, val)
	case int:
		return Int(key, val)
	case *int:
		return Intp(key, val)
	case int8:
		return Int8(key, val)
	case *int8:
		return Int8p(key, val)
	case int16:
		return Int16(key, val)
	case *int16:
		return Int16p(key, val)
	case int32:
		return Int32(key, val)
	case *int32:
		return Int32p(key, val)
	case int64:
		return Int64(key, val)
	case *int64:
		return Int64p(key, val)
	case uint:
		return Uint(key, val)
	case *uint:
		return Uintp(key, val)
	case uint8:
		return Uint8(key, val)
	case *uint8:
		return Uint8p(key, val)
	case uint16:
		return Uint16(key, val)
	case *uint16:
		return Uint16p(key, val)
	case uint32:
		return Uint32(key, val)
	case *uint32:
		return Uint32p(key, val)
	case uint64:
		return Uint64(key, val)
	case *uint64:
		return Uint64p(key, val)
	case float32:
		return Float32(key, val)
	case *float32:
		return Float32p(key, val)
	case float64:
		return Float64(key, val)
	case *float64:
		return Float64p(key, val)
	case complex64:
		return Complex64(key, val)
	case *complex64:
		return Complex64p(key, val)
	case complex128:
		return Complex128(key, val)
	case *complex128:
		return Complex128p(key, val)
	case uintptr:
		return Uintptr(key, val)
	case *uintptr:
		return Uintptrp(key, val)
	default:
		// data, err := json.Marshal(x)
		// if err != nil {

		// }

	}
	return Field{}
}
