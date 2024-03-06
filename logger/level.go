package logger

// Level 日志级别类型
type Level uint8

// 日志保存级别定义
const (
	Trace Level = 1 // 通知级别
	Debug Level = 2 // 调试级别
	Info  Level = 3 // 普通级别
	Warn  Level = 4 // 警告级别
	Error Level = 5 // 错误级别
	Fatal Level = 6 // 致命级别
)

// levelByteMap 日志级别字符映射
var levelByteMap = map[Level]byte{
	Trace: 'T',
	Debug: 'D',
	Info:  'I',
	Warn:  'W',
	Error: 'E',
	Fatal: 'F',
}

// levelStringMap 日志级别字符串映射
var levelStringMap = map[Level]string{
	Trace: "trace",
	Debug: "debug",
	Info:  "info",
	Warn:  "warning",
	Error: "error",
	Fatal: "fatal",
}

// Byte 获取日志级别对应的字符
func (l Level) Byte() byte {
	if c, ok := levelByteMap[l]; ok {
		return c
	}
	return ' '
}

// Byte 获取日志级别对应的字符串
func (l Level) String() string {
	if s, ok := levelStringMap[l]; ok {
		return s
	}
	return " "
}
