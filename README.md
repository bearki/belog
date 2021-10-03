# 关于
这是一个高度解耦的日志框架，支持多引擎同时输出，你可以发挥自己的想象力，随意的创建自己喜爱的引擎，在 `v0.1.0` 版本中，我们已经提供了两个简单的引擎实现 `console` 、 `file` ，你会注意到它们都实现了 `logger` 中的 `Engine` 接口，该接口主要有两个方法实现，只要实现了这两个方法的引擎均可挂载到 `logger` 中，在已实现的两个引擎中均有示例，你可以查看源码来编写自己的引擎。

# 安装
```shell
go get github.com/bearki/belog
```

# 快速使用
如果你想快速体验belog的特性，我们已内置了一个默认 `console` 实现，你可以直接使用以下方式输出你的日志内容
```go
package main

import (
	"github.com/bearki/belog"
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
	"github.com/bearki/belog"
	"github.com/bearki/belog/file"
	"github.com/bearki/belog/logger"
)

// 全局日志对象
var Log logger.Logger

func main() {
	// 初始化一个实例(可实例化任意引擎)
	mylog, err := belog.New(
		new(file.Engine), // 初始化文件引擎
		file.Options{
			LogPath:      "./logs/app.log", // 日志储存路径
			MaxSize:      128,              // 日志单文件大小
			SaveDay:      7,                // 日志保存天数
			Async:        true,             // 开启异步写入
			AsyncChanCap: 20,               // 异步缓存管道容量
		},
	)
	if err != nil { // 初始化失败将不能执行任何后续操作，否则会引起恐慌
		panic("belog init error: " + err.Error())
	}

	// 配置日志记录级别
	mylog.SetLevel(
		logger.LevelTrace,
		logger.LevelDebug,
		logger.LevelInfo,
		logger.LevelWarn,
		logger.LevelError,
		logger.LevelFatal,
	).OpenFileLine() // 开启行号记录

	/**********支持多引擎同时输出**********/
	/*************************************/
	// err = mylog.SetEngine( // 增加控制台引擎
	// 	new(console.Engine),
	// 	nil,
	// )
	// if err != nil { // 该错误不影响已添加的引擎
	// 	mylog.Error("add console engine to log error: %s", err.Error())
	// }
	/*************************************/

	// 实例对象记录日志
	mylog.Trace("this is a trace log")
	mylog.Debug("this is a debug log")
	mylog.Info("this is a info log")
	mylog.Warn("this is a warn log")
	mylog.Error("this is a error log")
	mylog.Fatal("this is a fatal log")
}
```

# 二次封装
你可以参考 `belog.go` 的方式将 `belog` 进行二次封装，这样通过自己的包即可记录日志，防止在过多的包中引入第三方包，便于后期的管理，值得注意的是，二次封装时需要配置函数栈层数，否则将会造成文件名及行数捕获不一致的问题，大多数情况下采用如下层级即可
```go
belogDefault, _ = belog.New(new(console.Engine), nil) // 初始化引擎
belogDefault.SetSkip(1)                         // 配置函数栈
```

# 自定义引擎
没有过多的繁琐操作，实现 `logger` 中的 `Engine` 接口即可完成引擎自定义工作，最简实现方式在 `console` 引擎中，仅包含这两个实现函数，我们认为这是非常容易识别的
### 引擎接口
```go
// Engine 引擎接口
type Engine interface {
	Init(options interface{}) (Engine, error)
	Print(t time.Time, lc BeLevelChar, file string, line int, logStr string)
}
```
### 引擎实现
```go
package console

import (
	"fmt"
	"sync"
	"time"

	"github.com/bearki/belog/logger"
)

// Engine 控制台引擎
type Engine struct {
	mutex sync.Mutex // 控制台输出锁
}

// Init 初始化控制台引擎
func (e *Engine) Init(options interface{}) (logger.Engine, error) {
	e = new(Engine)
	return e, nil
}

// printConsoleLog 打印控制台日志
func (e *Engine) Print(t time.Time, lc logger.BeLevelChar, file string, line int, logStr string) {
	// 加锁
	e.mutex.Lock()
	// 解锁
	defer e.mutex.Unlock()
	// 判断是否需要文件行号
	if len(file) > 0 {
		// 格式化打印
		fmt.Printf(
			"%s.%03d [%s] [%s:%d]  %s\n",
			t.Format("2006/01/02 15:04:05"),
			(t.UnixNano()/1e6)%t.Unix(),
			string(lc),
			file,
			line,
			logStr,
		)
	} else {
		// 格式化打印
		fmt.Printf(
			"%s.%03d [%s]  %s\n",
			t.Format("2006/01/02 15:04:05"),
			(t.UnixNano()/1e6)%t.Unix(),
			string(lc),
			logStr,
		)
	}
}
```
### 引擎挂载
```go
package mylog

import (
	"your_module_name/console"
)

func main() {
    logs, err = belog.New(new(console.Engine), nil)
    if err != nil {
        panic(err.Error())
    }
    logs.Info("log init success")
}
```
