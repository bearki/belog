package pool

// StructPtrPool 结构体指针复用池
//
// 注意：T为结构体，请勿传入指针类型
type StructPtrPool[T any] struct {
	sppc chan *T // 复用池管道
	siv  T       // 结构体初始值
}

// NewStructPtrPool 结构体指针复用池
//
// 注意：T为结构体，请勿传入指针类型
//
// @params maxNum 复用池容量
//
// @params initStruct 创建新结构体指针时需要使用的初始值
func NewStructPtrPool[T any](maxNum int, initStruct T) *StructPtrPool[T] {
	spp := &StructPtrPool[T]{
		sppc: make(chan *T, maxNum),
		siv:  initStruct,
	}
	if maxNum > 0 {
		for i := 0; i < 10; i++ {
			tmp := initStruct
			spp.sppc <- &tmp
		}
	}
	return spp
}

// Get 从复用池中获取一个结构体指针
func (bp *StructPtrPool[T]) Get() (s *T) {
	select {

	case s = <-bp.sppc:
		// 获取成功

	default:
		// 获取失败，创建一个新结构体指针
		s = new(T)
		*s = bp.siv

	}

	return
}

// Put 将结构体指针放回复用池
func (bp *StructPtrPool[T]) Put(b *T) bool {
	// 尝试将结构体指针放回复用池
	select {

	case bp.sppc <- b:
		// 放回成功
		return true

	default:
		// 放回失败
		return false

	}
}
