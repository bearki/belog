package convert

import (
	"reflect"
	"unsafe"
)

// StringToBytes 字符串转字节切片
func StringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	var b []byte
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}

// BytesToString 字节切片转字符串
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
