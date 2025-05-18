package eviction

type ByteView struct {
	B []byte
}

func (b ByteView) Len() int {
	return len(b.B)
}

func (b ByteView) String() string {
	return string(b.B)
}

func (b ByteView) ByteSlice() []byte {
	res := cloneByte(b.B)
	return res
}

func cloneByte(b []byte) []byte {
	res := make([]byte, len(b))
	copy(res, b)
	return res
}

func (b ByteView) Empty() bool {
	if b.Len() == 0 {
		return true
	}
	return false
}
