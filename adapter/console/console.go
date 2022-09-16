/**
 *@Title 控制台日志记录适配器
 *@Desc 控制台打印将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package console

import (
	"bytes"
	"os"
	"strconv"
	"time"

	"github.com/bearki/belog/v2/internal/convert"
	"github.com/bearki/belog/v2/logger"
)

// 控制台字体颜色字节表
var (
	// 灰色
	colorGrayStartBytes = []byte{27, 91, 57, 48, 109}
	// 蓝色
	colorBlueStartBytes = []byte{27, 91, 51, 52, 109}
	// 绿色
	colorGreenStartBytes = []byte{27, 91, 51, 50, 109}
	// 黄色
	colorYellowStartBytes = []byte{27, 91, 51, 51, 109}
	// 红色
	colorRedStartBytes = []byte{27, 91, 51, 49, 109}
	// 洋红色
	colorMagentaStartBytes = []byte{27, 91, 51, 53, 109}
	// 重置
	colorResetBytes = []byte{27, 91, 48, 109}
)

// GetLevelConsoleColorBytes 获取日志级别对应的控制台颜色字节表
func GetLevelConsoleColorBytes(l logger.Level) []byte {
	switch l {
	case logger.LevelTrace: // 通知级别(灰色)
		return colorGrayStartBytes
	case logger.LevelDebug: // 调试级别(蓝色)
		return colorBlueStartBytes
	case logger.LevelInfo: // 普通级别(绿色)
		return colorGreenStartBytes
	case logger.LevelWarn: // 警告级别(黄色)
		return colorYellowStartBytes
	case logger.LevelError: // 错误级别(红色)
		return colorRedStartBytes
	case logger.LevelFatal: // 紧急级别(洋红色)
		return colorMagentaStartBytes
	default:
		return nil
	}
}

// Adapter 控制台日志适配器
type Adapter struct{}

// New 创建控制台日志适配器
func New() logger.Adapter {
	return new(Adapter)
}

// Name 用于获取适配器名称
//
// 注意：请确保适配器名称不与其他适配器名称冲突
func (e *Adapter) Name() string {
	return "belog-console-adapter"
}

// Print 普通日志打印方法
//
// @params logTime 日记记录时间
//
// @params level 日志级别
//
// @params content 日志内容
func (e *Adapter) Print(logTime time.Time, level logger.Level, content []byte) {
	// 不带调用栈：
	// 2022/09/14 20:28:13.793 <colorStart>[T]<colorEnd>  this is a trace log\r\n
	// ++++++++++_++++++++++++_============+++==========__+++++++++++++++++++____
	//     10    1    12      1      5      3      4    2       content        2
	//
	// const initSize = 10 + 1 + 12 + 1 + 5 + 3 + 4 + 2 + 2 = 40

	// 计算需要的大小
	size := 40 + len(content)
	// 创建一个指定容量的切片，避免二次扩容
	logSlice := make([]byte, 0, size)

	// 追加格式化好的日期和时间
	logSlice = append(logSlice, convert.StringToBytes(logTime.Format("2006/01/02 15:04:05.000"))...) // 23个字节
	// 追加空格
	logSlice = append(logSlice, ' ') // 1个字节
	// 追加颜色开始
	logSlice = append(logSlice, GetLevelConsoleColorBytes(level)...) // 5个字节
	// 追加级别
	logSlice = append(logSlice, '[', level.GetLevelChar(), ']') // 3个字节
	// 追加颜色结束
	logSlice = append(logSlice, colorResetBytes...) // 4个字节
	// 追加空格
	logSlice = append(logSlice, ' ', ' ') // 2个字节
	// 追加日志内容
	logSlice = append(logSlice, content...) // len(content)个字节
	// 追加回车换行
	logSlice = append(logSlice, '\r', '\n') // 2个字节

	// 打印到标准输出
	os.Stdout.Write(logSlice)
}

// PrintStack 调用栈日志打印方法
//
// @params logTime 日记记录时间
//
// @params level 日志级别
//
// @params content 日志内容
//
// @params fileName 日志记录调用文件路径
//
// @params lineNo 日志记录调用文件行号
//
// @params methodName 日志记录调用函数名
func (e *Adapter) PrintStack(logTime time.Time, level logger.Level, content []byte, fileName []byte, lineNo int, methodName []byte) {
	// 带调用栈：
	// 2022/09/14 20:28:13.793 <colorStart>[T]<colorEnd> [belog_test.go:82] [PrintLog]  this is a trace log\r\n
	// ++++++++++_++++++++++++_============+++==========_++++++++++++++++++_++++++++++__+++++++++++++++++++____
	//     10    1    12      1      5      3      4    1   3+file+line    1 2+method  2   content           2
	//
	// const initSize = 10 + 1 + 12 + 1 + 5 + 3 + 4 + 1 + 3 + 1 + 2 + 2 + 2 = 47

	// 裁剪为基础文件名
	index := bytes.LastIndexByte(fileName, '/')
	if index > -1 && index+1 < len(fileName) {
		fileName = fileName[index+1:]
	}

	// 裁剪为基础函数名
	index = bytes.LastIndexByte(methodName, '/')
	if index > 0 && index+1 < len(methodName) {
		methodName = methodName[index+1:]
	}

	// 转换行号为切片
	lineNoBytes := convert.StringToBytes(strconv.Itoa(lineNo))
	// 计算需要的大小
	size := 47 + len(fileName) + len(lineNoBytes) + len(methodName) + len(content)
	// 创建一个指定容量的切片，避免二次扩容
	logSlice := make([]byte, 0, size)

	// 追加格式化好的日期和时间
	logSlice = append(logSlice, convert.StringToBytes(logTime.Format("2006/01/02 15:04:05.000"))...) // 23个字节
	// 追加空格
	logSlice = append(logSlice, ' ') // 1个字节
	// 追加颜色开始
	logSlice = append(logSlice, GetLevelConsoleColorBytes(level)...) // 5个字节
	// 追加级别
	logSlice = append(logSlice, '[', level.GetLevelChar(), ']') // 3个字节
	// 追加颜色结束
	logSlice = append(logSlice, colorResetBytes...) // 4个字节
	// 追加空格
	logSlice = append(logSlice, ' ') // 1个字节
	// 追加文件名和行号
	logSlice = append(logSlice, '[')            // 1个字节
	logSlice = append(logSlice, fileName...)    // len(fileName)个字节
	logSlice = append(logSlice, ':')            // 1个字节
	logSlice = append(logSlice, lineNoBytes...) // len(lineNo)个字节
	logSlice = append(logSlice, ']')            // 1个字节
	// 追加空格
	logSlice = append(logSlice, ' ') // 1个字节
	// 追加函数名
	logSlice = append(logSlice, '[')           // 1个字节
	logSlice = append(logSlice, methodName...) // len(methodName)个字节
	logSlice = append(logSlice, ']')           // 1个字节
	// 追加空格
	logSlice = append(logSlice, ' ', ' ') // 2个字节
	// 追加日志内容
	logSlice = append(logSlice, content...) // len(content)个字节
	// 追加回车换行
	logSlice = append(logSlice, '\r', '\n') // 2个字节

	// 打印到标准输出
	os.Stdout.Write(logSlice)
}

// Flush 日志缓存刷新
//
// 用于日志缓冲区刷新
// 接收到该通知后需要立即将缓冲区中的日志持久化
func (e *Adapter) Flush() {}
