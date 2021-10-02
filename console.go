/**
 *@Title 控制台日志记录引擎
 *@Desc 控制台打印将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// consoleEngine 控制台引擎
type consoleEngine struct {
}

// consolePrintLock 控制台打印锁
var consolePrintLock sync.Mutex

// 全局控制台颜色
const (
	colorReset  = "\033[0m"  // 重置
	colorGray   = "\033[37m" // 白色
	colorBlue   = "\033[34m" // 蓝色
	colorGreen  = "\033[32m" // 绿色
	colorYellow = "\033[33m" // 黄色
	colorRed    = "\033[31m" // 红色
	colorViolet = "\033[35m" // 紫色
)

// initConsoleEngine 初始化控制台引擎
func initConsoleEngine(options interface{}) *consoleEngine {
	consolelog := new(consoleEngine)
	return consolelog
}

// printConsoleLog 打印控制台日志
func (consolelog *consoleEngine) printConsoleLog(logStr string) {
	// 加锁
	consolePrintLock.Lock()
	// 解锁
	defer consolePrintLock.Unlock()
	// 提取日志级别
	compileRegex := regexp.MustCompile(`\[(T|D|I|W|E|F){1}\]`) // 正则表达式的分组，以括号()表示，每一对括号就是我们匹配到的一个文本，可以把他们提取出来。
	levelStr := compileRegex.FindString(logStr)
	// 根据级别赋值颜色
	var levelColorStr string
	switch levelStr {
	case `[T]`: // 通知级别(灰色)
		levelColorStr = colorGray + levelStr + colorReset
	case `[D]`: // 调试级别(蓝色)
		levelColorStr = colorBlue + levelStr + colorReset
	case `[I]`: // 普通级别(绿色)
		levelColorStr = colorGreen + levelStr + colorReset
	case `[W]`: // 警告级别(黄色)
		levelColorStr = colorYellow + levelStr + colorReset
	case `[E]`: // 错误级别(红色)
		levelColorStr = colorRed + levelStr + colorReset
	case `[F]`: // 紧急级别(紫色)
		levelColorStr = colorViolet + levelStr + colorReset
	}
	// 替换需要加颜色的部分
	logStr = strings.Replace(logStr, levelStr, levelColorStr, 1)
	// 打印到控制台
	fmt.Println(logStr)
}
