package field

import (
	"encoding/binary"

	"github.com/bearki/belog/v2/internal/convert"
)

var (
	intnFmtStartSym     = []byte{'"'}
	intnFmtIntervaltSym = []byte{'"', ':', ' '}
	intnFmtEndSym       = []byte{}
)

func IntnToBytes(val int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(val))
	return nil
}

// Intn 格式化~int类型字段信息
//
// 拼接格式  "index": 20
func Intn(name string, val int64) *Field {
	// 从对象池中取一个对象
	field := fieldStructPool.Get()
	nameBytes := convert.StringToBytes(name)
	valBytes := IntnToBytes(val)
	field.Size = len(intnFmtStartSym) + len(name) + len(intnFmtIntervaltSym) + len(valBytes) + len(intnFmtEndSym)
	field.StartSymBytes = intnFmtStartSym
	field.NameBytes = nameBytes
	field.IntervaltSymBytes = intnFmtIntervaltSym
	field.ValBytes = valBytes
	field.EndSymBytes = intnFmtEndSym
	return field
}
