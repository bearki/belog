package tools

import "unsafe"

// BytesPtrToStringPtr 字节切片指针转字符串指针(外部自己保证指针有效性)
func BytesPtrToStringPtr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
