package logger

import (
	"strconv"

	"github.com/bearki/belog/v2/pkg/tool"
)

// Field 字段序列化接口
type Field interface {
	Bytes() []byte // 获取字段信息的字节切片
}

type intf struct {
	content []byte
}

// Bytes 获取字段信息的字节切片
func (v intf) Bytes() []byte {
	return v.content
}

// Intf 格式化int类型字段信息
//
// 输入值:
//
//	{
//	  name: []byte("index"),
//	  val: []byte(20),
//	}
//
// 最终格式:
// []byte("\"index\"": 20")
func Intf(name string, val int) Field {
	valStr := strconv.Itoa(val)
	field := new(intf)
	field.content = make([]byte, 0, 4+len(name)+len(valStr))
	field.content = append(field.content, '"')
	field.content = append(field.content, tool.StringToBytes(name)...)
	field.content = append(field.content, '"', ':', ' ')
	field.content = append(field.content, tool.StringToBytes(valStr)...)
	return field
}
