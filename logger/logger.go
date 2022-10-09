/**
 *@Title belog核心代码
 *@Desc belog日志的主要实现都在这里了，欢迎大家指出需要改进的地方
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package logger

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/bearki/belog/v2/internal/pool"
	"github.com/bearki/belog/v2/level"
)

var (
	// 日志字节流对象池
	logBytesPool = pool.NewBytesPool(100, 0, 1024)
)

// belog 标准记录器
type belog struct {
	//
	// 缓存映射配置
	//

	levelMapRWMutex sync.RWMutex             // 日志级别配置读写锁
	levelMap        map[level.Level]struct{} // 需要记录的日志级别字符映射
	adaptersRWMutex sync.RWMutex             // 适配器配置读写锁
	adapters        map[string]Adapter       // 适配器缓存映射

	//
	// 功能配置
	//

	stackSkip          uint // 需要跳过的调用栈层数
	printCallStack     bool // 是否打印调用栈
	callStackFullPath  bool // 是否输出调用栈完整路径
	disabledJsonFormat bool // 是否禁用JSON序列化输出

	//
	// 序列化格式配置
	//

	timeFormat  string // 时间序列化格式
	levelFormat bool   // 日志级别输出格式

	//
	// JSON字段配置
	//

	timeJsonKey        string // 时间的JSON键名
	levelJsonKey       string // 日志级别的JSON键名
	stackJsonKey       string // 调用栈信息JSON键名
	stackFileJsonKey   string // 调用栈文件名JSON键名
	stackLineNoJsonKey string // 调用栈文件行号JSON键名
	stackMethodJsonKey string // 调用栈函数名JSON键名
	fieldsJsonKey      string // 字段集JSON键名
	messageJsonKey     string // 日志消息JSON键名
}

// New 初始化一个日志记录器实例
//
// @params adapter 日志适配器
//
// @return 日志记录器实例
func New(option Option, adapter ...Adapter) (Logger, error) {
	// 获取有效参数
	option = getValidOption(option)

	// 初始化日志记录器对象
	bl := &belog{
		stackSkip:          stackBaseSkip, // 初始为默认最小跳过层数
		levelMap:           nil,
		printCallStack:     option.PrintCallStack || option.CallStackFullPath,
		callStackFullPath:  option.CallStackFullPath,
		disabledJsonFormat: option.DisabledJsonFormat,
		timeFormat:         option.TimeFormat,
		levelFormat:        option.LevelFormat,
		timeJsonKey:        option.TimeJsonKey,
		levelJsonKey:       option.LevelJsonKey,
		stackJsonKey:       option.StackJsonKey,
		stackFileJsonKey:   option.StackFileJsonKey,
		stackLineNoJsonKey: option.StackLineNoJsonKey,
		stackMethodJsonKey: option.StackMethodJsonKey,
		fieldsJsonKey:      option.FieldsJsonKey,
		messageJsonKey:     option.MessageJsonKey,
	}

	// 默认开启全部级别的日志记录
	bl.SetLevel(
		level.Trace, level.Debug, level.Info,
		level.Warn, level.Error, level.Fatal,
	)

	// 预初始化（time包首次Format很慢）
	_ = time.Now().Format(bl.timeFormat)

	// 初始化适配器
	for _, v := range adapter {
		err := bl.SetAdapter(v)
		if err != nil {
			return nil, err
		}
	}

	// 返回标准记录器
	return &StandardBelog{
		belog: bl,
	}, nil
}

// SetAdapter 设置日志记录适配器
//
// @params adapter 适配器实例
//
// @return error 错误信息
func (b *belog) SetAdapter(adapter Adapter) error {
	// 适配器是否为空
	if adapter == nil {
		return errors.New("the address of `adapter` is adapter null pointer")
	}

	// 适配器名称不能为空
	if len(strings.TrimSpace(adapter.Name())) == 0 {
		return errors.New("the return value of `Name()` is empty")
	}

	// 加个写锁
	b.adaptersRWMutex.Lock()
	defer b.adaptersRWMutex.Unlock()

	// map为空需要初始化
	if b.adapters == nil {
		b.adapters = make(map[string]Adapter)
	}

	// 赋值适配器操作方法
	b.adapters[adapter.Name()] = adapter

	return nil
}

// levelIsExist 判断日志级别是否在缓存中
func (b *belog) levelIsExist(l level.Level) bool {
	// 加个读锁
	b.levelMapRWMutex.RLock()
	// 检查
	_, ok := b.levelMap[l]
	// 释放读锁
	b.levelMapRWMutex.RUnlock()
	// 返回结果
	return ok
}

// SetLevel 设置日志记录保存级别
//
// @params val 日志记录级别（会覆盖上一次的级别配置）
func (b *belog) SetLevel(ls ...level.Level) {
	// 加个写锁
	b.levelMapRWMutex.Lock()

	// 置空，用于覆盖后续输入的级别
	b.levelMap = nil
	// 初始化一下
	b.levelMap = make(map[level.Level]struct{})

	// 遍历输入的级别
	for _, l := range ls {
		b.levelMap[l] = struct{}{}
	}

	// 释放写锁
	b.levelMapRWMutex.Unlock()
}

// SetSkip 配置需要向上捕获的函数栈层数
//
// @params skip 需要跳过的函数栈层数
//
// @return 日志记录器实例
func (b *belog) SetSkip(skip uint) {
	b.stackSkip = stackBaseSkip + skip
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

// adapterPrintFunc 适配器打印方法
type adapterPrintFunc func(t time.Time, l level.Level, c []byte, fn string, ln int, mn string)

// filterAdapterPrint 筛选合适的适配器
func (b *belog) filterAdapterPrint() adapterPrintFunc {
	// 是否为单适配器输出
	if len(b.adapters) <= 1 {
		// 单适配器输出
		return b.singleAdapterPrint
	}
	// 多适配器并发输出
	return b.multipleAdapterPrint
}

// singleAdapterPrint 单适配器输出
func (b *belog) singleAdapterPrint(t time.Time, l level.Level, c []byte, fn string, ln int, mn string) {
	// 加个读锁
	b.adaptersRWMutex.RLock()

	// 是否需要输出调用栈
	if b.printCallStack {
		// 遍历所有适配器
		for _, adapter := range b.adapters {
			adapter.PrintStack(t, l, c, fn, ln, mn)
		}
	} else {
		// 遍历所有适配器
		for _, adapter := range b.adapters {
			adapter.Print(t, l, c)
		}
	}

	// 释放读锁
	b.adaptersRWMutex.RUnlock()
}

// multipleAdapterPrint 多适配器输出
func (b *belog) multipleAdapterPrint(t time.Time, l level.Level, c []byte, fn string, ln int, mn string) {
	// 加个读锁
	b.adaptersRWMutex.RLock()

	// 协程等待分组（WaitGroup会增加1个开销）
	var wg sync.WaitGroup

	// 是否需要输出调用栈
	if b.printCallStack {
		// 遍历所有适配器
		for _, adapter := range b.adapters {
			wg.Add(1)
			go func(a Adapter) {
				defer wg.Done()
				a.PrintStack(t, l, c, fn, ln, mn)
			}(adapter)
		}
	} else {
		// 遍历所有适配器
		for _, adapter := range b.adapters {
			wg.Add(1)
			go func(a Adapter) {
				defer wg.Done()
				a.Print(t, l, c)
			}(adapter)
		}
	}

	// 等待所有适配器完成日志记录
	wg.Wait()

	// 释放读锁
	b.adaptersRWMutex.RUnlock()
}
