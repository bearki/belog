/**
 *@Title 无输出日志记录适配器
 *@Desc 无输出打印将在这里完成
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package discard

import (
	"io"
	"time"

	"github.com/bearki/belog/v2/level"
	"github.com/bearki/belog/v2/logger"
)

// Adapter 无输出日志适配器
type Adapter struct {
}

// New 创建无输出日志适配器
func New() logger.Adapter {
	return &Adapter{}
}

// Name 用于获取适配器名称
// 注意：请确保适配器名称不与其他适配器名称冲突
func (e *Adapter) Name() string {
	return "belog-discard-adapter"
}

// Print 普通日志打印方法
//
//	@var logTime 日记记录时间
//	@var level 日志级别
//	@var content 日志内容
func (e *Adapter) Print(_ time.Time, _ level.Level, content []byte) {
	io.Discard.Write(content)
}

// PrintStack 调用栈日志打印方法
//
//	@var logTime 日记记录时间
//	@var level 日志级别
//	@var content 日志内容
//	@var fileName 日志记录调用文件路径
//	@var lineNo 日志记录调用文件行号
//	@var methodName 日志记录调用函数名
func (e *Adapter) PrintStack(_ time.Time, _ level.Level, content []byte, _ string, _ int, _ string) {
	io.Discard.Write(content)
}

// Flush 日志缓存刷新
// 用于日志缓冲区刷新
// 接收到该通知后需要立即将缓冲区中的日志持久化
func (e *Adapter) Flush() {}
