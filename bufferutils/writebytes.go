package bufferutils

type ByteWriter struct {
	Buf []byte
	Ptr int
}

func WriteBytes(buf []byte) ByteWriter {
	return ByteWriter{ buf, 0 }
}

func (b *ByteWriter) Reset() {
	b.Seek(0)
}

func (b *ByteWriter) Seek(offset int) {
	// Bounds check
	_ = b.Buf[offset]
	b.Ptr = offset
}

func (b *ByteWriter) WriteNext(target []byte) {
	copy(b.Next(len(target)), target)
}

func (b *ByteWriter) Next(count int) []byte {
	buf := b.Peek(count)
	b.Skip(count)
	return buf
}

func (b *ByteWriter) Peek(count int) []byte {
	_ = b.Buf[b.Ptr + count - 1]
	return b.Buf[b.Ptr:b.Ptr+count]
}

func (b *ByteWriter) Skip(count int) {
	b.Ptr += count
}

func (b *ByteWriter) Uint8(i uint8) {
	b.Buf[b.Ptr] = i
	b.Ptr++
}

func (b *ByteWriter) Uint16(i uint16) {
	bin.PutUint16(b.Next(2), i)
}

func (b *ByteWriter) Uint32(i uint32) {
	bin.PutUint32(b.Next(4), i)
}

func (b *ByteWriter) Uint64(i uint64) {
	bin.PutUint64(b.Next(8), i)
}

func (b *ByteWriter) Int8(i int8) {
	b.Uint8(uint8(i))
}

func (b *ByteWriter) Int16(i int16) {
	b.Uint16(uint16(i))
}

func (b *ByteWriter) Int32(i uint32) {
	b.Uint32(uint32(i))
}

func (b *ByteWriter) Int64(i uint64) {
	b.Uint64(uint64(i))
}
