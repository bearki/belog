/**
 *@Title 文件日志记录引擎
 *@Desc 文件日志的写入将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

import "sync"

// filePrintLock 控制台打印锁
var filePrintLock sync.Mutex

// printFileLog 记录日志到文件
func printFileLog(logStr string) {
	// 加锁
	filePrintLock.Lock()
	// 解锁
	defer filePrintLock.Unlock()
}
