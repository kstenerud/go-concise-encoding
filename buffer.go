package concise_encoding

type buffer struct {
	bytes             []byte
	lastAllocatedSize int
}

func (this *buffer) Bytes() []byte {
	return this.bytes
}

func (this *buffer) Grow(byteCount int) {
	length := len(this.bytes)
	growAmount := cap(this.bytes)
	if byteCount > growAmount {
		if byteCount > minBufferCap {
			growAmount = byteCount
		} else {
			growAmount = minBufferCap
		}
	}
	newCap := cap(this.bytes) + growAmount
	newBytes := make([]byte, length+byteCount, newCap)
	oldBytes := this.bytes
	copy(newBytes, oldBytes)
	this.bytes = newBytes
}

func (this *buffer) Allocate(byteCount int) []byte {
	length := len(this.bytes)
	if cap(this.bytes)-length < byteCount {
		this.Grow(byteCount)
	} else {
		this.bytes = this.bytes[:length+byteCount]
	}
	this.lastAllocatedSize = byteCount
	return this.bytes[length:]
}

func (this *buffer) TrimUnused(usedByteCount int) {
	unused := this.lastAllocatedSize - usedByteCount
	this.bytes = this.bytes[:len(this.bytes)-unused]
}

func (this *buffer) Shrink(amount int) {
	this.bytes = this.bytes[:len(this.bytes)-amount]
}
