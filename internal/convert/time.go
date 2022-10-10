package convert

import "time"

// TimeToInt64 时间类型转int64时间戳
func TimeToInt64(t time.Time) int64 {
	return t.UnixMicro()
}

// TimeFromInt64 int64时间戳转时间类型
func TimeFromInt64(v int64) time.Time {
	return time.UnixMicro(v)
}
