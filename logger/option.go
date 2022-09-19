package logger

// Option 日志记录器初始化参数
type Option struct {
	// 是否记录调用栈
	//
	// Default: false
	PrintCallStack bool

	// 是否禁用JSON序列化输出
	//
	// Default: false
	DisabledJsonFormat bool

	// 时间序列化格式
	//
	// 支持time.Format的所有格式
	//
	// Default: 2006/01/02 15:04:05.000
	TimeFormat string

	// 时间的JSON键名
	//
	// Default: time
	TimeJsonKey string

	// 日志级别的JSON键名
	//
	// Default: level
	LevelJsonKey string // 日志级别的JSON键名

	// 调用栈信息JSON键名
	//
	// Default: stack
	StackJsonKey string

	// 调用栈文件名JSON键名
	//
	// Default: file
	StackFileJsonKey string

	// 调用栈文件行号JSON键名
	//
	// Default: line
	StackLineNoJsonKey string

	// 调用栈函数名JSON键名
	//
	// Default: method
	StackMethodJsonKey string

	// 字段集JSON键名
	//
	// Default: fields
	FieldsJsonKey string

	// 日志消息JSON键名
	//
	// Default: message
	MessageJsonKey string
}

// 默认的参数配置
var defaultOption = Option{
	TimeFormat:         "2006/01/02 15:04:05.000",
	TimeJsonKey:        "time",
	LevelJsonKey:       "level",
	StackJsonKey:       "stack",
	StackFileJsonKey:   "file",
	StackLineNoJsonKey: "line",
	StackMethodJsonKey: "method",
	FieldsJsonKey:      "fields",
	MessageJsonKey:     "message",
}

// getValidOption 获取有效参数
func getValidOption(option Option) Option {
	// 时间序列化格式是否需要使用默认格式
	if len(option.TimeFormat) == 0 {
		option.TimeFormat = defaultOption.TimeFormat
	}

	// 时间的JSON键名是否使用默认值
	if len(option.TimeJsonKey) == 0 {
		option.TimeJsonKey = defaultOption.TimeJsonKey
	}

	// 日志级别的JSON键名是否使用默认值
	if len(option.LevelJsonKey) == 0 {
		option.LevelJsonKey = defaultOption.LevelJsonKey
	}

	// 调用栈信息JSON键名是否使用默认值
	if len(option.StackJsonKey) == 0 {
		option.StackJsonKey = defaultOption.StackJsonKey
	}

	// 调用栈文件名JSON键名是否使用默认值
	if len(option.StackFileJsonKey) == 0 {
		option.StackFileJsonKey = defaultOption.StackFileJsonKey
	}

	// 调用栈文件行号JSON键名是否使用默认值
	if len(option.StackLineNoJsonKey) == 0 {
		option.StackLineNoJsonKey = defaultOption.StackLineNoJsonKey
	}

	// 调用栈函数名JSON键名是否使用默认值
	if len(option.StackMethodJsonKey) == 0 {
		option.StackMethodJsonKey = defaultOption.StackMethodJsonKey
	}

	// 字段集JSON键名是否使用默认值
	if len(option.FieldsJsonKey) == 0 {
		option.FieldsJsonKey = defaultOption.FieldsJsonKey
	}

	// 日志消息JSON键名是否使用默认值
	if len(option.MessageJsonKey) == 0 {
		option.MessageJsonKey = defaultOption.MessageJsonKey
	}

	// OK,参数校验完成
	return option
}
