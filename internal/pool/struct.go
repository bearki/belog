package pool

// StructPool 对象池
type StructPool[T any] struct {
	sp chan *T // 对象池管道
	si T       // 对象初始值
}

// NewStructPool 创建一个对象池
func NewStructPool[T any](maxNum int, initStruct T) *StructPool[T] {
	return &StructPool[T]{
		sp: make(chan *T, maxNum), // 对象池管道
		si: initStruct,            // 对象初始化时的默认值
	}
}

// Get 从对象池中获取一个对象
func (bp *StructPool[T]) Get() (s *T) {
	select {

	case s = <-bp.sp:
		// 获取成功

	default:
		// 获取失败，创建一个新对象
		s = new(T)
		*s = bp.si

	}

	return
}

// Put 将当前对象放回对象池
func (bp *StructPool[T]) Put(b *T) bool {
	// 尝试将当前对象放回对象池
	select {

	case bp.sp <- b:
		// 放回成功
		return true

	default:
		// 放回失败
		return false

	}
}
