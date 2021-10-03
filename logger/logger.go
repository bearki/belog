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

// BeLevel 日志级别类型
type BeLevel uint8

// BeLevelChar 日志界别字符类型
type BeLevelChar byte

// Engine 引擎接口
type Engine interface {
	Init(options interface{}) (Engine, error)                                // 引擎初始化函数
	Print(t time.Time, lc BeLevelChar, file string, line int, logStr string) // 日志打印函数
}

// Logger 日志接口
type Logger interface {
	SetEngine(engine Engine, options interface{}) error // 引擎设置
	SetLevel(val ...BeLevel) Logger                     // 日志级别设置
	OpenFileLine() Logger                               // 行号打印开启
	SetSkip(skip uint) Logger                           // 函数栈配置
	Trace(format string, val ...interface{})            // 通知级别的日志
	Debug(format string, val ...interface{})            // 调试级别的日志
	Info(format string, val ...interface{})             // 普通级别的日志
	Warn(format string, val ...interface{})             // 警告级别的日志
	Error(format string, val ...interface{})            // 错误级别的日志
	Fatal(format string, val ...interface{})            // 致命级别的日志
}

// beEnginePrint 引擎打印日志输出方法类型
type beEnginePrint func(t time.Time, lc BeLevelChar, file string, line int, logStr string)

// beLog 记录器对象
type beLog struct {
	isFileLine bool                    // 是否开启文件行号记录
	skip       uint                    // 需要向上捕获的函数栈层数（该值会自动加2，以便于实例化用户可直接使用）【0-runtime.Caller函数的执行位置（在belog包内），1-Belog各级别方法实现位置（在belog包内），2-belog实例调用各级别日志函数位置，依次类推】
	engine     []beEnginePrint         // 引擎日志输出方法
	level      map[BeLevel]BeLevelChar // 需要记录的日志级别字符映射
}

// 日志保存级别定义
var (
	LevelTrace BeLevel = 1 // 通知级别
	LevelDebug BeLevel = 2 // 调试级别
	LevelInfo  BeLevel = 3 // 普通级别
	LevelWarn  BeLevel = 4 // 警告级别
	LevelError BeLevel = 5 // 错误级别
	LevelFatal BeLevel = 6 // 致命级别
)

// levelMap 日志级别字符映射
var levelMap = map[BeLevel]BeLevelChar{
	1: 'T',
	2: 'D',
	3: 'I',
	4: 'W',
	5: 'E',
	6: 'F',
}

// New 初始化一个日志记录器实例
// @params engine  belogEngine 必选的基础日志引擎
// @params options interface{} 引擎的配置参数
// @return         *beLog      日志记录器实例指针
func New(engine Engine, options interface{}) (Logger, error) {
	// 初始化日志记录器对象
	b := new(beLog)
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
func (b *beLog) SetEngine(engine Engine, options interface{}) error {
	// 初始化引擎
	e, err := engine.Init(options)
	if err != nil {
		return err
	}
	// 追加到引擎列表
	b.engine = append(b.engine, e.Print)
	return nil
}

// SetLevel 设置日志记录保存级别
// @params val ...belogLevel 任意数量的日志记录级别
func (b *beLog) SetLevel(val ...BeLevel) Logger {
	// 置空，用于覆盖后续输入的级别
	b.level = nil
	// 初始化一下
	b.level = make(map[BeLevel]BeLevelChar)
	// 遍历输入的级别
	for _, item := range val {
		b.level[item] = levelMap[item]
	}
	return b
}

// OpenFileLine 开启文件行号记录
func (b *beLog) OpenFileLine() Logger {
	b.isFileLine = true
	return b
}

// SetSkip 配置需要向上捕获的函数栈层数
func (b *beLog) SetSkip(skip uint) Logger {
	b.skip = 2 + skip
	return b
}

// Trace 通知级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (b *beLog) Trace(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelTrace]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelTrace])
}

// Debug 调试级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (b *beLog) Debug(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelDebug]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelDebug])
}

// Info 普通级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (b *beLog) Info(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelInfo]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelInfo])
}

// Warn 警告级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (b *beLog) Warn(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelWarn]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelWarn])
}

// Error 错误级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (b *beLog) Error(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelError]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelError])
}

// Fatal 致命级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (b *beLog) Fatal(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelFatal]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(logStr, b.level[LevelFatal])
}

// print 日志集中打印地，日志的真实记录地
func (b *beLog) print(logstr string, levelChar BeLevelChar) {
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
	for _, printFunc := range b.engine {
		wg.Add(1)
		go func(ouput beEnginePrint) {
			defer wg.Done()
			ouput(currTime, levelChar, file, line, logstr)
		}(printFunc)
	}
	// 等待所有引擎完成日志记录
	wg.Wait()
}
