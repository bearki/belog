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
