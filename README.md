# 关于
这是一个高度解耦的日志框架，支持多适配器同时输出，你可以发挥自己的想象力，随意的创建自己喜爱的适配器，在 `v0.1+` 版本中，我们已经提供了两个简单的适配器实现 `console` 、 `file` ，你会注意到它们都实现了 `logger` 中的 `Adapter` 接口，只要该接口的适配器均可挂载到 `logger` 中，在已实现的两个适配器中均有示例，你可以查看源码来编写自己的适配器。

# 安装
```shell
go get github.com/bearki/belog/v2
```

# 快速使用
如果你想快速体验 `belog` 的特性，我们已内置了一个默认 `console` 实例，你可以直接使用以下方式输出你的日志内容
```go
package main

import (
	"github.com/bearki/belog/v2"
)

func main() {
	// 直接使用方式（仅能输出到控制台）
	belog.Trace("this is a trace log")
	belog.Debug("this is a debug log")
	belog.Info("this is a info log")
	belog.Warn("this is a warn log")
	belog.Error("this is a error log")
	belog.Fatal("this is a fatal log")
}
```

# 实例使用
正常情况下我们一般会使用这种方式来使用 `belog` ，我们也推荐使用该方式来记录日志
```go
package main

import (
	"github.com/bearki/belog/v2"
	"github.com/bearki/belog/v2/adapter/file"
	"github.com/bearki/belog/v2/adapter/logger"
)

func main() {
	// 初始化文件日志适配器
	fielAdapter, err := file.New(file.Options{
		LogPath:      "../logs/test_new.log", // 日志储存路径
		MaxSize:      100,                    // 日志单文件大小
		MaxLines:     1000000,                // 单文件最大行数
		SaveDay:      7,                      // 日志保存天数
		Async:        true,                   // 开启异步写入(main函数提前结束会导致日志未写入)
		AsyncChanCap: 100,                    // 异步缓存管道容量
	})
	if err != nil {
		panic("file adapter create failed, " + err.Error())
	}

	// 初始化一个实例(可实例化任意适配器)
	mylog, err := belog.New(fielAdapter)
	if err != nil {
		panic("belog logger create failed, " + err.Error())
	}

	// 配置日志记录级别
	mylog.SetLevel(
		logger.LevelTrace,
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
	)

	// 开启调用栈打印
	mylog.PrintCallStack()

	// 实例对象记录日志
	for i := 0; i < 1000000; i++ {
		mylog.Trace("this is a trace log")
	}
}
```

# 二次封装
你可以参考 `belog.go` 的方式将 `belog` 进行二次封装，这样通过自己的包即可记录日志，防止在过多的包中引入第三方包，便于后期的管理，值得注意的是，二次封装时需要配置函数栈层数，否则将会造成文件名及行数捕获不一致的问题，大多数情况下采用如下层级即可
```go
belogDefault, _ = belog.New(console.New(), nil) // 初始化适配器
belogDefault.SetSkip(1)                         // 配置函数栈
```

# 自定义适配器
没有过多的繁琐操作，实现 `logger` 中的 `Adapter` 接口即可完成适配器自定义工作，最简实现方式在 `console` 适配器中
### 适配器接口
```go
// Adapter 适配器接口
type Adapter interface {
	// Name 用于获取适配器名称
	//
	// 注意：请确保适配器名称不与其他适配器名称冲突
	Name() string

	// Print 普通日志打印方法
	//
	// @params logTime 日记记录时间
	//
	// @params level 日志级别
	//
	// @params content 日志内容
	Print(logTime time.Time, level Level, content []byte)

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
	PrintStack(logTime time.Time, level Level, content []byte, fileName string, lineNo int, methodName string)
}
```
### 适配器实现
```go
/**
 *@Title 控制台日志记录适配器
 *@Desc 控制台打印将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package console

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bearki/belog/v2/logger"
	"github.com/bearki/belog/v2/pkg/tool"
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
	// 2022/09/14 20:28:13.793 [T]  this is a trace log
	// 日期(10) + 空格(1) + 时间(12) + 空格(1) + 颜色开始(5) + 级别(3) + 颜色结束(4) + 空格(2) + 日志内容(len(content)) + 回车换行(2)
	//
	// 计算需要的大小
	size := 41 + len(content)
	// 创建一个指定容量的切片，避免二次扩容
	logSlice := make([]byte, 0, size)
	// 追加格式化好的日期和时间
	logSlice = append(logSlice, tool.StringToBytes(logTime.Format("2006/01/02 15:04:05.000"))...) // 23个字节
	// 追加级别对应的颜色
	logSlice = append(logSlice, GetLevelConsoleColorBytes(level)...)
	// 追加级别
	logSlice = append(logSlice, ' ', '[', level.GetLevelChar(), ']') // 4个字节
	// 追加颜色结束
	logSlice = append(logSlice, colorResetBytes...) // 4个字节
	// 追加日志内容
	logSlice = append(logSlice, ' ', ' ')   // 2个字节
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
func (e *Adapter) PrintStack(logTime time.Time, level logger.Level, content []byte, fileName string, lineNo int, methodName string) {
	// 带调用栈：
	// 2022/09/14 20:28:13.793 [T] [belog_test.go:82] [PrintLog]  this is a trace log
	// 日期(10) + 空格(1) + 时间(12) + 空格(1) + 颜色开始(5) +  级别(3) + 颜色结束(4) + 空格(1) + 文件名和行数(len(fileName) + 3 + 行数(5)) + 空格(1) + 函数名(2+len(methodName)) + 空格(2) + 日志内容(len(fileName)) + 回车换行(2)
	//
	// 裁剪为基础文件名
	fileName = filepath.Base(fileName)
	// 计算需要的大小
	size := 51 + len(content) + len(fileName) + len(methodName)
	// 创建一个指定容量的切片，避免二次扩容
	logSlice := make([]byte, 0, size)
	// 追加格式化好的日期和时间
	logSlice = append(logSlice, tool.StringToBytes(logTime.Format("2006/01/02 15:04:05.000"))...) // 23个字节
	// 追加级别对应的颜色
	logSlice = append(logSlice, GetLevelConsoleColorBytes(level)...)
	// 追加级别
	logSlice = append(logSlice, ' ', '[', level.GetLevelChar(), ']') // 4个字节
	// 追加颜色结束
	logSlice = append(logSlice, colorResetBytes...) // 4个字节
	// 追加文件名和行号，len(strconv.FormatInt(int64(fileName), 10))大于5个字节时，logSlice会发生扩容
	logSlice = append(logSlice, ' ', '[')                                                    // 2个字节
	logSlice = append(logSlice, tool.StringToBytes(fileName)...)                             // len(fileName)个字节
	logSlice = append(logSlice, ':')                                                         // 1个字节
	logSlice = append(logSlice, tool.StringToBytes(strconv.FormatInt(int64(lineNo), 10))...) // 默认5个字节
	logSlice = append(logSlice, ']')                                                         // 1个字节
	// 追加函数名
	logSlice = append(logSlice, ' ', '[')                          // 2个字节
	logSlice = append(logSlice, tool.StringToBytes(methodName)...) // len(methodName)个字节
	logSlice = append(logSlice, ']')                               // 1个字节
	// 追加日志内容
	logSlice = append(logSlice, ' ', ' ')    // 2个字节
	logSlice = append(logSlice, fileName...) // len(content)个字节
	// 追加回车换行
	logSlice = append(logSlice, '\r', '\n') // 2个字节

	// 打印到标准输出
	os.Stdout.Write(logSlice)
}
```
### 适配器挂载
```go
package mylog

import (
	"${your_module_name}/console"
)

func main() {
    logs, err = belog.New(console.New(), nil)
    if err != nil {
        panic(err.Error())
    }
    logs.Info("log init success")
}
```
