package encoder

import (
	"math"
	"strconv"
	"time"

	"github.com/bearki/belog/v2/field"
	"github.com/bearki/belog/v2/internal/convert"
)

// appendFieldValue 追加字段值
func appendFieldValue(isJson bool, dst []byte, val field.Field) []byte {
	// 是否为JSON格式
	if isJson {
		// 追加键名
		dst = append(dst, '"')
		dst = append(dst, val.Key...)
		dst = append(dst, `": `...)
	} else {
		// 追加键名
		dst = append(dst, val.Key...)
		dst = append(dst, `: `...)
	}

	// 根据类型追加值
	switch true {

	// int8 <= type <= int64
	case field.IsValidRange(field.TypeInt8, val.ValType, field.TypeInt64):
		dst = strconv.AppendInt(dst, val.Integer, 10)

	// uint8 <= type <= uint64
	case field.IsValidRange(field.TypeUint8, val.ValType, field.TypeUint64):
		dst = strconv.AppendUint(dst, uint64(val.Integer), 10)

	// type == float32
	case val.ValType == field.TypeFloat32:
		dst = strconv.AppendFloat(dst, float64(math.Float32frombits(uint32(val.Integer))), 'E', -1, 32)

	// type == float64
	case val.ValType == field.TypeFloat64:
		dst = strconv.AppendFloat(dst, math.Float64frombits(uint64(val.Integer)), 'E', -1, 64)

	// type == complex64
	case val.ValType == field.TypeComplex64:
		dst = append(dst, strconv.FormatComplex(complex128(val.Interface.(complex64)), 'E', -1, 64)...)

	// type == complex128
	case val.ValType == field.TypeComplex128:
		dst = append(dst, strconv.FormatComplex(val.Interface.(complex128), 'E', -1, 128)...)

	// type == nil
	case val.ValType == field.TypeNull:
		dst = append(dst, val.String...)

	// type == bool
	case val.ValType == field.TypeBool:
		dst = strconv.AppendBool(dst, convert.IntToBool(int(val.Integer)))

	// type == string
	case val.ValType == field.TypeString:
		// 是否为JSON格式
		if isJson {
			dst = append(dst, '"')
			dst = append(dst, val.String...)
			dst = append(dst, '"')
		} else {
			dst = append(dst, val.String...)
		}

	// type == time
	case val.ValType == field.TypeTime:
		// 微秒时间戳
		dst = appendTime(isJson, dst, time.UnixMicro(val.Integer), val.String)
	}

	// 组装完成
	return dst
}

// AppendFieldAndMsg 将字段拼接为行格式
//
// @params dst 目标切片
//
// @params message 日志消息
//
// @params val 字段列表
//
// @return 序列化后的行格式字段字符串
//
// 返回示例，反引号内为实际内容:
// `k1: v1, k2: v2, ..., message`
func AppendFieldAndMsg(dst []byte, message string, val ...field.Field) []byte {
	// 遍历所有字段
	for _, v := range val {
		// 追加字段并序列化
		dst = appendFieldValue(false, dst, v)
		// 追加分隔符
		dst = append(dst, `, `...)
	}

	// 追加message内容
	dst = append(dst, convert.StringToBytes(message)...)

	// 返回组装好的数据
	return dst
}

// AppendFieldAndMsgJSON 将字段拼接为json格式
//
// @params dst 目标切片
//
// @params messageKey 消息的键名
//
// @params message 消息内容
//
// @params fieldsKey 包裹所有字段的键名
//
// @params val 字段列表
//
// @return 序列化后的JSON格式字段字符串
//
// 返回示例，反引号内为实际内容:
// `"fields": {"k1": "v1", ...}, "msg": "message"`
func AppendFieldAndMsgJSON(dst []byte, messageKey string, message string, fieldsKey string, val ...field.Field) []byte {
	// 追加字段集字段
	dst = append(dst, '"')
	dst = append(dst, fieldsKey...)
	dst = append(dst, `": {`...)
	// 遍历所有字段
	for i, v := range val {
		// 从第二个有效字段开始追加分隔符号
		if i > 0 {
			dst = append(dst, `, `...)
		}
		// 追加字段并序列化
		dst = appendFieldValue(true, dst, v)
	}
	// 追加字段结束括号
	dst = append(dst, `}, `...)

	// 追加message字段及其内容
	dst = append(dst, '"')
	dst = append(dst, messageKey...)
	dst = append(dst, `": "`...)
	dst = append(dst, convert.StringToBytes(message)...)
	dst = append(dst, '"')

	// 返回组装好的数据
	return dst
}
