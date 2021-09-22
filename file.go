/**
 *@Title 文件日志记录引擎
 *@Desc 文件日志的写入将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

import (
	"sync"
)

// fileEngine 文件引擎
type fileEngine struct {
	logPath    string // 日志文件保存路径
	isSplitDay bool   // 是否开启按日分割
	maxSize    uint16 // 单文件最大容量（单位：byte）
	saveDay    uint16 // 日志保存天数
}

// filePrintLock 控制台打印锁
var filePrintLock sync.Mutex

// initFileEngine 初始化文件引擎
func initFileEngine(option interface{}) *fileEngine {
	// 实例化文件引擎
	filelog := new(fileEngine)
	// // 判断路径是否为空
	// if len(logPath) <= 0 {
	// 	panic("beLog: `logPath` is null")
	// }
	// // 创建路径的文件夹部分
	// err := os.MkdirAll(filepath.Dir(logPath), 0755)
	// if err != nil {
	// 	panic(fmt.Sprintf("beLog: create dir error, %s", err.Error()))
	// }
	// // 判断是文件路径还是文件夹路径
	// filePath, err := os.Stat(logPath)
	// if err != nil {
	// 	// 判断文件是否存在
	// 	if os.IsNotExist(err) {
	// 		// 创建文件
	// 		os.Create("")
	// 	}
	// 	panic(fmt.Sprintf("beLog: %s", err.Error()))
	// }
	// if filePath.IsDir() {
	// 	panic(fmt.Sprintf("beLog: %s is dir, `logPath` it should be a file", logPath))
	// }
	return filelog
}

// OpenSplitDay 开启日志文件按日分割
func (filelog *fileEngine) OpenSplitDay() *fileEngine {
	filelog.isSplitDay = true
	return filelog
}

// SetMaxSize 配置单文件储存容量
// @params maxSize uint16 单文件最大容量（单位：MB）
func (filelog *fileEngine) SetMaxSize(maxSize uint16) *fileEngine {
	byteSize := 1024 * 1024
	filelog.maxSize = uint16(byteSize) * maxSize
	return filelog
}

// SetSaveDay 配置日志保存天数
// @params saveDay uint16 保存天数
func (filelog *fileEngine) SetSaveDay(saveDay uint16) *fileEngine {
	filelog.saveDay = saveDay
	return filelog
}

// printFileLog 记录日志到文件
func (filelog *fileEngine) printFileLog(logStr string) {
	// 加锁
	filePrintLock.Lock()
	// 解锁
	defer filePrintLock.Unlock()
}
