/**
 *@Title belog核心代码
 *@Desc belog日志的主要实现都在这里了，欢迎大家指出需要改进的地方
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package logger

import (
	"errors"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/bearki/belog/v3/pkg/pool"
)

// 需要跳过的最少调用栈层数
//
//	该值由belog内部自定义，外部无需关心
const stackBaseSkip uint = 4

// 日志字节流对象池
var logBytesPool = pool.NewBytesPool(100, 0, 1024)

// 标准记录器
type belog struct {
	minLevel          Level              // 需要记录的最小日志级别
	adaptersRWMutex   sync.RWMutex       // 适配器配置读写锁
	adapters          map[string]Adapter // 适配器缓存映射
	encoder           Encoder            // 编码器
	stackSkip         uint               // 需要跳过的调用栈层数
	enabledStackPrint bool               // 是否打印调用栈
}

// 获取调用栈信息
//
//	@param	skip	需要跳过的调用栈数量
//	@return	文件名字节切片
//	@return	行号
//	@return	函数名
func getCallStack(skip uint) (fn string, ln int, mn string) {
	// 获取调用栈信息
	pc, fn, ln, _ := runtime.Caller(int(skip))

	// 获取函数名字节切片
	if funcForPC := runtime.FuncForPC(pc); funcForPC != nil {
		mn = funcForPC.Name()
	}

	// OK
	return
}

// New 初始化一个日志记录器实例
//
//	@param	adapter	日志适配器
//	@return	日志记录器实例
func New(option Option, adapter ...Adapter) (Logger, error) {
	// 检查参数
	if option.Encoder == nil {
		return nil, errors.New("the log encoder field cannot be empty")
	}

	// 初始化日志记录器对象
	bl := &belog{
		minLevel:          Trace, // 默认最低级别
		encoder:           option.Encoder,
		stackSkip:         stackBaseSkip, // 初始为默认最小跳过层数
		enabledStackPrint: option.EnabledStackPrint,
	}

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
//	@param	adapter	适配器实例
//	@return	异常信息
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

// SetLevel 设置日志记录保存级别
//
//	@param	level	日志最小记录级别
func (b *belog) SetLevel(level Level) {
	// 赋值最小记录级别
	b.minLevel = level
}

// SetSkip 配置需要向上捕获的函数栈层数
//
//	@param	skip	需要跳过的函数栈层数
func (b *belog) SetSkip(skip uint) {
	b.stackSkip = stackBaseSkip + skip
}

// Flush 日志缓存刷新
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

// 筛选合适的适配器
func (b *belog) adapterPrint(t time.Time, l Level, c []byte) {
	// 是否为单适配器输出
	if len(b.adapters) <= 1 {
		// 单适配器输出
		b.singleAdapterPrint(t, l, c)
		return
	}
	// 多适配器并发输出
	b.multipleAdapterPrint(t, l, c)
}

// 单适配器输出
func (b *belog) singleAdapterPrint(t time.Time, l Level, c []byte) {
	// 加个读锁
	b.adaptersRWMutex.RLock()

	// 适配器是否为空
	if b.adapters == nil {
		return
	}

	// 遍历所有适配器
	for _, adapter := range b.adapters {
		adapter.Print(t, l, c)
	}

	// 释放读锁
	b.adaptersRWMutex.RUnlock()
}

// 多适配器输出
func (b *belog) multipleAdapterPrint(t time.Time, l Level, c []byte) {
	// 加个读锁
	b.adaptersRWMutex.RLock()

	// 适配器是否为空
	if b.adapters == nil {
		return
	}

	// 协程等待分组（WaitGroup会增加1个开销）
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

	// 释放读锁
	b.adaptersRWMutex.RUnlock()
}

// 筛选合适的调用栈适配器
func (b *belog) adapterPrintStack(t time.Time, l Level, c []byte, fn string, ln int, mn string) {
	// 是否为单适配器输出
	if len(b.adapters) <= 1 {
		// 单适配器输出
		b.singleAdapterPrintStack(t, l, c, fn, ln, mn)
		return
	}
	// 多适配器并发输出
	b.multipleAdapterPrintStack(t, l, c, fn, ln, mn)
}

// 单适配器输出调用栈
func (b *belog) singleAdapterPrintStack(t time.Time, l Level, c []byte, fn string, ln int, mn string) {
	// 加个读锁
	b.adaptersRWMutex.RLock()

	// 适配器是否为空
	if b.adapters == nil {
		return
	}

	// 遍历所有适配器
	for _, adapter := range b.adapters {
		adapter.PrintStack(t, l, c, fn, ln, mn)
	}

	// 释放读锁
	b.adaptersRWMutex.RUnlock()
}

// 多适配器输出调用栈
func (b *belog) multipleAdapterPrintStack(t time.Time, l Level, c []byte, fn string, ln int, mn string) {
	// 加个读锁
	b.adaptersRWMutex.RLock()

	// 适配器是否为空
	if b.adapters == nil {
		return
	}

	// 协程等待分组（WaitGroup会增加1个开销）
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

	// 释放读锁
	b.adaptersRWMutex.RUnlock()
}
