package logger

import (
	"strconv"

	"github.com/bearki/belog/v2/pkg/tool"
)

// Field 字段序列化接口
type Field interface {
	Bytes() []byte
}

type intf struct {
	name []byte
	val  []byte
}

func (v intf) Bytes() []byte {
	tmp := make([]byte, 0, len(v.name)+len(v.val)+3)
	tmp = append(tmp, '"')
	tmp = append(tmp, v.name...)
	tmp = append(tmp, '"', ':', ' ')
	tmp = append(tmp, v.val...)
	return tmp
}

func Intf(name string, val int) Field {
	return &intf{
		name: tool.StringToBytes(name),
		val:  tool.StringToBytes(strconv.Itoa(val)),
	}
}
