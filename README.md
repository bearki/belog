# BeLog 高性能日志库（结构化日志）
[![Go Version](https://img.shields.io/github/go-mod/go-version/bearki/belog)](https://go.dev/)
[![Go Doc](https://pkg.go.dev/badge/github.com/bearki/belog/v2.svg)](https://pkg.go.dev/github.com/bearki/belog/v2)
[![Latest Release](https://img.shields.io/github/v/release/bearki/belog)](https://github.com/bearki/belog/releases)

这是一个高度解耦的日志框架，支持多适配器同时输出，你可以发挥自己的想象力，随意的创建自己喜爱的适配器；我们已经提供了几个简单的适配器实现，你会注意到它们都实现了 `logger` 中的 `Adapter` 接口，只要该接口的适配器均可挂载到 `logger` 中，你可以查看我们的适配器源码来编写自己的适配器。

## 基准测试
1、测试静态字符串格式化
```txt
goos: windows
goarch: amd64
pkg: github.com/bearki/belog/v2/test
cpu: Intel(R) Core(TM) i5-10500 CPU @ 3.10GHz
BenchmarkBelogLoggerFormatStatic
BenchmarkBelogLoggerFormatStatic-12
 6137994               194.0 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/bearki/belog/v2/test 1.703s


> 测试运行完成时间: 2022/10/13 13:06:25 <
```

2、测试5个字段格式化
```txt
goos: windows
goarch: amd64
pkg: github.com/bearki/belog/v2/test
cpu: Intel(R) Core(TM) i5-10500 CPU @ 3.10GHz
BenchmarkBelogLoggerFormatFiveFields
BenchmarkBelogLoggerFormatFiveFields-12
 1696510               697.2 ns/op           320 B/op          1 allocs/op
PASS
ok      github.com/bearki/belog/v2/test 2.209s


> 测试运行完成时间: 2022/10/13 13:27:31 <
```

3、测试10个字段格式化
```txt
goos: windows
goarch: amd64
pkg: github.com/bearki/belog/v2/test
cpu: Intel(R) Core(TM) i5-10500 CPU @ 3.10GHz
BenchmarkBelogLoggerFormatTenFields
BenchmarkBelogLoggerFormatTenFields-12
  936278              1180 ns/op             640 B/op          1 allocs/op
PASS
ok      github.com/bearki/belog/v2/test 1.176s


> 测试运行完成时间: 2022/10/13 13:28:54 <
```

4、测试切片字段格式化
```txt
goos: windows
goarch: amd64
pkg: github.com/bearki/belog/v2/test
cpu: Intel(R) Core(TM) i5-10500 CPU @ 3.10GHz
BenchmarkBelogLoggerFormatSlice
BenchmarkBelogLoggerFormatSlice-12
 2113472               573.2 ns/op           248 B/op          3 allocs/op
PASS
ok      github.com/bearki/belog/v2/test 2.093s


> 测试运行完成时间: 2022/10/13 13:35:35 <
```

4、测试反射字段格式化
```txt
goos: windows
goarch: amd64
pkg: github.com/bearki/belog/v2/test
cpu: Intel(R) Core(TM) i5-10500 CPU @ 3.10GHz
BenchmarkBelogLoggerFormatInterface
BenchmarkBelogLoggerFormatInterface-12
  649614              1753 ns/op            1096 B/op         18 allocs/op
PASS
ok      github.com/bearki/belog/v2/test 1.218s


> 测试运行完成时间: 2022/10/13 13:41:16 <
```

## 安装
```shell
go get -u github.com/bearki/belog/v2
```

## 快速使用
如果你想快速体验 `BeLog` 的特性，我们已内置了一个默认 `Console` 实例，你可以直接使用以下方式输出你的日志内容
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

## 实例使用
正常情况下我们一般会使用这种方式来使用 `BeLog` ，我们也推荐使用该方式来记录日志
```go
package main

import (
	"github.com/bearki/belog/v2"
	"github.com/bearki/belog/v2/adapter/file"
	"github.com/bearki/belog/v2/field"
	"github.com/bearki/belog/v2/logger"
)

func main() {
	// 初始化文件日志适配器
	fileAdapter, err := file.New(file.Options{
		LogPath:      "../logs/test_new.log", // 日志储存路径
		MaxSize:      200,                    // 日志单文件大小
		MaxLines:     1000000,                // 单文件最大行数
		SaveDay:      7,                      // 日志保存天数
		Async:        true,                   // 开启异步写入(main函数提前结束会导致日志未写入)
		AsyncChanCap: 100,                    // 异步缓存管道容量
	})
	if err != nil {
		b.Fatalf("file adapter create failed, %s\r\n", err)
	}

	// 配置日志记录器参数
	opt := logger.Option{
		EnabledStackPrint: true,
		Encoder:           logger.NewJsonEncoder(logger.DefaultJsonOption),
	}

	// 初始化一个实例(可实例化任意适配器)
	l, err := belog.New(opt, fileAdapter)
	if err != nil {
		b.Fatalf("belog logger create failed, %s\r\n", err)
	}

	// 程序结束需要刷新缓冲区
	defer l.Flush()

	// 开始测试
	for i := 0; i < b.N; i++ {
		l.Trace(
			"this is a trace log",
			field.Bool("key0", i%2 == 0),
			field.Int8("key1", 1),
			field.Ints("int10", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			field.Interface("any", map[string]int{"1": 1, "2": 2, "3": 3}),
		)
	}
}
```

## 二次封装
你可以参考 `belog.go` 的方式将 `BeLog` 进行二次封装，这样通过自己的包即可记录日志，防止在过多的包中引入第三方包，便于后期的管理，值得注意的是，二次封装时需要配置函数栈层数，否则将会造成文件名及行数捕获不一致的问题，大多数情况下采用如下层级即可
```go
DefaultLog.SetSkip(1)
```

## 自定义适配器
没有过多的繁琐操作，实现 `logger/interface.go` 中的 `Adapter` 接口即可完成适配器自定义工作，最简实现方式在 `discard` 适配器中
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
	Print(logTime time.Time, level level.Level, content []byte)

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
	PrintStack(logTime time.Time, level level.Level, content []byte, fileName string, lineNo int, methodName string)

	// Flush 日志缓存刷新
	//
	// 用于日志缓冲区刷新,
	// 接收到该通知后需要立即将缓冲区中的日志持久化,
	// 因为程序很有可能将在短时间内退出
	Flush()
}
```
### 适配器实现
```go
package discard

import (
	"io"
	"time"

	"github.com/bearki/belog/v2/level"
	"github.com/bearki/belog/v2/logger"
)

// Adapter 无输出日志适配器
type Adapter struct {
}

// New 创建无输出日志适配器
func New() logger.Adapter {
	return &Adapter{}
}

// Name 用于获取适配器名称
//
// 注意：请确保适配器名称不与其他适配器名称冲突
func (e *Adapter) Name() string {
	return "belog-discard-adapter"
}

// Print 普通日志打印方法
//
// @params logTime 日记记录时间
//
// @params level 日志级别
//
// @params content 日志内容
func (e *Adapter) Print(_ time.Time, _ level.Level, content []byte) {
	io.Discard.Write(content)
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
func (e *Adapter) PrintStack(_ time.Time, _ level.Level, content []byte, _ string, _ int, _ string) {
	io.Discard.Write(content)
}

// Flush 日志缓存刷新
//
// 用于日志缓冲区刷新
// 接收到该通知后需要立即将缓冲区中的日志持久化
func (e *Adapter) Flush() {}
```
### 适配器挂载
```go
package mylog

import (
	"${your_module_name}/discard"

	"github.com/bearki/belog/v2"
)

func main() {
    l, err := belog.New(logger.DefaultOption, discard.New())
	if err != nil {
		fmt.Printf("belog logger create failed, %s\r\n", err)
		return
	}
    l.Info("log init success")
}
```
