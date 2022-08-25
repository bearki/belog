/**
 *@Title belog核心代码
 *@Desc belog日志的主要实现都在这里了，欢迎大家指出需要改进的地方
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Logger 日志接口
type Logger interface {
	SetEngine(engine Engine, options interface{}) error // 引擎设置
	SetLevel(val ...Level) Logger                       // 日志级别设置
	OpenFileLine() Logger                               // 行号打印开启
	SetSkip(skip uint) Logger                           // 函数栈配置
	Trace(format string, val ...interface{})            // 通知级别的日志
	Debug(format string, val ...interface{})            // 调试级别的日志
	Info(format string, val ...interface{})             // 普通级别的日志
	Warn(format string, val ...interface{})             // 警告级别的日志
	Error(format string, val ...interface{})            // 错误级别的日志
	Fatal(format string, val ...interface{})            // 致命级别的日志
}

// belogPrint 引擎打印日志输出方法类型
type belogPrint func(t time.Time, lc LevelChar, file string, line int, logStr string)

// belog 记录器对象
type belog struct {
	// 是否开启文件行号记录
	isFileLine bool
	// 需要向上捕获的函数栈层数（
	// 该值会自动加2，以便于实例化用户可直接使用
	// 【示例】
	// 0：runtime.Caller函数的执行位置（在belog包内）
	// 1：belog各级别方法实现位置（在belog包内）
	// 2：belog实例调用各级别日志函数位置，依此类推】
	skip uint
	// 引擎日志输出方法，每一个实例的输出句柄将会缓存在该map中
	engine map[string]belogPrint
	// 需要记录的日志级别字符映射
	level map[Level]LevelChar
}

// New 初始化一个日志记录器实例
// 	@params engine  belogEngine 必选的基础日志引擎
// 	@params options interface{} 引擎的配置参数
// 	@return         *beLog      日志记录器实例指针
func New(engine Engine, options interface{}) (Logger, error) {
	// 初始化日志记录器对象
	b := new(belog)
	// 初始化引擎
	err := b.SetEngine(engine, options)
	if err != nil {
		return nil, err
	}
	// 默认不需要跳过函数堆栈
	b.SetSkip(0)
	// 默认开启全部级别的日志记录
	b.SetLevel(LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal)
	// 返回日志示例指针
	return b, nil
}

// SetEngine 设置日志记录引擎
// 	@params engine  Engine      引擎对象
// 	@params options interface{} 引擎参数
// 	@return         error       错误信息
func (b *belog) SetEngine(engine Engine, options interface{}) error {
	// 初始化引擎
	e, err := engine.Init(options)
	if err != nil {
		return err
	}
	// map为空需要初始化
	if b.engine == nil {
		b.engine = make(map[string]belogPrint)
	}
	// 赋值引擎
	b.engine[fmt.Sprintf("%T", e)] = e.Print
	return nil
}

// SetLevel 设置日志记录保存级别
// 	@params val ...belogLevel 任意数量的日志记录级别
func (b *belog) SetLevel(val ...Level) Logger {
	// 置空，用于覆盖后续输入的级别
	b.level = nil
	// 初始化一下
	b.level = make(map[Level]LevelChar)
	// 遍历输入的级别
	for _, item := range val {
		b.level[item] = levelMap[item]
	}
	return b
}

// OpenFileLine 开启文件行号记录
func (b *belog) OpenFileLine() Logger {
	b.isFileLine = true
	return b
}

// SetSkip 配置需要向上捕获的函数栈层数
func (b *belog) SetSkip(skip uint) Logger {
	b.skip = 2 + skip
	return b
}

// print 日志集中打印地，日志的真实记录地
func (b *belog) print(logstr string, levelChar LevelChar) {
	// 统一当前时间
	currTime := time.Now()
	// 文件行号
	var file string
	var line int
	// 是否需要打印文件行数
	if b.isFileLine {
		// 捕获函数栈文件名及执行行数
		_, file, line, _ = runtime.Caller(int(b.skip))
		// 提取文件名部分
		file = filepath.Base(file)
	}
	// 异步等待组
	var wg sync.WaitGroup
	// 遍历引擎，执行输出
	for _, out := range b.engine {
		wg.Add(1)
		go func(ouput belogPrint) {
			defer wg.Done()
			ouput(currTime, levelChar, file, line, logstr)
		}(out)
	}
	// 等待所有引擎完成日志记录
	wg.Wait()
}

// Trace 通知级别的日志
// 	@params format string         序列化格式
// 	@params val    ...interface{} 待序列化内容
func (b *belog) Trace(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelTrace]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelTrace])
}

// Debug 调试级别的日志
// 	@params format string         序列化格式
// 	@params val    ...interface{} 待序列化内容
func (b *belog) Debug(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelDebug]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelDebug])
}

// Info 普通级别的日志
// 	@params format string         序列化格式
// 	@params val    ...interface{} 待序列化内容
func (b *belog) Info(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelInfo]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelInfo])
}

// Warn 警告级别的日志
// 	@params format string         序列化格式
// 	@params val    ...interface{} 待序列化内容
func (b *belog) Warn(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelWarn]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelWarn])
}

// Error 错误级别的日志
// 	@params format string         序列化格式
// 	@params val    ...interface{} 待序列化内容
func (b *belog) Error(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelError]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelError])
}

// Fatal 致命级别的日志
// 	@params format string         序列化格式
// 	@params val    ...interface{} 待序列化内容
func (b *belog) Fatal(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelFatal]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelFatal])
}
