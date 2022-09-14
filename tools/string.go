package tools

import (
	"reflect"
	"unsafe"
)

// StringPtrToBytesPtr 字符串指针转字节切片指针(外部自己保证指针有效性)
func StringPtrToBytesPtr(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	var b []byte
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}
