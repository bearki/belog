package field

import "github.com/bearki/belog/v2/internal/pool"

// 普通类型键值对序列化符号（包含：布尔型、（有、无）符号整形、浮点型、null...等无需使用符号包裹的类型）
var (
	normalValPrefix = [0]byte{}
	normalValSuffix = [0]byte{}
)

// 字符串类型键值对序列化符号
var (
	stringValPrefix = [1]byte{'"'}
	stringValSuffix = [1]byte{'"'}
)

// Field 键值对序列化结构体
type Field struct {
	KeyBytes       []byte                // 键的字节流
	ValPrefixBytes []byte                // 值的前缀字节流
	ValBytes       []byte                // 值的字节流
	ValSuffixBytes []byte                // 值的后缀字节流
	valBytesPut    pool.BytesPoolPutFunc // 值的字节切片回收到复用池的方法
}

// Put 将字段引用的底层字节数组回收到复用池
func (v Field) Put() {
	// 回收值使用的底层字节数组
	if v.valBytesPut != nil {
		v.valBytesPut(v.ValBytes)
	}
}

var (
	// 8个容量字节切片复用池
	eightCapBytesPool = pool.NewBytesPool(100, 0, 8)
	// 16个容量的字节切片复用池
	sixteenCapBytesPool = pool.NewBytesPool(100, 0, 16)
)
