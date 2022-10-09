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

	"github.com/bearki/belog/v2/level"
	"github.com/bearki/belog/v2/logger"
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

// GetLevelConsoleColorBytes 获取日志级别对应的控制台颜色字节表
func GetLevelConsoleColorBytes(l level.Level) []byte {
	switch l {
	case level.Trace: // 通知级别(灰色)
		return colorGrayStartBytes[:]
	case level.Debug: // 调试级别(蓝色)
		return colorBlueStartBytes[:]
	case level.Info: // 普通级别(绿色)
		return colorGreenStartBytes[:]
	case level.Warn: // 警告级别(黄色)
		return colorYellowStartBytes[:]
	case level.Error: // 错误级别(红色)
		return colorRedStartBytes[:]
	case level.Fatal: // 紧急级别(洋红色)
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
func New(op Option) logger.Adapter {
	adapter := &Adapter{
		disabledBuffer: op.DisabledBuffer,
		disabledColor:  op.DisabledColor,
	}
	if !op.DisabledBuffer {
		adapter.write = bufio.NewWriter(os.Stdout)
	}
	return adapter
}

// Name 用于获取适配器名称
//
// 注意：请确保适配器名称不与其他适配器名称冲突
func (e *Adapter) Name() string {
	return "belog-console-adapter"
}

// Print 普通日志打印方法
//
// @params logTime 日记记录时间
//
// @params level 日志级别
//
// @params content 日志内容
func (e *Adapter) Print(_ time.Time, level level.Level, content []byte) {
	// 是否禁用颜色
	if !e.disabledColor {
		oldBytes := []byte{'[', level.Byte(), ']'}
		newBytes := GetLevelConsoleColorBytes(level)
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

	// 是否禁用缓冲区
	if e.disabledBuffer {
		_, _ = os.Stdout.Write(content)
	}

	// 使用缓冲区输出
	if e.write != nil {
		_, _ = e.write.Write(content)
	} else {
		os.Stdout.WriteString(e.Name() + " writer is nil pointer")
		_, _ = os.Stdout.Write(content)
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
func (e *Adapter) PrintStack(_ time.Time, level level.Level, content []byte, _ string, _ int, _ string) {
	// 是否禁用颜色
	if !e.disabledColor {
		oldBytes := []byte{'[', level.Byte(), ']'}
		newBytes := GetLevelConsoleColorBytes(level)
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

	// 是否禁用缓冲区
	if e.disabledBuffer {
		_, _ = os.Stdout.Write(content)
	}

	// 使用缓冲区输出
	if e.write != nil {
		_, _ = e.write.Write(content)
	} else {
		os.Stdout.WriteString(e.Name() + " writer is nil pointer")
		_, _ = os.Stdout.Write(content)
	}
}

// Flush 日志缓存刷新
//
// 用于日志缓冲区刷新
// 接收到该通知后需要立即将缓冲区中的日志持久化
func (e *Adapter) Flush() {
	if e.write != nil {
		_ = e.write.Flush()
	} else {
		os.Stdout.WriteString(e.Name() + " writer is nil pointer")
	}
}
