package encoder

import (
	"strconv"
	"strings"
)

// 追加行格式的调用栈
//
//	@param	dst			目标切片
//	@param	fullPath	是否保留完整路径
//	@param	fn			完整文件名
//	@param	ln			行号
//	@param	mn			完整函数名
//	@return	序列化后的调用栈字符串
//
// 返回示例: [test.go:100] [test.TestLogger]
func appendStack(dst []byte, fullPath bool, fn string, ln int, mn string) []byte {
	if !fullPath {
		// 裁剪为基础文件名
		index := strings.LastIndexByte(fn, '/')
		if index > -1 && index+1 < len(fn) {
			fn = fn[index+1:]
		}

		// 裁剪为基础函数名
		index = strings.LastIndexByte(mn, '/')
		if index > 0 && index+1 < len(mn) {
			mn = mn[index+1:]
		}
	}

	// 追加内容
	dst = append(dst, '[')
	dst = append(dst, fn...)
	dst = append(dst, ':')
	dst = strconv.AppendInt(dst, int64(ln), 10)
	dst = append(dst, `] [`...)
	dst = append(dst, mn...)
	dst = append(dst, `]`...)

	// OK
	return dst
}

// 追加JSON格式的调用栈
//
//	@param	dst			目标切片
//	@param	fullPath	是否保留完整路径
//	@param	stackKey	调用栈信息键名
//	@param	fnKey		文件名的JSON键名
//	@param	fn			完整文件名
//	@param	lnKey		行号的JSON键名
//	@param	ln			行号
//	@param	mnKey		函数名的JSON键名
//	@param	mn			完整函数名
//	@return	序列化后的调用栈字符串
//
// 返回示例: "stack": {"file": "test.go", "line": 100, "method": "test.TestLogger"}
func appendStackJSON(dst []byte, fullPath bool, stackKey string, fnKey string, fn string, lnKey string, ln int, mnKey string, mn string) []byte {
	if !fullPath {
		// 裁剪为基础文件名
		index := strings.LastIndexByte(fn, '/')
		if index > -1 && index+1 < len(fn) {
			fn = fn[index+1:]
		}

		// 裁剪为基础函数名
		index = strings.LastIndexByte(mn, '/')
		if index > 0 && index+1 < len(mn) {
			mn = mn[index+1:]
		}
	}

	// 追加内容
	dst = append(dst, '"')
	dst = append(dst, stackKey...)
	dst = append(dst, `": {"`...)
	dst = append(dst, fnKey...)
	dst = append(dst, `": "`...)
	dst = append(dst, fn...)
	dst = append(dst, `", "`...)
	dst = append(dst, lnKey...)
	dst = append(dst, `": `...)
	dst = strconv.AppendInt(dst, int64(ln), 10)
	dst = append(dst, `, "`...)
	dst = append(dst, mnKey...)
	dst = append(dst, `": "`...)
	dst = append(dst, mn...)
	dst = append(dst, `"}`...)

	// OK
	return dst
}
