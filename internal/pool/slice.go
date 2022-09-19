package pool

// BytesPoolPutFunc 字节切片回收到复用池的方法类型
type BytesPoolPutFunc func([]byte) bool

// BytesPool 字节切片复用池
type BytesPool struct {
	bpc chan []byte // 字节切片管道
	bsl int         // 字节切片初始长度
	bsc int         // 字节切片初始容量
}

// NewBytesPool 初始化一个字节切片复用池
//
// @params maxNum 字节切片复用池容量
//
// @params initLen 字节切片初始长度
//
// @params initCap 字节切片初始容量
//
// @return 字节切片复用池指针
func NewBytesPool(maxNum int, initLen int, initCap int) *BytesPool {
	bp := &BytesPool{
		bpc: make(chan []byte, maxNum),
		bsl: initLen,
		bsc: initCap,
	}
	if maxNum > 0 {
		for i := 0; i < 10; i++ {
			bp.bpc <- make([]byte, initLen, initCap)
		}
	}
	return bp
}

// Get 从复用池中获取一个字节切片（会将len置为0）
func (bp *BytesPool) Get() (b []byte) {
	select {

	case b = <-bp.bpc:
		// 清空原内容
		b = b[:0]

	default:
		// 创建一个新对象
		if bp.bsc > 0 {
			b = make([]byte, bp.bsl, bp.bsc)
		} else {
			b = make([]byte, bp.bsl)
		}

	}

	return
}

// Put 将字节切片放回复用池（会将len置为0）
func (bp *BytesPool) Put(b []byte) bool {
	// 清空原内容
	b = b[:0]

	// 尝试放回复用池
	select {

	case bp.bpc <- b:
		// 放回复用池成功
		return true

	default:
		// 放回复用池失败
		return false

	}
}
