package encoder

import (
	"runtime"
	"strconv"
	"strings"
)

// GetCallStack 获取调用栈信息
//
// @params skip 需要跳过的调用栈数量
//
// @return fn 文件名字节切片
//
// @return ln 行号
//
// @return mn 函数名
func GetCallStack(skip uint) (fn string, ln int, mn string) {
	// 获取调用栈信息
	pc, fn, ln, _ := runtime.Caller(int(skip))

	// 获取函数名字节切片
	if funcForPC := runtime.FuncForPC(pc); funcForPC != nil {
		mn = funcForPC.Name()
	}

	// OK
	return
}

// AppendStack 追加行格式的调用栈
//
// @params dst 目标切片
//
// @params fullPath 是否保留完整路径
//
// @params fn 完整文件名
//
// @params ln 行号
//
// @params mn 完整函数名
//
// @return 序列化后的调用栈字符串
//
// 返回示例，反引号内为实际内容:
// `[test.go:100] [test.TestLogger]`
func AppendStack(dst []byte, fullPath bool, fn string, ln int, mn string) []byte {
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

// AppendStackJSON 追加JSON格式的调用栈
//
// @params dst 目标切片
//
// @params fullPath 是否保留完整路径
//
// @params stackKey 调用栈信息键名
//
// @params fnKey 文件名的JSON键名
//
// @params fn 完整文件名
//
// @params lnKey 行号的JSON键名
//
// @params ln 行号
//
// @params mnKey 函数名的JSON键名
//
// @params mn 完整函数名
//
// @return 序列化后的调用栈字符串
//
// 返回示例，反引号内为实际内容:
// `"stack": {"file": "test.go", "line": 100, "method": "test.TestLogger"}`
func AppendStackJSON(dst []byte, fullPath bool, stackKey string, fnKey string, fn string, lnKey string, ln int, mnKey string, mn string) []byte {
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
