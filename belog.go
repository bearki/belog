/**
 *@Title belog默认实例
 *@Desc 该实例可通过控制台打印日志
 *@Author Bearki
 *@DateTime 2021/09/21 19:16
 */

package belog

import (
	"github.com/bearki/belog/v2/adapter/console"
	"github.com/bearki/belog/v2/logger"
)

// DefaultLog 默认实例(控制台适配器记录日志)
var DefaultLog, _ = logger.New(
	logger.Option{},
	console.New(console.Option{}),
)

// New 初始化一个日志记录器实例
//
// @params option 日志记录器初始化参数
//
// @params adapter 日志适配器
//
// @return 日志记录器实例
func New(option logger.Option, adapter ...logger.Adapter) (logger.Logger, error) {
	// 返回日志示例指针
	return logger.New(option, adapter...)
}
