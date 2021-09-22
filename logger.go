/**
 *@Title belog核心代码
 *@Desc belog日志的主要实现都在这里了，欢迎大家指出需要改进的地方
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Logger 日志接口
type Logger interface {
	Trace(format string, val ...interface{}) // 通知级别的日志
	Debug(format string, val ...interface{}) // 调试级别的日志
	Info(format string, val ...interface{})  // 普通级别的日志
	Warn(format string, val ...interface{})  // 警告级别的日志
	Error(format string, val ...interface{}) // 错误级别的日志
	Fatal(format string, val ...interface{}) // 致命级别的日志
}

// 日志引擎类型
type belogEngine uint8

// 日志级别类型
type belogLevel uint8

// 日志级别字符类型
type belogLevelChar byte

// 日志输出引擎方法类型
type printFuncEngine func(logStr string)

// BeLog 记录器对象
type BeLog struct {
	isFileLine bool                            // 是否开启文件行号记录
	skip       uint                            // 需要向上捕获的函数栈层数（该值会自动加2，以便于实例化用户可直接使用）【0-runtime.Caller函数的执行位置（在belog包内），1-Belog各级别方法实现位置（在belog包内），2-belog实例调用各级别日志函数位置，依次类推】
	engine     map[belogEngine]printFuncEngine // 输出引擎方法类型映射
	level      map[belogLevel]belogLevelChar   // 需要记录的日志级别字符映射
}

// 记录引擎定义
var (
	EngineConsole belogEngine = 1 // 控制台引擎
	EngineFile    belogEngine = 2 // 文件引擎
)

// 日志保存级别定义
var (
	LevelTrace belogLevel = 1 // 通知级别
	LevelDebug belogLevel = 2 // 调试级别
	LevelInfo  belogLevel = 3 // 普通级别
	LevelWarn  belogLevel = 4 // 警告级别
	LevelError belogLevel = 5 // 错误级别
	LevelFatal belogLevel = 6 // 致命级别
)

// 日志级别映射
var levelMap = map[belogLevel]belogLevelChar{
	1: 'T',
	2: 'D',
	3: 'I',
	4: 'W',
	5: 'E',
	6: 'F',
}

// New 初始化一个日志记录器实例
// @params engine belogEngine 必选的基础日志引擎
// @params option interface{} 引擎的配置参数
// @return        *BeLog      日志记录器实例指针
// @return        error       初始化时发生的错误信息
func New(engine belogEngine, option interface{}) *BeLog {
	// 初始化日志记录器对象
	belog := new(BeLog)
	// 初始化引擎
	belog.SetEngine(engine, option)
	// 默认不需要跳过函数堆栈
	belog.SetSkip(0)
	// 默认开启全部级别的日志记录
	belog.SetLevel(LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal)
	// 返回日志示例指针
	return belog
}

// SetEngine 设置日志记录引擎
// @params val ...belogEngine 任意数量的日志记录引擎
func (belog *BeLog) SetEngine(engine belogEngine, option interface{}) *BeLog {
	// 判断引擎是否为空
	if belog.engine == nil {
		// 初始化一下
		belog.engine = make(map[belogEngine]printFuncEngine)
	}
	// 使用不同引擎进行初始化
	switch engine {
	case EngineFile: // 文件引擎
		// 初始化文件引擎
		fileEngine := initFileEngine(option)
		// 赋值引擎输出方法映射
		belog.engine[engine] = fileEngine.printFileLog
	/*--------------------------可以在此之后添加更多的引擎---------------------------*/

	/*--------------------------可以在此之前添加更多的引擎---------------------------*/
	case EngineConsole: // 控制台引擎（和默认引擎一致）
		fallthrough
	default: // 默认引擎（和控制台引擎一致）
		// 初始化控制台引擎
		consolelog := initConsoleEngine(option)
		// 赋值引擎输出方法映射
		belog.engine[engine] = consolelog.printConsoleLog
	}
	return belog
}

// SetLevel 设置日志记录保存级别
// @params val ...belogLevel 任意数量的日志记录级别
func (belog *BeLog) SetLevel(val ...belogLevel) *BeLog {
	// 置空，用于覆盖后续输入的级别
	belog.level = nil
	// 初始化一下
	belog.level = make(map[belogLevel]belogLevelChar)
	// 遍历输入的级别
	for _, item := range val {
		belog.level[item] = levelMap[item]
	}
	return belog
}

// OpenFileLine 开启文件行号记录
func (belog *BeLog) OpenFileLine() *BeLog {
	belog.isFileLine = true
	return belog
}

// SetSkip 配置需要向上捕获的函数栈层数
func (belog *BeLog) SetSkip(skip uint) *BeLog {
	belog.skip = 2 + skip
	return belog
}

// print 日志集中打印地，日志的真实记录地
func (belog *BeLog) print(logstr string, levelChar belogLevelChar) {
	// 统一当前时间
	currTime := time.Now()
	// 是否需要打印文件行数
	if belog.isFileLine {
		// 捕获函数栈文件名及执行行数
		_, file, line, _ := runtime.Caller(int(belog.skip))
		// 格式化
		logstr = fmt.Sprintf(
			"%s.%03d [%s] [%s:%d]  %s\n",
			currTime.Format("2006/01/02 15:04:05"),
			(currTime.UnixNano()/1e6)%currTime.Unix(),
			string(levelChar),
			filepath.Base(file),
			line,
			logstr,
		)
	} else {
		// 格式化
		logstr = fmt.Sprintf(
			"%s.%03d [%s]  %s\n",
			currTime.Format("2006/01/02 15:04:05"),
			(currTime.UnixNano()/1e6)%currTime.Unix(),
			string(levelChar),
			logstr,
		)
	}
	// 异步等待组
	var wg sync.WaitGroup
	// 遍历引擎，执行输出
	for _, printFunc := range belog.engine {
		wg.Add(1)
		go func(ouput printFuncEngine) {
			defer wg.Done()
			ouput(logstr)
		}(printFunc)
	}
	// 等待所有引擎完成日志记录
	wg.Wait()
}

// Trace 通知级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (belog *BeLog) Trace(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := belog.level[LevelTrace]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	belog.print(logStr, belog.level[LevelTrace])
}

// Debug 调试级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (belog *BeLog) Debug(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := belog.level[LevelDebug]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	belog.print(logStr, belog.level[LevelDebug])
}

// Info 普通级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (belog *BeLog) Info(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := belog.level[LevelInfo]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	belog.print(logStr, belog.level[LevelInfo])
}

// Warn 警告级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (belog *BeLog) Warn(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := belog.level[LevelWarn]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	belog.print(logStr, belog.level[LevelWarn])
}

// Error 错误级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (belog *BeLog) Error(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := belog.level[LevelError]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	belog.print(logStr, belog.level[LevelError])
}

// Fatal 致命级别的日志
// @params format string         序列化格式
// @params val    ...interface{} 待序列化内容
func (belog *BeLog) Fatal(format string, val ...interface{}) {
	// 判断当前级别日志是否需要记录
	if _, ok := belog.level[LevelFatal]; !ok { // 当前级别日志不需要记录
		return
	}
	// 执行日志记录
	logStr := fmt.Sprintf(format, val...)
	belog.print(logStr, belog.level[LevelFatal])
}
