package encoder

import "time"

// BaseOption 编码器基础参数
type BaseOption struct {
	// 时间序列化格式
	//
	// 支持time.Format的所有格式和以下内置格式
	//
	// Unix(秒级时间戳),
	// UnixMilli(毫秒级时间戳),
	// UnixMicro(微秒级时间戳),
	// UnixNano(纳秒级时间戳)
	//
	// Default: UnixMilli
	TimeFormat string

	// 日志级别输出格式
	//
	// true: 使用日志级别完整字符串
	// false: 使用日志级别首字符
	//
	// Default: false
	LevelFormat bool

	// 调用栈文件路径格式
	//
	// true: 使用完整路径
	// false: 使用Base文件名
	//
	// Default: false
	StackFileFormat bool
}

// DefaultBaseOption 编码器默认基础参数
var DefaultBaseOption = BaseOption{
	TimeFormat:      TimeFormatUnixMilli,
	LevelFormat:     false,
	StackFileFormat: false,
}

// checkBaseOptionValid 检查编码器基础参数有效性
func checkBaseOptionValid(opt BaseOption) BaseOption {
	// 检验时间格式
	if len(opt.TimeFormat) == 0 {
		opt.TimeFormat = DefaultBaseOption.TimeFormat
	}

	// 预初始化（time包首次Format很慢）
	_ = time.Now().Format(opt.TimeFormat)

	// 校验完成
	return opt
}
