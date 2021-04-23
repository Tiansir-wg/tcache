package main

// 缓存值的抽象
type ByteView struct {
	// 表示真实的缓存值，使用byte切片使其支持任意数据
	b []byte
}

// 实现了lru.Value接口
func (bv ByteView) Len() int {
	return len(bv.b)
}

// 返回对应的切片副本，防止缓存值被修改
func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

func (bv ByteView) String() string {
	return string(bv.b)
}

func cloneBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
