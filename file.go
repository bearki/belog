/**
 *@Title 文件日志记录引擎
 *@Desc 文件日志的写入将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileEngineOption 文件引擎参数
type FileEngineOption struct {
	LogPath    string // 日志文件储存路径
	IsSplitDay bool   // 是否开启按日分割
	MaxSize    uint16 // 单文件最大容量（单位：MB）
	SaveDay    uint16 // 日志保存天数
}

// fileEngine 文件引擎
type fileEngine struct {
	logPath     string // 日志文件保存路径（软链）
	currLogPath string // 日志文件源文件路径
	isSplitDay  bool   // 是否开启按日分割
	maxSize     uint16 // 单文件最大容量（单位：byte）
	saveDay     uint16 // 日志保存天数
}

// fileWriteChan 日志写入缓冲管道
var fileWriteChan = make(chan string, 20)

// initFileEngine 初始化文件引擎
func initFileEngine(option interface{}) *fileEngine {
	// 类型断言参数
	data, ok := option.(FileEngineOption)
	if !ok {
		panic("BeLog: file log option is nil")
	}
	// 判断路径是否为空
	if len(data.LogPath) <= 0 {
		panic("BeLog: `logPath` is null")
	}
	// 创建路径的文件夹部分
	err := os.MkdirAll(filepath.Dir(data.LogPath), 0755)
	if err != nil {
		panic(fmt.Sprintf("BeLog: create dir error, %s", err.Error()))
	}
	// 实例化文件引擎
	filelog := new(fileEngine)

	// 赋值软链文件名
	filelog.logPath = data.LogPath
	// 默认日志文件路径和软链一致
	filelog.currLogPath = data.LogPath
	// 赋值文件大小限制
	filelog.setMaxSize(data.MaxSize)
	// 赋值日志保存天数
	filelog.setSaveDay(data.SaveDay)

	// 截取文件名后缀
	logExt := filepath.Ext(data.LogPath)
	// 判断是否开启按日分割（需要开启软链）
	if data.IsSplitDay {
		// 开启按日分割
		filelog.openSplitDay()
		// 以当前日期命名文件名
		filelog.currLogPath = data.LogPath[:len(logExt)]
	}
	// 循环取文件名(最多允许存在999个日志文件)
	for i := 1; i <= 999; i++ {
		currPath := fmt.Sprintf("%s.%s.%03d%s", filelog.currLogPath, time.Now().Format("2006-01-02"), i, logExt)
		file, err := os.Stat(currPath)
		if err != nil {
			if os.IsNotExist(err) {
				filelog.currLogPath = currPath
				break
			}
			panic("BeLog: get file error: %s")
		}
		// 判断文件大小是否超过限制
		if file.Size() >= int64(filelog.maxSize) {
			continue
		}
		// 赋值当前文件的路径
		filelog.currLogPath = currPath
		break
	}
	// 开启监听写入日志
	go filelog.writeFile()
	// 返回文件日志实例
	return filelog
}

// openSplitDay 开启日志文件按日分割
func (filelog *fileEngine) openSplitDay() *fileEngine {
	filelog.isSplitDay = true
	return filelog
}

// setMaxSize 配置单文件储存容量
// @params maxSize uint16 单文件最大容量（单位：MB）
func (filelog *fileEngine) setMaxSize(maxSize uint16) *fileEngine {
	byteSize := 1024 * 1024
	filelog.maxSize = uint16(byteSize) * maxSize
	return filelog
}

// setSaveDay 配置日志保存天数
// @params saveDay uint16 保存天数
func (filelog *fileEngine) setSaveDay(saveDay uint16) *fileEngine {
	filelog.saveDay = saveDay
	return filelog
}

// printFileLog 记录日志到文件
func (filelog *fileEngine) printFileLog(logStr string) {
	// 写入到管道中
	fileWriteChan <- logStr
}

func (filelog *fileEngine) writeFile() {
	// 创建或追加文件
	file, err := os.OpenFile(filelog.currLogPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		panic("BeLog: open file error: " + err.Error())
	}
	// 结束时关闭文件
	defer file.Close()
	// 创建软链
	err = os.Symlink(filelog.currLogPath, filelog.logPath)
	if err != nil {
		panic("BeLog: create file link error: " + err.Error())
	}
	// 监听文件写入
	for logStr := range fileWriteChan {
		// 写入到文件中
		_, _ = file.WriteString(logStr)
	}
}
