/**
 *@Title belog核心代码
 *@Desc belog日志的主要实现都在这里了，欢迎大家指出需要改进的地方
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package logger

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/bearki/belog/v2/pkg/tool"
)

// Logger 日志接口
type Logger interface {
	SetAdapter(Adapter) error     // 适配器设置
	SetLevel(...Level) Logger     // 日志级别设置
	PrintCallStack() Logger       // 开启调用栈打印
	SetSkip(uint) Logger          // 函数栈配置
	Trace(string, ...interface{}) // 通知级别的日志
	Debug(string, ...interface{}) // 调试级别的日志
	Info(string, ...interface{})  // 普通级别的日志
	Warn(string, ...interface{})  // 警告级别的日志
	Error(string, ...interface{}) // 错误级别的日志
	Fatal(string, ...interface{}) // 致命级别的日志
}

// belogPrint 适配器普通日志打印方法类型
type belogPrint func(time.Time, Level, []byte)

// belogPrint 适配器调用栈日志打印方法类型
type belogPrintStack func(time.Time, Level, []byte, string, int, string)

// belog 记录器对象
type belog struct {
	// 是否开启调用栈记录
	printCallStack bool

	// 需要向上捕获的函数栈层数（
	// 该值会自动加2，以便于实例化用户可直接使用
	// 【示例】
	// 0：runtime.Caller函数的执行位置（在belog包内）
	// 1：belog各级别方法实现位置（在belog包内）
	// 2：belog实例调用各级别日志函数位置，依此类推】
	skip uint

	// 需要记录的日志级别字符映射
	level map[Level]LevelChar

	// 适配器普通日志打印方法，每一个实例的输出句柄将会缓存在该map中
	adapterPrints map[string]belogPrint
	// 适配器调用栈日志打印方法，每一个实例的输出句柄将会缓存在该map中
	adapterPrintStacks map[string]belogPrintStack
}

// New 初始化一个日志记录器实例
//
// @params adapter 日志适配器
//
// @return 日志记录器实例
func New(adapter Adapter) (Logger, error) {
	// 初始化日志记录器对象
	b := new(belog)
	// 初始化适配器
	err := b.SetAdapter(adapter)
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

// SetAdapter 设置日志记录适配器
//
// @params adapter 适配器实例
//
// @return error 错误信息
func (b *belog) SetAdapter(adapter Adapter) error {
	// 适配器是否为空
	if adapter == nil {
		return errors.New("the address of `adapter` is a null pointer")
	}

	// 适配器名称不能为空
	if len(strings.TrimSpace(adapter.Name())) == 0 {
		return errors.New("the return value of `Name()` is empty")
	}

	// map为空需要初始化
	if b.adapterPrints == nil {
		b.adapterPrints = make(map[string]belogPrint)
	}
	if b.adapterPrintStacks == nil {
		b.adapterPrintStacks = make(map[string]belogPrintStack)
	}

	// 赋值适配器
	b.adapterPrints[adapter.Name()] = adapter.Print
	b.adapterPrintStacks[adapter.Name()] = adapter.PrintStack

	return nil
}

// SetLevel 设置日志记录保存级别
//
// @params val 日志记录级别（会覆盖上一次的级别配置）
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

// PrintCallStack 是否记录调用栈
//
// @return 日志记录器实例
func (b *belog) PrintCallStack() Logger {
	b.printCallStack = true
	return b
}

// SetSkip 配置需要向上捕获的函数栈层数
//
// @params skip 需要跳过的函数栈层数
//
// @return 日志记录器实例
func (b *belog) SetSkip(skip uint) Logger {
	b.skip = 2 + skip
	return b
}

// print 日志打印方法
//
// @params levelChar 日志级别
//
// @params content 日志内容
func (b *belog) print(level Level, content []byte) {
	// 获取日志记录时间
	currTime := time.Now()

	// 协程等待分组
	var wg sync.WaitGroup

	// 是否需要打印调用栈
	if b.printCallStack {
		// 捕获函数栈文件名及执行行数
		methodPtr, fileName, lineNo, _ := runtime.Caller(int(b.skip))
		// 遍历适配器，执行输出
		for _, out := range b.adapterPrintStacks {
			wg.Add(1)
			go func(ouput belogPrintStack) {
				defer wg.Done()
				ouput(currTime, level, content, fileName, lineNo, runtime.FuncForPC(methodPtr).Name())
			}(out)
		}

	} else {
		// 遍历适配器，执行输出
		for _, out := range b.adapterPrints {
			wg.Add(1)
			go func(ouput belogPrint) {
				defer wg.Done()
				ouput(currTime, level, content)
			}(out)
		}
	}

	// 等待所有适配器完成日志记录
	wg.Wait()
}

// Trace 通知级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Trace(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelTrace]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(LevelTrace, tool.StringToBytes(logStr))
}

// Debug 调试级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Debug(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelDebug]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(LevelDebug, tool.StringToBytes(logStr))
}

// Info 普通级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Info(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelInfo]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(LevelInfo, tool.StringToBytes(logStr))
}

// Warn 警告级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Warn(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelWarn]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(LevelWarn, tool.StringToBytes(logStr))
}

// Error 错误级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Error(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelError]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(LevelError, tool.StringToBytes(logStr))
}

// Fatal 致命级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Fatal(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[LevelFatal]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	b.print(LevelFatal, tool.StringToBytes(logStr))
}
