/**
 *@Title 文件日志记录适配器
 *@Desc 文件日志的写入将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package file

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bearki/belog/v2/logger"
	"github.com/bearki/belog/v2/pkg/tool"
)

// Options 文件日志适配器参数
type Options struct {

	// 日志文件储存路径
	//
	// Default: ${work_dir}/app.log
	LogPath string

	// 单文件最大容量
	//
	// Unit:MB, Default:4, Min:1, Max:10000
	MaxSize uint16

	// 单文件最大行数
	//
	// Unit:行, Default:100000, Min:100, Max:100000000
	MaxLines uint64

	// 日志保存天数
	//
	// Unit: 天, Default: 30
	SaveDay uint16

	// 注意：
	//
	// 该选项可用于选择是否开启日志文件异步写入
	//
	// 我建议您开启异步写入而不是使用同步写入，
	// 同步写入相较于异步写入将会损失50%左右的性能
	//
	// 你可以使用Async=true && AsyncChanCap=1的组合来实现与同步写入相似的功能，
	// 这种组合的方式能使性能得到保证
	//
	// Default: false
	Async bool

	// 异步写入管道容量
	//
	// Default: 1, Min: 1, Max: 100
	AsyncChanCap uint
}

// Adapter 文件日志适配器
type Adapter struct {

	// 外部传入字段

	logPath        string      // 日志文件保存路径（默认：app.log）
	maxSize        uint64      // 单文件最大容量（单位：byte, 默认：4MB）
	maxLines       uint64      // 单文件最大保存行数（默认：10万行）
	saveDay        uint16      // 日志最大保存天数（默认：30天）
	fileWriteAsync bool        // 日志写入是否为异步（默认：false）
	fileWriteChan  chan []byte // 日志写入缓冲管道（默认：1）

	// 内部字段

	logPathFormat string    // 日志文件路径格式
	currLogPath   string    // 日志文件源文件路径
	currTime      time.Time // 当前日志文件使用的日期
	currIndex     uint32    // 当前日志文件分割后缀标识
	currSize      uint64    // 当前日志文件大小（单位：byte）
	currLines     uint64    // 当前日志文件行数
}

// validity 判断参数有效性
func (p *Options) validity() {
	// 转换路径为当前系统格式
	p.LogPath = filepath.Join(p.LogPath)
	// 判断路径是否为空
	if len(p.LogPath) <= 0 {
		p.LogPath = "app.log"
		log.Println("`logPath` is null, use the default value `app.log`")
	}
	// 判断日志大小限制
	if p.MaxSize < 1 || p.MaxSize > 10000 {
		p.MaxSize = 4
		log.Println("file log size min value is 1(MB),max value is 10000(MB), use the default value 4(MB)")
	}
	// 判断日志行数限制
	if p.MaxLines < 100 || p.MaxLines > 100000000 {
		p.MaxLines = 100000
		log.Println("file log lines min value is 100,max value is 100000000, use the default value 100000")
	}
	// 判断日志保存天数
	if p.SaveDay < 1 || p.SaveDay > 1000 {
		p.SaveDay = 30
		log.Println("file log save day min value is 1,max value is 1000, use the default value 30")
	}
	// 判断异步模式下的管道容量
	if p.Async {
		if p.AsyncChanCap > 100 {
			p.AsyncChanCap = 1
			log.Println("async log channel cap error, min 1, max 100, use the default value 1")
		}
	} else {
		// 非异步情况下管道容量为1
		p.AsyncChanCap = 1
	}
}

// New 创建文件日志适配器
//
// @params options 文件日志适配器参数
//
// @return 文件日志适配器实例
//
// @return 错误信息
func New(options Options) (logger.Adapter, error) {
	// 判断参数有效性
	options.validity()
	// 创建路径的文件夹部分
	err := os.MkdirAll(filepath.Dir(options.LogPath), 0755)
	if err != nil {
		return nil, err
	}

	// 预处理一些变量
	// 分割文件夹与文件名部分
	logDir, logFile := filepath.Split(options.LogPath)
	// 截取文件名后缀
	logExt := filepath.Ext(logFile)
	// 日志文件名（不含后缀）
	logName := strings.TrimSuffix(logFile, logExt)
	// 定义MB的字节大小
	MB := uint64(1024 * 1024)

	// 实例化文件日志适配器
	e := new(Adapter)
	// 赋值软链文件名
	e.logPath = options.LogPath
	// 赋值文件大小限制
	e.maxSize = uint64(options.MaxSize) * MB
	// 赋值文件最大行数
	e.maxLines = options.MaxLines
	// 赋值日志保存天数
	e.saveDay = options.SaveDay
	// 赋值日志文件路径生成格式
	e.logPathFormat = filepath.Join(logDir, logName+".%s.%03d"+logExt)
	// 赋值是否为异步写入
	e.fileWriteAsync = options.Async
	// 初始化日志写入管道容量
	e.fileWriteChan = make(chan []byte, options.AsyncChanCap)
	// 筛选出合适的下标日志文件
	if err = e.selectAvailableFile(); err != nil {
		return nil, err
	}

	// 异步执行一次过期日志文件删除
	go e.deleteTimeoutLogFile()
	// 异步死循环监听文件写入
	go func() {
		for {
			// 阻塞开启监听写入日志
			e.writeFile()
		}
	}()

	// 返回文件日志实例
	return e, nil
}

// Name 用于获取适配器名称
//
// 注意：请确保适配器名称不与其他适配器名称冲突
func (e *Adapter) Name() string {
	return "belog-file-adapter"
}

// Print 普通日志打印方法
//
// @params logTime 日记记录时间
//
// @params level 日志级别
//
// @params content 日志内容
func (e *Adapter) Print(logTime time.Time, level logger.Level, content []byte) {
	// 不带调用栈：
	// 2022/09/14 20:28:13.793 [T]  this is a trace log
	// 日期(10) + 空格(1) + 时间(12) + 空格(1) + 级别(3) + 空格(2) + 日志内容(len(content)) + 回车换行(2)
	//
	// 计算需要的大小
	size := 32 + len(content)
	// 创建一个指定容量的切片，避免二次扩容
	logSlice := make([]byte, 0, size)
	// 追加格式化好的日期和时间
	logSlice = append(logSlice, tool.StringToBytes(logTime.Format("2006/01/02 15:04:05.000"))...) // 23个字节
	// 追加级别
	logSlice = append(logSlice, ' ', '[', level.GetLevelChar(), ']') // 4个字节
	// 追加日志内容
	logSlice = append(logSlice, ' ', ' ')   // 2个字节
	logSlice = append(logSlice, content...) // len(content)个字节
	// 追加回车换行
	logSlice = append(logSlice, '\r', '\n') // 2个字节

	// 写入到管道中
	e.fileWriteChan <- logSlice
	// 同步时需要等待日志缓冲区被清空
	if !e.fileWriteAsync {
		// 空会在写入时进行判断
		e.fileWriteChan <- nil
	}
}

// PrintStack 调用栈日志打印方法
//
// @params logTime 日记记录时间
//
// @params level 日志级别
//
// @params content 日志内容
//
// @params fileName 日志记录调用文件路径
//
// @params lineNo 日志记录调用文件行号
//
// @params methodName 日志记录调用函数名
func (e *Adapter) PrintStack(logTime time.Time, level logger.Level, content []byte, fileName string, lineNo int, methodName string) {
	// 带调用栈：
	// 2022/09/14 20:28:13.793 [T] [belog_test.go:82] [PrintLog]  this is a trace log
	// 日期(10) + 空格(1) + 时间(12) + 空格(1) + 级别(3) + 空格(1) + 文件名和行数(len(fileName) + 3 + 行数(5)) + 空格(1) + 函数名(2+len(methodName)) + 空格(2) + 日志内容(len(fileName)) + 回车换行(2)
	//
	// 裁剪为基础文件名
	fileName = filepath.Base(fileName)
	// 计算需要的大小
	size := 43 + len(content) + len(fileName) + len(methodName)
	// 创建一个指定容量的切片，避免二次扩容
	logSlice := make([]byte, 0, size)
	// 追加格式化好的日期和时间
	logSlice = append(logSlice, tool.StringToBytes(logTime.Format("2006/01/02 15:04:05.000"))...) // 23个字节
	// 追加级别
	logSlice = append(logSlice, ' ', '[', level.GetLevelChar(), ']') // 4个字节
	// 追加文件名和行号，len(strconv.FormatInt(int64(fileName), 10))大于5个字节时，logSlice会发生扩容
	logSlice = append(logSlice, ' ', '[')                                                    // 2个字节
	logSlice = append(logSlice, tool.StringToBytes(fileName)...)                             // len(fileName)个字节
	logSlice = append(logSlice, ':')                                                         // 1个字节
	logSlice = append(logSlice, tool.StringToBytes(strconv.FormatInt(int64(lineNo), 10))...) // 默认5个字节
	logSlice = append(logSlice, ']')                                                         // 1个字节
	// 追加函数名
	logSlice = append(logSlice, ' ', '[')                          // 2个字节
	logSlice = append(logSlice, tool.StringToBytes(methodName)...) // len(methodName)个字节
	logSlice = append(logSlice, ']')                               // 1个字节
	// 追加日志内容
	logSlice = append(logSlice, ' ', ' ')   // 2个字节
	logSlice = append(logSlice, content...) // len(content)个字节
	// 追加回车换行
	logSlice = append(logSlice, '\r', '\n') // 2个字节

	// 写入到管道中
	e.fileWriteChan <- logSlice
	// 同步时需要等待日志缓冲区被清空
	if !e.fileWriteAsync {
		// 空会在写入时进行判断
		e.fileWriteChan <- nil
	}
}

// getFileLines 获取文件当前总行数
func getFileLines(file *os.File) uint64 {
	// 行数统计
	var lines uint64 = 0
	// 获取文件行数
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		lines++
	}
	// 返回行数
	return lines
}

// selectAvailableFile 选择一个可用的文件
func (e *Adapter) selectAvailableFile() error {
	// 赋值当前日期
	e.currTime = time.Now()
	// 循环取文件名(单日最多允许存在999个日志文件)
	for i := 1; i <= 999; i++ {
		ok, err := func(j int) (bool, error) {
			// 赋值文件分割后缀标识
			e.currIndex = uint32(i)
			// 以当前日期命名文件名
			e.currLogPath = fmt.Sprintf(e.logPathFormat, e.currTime.Format("2006-01-02"), i)

			// 判断拼接出来的文件状态
			file, err := os.Open(e.currLogPath)
			if err != nil {
				if os.IsNotExist(err) {
					// 文件不存在，则该文件可用，跳出循环
					return true, nil
				}
				// 文件存在，但获取信息错误
				return false, err
			}
			defer file.Close()

			// 判断文件大小是否超过限制
			fileStat, _ := file.Stat()
			if fileStat.Size() >= int64(e.maxSize) {
				// 超过了限制，递增后缀标识
				return false, nil
			}

			// 判断文件是否超过了最大行数
			lines := getFileLines(file)
			if lines >= e.maxLines {
				// 超过了限制，递增后缀标识
				return false, nil
			}

			// 文件可用
			return true, nil
		}(i)

		// 是否异常
		if err != nil {
			return err
		}

		// 文件可用
		if ok {
			return nil
		}
	}

	// 全部文件不可用
	return errors.New("files from 0 to 999 are not available")
}

// deleteTimeoutLogFile 过期日志文件删除
func (e *Adapter) deleteTimeoutLogFile() {
	// 获取日志储存文件夹部分
	logDirPath := filepath.Dir(e.logPath)
	// 打开文件夹
	logDir, err := os.ReadDir(logDirPath)
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

// fileSplit 日志文件是否需要分割
func (e *Adapter) fileSplit() bool {
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
		return true
	}

	// 判断容量或行数是否超过了
	if e.currSize >= e.maxSize || e.currLines >= e.maxLines {
		// 后缀标识加1
		e.currIndex++
		// 拼接后缀标识及文件后缀
		e.currLogPath = fmt.Sprintf(e.logPathFormat, e.currTime.Format("2006-01-02"), e.currIndex)
		// 通知执行文件分割
		return true
	}

	// 不分隔文件
	return false
}

// writeFile 写入日志到文件中
func (e *Adapter) writeFile() {
	// 函数结束时的操作
	defer func() {
		// 异步执行一次过期日志文件删除
		go e.deleteTimeoutLogFile()
	}()

	// 移除软链
	err := os.RemoveAll(e.logPath)
	if err != nil {
		log.Fatalln("remove file error: " + err.Error())
	}

	// 创建或追加文件
	file, err := os.OpenFile(e.currLogPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		log.Fatalln("open file error: " + err.Error())
	}
	defer func() {
		file.Sync()  // 同步IO底层缓存到磁盘
		file.Close() // 关闭文件句柄
	}()

	// 赋值当前文件大小
	fileStat, _ := file.Stat()
	e.currSize = uint64(fileStat.Size())
	// 赋值当前文件行数
	e.currLines = getFileLines(file)

	// 创建写入缓冲区
	writer := bufio.NewWriter(file)
	defer func() {
		writer.Flush() // 结束时刷新到文件中
	}()

	// 创建软链
	err = os.Link(e.currLogPath, e.logPath)
	if err != nil {
		log.Fatalln("create file link error: " + err.Error())
	}

	// 创建当前时间在指定时间后接收到信号的管道
	specifiedTimeAfter := time.After(time.Hour * 1)

	// 监听文件写入或重新打开新的日志文件
	for {
		select {

		// 是否持久监听4个小时了
		case <-specifiedTimeAfter:
			// 判断一下是否需要分隔文件了
			if e.fileSplit() {
				// 结束当前文件的写入
				return
			}

		// 写日志
		case logStr := <-e.fileWriteChan:
			// 为空时执行跳过
			if logStr == nil {
				break
			}

			// 写入到缓冲区
			count, err := writer.Write(logStr)
			if err != nil {
				log.Println(err.Error())
			}

			// 增加当前文件已写入的大小
			e.currSize += uint64(count)
			// 增加当前文件已写入的行数
			e.currLines++

			// 判断大小或行数是否超过
			if e.currSize >= e.maxSize || e.currLines >= e.maxLines {
				// 判断一下是否需要分隔文件了
				if e.fileSplit() {
					// 结束当前文件的写入
					return
				}
			}
		}
	}
}
