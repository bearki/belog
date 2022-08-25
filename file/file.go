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
	"log"
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
	logPath       string      // 日志文件保存路径
	maxSize       uint64      // 单文件最大容量（单位：byte）
	saveDay       uint16      // 日志保存天数
	logPathFormat string      // 日志文件路径格式
	currLogPath   string      // 日志文件源文件路径
	currTime      time.Time   // 当前日志文件使用的日期
	currIndex     uint32      // 当前文件分割后缀标识
	fileWriteChan chan string // 日志写入缓冲管道
}

// 创建一个引擎
func New() *Engine {
	return new(Engine)
}

// validity 判断参数有效性
func (p *Options) validity() error {
	// 转换路径为当前系统格式
	p.LogPath = filepath.Join(p.LogPath)
	// 判断路径是否为空
	if len(p.LogPath) <= 0 {
		return fmt.Errorf("`logPath` is null")
	}
	// 判断日志保存天数
	if p.SaveDay < 1 || p.SaveDay > 1000 {
		return fmt.Errorf("file log save day min value is 1,max value is 1000")
	}
	// 判断日志大小限制
	if p.MaxSize < 1 || p.MaxSize > 10000 {
		return fmt.Errorf("file log size min value is 1(MB),max value is 10000(MB)")
	}
	// 判断异步模式下的管道容量
	if p.Async {
		if p.AsyncChanCap > 99999 {
			return fmt.Errorf("async log channel cap error, min 1, max 99999")
		}
	} else {
		// 非异步情况下管道容量为1
		p.AsyncChanCap = 1
	}
	return nil
}

// initFileEngine 初始化文件引擎
func (e *Engine) Init(options interface{}) (logger.Engine, error) {
	// 类型断言参数
	data, ok := options.(Options)
	if !ok {
		return nil, fmt.Errorf("file log optionsis nil")
	}

	// 判断参数有效性
	if err := data.validity(); err != nil {
		return nil, err
	}

	// 创建路径的文件夹部分
	err := os.MkdirAll(filepath.Dir(data.LogPath), 0755)
	if err != nil {
		return nil, err
	}

	// 实例化文件引擎
	e = new(Engine)
	// 赋值软链文件名
	e.logPath = data.LogPath
	// 分割文件夹与文件名部分
	logDir, logFile := filepath.Split(e.logPath)
	// 截取文件名后缀
	logExt := filepath.Ext(logFile)
	// 日志文件名（不含后缀）
	logName := strings.TrimSuffix(logFile, logExt)
	// 赋值日志文件路径生成格式
	e.logPathFormat = filepath.Join(logDir, logName+".%s.%03d"+logExt)
	// 初始化写入管道容量
	e.fileWriteChan = make(chan string, data.AsyncChanCap)
	// 赋值文件大小限制
	MB := 1024 * 1024
	e.maxSize = uint64(MB) * uint64(data.MaxSize)
	// 赋值日志保存天数
	e.saveDay = data.SaveDay
	// 赋值当前日期
	e.currTime = time.Now()
	// 筛选出合适的下标日志文件
	if err = e.selectAvailableFile(); err != nil {
		return nil, err
	}

	// 异步监听过期文件删除
	go e.deleteTimeoutLogFile()
	// 创建一个文件分割通信管道
	splitFile := make(chan struct{}, 1)
	// 异步监听文件切割
	go e.listenLogFileSplit(splitFile)
	// 异步处理日志文件写入
	go func() {
		for {
			// 阻塞开启监听写入日志
			e.writeFile(splitFile)
		}
	}()

	// 返回文件日志实例
	return e, nil
}

// selectAvailableFile 选择一个可用的文件
func (e *Engine) selectAvailableFile() error {
	// 循环取文件名(单日最多允许存在999个日志文件)
	for i := 1; i <= 999; i++ {
		// 赋值文件分割后缀标识
		e.currIndex = uint32(i)
		// 以当前日期命名文件名
		e.currLogPath = fmt.Sprintf(e.logPathFormat, e.currTime.Format("2006-01-02"), i)
		// 判断拼接出来的文件状态
		file, err := os.Stat(e.currLogPath)
		if err != nil {
			if os.IsNotExist(err) {
				// 文件不存在，则该文件可用，跳出循环
				break
			}
			// 文件存在，但获取信息错误
			return err
		}
		// 判断文件大小是否超过限制
		if file.Size() >= int64(e.maxSize) {
			// 超过了限制，递增后缀标识
			continue
		}
		// 文件可用
		break
	}
	return nil
}

// listenLogFileSplit 监听日志文件分割
func (e *Engine) listenLogFileSplit(splitFile chan<- struct{}) {
	// 定义一个10秒间隔的定时器
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	// 死循环监听吧
	for range ticker.C {
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
			// 通知执行文件分割
			splitFile <- struct{}{}
			// 开启下一个监听
			continue
		}

		// 获取文件信息
		fileinfo, err := os.Stat(e.currLogPath)
		if err != nil {
			// 开启下一个监听
			continue
		}

		// 判断容量是否超过了
		if fileinfo.Size() >= int64(e.maxSize) {
			// 后缀标识加1
			e.currIndex++
			// 拼接后缀标识及文件后缀
			e.currLogPath = fmt.Sprintf(e.logPathFormat, e.currTime.Format("2006-01-02"), e.currIndex)
			// 通知执行文件分割
			splitFile <- struct{}{}
			// 开启下一个监听
			continue
		}
	}
}

// writeFile 写入日志到文件中
func (e *Engine) writeFile(splitFile <-chan struct{}) {
	// 移除软链
	err := os.RemoveAll(e.logPath)
	if err != nil {
		log.Fatalln("remove file error: " + err.Error())
	}

	// 创建或追加文件
	fileObj, err := os.OpenFile(e.currLogPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		log.Fatalln("open file error: " + err.Error())
	}
	// 结束时关闭文件
	defer fileObj.Close()

	// 创建软链
	err = os.Link(e.currLogPath, e.logPath)
	if err != nil {
		log.Fatalln("create file link error: " + err.Error())
	}

	// 监听文件写入或重新打开新的日志文件
	for {
		select {
		case logStr := <-e.fileWriteChan: // 写日志
			// 写入到文件中
			_, err = fileObj.WriteString(logStr)
			if err != nil {
				log.Println(err.Error())
			}
		case <-splitFile: // 需要分割新文件了
			return
		}
	}
}

// listenLogFileDelete 监听日志文件删除
func (e *Engine) deleteTimeoutLogFile() {
	// 定义一个1分钟间隔的定时器
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	// 死循环监听吧
	for range ticker.C {
		// 获取日志储存文件夹部分
		logDirPath := filepath.Dir(e.logPath)
		// 打开文件夹
		logDir, err := ioutil.ReadDir(logDirPath)
		if err != nil {
			log.Println(err.Error())
			return
		}
		// 获取当天整点时间
		currDateStr := time.Now().Format("2006-01-02")
		// 再解析成时间类型
		currDate, err := time.Parse("2006-01-02", currDateStr)
		if err != nil {
			log.Println(err.Error())
			return
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
					err = os.Remove(filepath.Join(logDirPath, item.Name()))
					if err != nil {
						log.Println(err.Error())
					}
				}
			}
		}
	}
}

// printFileLog 记录日志到文件
func (e *Engine) Print(t time.Time, lc logger.LevelChar, file string, line int, logStr string) {
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
}
