package logger

import "github.com/bearki/belog/v2/internal/encoder"

// 导出时间格式
const (
	TimeFormat1         = encoder.TimeFormat1
	TimeFormat2         = encoder.TimeFormat2
	TimeFormat3         = encoder.TimeFormat3
	TimeFormat4         = encoder.TimeFormat4
	TimeFormat5         = encoder.TimeFormat5
	TimeFormat6         = encoder.TimeFormat6
	TimeFormat7         = encoder.TimeFormat7
	TimeFormat8         = encoder.TimeFormat8
	TimeFormatUnix      = encoder.TimeFormatUnix
	TimeFormatUnixMilli = encoder.TimeFormatUnixMilli
	TimeFormatUnixMicro = encoder.TimeFormatUnixMicro
	TimeFormatUnixNano  = encoder.TimeFormatUnixNano
)

// 导出编码器参数类型

type (
	EncoderBaseOption   = encoder.BaseOption          // 编码器基础参数类型
	NormalEncoderOption = encoder.NormalEncoderOption // 普通编码器参数类型
	JsonEncoderOption   = encoder.JsonEncoderOption   // JSON编码器参数类型
)

// 导出编码器默认参数

var (
	DefaultBaseOption   = encoder.DefaultBaseOption   // 所有编码器的基础默认参数
	DefaultNormalOption = encoder.DefaultNormalOption // 普通编码器默认参数
	DefaultJsonOption   = encoder.DefaultJsonOption   // JSON编码器默认参数
)

// NewNormalEncoder 创建一个普通格式编码器
func NewNormalEncoder(opt NormalEncoderOption) Encoder {
	return encoder.NewNormalEncoder(opt)
}

// NewJsonEncoder 创建一个JSON格式编码器
func NewJsonEncoder(opt JsonEncoderOption) Encoder {
	return encoder.NewJsonEncoder(opt)
}
