package logger

import "github.com/bearki/belog/v2/internal/encoder"

// 导出编码器参数类型

type (
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
