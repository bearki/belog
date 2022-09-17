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

	"github.com/bearki/belog/v2/internal/convert"
	"github.com/bearki/belog/v2/internal/pool"
	"github.com/bearki/belog/v2/logger/field"
)

// EncoderFunc 格式化编码器类型
type EncoderFunc func(string, ...interface{}) string

// 默认编码器
var defaultEncoder EncoderFunc = fmt.Sprintf

// 日志字节流对象池
var logBytesPool = pool.NewBytesPool(100, 0, 1024)

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

	// 适配器缓存映射
	adapters map[string]Adapter

	// 格式化编码器
	//
	// 默认的格式化编码器为: fmt.Sprintf()
	encoder EncoderFunc
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
	if b.adapters == nil {
		b.adapters = make(map[string]Adapter)
	}

	// 赋值适配器操作方法
	b.adapters[adapter.Name()] = adapter

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

// SetEncoder 设置日志格式化编码器
func (b *belog) SetEncoder(encoder EncoderFunc) Logger {
	b.encoder = encoder
	return b
}

// PrintCallStack 是否记录调用栈
//
// 注意：开启调用栈打印将会损失部分性能
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
	b.skip = 3 + skip
	return b
}

// Flush 日志缓存刷新
//
// 用于日志缓冲区刷新，
// 建议在程序正常退出时调用一次日志刷新，
// 以保证日志能完整的持久化
func (b *belog) Flush() {
	// 协程等待组
	var wg sync.WaitGroup
	// 遍历适配器
	for _, adapter := range b.adapters {
		wg.Add(1)
		go func(a Adapter) {
			defer wg.Done()
			a.Flush()
		}(adapter)
	}
	// 等待所有协程结束
	wg.Wait()
}

// print 日志打印方法
//
// @params levelChar 日志级别
//
// @params content 日志内容
func (b *belog) print(level Level, content []byte) {
	// 获取日志记录时间
	currTime := time.Now()

	// 协程等待分组（WaitFroup会增加1个开销）
	// var wg sync.WaitGroup

	// 等待所有适配器完成日志记录
	// wg.Wait()

	// return

	// 是否需要打印调用栈
	if b.printCallStack {
		// 捕获函数栈文件名及执行行数
		methodPtr, fileName, lineNo, _ := runtime.Caller(int(b.skip))
		// 遍历适配器，执行输出
		for _, adapter := range b.adapters {
			adapter.PrintStack(
				currTime,
				level,
				content,
				convert.StringToBytes(fileName),
				lineNo,
				convert.StringToBytes(runtime.FuncForPC(methodPtr).Name()),
			)
			// （异步会增加开销）
			// wg.Add(1)
			// go func(a Adapter) {
			// 	defer wg.Done()
			// a.PrintStack(
			// 	currTime,
			// 	level,
			// 	content,
			// 	convert.StringToBytes(fileName),
			// 	lineNo,
			// 	convert.StringToBytes(runtime.FuncForPC(methodPtr).Name()),
			// )
			// }(adapter)
		}
	} else {
		// 遍历适配器，执行输出
		for _, adapter := range b.adapters {
			adapter.Print(
				currTime,
				level,
				content,
			)
			// （异步会增加开销）
			// wg.Add(1)
			// go func(a Adapter) {
			// 	defer wg.Done()
			// 	a.Print(
			// 		currTime,
			// 		level,
			// 		content,
			// 	)
			// }(adapter)
		}
	}

	// // 等待所有适配器完成日志记录
	// wg.Wait()
}

// preCheck 常规日志前置判断
//
// @params level 日志级别
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) preCheck(level Level, format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[level]; !ok {
		// 当前级别日志不需要记录
		return
	}
	var content []byte
	// 是否使用默认编码器
	if b.encoder == nil {
		// 仅静态字符串时开启优化打印
		if len(val) == 0 {
			content = convert.StringToBytes(format)
		} else {
			// 使用默认编码器
			content = convert.StringToBytes(defaultEncoder(format, val...))
		}
	} else {
		// 使用自定义编码器
		content = convert.StringToBytes(b.encoder(format, val...))
	}
	b.print(LevelTrace, content)
}

// Trace 通知级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Trace(format string, val ...interface{}) {
	b.preCheck(LevelTrace, format, val...)
}

// Debug 调试级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Debug(format string, val ...interface{}) {
	b.preCheck(LevelDebug, format, val...)
}

// Info 普通级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Info(format string, val ...interface{}) {
	b.preCheck(LevelInfo, format, val...)
}

// Warn 警告级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Warn(format string, val ...interface{}) {
	b.preCheck(LevelWarn, format, val...)
}

// Error 错误级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Error(format string, val ...interface{}) {
	b.preCheck(LevelError, format, val...)
}

// Fatal 致命级别的日志
//
// @params format 序列化格式或内容
//
// @params val 待序列化内容
func (b *belog) Fatal(format string, val ...interface{}) {
	b.preCheck(LevelFatal, format, val...)
}

// preCheck 高性能日志前置判断和序列化
//
// @params level 日志级别
//
// @params message 日志消息
//
// @params val 字段信息
func (b *belog) preCheckf(level Level, message string, val ...field.Field) {
	// 判断当前级别日志是否需要记录
	if _, ok := b.level[level]; !ok {
		// 当前级别日志不需要记录
		return
	}

	// 结构化为：
	//
	// {"fields": {"k1": v1, "k2": "v2", ...}, "message": "this is test message"}
	// ++++++++++++=========================+++++++++++++++====================++
	//      12         sum(len(field)...)          15          len(message)     2
	//
	// const count = 12 + 15 + 2 = 29

	// 计算Byte大小
	// size := 1024
	size := 29 + len(message) // 除去字段的大小
	for _, v := range val {
		size += len(v.KeyBytes) + len(v.ValPrefixBytes) + len(v.ValSuffixBytes)
	}

	// 从对象池中获取一个日志字节流对象
	logBytes := logBytesPool.Get()
	// 判读对象容量是否足够
	if cap(logBytes) < size {
		logBytes = make([]byte, 0, size)
	}

	// 追加json开始括号
	logBytes = append(logBytes, '{')

	// 追加字段集字段
	logBytes = append(logBytes, `"fields": `...)
	// 追加字段开始括号
	logBytes = append(logBytes, '{')
	// 是否需要追加分隔符了
	appendDelimiter := false
	// 遍历所有字段
	for _, v := range val {
		// 从第二个有效字段开始追加分隔符号
		if appendDelimiter {
			logBytes = append(logBytes, ',', ' ')
		}

		// 追加字段并序列化
		logBytes = append(logBytes, '"')
		logBytes = append(logBytes, v.KeyBytes...)
		logBytes = append(logBytes, '"')
		logBytes = append(logBytes, ':', ' ')
		logBytes = append(logBytes, v.ValPrefixBytes...)
		logBytes = append(logBytes, v.ValBytes...)
		logBytes = append(logBytes, v.ValSuffixBytes...)

		// 已经填充了一个有效字段了
		if !appendDelimiter {
			// 下一次需要追加分隔符
			appendDelimiter = true
		}

		// 回收到复用池
		v.Put()
	}
	// 追加字段结束括号
	logBytes = append(logBytes, '}')

	// 追加分隔符
	logBytes = append(logBytes, ',', ' ')

	// 追加message字段及其内容
	logBytes = append(logBytes, `"message": `...)
	logBytes = append(logBytes, '"')
	logBytes = append(logBytes, convert.StringToBytes(message)...)
	logBytes = append(logBytes, '"')

	// 追加json结尾括号
	logBytes = append(logBytes, '}')

	// 执行打印
	b.print(level, logBytes)

	// 字节流对象放回对象池
	logBytesPool.Put(logBytes)
}

// Tracef 通知级别的日志（高性能序列化）
func (b *belog) Tracef(message string, val ...field.Field) {
	b.preCheckf(LevelTrace, message, val...)
}

// Debugf 调试级别的日志（高性能序列化）
func (b *belog) Debugf(message string, val ...field.Field) {
	b.preCheckf(LevelDebug, message, val...)
}

// Infof 普通级别的日志（高性能序列化）
func (b *belog) Infof(message string, val ...field.Field) {
	b.preCheckf(LevelInfo, message, val...)
}

// Warnf 警告级别的日志（高性能序列化）
func (b *belog) Warnf(message string, val ...field.Field) {
	b.preCheckf(LevelWarn, message, val...)
}

// Errorf 错误级别的日志（高性能序列化）
func (b *belog) Errorf(message string, val ...field.Field) {
	b.preCheckf(LevelError, message, val...)
}

// Fatalf 致命级别的日志（高性能序列化）
func (b *belog) Fatalf(message string, val ...field.Field) {
	b.preCheckf(LevelFatal, message, val...)
}
