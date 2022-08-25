package logger

import "time"

// Engine 引擎接口
type Engine interface {
	Init(interface{}) (Engine, error)                // 引擎初始化函数
	Print(time.Time, LevelChar, string, int, string) // 日志打印函数
}
