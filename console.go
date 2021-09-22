/**
 *@Title 控制台日志记录引擎
 *@Desc 控制台打印将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

import (
	"fmt"
	"sync"
)

// consoleEngine 控制台引擎
type consoleEngine struct {
}

// consolePrintLock 控制台打印锁
var consolePrintLock sync.Mutex

// initConsoleEngine 初始化控制台引擎
func initConsoleEngine(optin interface{}) *consoleEngine {
	consolelog := new(consoleEngine)
	return consolelog
}

// printConsoleLog 打印控制台日志
func (consolelog *consoleEngine) printConsoleLog(logStr string) {
	// 加锁
	consolePrintLock.Lock()
	// 解锁
	defer consolePrintLock.Unlock()
	// 打印到控制台
	fmt.Println(logStr)
}
