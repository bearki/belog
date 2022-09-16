package pool

// BytesPool 字节流对象池
type BytesPool struct {
	bp chan []byte // 字节流对象管道
	bl int         // 字节流对象初始长度
	bc int         // 字节流对象初始容量
}

// NewBytesPool 初始化一个字节流对象池
func NewBytesPool(maxNum int, len int, cap int) (bp *BytesPool) {
	return &BytesPool{
		bp: make(chan []byte, maxNum),
		bl: len,
		bc: cap,
	}
}

// Get 从对象池中获取一个对象
func (bp *BytesPool) Get() (b []byte) {
	select {

	case b = <-bp.bp:
		// 清空原内容
		b = b[:0]

	default:
		// 创建一个新对象
		if bp.bc > 0 {
			b = make([]byte, bp.bl, bp.bc)
		} else {
			b = make([]byte, bp.bl)
		}

	}

	return
}

// Put 将对象放回对象池
func (bp *BytesPool) Put(b []byte) bool {
	// 清空原内容
	b = b[:0]

	// 尝试放回对象池
	select {

	case bp.bp <- b:
		// 放回对象池成功
		return true

	default:
		// 放回对象池失败
		return false

	}
}
