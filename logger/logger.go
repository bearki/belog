/**
 *@Title belog核心代码
 *@Desc belog日志的主要实现都在这里了，欢迎大家指出需要改进的地方
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package logger

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strconv"
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
	// 是否禁用JSON序列化输出
	disabledJsonFormat bool

	// 时间格式化样式
	//
	// 默认: 2006/01/02 15:04:05.000
	timeFormat string

	timeJsonKey        string
	levelJsonKey       string
	stackJsonKey       string
	stackFileJsonKey   string
	stackLineNoJsonKey string
	stackMethodJsonKey string
	fieldsJsonKey      string
	messageJsonKey     string

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
	// 赋值默认的时间样式
	b.timeFormat = "2006/01/02 15:04:05.000"
	// json字段key
	b.timeJsonKey = "time"
	b.levelJsonKey = "level"
	b.stackJsonKey = "stack"
	b.stackFileJsonKey = "file"
	b.stackLineNoJsonKey = "line"
	b.stackMethodJsonKey = "method"
	b.fieldsJsonKey = "fields"
	b.messageJsonKey = "message"
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

// singleAdaptersPrint 单适配器输出
func (b *belog) singleAdaptersPrint(t time.Time, l Level, c []byte) {
	// 遍历所有适配器
	for _, adapter := range b.adapters {
		// 不打印调用栈信息
		adapter.Print(t, l, c)
	}
}

// multipleAdaptersPrint 多适配器输出
func (b *belog) multipleAdaptersPrint(t time.Time, l Level, c []byte) {
	// 协程等待分组（WaitFroup会增加1个开销）
	var wg sync.WaitGroup
	// 遍历所有适配器
	for _, adapter := range b.adapters {
		wg.Add(1)
		go func(a Adapter) {
			defer wg.Done()
			a.Print(t, l, c)
		}(adapter)
	}
	// 等待所有适配器完成日志记录
	wg.Wait()
}

// singleAdaptersStackPrint 单适配器含调用栈输出
func (b *belog) singleAdaptersStackPrint(t time.Time, l Level, c []byte, fn []byte, ln int, mn []byte) {
	// 遍历所有适配器
	for _, adapter := range b.adapters {
		adapter.PrintStack(t, l, c, fn, ln, mn)
	}
}

// multipleAdaptersStackPrint 多适配器含调用栈输出
func (b *belog) multipleAdaptersStackPrint(t time.Time, l Level, c []byte, fn []byte, ln int, mn []byte) {
	// 协程等待分组（WaitFroup会增加1个开销）
	var wg sync.WaitGroup
	// 遍历所有适配器
	for _, adapter := range b.adapters {
		wg.Add(1)
		go func(a Adapter) {
			defer wg.Done()
			a.PrintStack(t, l, c, fn, ln, mn)
		}(adapter)
	}
	// 等待所有适配器完成日志记录
	wg.Wait()
}

// format 序列化行格式日志
//
// 传入的content格式:
//
// k1: v1, k2: v2, ..., message
//
// 打印格式：
//
// 2022/09/14 20:28:13.793 [T]  k1: v1, k2: v2, ..., message\r\n
func (b *belog) format(currTime time.Time, level Level, content []byte) {
	// 2022/09/14 20:28:13.793 [T]  k1: v1, k2: v2, ..., message\r\n
	// +++++++++++++++++++++++======++++++++++++++++++++++++++++====
	//        len(tf)           6               len(c)            2
	//
	// const count = 8

	// 计算大小
	size := 8 + len(b.timeFormat) + len(content)
	// 从对象池中获取一个日志字节流对象
	logBytes := logBytesPool.Get()
	defer func() {
		// 回收切片
		logBytesPool.Put(logBytes)
	}()
	// 判读对象容量是否足够
	if cap(logBytes) < size {
		logBytes = make([]byte, 0, size)
	}

	// 开始追加内容
	logBytes = currTime.AppendFormat(logBytes, b.timeFormat)
	logBytes = append(logBytes, ` [`...)
	logBytes = append(logBytes, level.GetLevelChar())
	logBytes = append(logBytes, `]  `...)
	logBytes = append(logBytes, content...)
	logBytes = append(logBytes, "\r\n"...)

	// 是否为多适配器输出
	if len(b.adapters) > 1 {
		// 多适配器并发输出
		b.multipleAdaptersPrint(currTime, level, logBytes)
	} else {
		// 单适配器输出
		b.singleAdaptersPrint(currTime, level, logBytes)
	}
}

// formatJSON 序列化JSON格式日志
//
// 传入的content格式:
//
// "fields": {"k1": "v1", ...}, "msg": "message"
//
// 打印格式：
//
// {"$(timeKey)": "$(2006/01/02 15:04:05.000)", "$(levelKey)": "D", $(content)}\r\n
func (b *belog) formatJSON(currTime time.Time, level Level, content []byte) {
	// = 固定长度
	// + 动态长度
	//
	// {"$(timeKey)": "$(2006/01/02 15:04:05.000)", "$(levelKey)": "D", $(content)}\r\n
	// ==++++++++++====++++++++++++++++++++++++++====+++++++++++========++++++++++=====
	// 2   len(tk)   4         len(tf)             4   len(lk)     8      len(c)    3
	//
	// const count = 2 + 4 + 4 + 8 + 3 = 21

	// 计算大小
	size := 21 + len(b.timeJsonKey) + len(b.timeFormat) + len(b.levelJsonKey) + len(content)
	// 从对象池中获取一个日志字节流对象
	logBytes := logBytesPool.Get()
	defer func() {
		// 回收切片
		logBytesPool.Put(logBytes)
	}()
	// 判读对象容量是否足够
	if cap(logBytes) < size {
		logBytes = make([]byte, 0, size)
	}

	// 开始追加内容
	logBytes = append(logBytes, `{"`...)
	logBytes = append(logBytes, b.timeJsonKey...)
	logBytes = append(logBytes, `": "`...)
	logBytes = currTime.AppendFormat(logBytes, b.timeFormat)
	logBytes = append(logBytes, `", "`...)
	logBytes = append(logBytes, b.levelJsonKey...)
	logBytes = append(logBytes, `": "`...)
	logBytes = append(logBytes, level.GetLevelChar())
	logBytes = append(logBytes, `", `...)
	logBytes = append(logBytes, content...)
	logBytes = append(logBytes, "}\r\n"...)

	// 是否为多适配器输出
	if len(b.adapters) > 1 {
		// 多适配器并发输出
		b.multipleAdaptersPrint(currTime, level, logBytes)
	} else {
		// 单适配器输出
		b.singleAdaptersPrint(currTime, level, logBytes)
	}
}

// formatStack 序列化带调用栈的行格式日志
//
// 传入的content格式:
//
// k1: v1, k2: v2, ..., message
//
// 打印格式：
//
// 2022/09/14 20:28:13.793 [T] [belog_test.go:10000000000] [PrintLogTest]  k1: v1, k2: v2, ..., message\r\n
func (b *belog) formatStack(currTime time.Time, level Level, content []byte) {
	// 2022/09/14 20:28:13.793 [T] [belog_test.go:10000000000] [PrintLogTest]  k1: v1, k2: v2, ..., message\r\n
	// +++++++++++++++++++++++======+++++++++++++=+++++++++++===++++++++++++===++++++++++++++++++++++++++++====
	//        len(tf)           6      len(fn)   1  len(ln)   3    len(mn)   3              len(c)           2
	//
	// const count = 6 + 1 + 3 + 3 + 2 = 15

	// 获取调用栈信息
	pc, fileName, lineNo, _ := runtime.Caller(int(b.skip))
	// 转换文件名
	fn := convert.StringToBytes(fileName)
	// 裁剪为基础文件名
	index := bytes.LastIndexByte(fn, '/')
	if index > -1 && index+1 < len(fn) {
		fn = fn[index+1:]
	}
	// 整形行号转字节切片
	ln := make([]byte, 0, 10)
	ln = strconv.AppendInt(ln, int64(lineNo), 10)
	// 获取函数名字节切片
	var mn []byte
	if rfunc := runtime.FuncForPC(pc); rfunc != nil {
		mn = convert.StringToBytes(rfunc.Name())
		// 裁剪为基础函数名
		index = bytes.LastIndexByte(mn, '/')
		if index > 0 && index+1 < len(mn) {
			mn = mn[index+1:]
		}
	}

	// 计算大小
	size := 15 + len(b.timeFormat) + len(fn) + len(ln) + len(mn) + len(content)
	// 从对象池中获取一个日志字节流对象
	logBytes := logBytesPool.Get()
	defer func() {
		// 回收切片
		logBytesPool.Put(logBytes)
	}()
	// 判读对象容量是否足够
	if cap(logBytes) < size {
		logBytes = make([]byte, 0, size)
	}

	// 开始追加内容
	logBytes = currTime.AppendFormat(logBytes, b.timeFormat)
	logBytes = append(logBytes, ` [`...)
	logBytes = append(logBytes, level.GetLevelChar())
	logBytes = append(logBytes, `] [`...)
	logBytes = append(logBytes, fn...)
	logBytes = append(logBytes, `:`...)
	logBytes = append(logBytes, ln...)
	logBytes = append(logBytes, `] [`...)
	logBytes = append(logBytes, mn...)
	logBytes = append(logBytes, `]  `...)
	logBytes = append(logBytes, content...)
	logBytes = append(logBytes, "\r\n"...)

	// 是否为多适配器输出
	if len(b.adapters) > 1 {
		// 含调用栈多适配器并发输出
		b.multipleAdaptersStackPrint(currTime, level, logBytes, fn, lineNo, mn)
	} else {
		// 含调用栈单适配器输出
		b.singleAdaptersStackPrint(currTime, level, logBytes, fn, lineNo, mn)
	}
}

// formatStackJSON 序列化带调用栈的JSON格式日志
//
// 传入的content格式:
//
// "fields": {"k1": "v1", ...}, "msg": "message"
//
// 打印格式：
//
// {"$(timeKey)": "$(2006/01/02 15:04:05.000)", "$(levelKey)": "D", "$(stackKey)": {"$(fileKey)": "xxxxxxx", "$(lineNoKey)": 1000000, "$(methodKey)": "xxxxxxx"}, $(content)}\r\n
func (b *belog) formatStackJSON(currTime time.Time, level Level, content []byte) {
	// = 固定长度
	// + 动态长度
	//
	// {"$(timeKey)": "$(2006/01/02 15:04:05.000)", "$(levelKey)": "D", "$(stackKey)": {"$(fileKey)": "xxxxxxx", "$(lineNoKey)": 1000000, "$(methodKey)": "xxxxxxx"}, $(content)}\r\n
	// ==++++++++++====++++++++++++++++++++++++++====+++++++++++=========+++++++++++=====++++++++++====+++++++====++++++++++++===+++++++===++++++++++++====+++++++====++++++++++=====
	// 2   len(tk)   4         len(tf)             4   len(lk)      9      len(sk)    5   len(fk)   4  len(fn)  4   len(lnk)   3 len(ln) 3    len(mk)   4  len(mn)  4   len(c)    3
	//
	// const count = 2 + 4 + 4 + 9 + 5 + 4 + 4 + 3 + 3 + 4 + 4 + 3 = 49

	// 获取调用栈信息
	pc, fileName, lineNo, _ := runtime.Caller(int(b.skip))
	// 转换文件名
	fn := convert.StringToBytes(fileName)
	// 裁剪为基础文件名
	index := bytes.LastIndexByte(fn, '/')
	if index > -1 && index+1 < len(fn) {
		fn = fn[index+1:]
	}
	// 整形行号转字节切片
	ln := make([]byte, 0, 10)
	ln = strconv.AppendInt(ln, int64(lineNo), 10)
	// 获取函数名字节切片
	var mn []byte
	if rfunc := runtime.FuncForPC(pc); rfunc != nil {
		mn = convert.StringToBytes(rfunc.Name())
		// 裁剪为基础函数名
		index = bytes.LastIndexByte(mn, '/')
		if index > 0 && index+1 < len(mn) {
			mn = mn[index+1:]
		}
	}

	// 计算大小
	size := 49 + len(b.timeJsonKey) + len(b.timeFormat) + len(b.levelJsonKey) +
		len(b.stackJsonKey) + len(b.stackFileJsonKey) + len(fn) + len(b.stackLineNoJsonKey) +
		len(ln) + len(b.stackMethodJsonKey) + len(mn) + len(content)
	// 从对象池中获取一个日志字节流对象
	logBytes := logBytesPool.Get()
	defer func() {
		// 回收切片
		logBytesPool.Put(logBytes)
	}()
	// 判读对象容量是否足够
	if cap(logBytes) < size {
		logBytes = make([]byte, 0, size)
	}

	// 开始追加内容
	logBytes = append(logBytes, `{"`...)
	logBytes = append(logBytes, b.timeJsonKey...)
	logBytes = append(logBytes, `": "`...)
	logBytes = currTime.AppendFormat(logBytes, b.timeFormat)
	logBytes = append(logBytes, `", "`...)
	logBytes = append(logBytes, b.levelJsonKey...)
	logBytes = append(logBytes, `": "`...)
	logBytes = append(logBytes, level.GetLevelChar())
	logBytes = append(logBytes, `", "`...)
	logBytes = append(logBytes, b.stackJsonKey...)
	logBytes = append(logBytes, `": {"`...)
	logBytes = append(logBytes, b.stackFileJsonKey...)
	logBytes = append(logBytes, `": "`...)
	logBytes = append(logBytes, fn...)
	logBytes = append(logBytes, `", "`...)
	logBytes = append(logBytes, b.stackLineNoJsonKey...)
	logBytes = append(logBytes, `": `...)
	logBytes = append(logBytes, ln...)
	logBytes = append(logBytes, `, "`...)
	logBytes = append(logBytes, b.stackMethodJsonKey...)
	logBytes = append(logBytes, `": "`...)
	logBytes = append(logBytes, mn...)
	logBytes = append(logBytes, `"}, `...)
	logBytes = append(logBytes, content...)
	logBytes = append(logBytes, "}\r\n"...)

	// 是否为多适配器输出
	if len(b.adapters) > 1 {
		// 含调用栈多适配器并发输出
		b.multipleAdaptersStackPrint(currTime, level, logBytes, fn, lineNo, mn)
	} else {
		// 含调用栈单适配器输出
		b.singleAdaptersStackPrint(currTime, level, logBytes, fn, lineNo, mn)
	}
}

// preCheck 常规日志前置判断和序列化
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

	// 获取当前时间
	currTime := time.Now()
	// 数据容器
	var logBytes []byte

	// 是否禁用JSON格式输出
	if b.disabledJsonFormat {
		// 使用行格式输出

		// 静态字符串时开启优化打印
		if len(val) == 0 {
			logBytes = convert.StringToBytes(format)
		} else {
			logBytes = convert.StringToBytes(fmt.Sprintf(format, val...))
		}

		// 是否需要打印调用栈
		if b.printCallStack {
			// 执行含调用栈输出行打印
			b.formatStack(currTime, level, logBytes)
		} else {
			// 执行不含调用栈输出行打印
			b.format(currTime, level, logBytes)
		}

	} else {
		// 使用json格式输出

		// 静态字符串时开启优化打印
		if len(val) == 0 {
			b.preCheckf(level, format)
			return
		}

		// 是否满足(message string, val1 field.Field, val2 field.Field...)
		var fields []field.Field
		for _, v := range val {
			// 是否为Field类型
			if item, ok := v.(field.Field); ok {
				fields = append(fields, item)
			}
		}
		if len(fields) == len(val) {
			// 使用高性能格式化输出
			b.preCheckf(level, format, fields...)
			return
		}

		// 是否满足(message string, key1 string, val1 interface{}, key2 string, val2 interface{}...)
		fields = fields[:]
		key := ""
		for i, v := range val {
			if i%2 == 0 {
				if item, ok := v.(string); ok {
					key = item
				} else {
					// 一旦匹配到偶数项不是字符串就立即跳出
					break
				}
			} else {
				fields = append(fields, field.Interface(key, v))
			}
		}
		if len(fields) == len(val) {
			// 使用高性能格式化输出
			b.preCheckf(level, format, fields...)
			return
		}

		// 格式化为json数据

		// 使用标准库的JSON输出

		// 是否需要打印调用栈
		if b.printCallStack {
			// 执行含调用栈输出JSON打印
			b.formatStackJSON(currTime, level, logBytes)
		} else {
			// 执行不含调用栈输出JSON打印
			b.formatJSON(currTime, level, logBytes)
		}

	}
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

	// 获取当前时间
	currTime := time.Now()

	// 是否禁用json序列化格式输出
	if b.disabledJsonFormat {
		// 结构化为:
		//
		// = 固定长度
		// + 动态长度
		//
		// k1: v1, k2: v2, ..., this is test message
		// +++++++++++++++++++++++++++++++++++++++++
		//
		// 全是动态长度

		// 计算byte大小
		size := len(message)
		for _, v := range val {
			// key: 10,  => "": => 4个字节
			// 每个键和值之间有一个冒号和一个空格 => 2个字节
			// 每个值后面有一个逗号和一个空格 => 2个字节
			size += 4 + len(v.KeyBytes) + len(v.ValBytes)
		}

		// 从对象池中获取一个日志字节流对象
		logBytes := logBytesPool.Get()
		// 用完后字节切片放回复用池
		defer func() {
			logBytesPool.Put(logBytes)
		}()

		// 判读对象容量是否足够
		if cap(logBytes) < size {
			logBytes = make([]byte, 0, size)
		}

		// 将字段和消息拼接为行格式
		logBytes = field.Append(logBytes, message, val...)

		// 是否需要打印调用栈
		if b.printCallStack {
			// 执行含调用栈输出行打印
			b.formatStack(currTime, level, logBytes)
		} else {
			// 执行不含调用栈输出行打印
			b.format(currTime, level, logBytes)
		}

	} else {
		// 结构化为:
		//
		// = 固定长度
		// + 动态长度
		//
		// "$(fields)": {$("k1": v1, "k2": "v2", ...)}, "$(message)": "this is test message"
		// =+++++++++====++++++++++++++++++++++++++++====++++++++++====++++++++++++++++++++=
		// 1 len(fk)   4      sum(len(field)...)       4   len(mk)   4        len(m)       1
		//
		// const count = 1 + 4 + 4 + 4 + 1 = 14

		// 计算Byte大小
		size := 14 + len(b.fieldsJsonKey) + len(b.messageJsonKey) + len(message)
		for _, v := range val {
			// "key": 10 => "": => 4个字节
			// 每个键有两个双引号包裹 => 2个字节
			// 每个键和值之间有一个冒号和一个空格 => 2个字节
			size += 4 + len(v.KeyBytes) + len(v.ValPrefixBytes) + len(v.ValBytes) + len(v.ValSuffixBytes)
		}

		// 从对象池中获取一个日志字节流对象
		logBytes := logBytesPool.Get()
		// 用完后字节切片放回复用池
		defer func() {
			logBytesPool.Put(logBytes)
		}()

		// 判读对象容量是否足够
		if cap(logBytes) < size {
			logBytes = make([]byte, 0, size)
		}

		// 将字段和消息拼接为json格式
		logBytes = field.AppendJSON(logBytes, b.fieldsJsonKey, b.messageJsonKey, false, message, val...)

		// 是否需要打印调用栈
		if b.printCallStack {
			// 执行含调用栈输出JSON打印
			b.formatStackJSON(currTime, level, logBytes)
		} else {
			// 执行不含调用栈输出JSON打印
			b.formatJSON(currTime, level, logBytes)
		}

	}
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
