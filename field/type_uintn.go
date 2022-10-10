/**
 * @Title 整型及整型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc uint64|uint32|uint|uint16|uint8|byte|uintptr
 * @Desc *uint64|*uint32|*uint|*uint16|*uint8|*byte|uintptr
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

//------------------------------ 值类型转换 ------------------------------//

// Byte 格式化byte类型字段信息
func Byte(key string, val byte) Field {
	return intn(key, int64(val), TypeByte)
}

// Uint8 格式化uint8类型字段信息
func Uint8(key string, val uint8) Field {
	return intn(key, int64(val), TypeUint8)
}

// Uint16 格式化uint16类型字段信息
func Uint16(key string, val uint16) Field {
	return intn(key, int64(val), TypeUint16)
}

// Uint 格式化uint类型字段信息
func Uint(key string, val uint) Field {
	return intn(key, int64(val), TypeUint)
}

// Uintptr 格式化uintptr类型字段信息
func Uintptr(key string, val uintptr) Field {
	return intn(key, int64(val), TypeUintptr)
}

// Uint32 格式化uint32类型字段信息
func Uint32(key string, val uint32) Field {
	return intn(key, int64(val), TypeUint32)
}

// Uint64 格式化uint64类型字段信息
func Uint64(key string, val uint64) Field {
	return intn(key, int64(val), TypeUint64)
}

//------------------------------ 指针类型转换 ------------------------------//

// Bytep 格式化*byte类型字段信息
func Bytep(key string, valp *byte) Field {
	if valp == nil {
		return nullField(key)
	}
	return Byte(key, *valp)
}

// Uint8p 格式化*uint8类型字段信息
func Uint8p(key string, valp *uint8) Field {
	if valp == nil {
		return nullField(key)
	}
	return Uint8(key, *valp)
}

// Uint16p 格式化*uint16类型字段信息
func Uint16p(key string, valp *uint16) Field {
	if valp == nil {
		return nullField(key)
	}
	return Uint16(key, *valp)
}

// Uintp 格式化*uint类型字段信息
func Uintp(key string, valp *uint) Field {
	if valp == nil {
		return nullField(key)
	}
	return Uint(key, *valp)
}

// Uintptrp 格式化*uintptr类型字段信息
func Uintptrp(key string, valp *uintptr) Field {
	if valp == nil {
		return nullField(key)
	}
	return Uintptr(key, *valp)
}

// Uint32p 格式化*uint32类型字段信息
func Uint32p(key string, valp *uint32) Field {
	if valp == nil {
		return nullField(key)
	}
	return Uint32(key, *valp)
}

// Uint64p 格式化*uint64类型字段信息
func Uint64p(key string, valp *uint64) Field {
	if valp == nil {
		return nullField(key)
	}
	return Uint64(key, *valp)
}

//------------------------------ 切片类型转换 ------------------------------//

// Bytes 格式化byte类型字段信息
func Bytes(key string, vals []byte) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeBytes, Interface: vals}
}

// Uint8s 格式化uint8类型字段信息
func Uint8s(key string, vals []uint8) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeUint8s, Interface: vals}
}

// Uint16s 格式化uint16类型字段信息
func Uint16s(key string, vals []uint16) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeUint16s, Interface: vals}
}

// Uints 格式化uint类型字段信息
func Uints(key string, vals []uint) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeUints, Interface: vals}
}

// Uintptrs 格式化uintptr类型字段信息
func Uintptrs(key string, vals []uintptr) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeUintptrs, Interface: vals}
}

// Uint32s 格式化uint32类型字段信息
func Uint32s(key string, vals []uint32) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeUint32s, Interface: vals}
}

// Uint64s 格式化uint64类型字段信息
func Uint64s(key string, vals []uint64) Field {
	if vals == nil {
		return nullField(key)
	}
	return Field{Key: key, Type: TypeUint64s, Interface: vals}
}
