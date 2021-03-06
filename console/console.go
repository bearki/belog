/**
 *@Title 控制台日志记录引擎
 *@Desc 控制台打印将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package console

import (
	"fmt"
	"sync"
	"time"

	"github.com/bearki/belog/logger"
	"github.com/gookit/color"
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
	// 根据级别赋值颜色
	var levelStr string
	switch lc {
	case 'T': // 通知级别(灰色)
		levelStr = color.Gray.Sprint("[T]")
	case 'D': // 调试级别(蓝色)
		levelStr = color.Blue.Sprint("[D]")
	case 'I': // 普通级别(绿色)
		levelStr = color.Green.Sprint("[I]")
	case 'W': // 警告级别(黄色)
		levelStr = color.Yellow.Sprint("[W]")
	case 'E': // 错误级别(红色)
		levelStr = color.Red.Sprint("[E]")
	case 'F': // 紧急级别(洋红色)
		levelStr = color.Magenta.Sprint("[F]")
	}
	// 判断是否需要文件行号
	if len(file) > 0 {
		// 格式化打印
		fmt.Printf(
			"%s.%03d %s [%s:%d]  %s\n",
			t.Format("2006/01/02 15:04:05"),
			(t.UnixNano()/1e6)%t.Unix(),
			levelStr,
			file,
			line,
			logStr,
		)
	} else {
		// 格式化打印
		fmt.Printf(
			"%s.%03d %s  %s\n",
			t.Format("2006/01/02 15:04:05"),
			(t.UnixNano()/1e6)%t.Unix(),
			levelStr,
			logStr,
		)
	}
}
