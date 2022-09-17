package convert

import "unsafe"

// BytesAppendSpace 为字节切片追加空格
func BytesAppendSpace(s *[]byte, num uint) {
	// 计算可用容量
	hcap := uint(cap(*s) - len(*s))
	// 当可用容量小于期望的填充长度时以可用容量为准
	if hcap < num {
		num = hcap
	}
	// 开始填充
	var i uint
	for i = 0; i < num; i++ {
		*s = append(*s, ' ')
	}
}

// BytesToString 字节切片转字符串
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
