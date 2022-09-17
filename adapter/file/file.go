/**
 *@Title 文件日志适配器
 *@Desc 文件日志的写入将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bearki/belog/v2/internal/convert"
	"github.com/bearki/belog/v2/internal/pool"
	"github.com/bearki/belog/v2/logger"
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

	logPathFormat    string        // 日志文件路径格式
	currLogPath      string        // 日志文件源文件路径
	currTime         time.Time     // 当前日志文件使用的日期
	currIndex        uint32        // 当前日志文件分割后缀标识
	currSize         uint64        // 当前日志文件大小（单位：byte）
	currLines        uint64        // 当前日志文件行数
	flushMutex       sync.Mutex    // 刷新操作锁
	flushStartSignal chan struct{} // 刷新开始信号
	flushOverSignal  chan struct{} // 刷新结束信号

	logBytesPool *pool.BytesPool // 日志字节流对象池
}

// printWarningMsg 打印警告信息
func printWarningMsg(msg string) {
	// _, _ = os.Stderr.WriteString(msg + "\r\n")
	_, _ = os.Stdout.WriteString(msg + "\r\n")
}

// validity 判断参数有效性
func (p *Options) validity() {
	// 转换路径为当前系统格式
	p.LogPath = filepath.Join(p.LogPath)
	// 判断路径是否为空
	if len(p.LogPath) <= 0 {
		p.LogPath = "app.log"
		printWarningMsg("`logPath` is null, use the default value `app.log`")
	}
	// 判断日志大小限制
	if p.MaxSize < 1 || p.MaxSize > 10000 {
		p.MaxSize = 4
		printWarningMsg("file log size min value is 1(MB),max value is 10000(MB), use the default value 4(MB)")
	}
	// 判断日志行数限制
	if p.MaxLines < 100 || p.MaxLines > 100000000 {
		p.MaxLines = 100000
		printWarningMsg("file log lines min value is 100,max value is 100000000, use the default value 100000")
	}
	// 判断日志保存天数
	if p.SaveDay < 1 || p.SaveDay > 1000 {
		p.SaveDay = 30
		printWarningMsg("file log save day min value is 1,max value is 1000, use the default value 30")
	}
	// 判断异步模式下的管道容量
	if p.Async {
		if p.AsyncChanCap > 100 {
			p.AsyncChanCap = 1
			printWarningMsg("async log channel cap error, min 1, max 100, use the default value 1")
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
	e.logPathFormat = filepath.Join(logDir, logName+".%s.%d"+logExt)
	// 赋值是否为异步写入
	e.fileWriteAsync = options.Async
	// 初始化日志写入管道容量
	e.fileWriteChan = make(chan []byte, options.AsyncChanCap)
	// 筛选出合适的下标日志文件
	if err = e.selectAvailableFile(); err != nil {
		return nil, err
	}
	// 初始化刷新信号管道
	e.flushStartSignal = make(chan struct{}, 1)
	e.flushOverSignal = make(chan struct{}, 1)
	// 日志字节流对象池
	e.logBytesPool = pool.NewBytesPool(100, 0, 1024)

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

// format 格式化输出内容
//
// @params t 日记记录时间
//
// @params l 日志级别
//
// @params c 日志内容
//
// @params f 日志记录调用文件路径
//
// @params n 日志记录调用文件行号
//
// @params m 日志记录调用函数名
func (e *Adapter) format(ps bool, t time.Time, l logger.Level, c []byte, f []byte, n int, m []byte) {
	// 不带调用栈：
	// 2022/09/14 20:28:13.793 [T]  this is a trace log\r\n
	// ++++++++++_++++++++++++_+++__+++++++++++++++++++____
	//     10    1    12      1 3 2   content           2
	//
	// const initSize = 10 + 1 + 12 + 1 + 3 + 2 + 2 = 31
	//
	// 带调用栈：
	// 2022/09/14 20:28:13.793 [T] [belog_test.go:82] [PrintLog]  this is a trace log\r\n
	// ++++++++++_++++++++++++_+++_++++++++++++++++++_++++++++++__+++++++++++++++++++____
	//     10    1    12      1 3 1   3+file+line    1 2+method  2   content           2
	//
	// const initSize = 10 + 1 + 12 + 1 + 3 + 1 + 3 + 1 + 2 + 2 + 2 = 38

	// 预留行号切片
	var lineNoBytes []byte

	// 计算需要的大小
	size := 31 + len(c)
	if ps {
		// 裁剪为基础文件名
		index := bytes.LastIndexByte(f, '/')
		if index > -1 && index+1 < len(f) {
			f = f[index+1:]
		}

		// 裁剪为基础函数名
		index = bytes.LastIndexByte(m, '/')
		if index > 0 && index+1 < len(m) {
			m = m[index+1:]
		}

		// 转换行号为切片
		lineNoBytes := convert.StringToBytes(strconv.Itoa(n))
		// 重新计算需要的大小
		size = 38 + len(f) + len(lineNoBytes) + len(m) + len(c)
	}

	// 从对象池获取切片
	logSlice := e.logBytesPool.Get()
	// 检查是否需要扩容
	if cap(logSlice) < size {
		// 创建一个指定容量的切片，避免二次扩容
		logSlice = make([]byte, 0, size)
	}

	// 追加格式化好的日期和时间
	logSlice = t.AppendFormat(logSlice, "2006/01/02 15:04:05.000") // 23个字节
	// 追加空格
	logSlice = append(logSlice, ' ') // 1个字节
	// 追加级别
	logSlice = append(logSlice, '[', l.GetLevelChar(), ']') // 3个字节
	// 是否需要记录调用栈
	if ps {
		// 追加空格
		logSlice = append(logSlice, ' ') // 1个字节
		// 追加文件名和行号
		logSlice = append(logSlice, '[')            // 1个字节
		logSlice = append(logSlice, f...)           // len(fileName)个字节
		logSlice = append(logSlice, ':')            // 1个字节
		logSlice = append(logSlice, lineNoBytes...) // len(lineNo)个字节
		logSlice = append(logSlice, ']')            // 1个字节
		// 追加空格
		logSlice = append(logSlice, ' ') // 1个字节
		// 追加函数名
		logSlice = append(logSlice, '[')  // 1个字节
		logSlice = append(logSlice, m...) // len(methodName)个字节
		logSlice = append(logSlice, ']')  // 1个字节
	}
	// 追加空格
	logSlice = append(logSlice, ' ', ' ') // 2个字节
	// 追加日志内容
	logSlice = append(logSlice, c...) // len(content)个字节
	// 追加回车换行
	logSlice = append(logSlice, '\r', '\n') // 2个字节

	// 发送到管道
	e.fileWriteChan <- logSlice
}

// Print 普通日志打印方法
//
// @params t 日记记录时间
//
// @params l 日志级别
//
// @params c 日志内容
func (e *Adapter) Print(t time.Time, l logger.Level, c []byte) {
	// 执行格式化并推送到管道
	e.format(false, t, l, c, nil, 0, nil)
}

// PrintStack 调用栈日志打印方法
//
// @params t 日记记录时间
//
// @params l 日志级别
//
// @params c 日志内容
//
// @params f 日志记录调用文件路径
//
// @params n 日志记录调用文件行号
//
// @params m 日志记录调用函数名
func (e *Adapter) PrintStack(t time.Time, l logger.Level, c []byte, f []byte, n int, m []byte) {
	// 执行格式化并推送到管道
	e.format(true, t, l, c, f, n, m)
}

// Flush 日志缓存刷新
//
// 用于日志缓冲区刷新
// 接收到该通知后需要立即将缓冲区中的日志持久化
func (e *Adapter) Flush() {
	// 加锁
	e.flushMutex.Lock()
	defer e.flushMutex.Unlock()

	// 发送刷新开始信号
	e.flushStartSignal <- struct{}{}

	// 阻塞，直到刷新完成
	<-e.flushOverSignal
}

// openFileGetLines 打开文件并获取文件总行数
//
// @params fileName 文件路径
//
// @params flag 文件打开模式
//
// @params closeFile 结束时是否关闭文件
//
// @return 文件句柄
//
// @return 总行数
//
// @return 异常信息
func openFileGetLines(fileName string, flag int, closeFile bool) (*os.File, uint64, error) {
	// 打开文件
	file, err := os.OpenFile(fileName, flag, 0666)
	if err != nil {
		return nil, 0, err
	}
	if closeFile {
		defer file.Close()
	}
	// 行数统计
	var lines uint64 = 0
	// 获取文件行数
	reader := bufio.NewScanner(file)
	for reader.Scan() {
		lines++
	}
	// 返回行数
	return file, lines, nil
}

// selectAvailableFile 选择一个可用的文件
func (e *Adapter) selectAvailableFile() error {
	// 赋值当前日期
	e.currTime = time.Now()
	// 循环取文件名
	for i := 1; i <= math.MaxInt32; i++ {
		// 赋值文件分割后缀标识
		e.currIndex = uint32(i)
		// 以当前日期命名文件名
		e.currLogPath = fmt.Sprintf(e.logPathFormat, e.currTime.Format("2006-01-02"), i)

		// 先获取文件信息
		fileInfo, err := os.Stat(e.currLogPath)
		if err != nil {
			if os.IsNotExist(err) {
				// 文件不存在，则该文件可用，跳出循环
				return nil
			}
			// 文件存在，但获取信息错误
			return err
		}

		// 判断文件大小是否超过限制
		if fileInfo.Size() >= int64(e.maxSize) {
			// 超过了限制，递增后缀标识
			continue
		}

		// 判断文件是否超过了最大行数
		_, lines, err := openFileGetLines(e.currLogPath, os.O_RDONLY, true)
		if err != nil {
			// 文件异常
			return err
		}
		if lines >= e.maxLines {
			// 超过了限制，递增后缀标识
			continue
		}

		// 文件可用
		return nil
	}

	// 全部文件不可用
	return errors.New("no available files found")
}

// deleteTimeoutLogFile 过期日志文件删除
func (e *Adapter) deleteTimeoutLogFile() {
	// 获取日志储存文件夹部分
	logDirPath := filepath.Dir(e.logPath)
	// 打开文件夹
	logDir, err := os.ReadDir(logDirPath)
	if err != nil {
		printWarningMsg(err.Error())
		return
	}

	// 获取当天整点时间
	currDateStr := time.Now().Format("2006-01-02")
	// 再解析成时间类型
	currDate, err := time.Parse("2006-01-02", currDateStr)
	if err != nil {
		printWarningMsg(err.Error())
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
					printWarningMsg(err.Error())
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

	// 移除硬连接
	err := os.RemoveAll(e.logPath)
	if err != nil {
		log.Fatalln("remove file error: " + err.Error())
	}

	// 创建或追加文件，并获取文件总行数
	file, lines, err := openFileGetLines(e.currLogPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, false)
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
	e.currLines = lines

	// 创建写入缓冲区
	writer := bufio.NewWriter(file)
	defer func() {
		writer.Flush() // 结束时刷新到文件中
	}()

	// 创建硬连接
	err = os.Link(e.currLogPath, e.logPath)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatalln("create file link error: " + err.Error())
	}

	// 阻塞，执行监听写入
	e.listenBufioWrite(file, writer)
}

// listenBufioWrite 监听日志并通过bufio写入
func (e *Adapter) listenBufioWrite(file *os.File, writer *bufio.Writer) {
	// 计算距离第二天凌晨0点还有多少时间
	tmpTime := e.currTime.AddDate(0, 0, 1)
	zeroTime := time.Date(tmpTime.Year(), tmpTime.Month(), tmpTime.Day(), 0, 0, 0, 0, time.Local)

	// 分割文件的信号管道，
	fileSplitChan := time.After(zeroTime.Sub(e.currTime))

	// 在指定时间间隔强制刷新一次缓冲区
	specifiedTimeAfter := time.NewTicker(time.Minute * 5)
	defer specifiedTimeAfter.Stop()

	// 预声明
	var count int
	var err error

	// 监听文件写入或重新打开新的日志文件
	for {
		select {

		// 是否监听到日志来了
		case logBytes := <-e.fileWriteChan:
			// 检查内容
			if logBytes == nil {
				// 跳过
				break
			}

			// 判断写入模式
			if e.fileWriteAsync {
				// 异步模式使用缓冲区写入
				count, err = writer.Write(logBytes)
			} else {
				// 使用文件句柄直接写入
				count, err = file.Write(logBytes)
			}
			if err != nil {
				printWarningMsg(err.Error())
			}

			// 将字节流对象放回对象池
			e.logBytesPool.Put(logBytes)

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

		// 是否接收到日志刷新信号
		case <-e.flushStartSignal:
			// 管道中是否还有内容
			num := len(e.fileWriteChan)
			for i := 0; i < num; i++ {
				logStr := <-e.fileWriteChan
				if logStr == nil {
					continue
				}
				if _, err := writer.Write(logStr); err != nil {
					printWarningMsg(err.Error())
				}
			}

			// 缓冲区内有内容时执行缓冲区刷新
			if writer.Buffered() > 0 {
				writer.Flush()
				file.Sync()
			}

			// 发送刷新完成信号
			e.flushOverSignal <- struct{}{}

		// 是否到达凌晨0点了
		case <-fileSplitChan:
			if e.fileSplit() || writer.Buffered() > 0 {
				return
			}

		// 是否需要强制刷新了
		case <-specifiedTimeAfter.C:
			if e.fileSplit() || writer.Buffered() > 0 {
				return
			}
		}
	}
}
