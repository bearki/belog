package field

import (
	"sync"

	"github.com/bearki/belog/v2/internal/pool"
)

// fieldBytesPool 字段结构对象池
var fieldStructPool = pool.NewStructPool(1000, Field{})

func init() {
	// 预放置5个对象在对象池中
	for i := 0; i < 5; i++ {
		fieldStructPool.Put(&Field{})
	}
}

// Field 字段序列化
type Field struct {
	Size              int
	StartSymBytes     []byte
	NameBytes         []byte
	IntervaltSymBytes []byte
	ValBytes          []byte
	EndSymBytes       []byte
	allowPutMutex     sync.Mutex
	allowPut          bool
}

// Put 将对象放回到对象池中
func (f *Field) Put() {
	// // 异步操作（异步会增加开销）
	// go func() {
	// 是否允许被放回对象池
	if !f.allowPut {
		return
	}

	// 尝试互斥操作
	if !f.allowPutMutex.TryLock() {
		return
	}
	defer f.allowPutMutex.Unlock()

	// 是否允许被放回对象池
	if !f.allowPut {
		return
	}

	// 将对象放回对象池
	f.allowPut = !fieldStructPool.Put(f)
	// }()
}
