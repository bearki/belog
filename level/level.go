package level

// Level 日志级别类型
type Level uint8

// Char 日志级别字符类型
type Char = byte

// 日志保存级别定义
const (
	Trace Level = 1 // 通知级别
	Debug Level = 2 // 调试级别
	Info  Level = 3 // 普通级别
	Warn  Level = 4 // 警告级别
	Error Level = 5 // 错误级别
	Fatal Level = 6 // 致命级别
)

// levelMap 日志级别字符映射
var levelMap = map[Level]Char{
	Trace: 'T',
	Debug: 'D',
	Info:  'I',
	Warn:  'W',
	Error: 'E',
	Fatal: 'F',
}

// GetLevelChar 获取日志级别对应的字符
func (l Level) GetLevelChar() Char {
	if c, ok := levelMap[l]; ok {
		return c
	}
	return ' '
}
