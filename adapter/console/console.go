/**
 *@Title 控制台日志记录适配器
 *@Desc 控制台打印将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package console

import (
	"bufio"
	"bytes"
	"os"
	"time"

	"github.com/bearki/belog/v3/logger"
)

// 控制台字体颜色字节表
var (
	// 灰色
	colorGrayStartBytes = [...]byte{27, 91, 57, 48, 109}
	// 蓝色
	colorBlueStartBytes = [...]byte{27, 91, 51, 52, 109}
	// 绿色
	colorGreenStartBytes = [...]byte{27, 91, 51, 50, 109}
	// 黄色
	colorYellowStartBytes = [...]byte{27, 91, 51, 51, 109}
	// 红色
	colorRedStartBytes = [...]byte{27, 91, 51, 49, 109}
	// 洋红色
	colorMagentaStartBytes = [...]byte{27, 91, 51, 53, 109}
	// 重置
	colorResetBytes = [...]byte{27, 91, 48, 109}
)

// 获取日志级别对应的控制台颜色字节表
func getLevelConsoleColorBytes(level logger.Level) []byte {
	switch level {
	case logger.Trace: // 通知级别(灰色)
		return colorGrayStartBytes[:]
	case logger.Debug: // 调试级别(蓝色)
		return colorBlueStartBytes[:]
	case logger.Info: // 普通级别(绿色)
		return colorGreenStartBytes[:]
	case logger.Warn: // 警告级别(黄色)
		return colorYellowStartBytes[:]
	case logger.Error: // 错误级别(红色)
		return colorRedStartBytes[:]
	case logger.Fatal: // 紧急级别(洋红色)
		return colorMagentaStartBytes[:]
	default:
		return nil
	}
}

// 控制台日志适配器初始化参数
type Option struct {
	DisabledBuffer bool // 禁用缓冲区输出
	DisabledColor  bool // 禁用颜色输出
}

// Adapter 控制台日志适配器
type Adapter struct {
	disabledBuffer bool          // 禁用缓冲区输出
	disabledColor  bool          // 禁用颜色输出
	write          *bufio.Writer // 写入缓冲器
}

// New 创建控制台日志适配器
//
//	@param	opt	适配器参数
//	@return	适配器实例
func New(opt Option) logger.Adapter {
	adapter := &Adapter{
		disabledBuffer: opt.DisabledBuffer,
		disabledColor:  opt.DisabledColor,
	}
	if !opt.DisabledBuffer {
		adapter.write = bufio.NewWriter(os.Stdout)
	}
	return adapter
}

// Name 用于获取适配器名称
//
//	注意：请确保适配器名称不与其他适配器名称冲突
func (e *Adapter) Name() string {
	return "belog-console-adapter"
}

// Print 普通日志打印方法
//
//	@param	logTime	日记记录时间
//	@param	level	日志级别
//	@param	content	日志内容
func (e *Adapter) Print(_ time.Time, level logger.Level, content []byte) {
	// 是否禁用颜色
	if !e.disabledColor {
		oldBytes := []byte{'[', level.Byte(), ']'}
		if !bytes.Contains(content, oldBytes) {
			oldBytes = make([]byte, 0, len(level.String())+2)
			oldBytes = append(oldBytes, '[')
			oldBytes = append(oldBytes, level.String()...)
			oldBytes = append(oldBytes, ']')
		}
		newBytes := getLevelConsoleColorBytes(level)
		newBytes = append(newBytes, oldBytes...)
		newBytes = append(newBytes, colorResetBytes[:]...)
		// 替换颜色
		content = bytes.Replace(
			content,
			oldBytes,
			newBytes,
			1,
		)
	}

	// 是否禁用缓冲区 或 缓冲器输出器为空
	if e.disabledBuffer || e.write == nil {
		_, _ = os.Stdout.Write(content)
		return
	}

	// 使用缓冲区输出
	_, _ = e.write.Write(content)
}

// PrintStack 调用栈日志打印方法
//
//	@param	logTime		日记记录时间
//	@param	level		日志级别
//	@param	content		日志内容
//	@param	fileName	日志记录调用文件路径
//	@param	lineNo		日志记录调用文件行号
//	@param	methodName	日志记录调用函数名
func (e *Adapter) PrintStack(_ time.Time, level logger.Level, content []byte, _ string, _ int, _ string) {
	// 是否禁用颜色
	if !e.disabledColor {
		oldBytes := []byte{'[', level.Byte(), ']'}
		if !bytes.Contains(content, oldBytes) {
			oldBytes = make([]byte, 0, len(level.String())+2)
			oldBytes = append(oldBytes, '[')
			oldBytes = append(oldBytes, level.String()...)
			oldBytes = append(oldBytes, ']')
		}
		newBytes := getLevelConsoleColorBytes(level)
		newBytes = append(newBytes, oldBytes...)
		newBytes = append(newBytes, colorResetBytes[:]...)
		// 替换颜色
		content = bytes.Replace(
			content,
			oldBytes,
			newBytes,
			1,
		)
	}

	// 是否禁用缓冲区 或 缓冲器输出器为空
	if e.disabledBuffer || e.write == nil {
		_, _ = os.Stdout.Write(content)
		return
	}

	// 使用缓冲区输出
	_, _ = e.write.Write(content)
}

// Flush 日志缓存刷新
//
//	注意：用于日志缓冲区刷新，接收到该通知后需要立即将缓冲区中的日志持久化
func (e *Adapter) Flush() {
	if e.write != nil {
		_ = e.write.Flush()
	}
}
