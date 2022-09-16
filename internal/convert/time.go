package convert

import (
	"errors"
	"time"
)

// TimeToBytes 时间类型转字节流
// 格式：2022/09/16 19:59:04.580
func TimeToBytes(s []byte, t time.Time) error {
	if len(s) != 23 {
		return errors.New("time format to bytes is 23 byte")
	}
	b, e := t.MarshalText()
	if e != nil {
		return e
	}
	copy(s, b[:23])
	s[4] = '/'
	s[7] = '/'
	s[10] = ' '
	return nil
}
