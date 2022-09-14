package logger

// Level 日志级别类型
type Level uint8

// LevelChar 日志级别字符类型
type LevelChar = byte

// 日志保存级别定义
var (
	LevelTrace Level = 1 // 通知级别
	LevelDebug Level = 2 // 调试级别
	LevelInfo  Level = 3 // 普通级别
	LevelWarn  Level = 4 // 警告级别
	LevelError Level = 5 // 错误级别
	LevelFatal Level = 6 // 致命级别
)

// levelMap 日志级别字符映射
var levelMap = map[Level]LevelChar{
	LevelTrace: 'T',
	LevelDebug: 'D',
	LevelInfo:  'I',
	LevelWarn:  'W',
	LevelError: 'E',
	LevelFatal: 'F',
}
