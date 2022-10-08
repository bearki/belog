/**
 * @Title 整型及整型指针的键值对序列化
 * @Desc 支持对以下类型进行序列化:
 * @Desc uint64|uint32|uint|uint16|uint8|byte
 * @Desc *uint64|*uint32|*uint|*uint16|*uint8|*byte
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

import (
	"unsafe"
)

// uintn 组装~uint类型的字段结构
func uintn(name string, value int64, vt Type) Field {
	return Field{
		Key:     name,
		ValType: vt,
		Integer: value,
	}
}

//------------------------------ 值类型转换 ------------------------------//

// Byte 格式化byte类型字段信息
func Byte(name string, value byte) Field {
	return uintn(name, int64(value), TypeByte)
}

// Uint8 格式化uint8类型字段信息
func Uint8(name string, value uint8) Field {
	return uintn(name, int64(value), TypeUint8)
}

// Uint16 格式化uint16类型字段信息
func Uint16(name string, value uint16) Field {
	return uintn(name, int64(value), TypeUint16)
}

// Uint 格式化uint类型字段信息
func Uint(name string, value uint) Field {
	return uintn(name, int64(value), TypeUint)
}

// Uintptr 格式化uintptr类型字段信息
func Uintptr(name string, value uintptr) Field {
	return uintn(name, int64(value), TypeUintptr)
}

// Uint32 格式化uint32类型字段信息
func Uint32(name string, value uint32) Field {
	return uintn(name, int64(value), TypeUint32)
}

// Uint64 格式化uint64类型字段信息
func Uint64(name string, value uint64) Field {
	return uintn(name, int64(value), TypeUint64)
}

//------------------------------ 指针类型转换 ------------------------------//

// Bytep 格式化*byte类型字段信息
func Bytep(name string, valuep *byte) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}

// Uint8p 格式化*uint8类型字段信息
func Uint8p(name string, valuep *uint8) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}

// Uint16p 格式化*uint16类型字段信息
func Uint16p(name string, valuep *uint16) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}

// Uintp 格式化*uint类型字段信息
func Uintp(name string, valuep *uint) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}

// Uintptrp 格式化*uintptr类型字段信息
func Uintptrp(name string, valuep *uintptr) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}

// Uint32p 格式化*uint32类型字段信息
func Uint32p(name string, valuep *uint32) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}

// Uint64p 格式化*uint64类型字段信息
func Uint64p(name string, valuep *uint64) Field {
	if f, ok := CheckPtr(name, unsafe.Pointer(valuep)); !ok {
		return f
	}
	return Uint64(name, uint64(*valuep))
}
