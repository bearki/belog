package logger

// Option 日志记录器初始化参数
type Option struct {
	// 是否记录调用栈
	//
	// Default: false
	EnabledStackPrint bool

	// Encoder 日志内容编码器
	Encoder Encoder
}

// 默认的参数配置
var defaultOption = Option{
	EnabledStackPrint: false,
	Encoder:           NewJsonEncoder(DefaultJsonOption),
}

// DefaultOption 默认配置
var DefaultOption = defaultOption

// checkOptionValid 获取有效参数
func checkOptionValid(opt Option) Option {
	if opt.Encoder == nil {
		opt.Encoder = defaultOption.Encoder
	}
	// OK,参数校验完成
	return opt
}
