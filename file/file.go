/**
 *@Title 文件日志记录引擎
 *@Desc 文件日志的写入将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bearki/belog/logger"
)

// Options 文件引擎参数
type Options struct {
	LogPath      string // 日志文件储存路径
	MaxSize      uint16 // 单文件最大容量（单位：MB）
	SaveDay      uint16 // 日志保存天数
	Async        bool   // 是否异步
	AsyncChanCap uint   // 异步通道容量
}

// Engine 文件引擎
type Engine struct {
	logPath       string        // 日志文件保存路径
	maxSize       uint64        // 单文件最大容量（单位：byte）
	saveDay       uint16        // 日志保存天数
	logPathFormat string        // 日志文件路径格式
	currLogPath   string        // 日志文件源文件路径
	currTime      time.Time     // 当前日志文件使用的日期
	currIndex     uint32        // 当前文件分割后缀标识
	currSize      uint64        // 当前日志文件大小
	async         bool          // 是否异步写入日志
	fileWriteChan chan string   // 日志写入缓冲管道
	isSplitFile   chan struct{} // 监听是否需要分割文件
}

// initFileEngine 初始化文件引擎
func (e *Engine) Init(options interface{}) (logger.Engine, error) {
	// 类型断言参数
	data, ok := options.(Options)
	if !ok {
		return nil, fmt.Errorf("file log optionsis nil")
	}
	// 判断路径是否为空
	if len(data.LogPath) <= 0 {
		return nil, fmt.Errorf("`logPath` is null")
	}
	// 转换路径为当前系统格式
	data.LogPath = filepath.Join(data.LogPath)
	// 分割文件夹与文件名部分
	logDir, logFile := filepath.Split(data.LogPath)
	// 截取文件名后缀
	logExt := filepath.Ext(logFile)
	// 日志文件名（不含后缀）
	logName := strings.TrimSuffix(logFile, logExt)
	// 创建路径的文件夹部分
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return nil, err
	}

	// 实例化文件引擎
	e = new(Engine)
	// 初始化是否需要分割文件监听通道
	e.isSplitFile = make(chan struct{}, 1)
	// 赋值软链文件名
	e.logPath = data.LogPath
	// 赋值日志文件路径生成格式
	e.logPathFormat = filepath.Join(logDir, logName+".%s.%03d"+logExt)
	// 赋值是否异步写入
	e.async = data.Async
	if e.async {
		// 开启异步写入，判断容量有效性
		if data.AsyncChanCap > 99999 {
			return nil, fmt.Errorf("async log channel cap error, min 1, max 99999")
		}
		// 容量有效，初始化管道
		e.fileWriteChan = make(chan string, data.AsyncChanCap)
	} else {
		// 不开启异步写入，默认缓冲容量为1
		e.fileWriteChan = make(chan string, 1)
	}
	// 赋值文件大小限制
	if data.MaxSize < 1 || data.MaxSize > 10000 {
		return nil, fmt.Errorf("file log size min value is 1(MB),max value is 10000(MB)")
	}
	MB := 1024 * 1024
	e.maxSize = uint64(MB) * uint64(data.MaxSize)
	// 赋值日志保存天数
	if data.SaveDay < 1 || data.SaveDay > 1000 {
		return nil, fmt.Errorf("file log save day min value is 1,max value is 1000")
	}
	e.saveDay = data.SaveDay
	// 赋值当前日期
	e.currTime = time.Now()

	// 循环取文件名(单日最多允许存在999个日志文件)
	for i := 1; i <= 999; i++ {
		// 赋值文件分割后缀标识
		e.currIndex = uint32(i)
		// 以当前日期命名文件名
		e.currLogPath = fmt.Sprintf(e.logPathFormat, e.currTime.Format("2006-01-02"), i)
		// 判断拼接出来的文件状态
		file, err := os.Stat(e.currLogPath)
		if err != nil {
			if os.IsNotExist(err) { // 文件不存在，则该文件可用，跳出循环
				break
			}
			// 文件存在，但获取信息错误
			return nil, err
		}
		// 判断文件大小是否超过限制
		if file.Size() >= int64(e.maxSize) { // 超过了限制，递增后缀标识
			continue
		}
		// 文件可用
		break
	}
	// 开启监听写入日志
	go e.writeFile()
	// 异步监听文件切割
	go e.listenLogFileSplit()
	// 异步监听文件删除
	go e.listenLogFileDelete()
	// 返回文件日志实例
	return e, nil
}

// printFileLog 记录日志到文件
func (e *Engine) Print(t time.Time, lc logger.BeLevelChar, file string, line int, logStr string) {
	// 判断是否需要文件行号
	if len(file) > 0 {
		// 格式化打印
		logStr = fmt.Sprintf(
			"%s.%03d [%s] [%s:%d]  %s\r\n",
			t.Format("2006/01/02 15:04:05"),
			(t.UnixNano()/1e6)%t.Unix(),
			string(lc),
			file,
			line,
			logStr,
		)
	} else {
		// 格式化打印
		logStr = fmt.Sprintf(
			"%s.%03d [%s]  %s\r\n",
			t.Format("2006/01/02 15:04:05"),
			(t.UnixNano()/1e6)%t.Unix(),
			string(lc),
			logStr,
		)
	}
	// 写入到管道中
	e.fileWriteChan <- logStr
	// 判断是否开启异步写入
	if !e.async { // 阻塞，直到写入完成
		for len(e.fileWriteChan) > 0 {
		}
	}
}

// listenLogFileSplit 监听日志文件分割
func (e *Engine) listenLogFileSplit() {
	// 定义一个10秒间隔的定时器
	ticker := time.Tick(time.Second * 10)
	// 死循环监听吧
	for range ticker {
		// 获取当前时间
		currTime := time.Now()
		// 比对一下当前日志文件的日期和当前日期是否是同一天
		if currTime.Day() != e.currTime.Day() { // 不是同一天
			// 赋值新日期
			e.currTime = currTime
			// 赋值新后缀标识
			e.currIndex = 1
			// 拼接后缀标识及文件后缀
			e.currLogPath = fmt.Sprintf(e.logPathFormat, e.currTime.Format("2006-01-02"), e.currIndex)
			// 通知文件写入函数需要切换文件了，可以结束了
			e.isSplitFile <- struct{}{}
			// 当切割管道未被释放时不允许重新调用文件写入函数
			for len(e.isSplitFile) > 0 {
			}
			// 异步重新打开文件写入函数
			go e.writeFile()
			// 开始下一个循环监听
			continue
		}

		// 获取文件信息
		fileinfo, err := os.Stat(e.currLogPath)
		if err != nil {
			// 开始下一个循环监听
			continue
		}
		// 判断容量是否超过了
		if e.currSize >= e.maxSize || fileinfo.Size() >= int64(e.maxSize) {
			// 后缀标识加1
			e.currIndex++
			// 拼接后缀标识及文件后缀
			e.currLogPath = fmt.Sprintf(e.logPathFormat, e.currTime.Format("2006-01-02"), e.currIndex)
			// 通知文件写入函数需要切换文件了，可以结束了
			e.isSplitFile <- struct{}{}
			// 当切割管道未被释放时不允许重新调用文件写入函数
			for {
				if len(e.isSplitFile) == 0 {
					break
				}
			}
			// 异步重新打开文件写入函数
			go e.writeFile()
		}
	}
}

// listenLogFileDelete 监听日志文件删除
func (e *Engine) listenLogFileDelete() {
	// 定义一个1分钟间隔的定时器
	ticker := time.Tick(time.Minute)
	// 死循环监听吧
	for range ticker {
		// 获取日志储存文件夹部分
		logDirPath := filepath.Dir(e.logPath)
		// 打开文件夹
		logDir, err := ioutil.ReadDir(logDirPath)
		if err != nil {
			continue
		}
		// 获取当天整点时间
		currDateStr := time.Now().Format("2006-01-02")
		// 再解析成时间类型
		currDate, err := time.Parse("2006-01-02", currDateStr)
		if err != nil {
			continue
		}
		// 初始化正则
		re := regexp.MustCompile(`[0-9]{4}-[0-9]{2}-[0-9]{2}`)
		// 遍历文件夹
		for _, item := range logDir {
			// 当前路径不是文件夹并且文件名不是当前正在使用的日志文件名
			if !item.IsDir() && filepath.Base(e.logPath) != item.Name() {
				// 获取文件名中的时间部分
				fileDateStr := re.FindString(item.Name())
				// 解析成时间类型
				fileDate, err := time.Parse("2006-01-02", fileDateStr)
				if err != nil {
					continue
				}
				// 比对两个时间是否大于指定的保存天数
				if currDate.Sub(fileDate).Hours() >= float64(24*e.saveDay) {
					// 删除这个文件
					_ = os.Remove(filepath.Join(logDirPath, item.Name()))
				}
			}
		}
	}
}

// writeFile 写入日志到文件中
func (e *Engine) writeFile() {
	// 移除软链
	// err = os.Remove(e.logPath)
	// 当路径不存在时RemoveAll return nil
	err := os.RemoveAll(e.logPath)
	if err != nil {
		panic("remove file error: " + err.Error())
	}
	// 创建或追加文件
	fileObj, err := os.OpenFile(e.currLogPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		panic("open file error: " + err.Error())
	}
	// 结束时关闭文件
	defer fileObj.Close()
	// 创建软链
	err = os.Link(e.currLogPath, e.logPath)
	if err != nil {
		panic("create file link error: " + err.Error())
	}
	// 监听文件写入或重新打开新的日志文件
	for {
		select {
		case logStr := <-e.fileWriteChan: // 写日志
			// 写入到文件中
			_, _ = fileObj.WriteString(logStr)
			// 获取文件当前大小
			fileinfo, err := fileObj.Stat()
			if err == nil {
				// 赋值当前文件大小
				e.currSize = uint64(fileinfo.Size())
			}
		case <-e.isSplitFile: // 需要分割新文件了
			// 结束掉该函数，分割监听函数将会重新打开该函数
			return
		}
	}
}
