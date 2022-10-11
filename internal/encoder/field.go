package encoder

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bearki/belog/v2/field"
	"github.com/bearki/belog/v2/internal/convert"
)

// 追加字段普通类型值
func appendFieldValue(isJson bool, dst []byte, val field.Field) []byte {
	// 根据类型追加值
	switch val.Type {

	// Type == time，使用字段为Integer
	case field.TypeTime:
		// 微秒时间戳
		dst = appendTimeValue(isJson, dst, convert.TimeFromInt64(val.Integer), val.String)

	// Type == ~intn，使用字段为Integer
	case field.TypeInt8, field.TypeInt16, field.TypeInt, field.TypeInt32, field.TypeInt64, field.TypeDuration:
		dst = strconv.AppendInt(dst, val.Integer, 10)

	// Type == ~uintn，使用字段为Integer
	case field.TypeUint8, field.TypeUint16, field.TypeUint, field.TypeUint32, field.TypeUint64, field.TypeByte, field.TypeUintptr:
		dst = strconv.AppendUint(dst, uint64(val.Integer), 10)

	// Type == float32，使用字段为Integer
	case field.TypeFloat32:
		dst = strconv.AppendFloat(dst, float64(convert.Float32FromInt64(val.Integer)), 'E', -1, 32)

	// Type == float64，使用字段为Integer
	case field.TypeFloat64:
		dst = strconv.AppendFloat(dst, convert.Float64FromInt64(val.Integer), 'E', -1, 64)

	// Type == complex64，使用字段为Interface
	case field.TypeComplex64:
		dst = append(dst, strconv.FormatComplex(complex128(val.Interface.(complex64)), 'E', -1, 64)...)

	// Type == complex128，使用字段为Interface
	case field.TypeComplex128:
		dst = append(dst, strconv.FormatComplex(val.Interface.(complex128), 'E', -1, 128)...)

	// Type == nil，使用字段为String
	case field.TypeNull:
		dst = append(dst, val.String...)

	// Type == bool，使用字段为Integer
	case field.TypeBool:
		dst = strconv.AppendBool(dst, convert.BoolFromInt64(val.Integer))

	// Type == field.Objecter，使用字段为Interface
	case field.TypeObjecter:
		obj := val.Interface.(field.Objecter)
		// 是否为JSON格式
		if isJson {
			dst = append(dst, obj.ToJSON()...)
		} else {
			dst = append(dst, obj.ToString()...)
		}

	// string <= type == error，使用字段为String
	case field.TypeString, field.TypeError:
		// 是否为JSON格式
		if isJson {
			dst = append(dst, '"')
			dst = append(dst, val.String...)
			dst = append(dst, '"')
		} else {
			dst = append(dst, val.String...)
		}

	}

	// 追加完成
	return dst
}

// 追加字段切切片类型值
func appendFieldValues(isJson bool, dst []byte, val field.Field) []byte {
	// 预构建一个复用字段
	f := field.Field{
		Key: val.Key,
	}

	// 根据类型追加值，所有切片类型均使用Interface字段
	switch val.Type {

	// Type == []time
	case field.TypeTimes:
		f.Type = field.TypeTime
		f.String = val.String
		tmps := val.Interface.([]time.Time)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = convert.TimeToInt64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []int8
	case field.TypeInt8s:
		f.Type = field.TypeInt8
		tmps := val.Interface.([]int8)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []int16
	case field.TypeInt16s:
		f.Type = field.TypeInt16
		tmps := val.Interface.([]int16)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []int
	case field.TypeInts:
		f.Type = field.TypeInt
		tmps := val.Interface.([]int)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []int32
	case field.TypeInt32s:
		f.Type = field.TypeInt32
		tmps := val.Interface.([]int32)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []int64
	case field.TypeInt64s:
		f.Type = field.TypeInt64
		tmps := val.Interface.([]int64)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = v
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []time.Duration
	case field.TypeDurations:
		f.Type = field.TypeDuration
		tmps := val.Interface.([]time.Duration)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []uint8
	case field.TypeUint8s:
		f.Type = field.TypeUint8
		tmps := val.Interface.([]uint8)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []uint16
	case field.TypeUint16s:
		f.Type = field.TypeUint16
		tmps := val.Interface.([]uint16)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []uint
	case field.TypeUints:
		f.Type = field.TypeUint
		tmps := val.Interface.([]uint)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []uint32
	case field.TypeUint32s:
		f.Type = field.TypeUint32
		tmps := val.Interface.([]uint32)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []uint64
	case field.TypeUint64s:
		f.Type = field.TypeUint64
		tmps := val.Interface.([]uint64)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []byte
	case field.TypeBytes:
		f.Type = field.TypeByte
		tmps := val.Interface.([]byte)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []uintptr
	case field.TypeUintptrs:
		f.Type = field.TypeUintptr
		tmps := val.Interface.([]uintptr)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = int64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []float32
	case field.TypeFloat32s:
		f.Type = field.TypeFloat32
		tmps := val.Interface.([]float32)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = convert.Float32ToInt64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []float64
	case field.TypeFloat64s:
		f.Type = field.TypeFloat64
		tmps := val.Interface.([]float64)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = convert.Float64ToInt64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []complex64
	case field.TypeComplex64s:
		f.Type = field.TypeComplex64
		tmps := val.Interface.([]complex64)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Interface = v
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []complex128
	case field.TypeComplex128s:
		f.Type = field.TypeComplex128
		tmps := val.Interface.([]complex128)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Interface = v
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []bool
	case field.TypeBools:
		f.Type = field.TypeBool
		tmps := val.Interface.([]bool)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.Integer = convert.BoolToInt64(v)
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []string
	case field.TypeStrings:
		f.Type = field.TypeString
		tmps := val.Interface.([]string)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			f.String = v
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []error
	case field.TypeErrors:
		tmps := val.Interface.([]error)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			if v == nil {
				f.Type = field.TypeNull
				f.String = "null"
			} else {
				f.Type = field.TypeError
				f.String = v.Error()
			}
			dst = appendFieldValue(isJson, dst, f)
		}

	// Type == []Objecter
	case field.TypeObjecters:
		tmps := val.Interface.([]field.Objecter)
		for i, v := range tmps {
			if i > 0 {
				dst = append(dst, `, `...)
			}
			if v == nil {
				f.Type = field.TypeNull
				f.String = "null"
			} else {
				f.Type = field.TypeObjecter
				f.Interface = v
			}
			dst = appendFieldValue(isJson, dst, f)
		}
	}

	// 追加完成
	return dst
}

// appendField 追加字段
func appendField(isJson bool, dst []byte, val field.Field) []byte {
	// 是否为JSON格式
	if isJson {
		// 追加键名
		dst = append(dst, '"')
		dst = append(dst, val.Key...)
		dst = append(dst, `": `...)
	} else {
		// 追加键名
		dst = append(dst, val.Key...)
		dst = append(dst, `:`...)
	}

	switch true {

	// 普通类型
	case field.NormalTypeStart < val.Type && val.Type < field.NormalTypeEnd:
		// 追加字段单个值
		dst = appendFieldValue(isJson, dst, val)

	// 切片类型
	case field.SliceTypeStart < val.Type && val.Type < field.SliceTypeEnd:
		// 在值的前面追加中括号
		dst = append(dst, '[')
		// 追加字段多个值
		dst = appendFieldValues(isJson, dst, val)
		// 在值的后面追加中括号
		dst = append(dst, ']')

	// 未知类型，走反射
	default:
		// 是否为JSON格式
		if isJson {
			tmp, _ := json.Marshal(val.Interface)
			dst = append(dst, tmp...)
		} else {
			dst = append(dst, fmt.Sprintf("%+v", val.Interface)...)
		}

	}

	// 组装完成
	return dst
}

// appendFieldAndMsg 将字段拼接为行格式
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
func appendFieldAndMsg(dst []byte, message string, val ...field.Field) []byte {
	// 遍历所有字段
	for _, v := range val {
		// 追加字段并序列化
		dst = appendField(false, dst, v)
		// 追加分隔符
		dst = append(dst, `, `...)
	}

	// 追加message内容
	dst = append(dst, convert.StringToBytes(message)...)

	// 返回组装好的数据
	return dst
}

// appendFieldAndMsgJSON 将字段拼接为json格式
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
func appendFieldAndMsgJSON(dst []byte, messageKey string, message string, fieldsKey string, val ...field.Field) []byte {
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
		dst = appendField(true, dst, v)
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
