package bufferutils

type ByteReader struct {
	Buf []byte
	Ptr int
}

func ReadBytes(buf []byte) ByteReader {
	return ByteReader{ buf, 0 }
}

func (b *ByteReader) Reset() {
	b.Seek(0)
}

func (b *ByteReader) Seek(offset int) {
	_ = b.Buf[offset] // Bounds check
	b.Ptr = offset
}

func (b *ByteReader) CopyNext(target []byte) {
	copy(target, b.Next(len(target)))
}

func (b *ByteReader) Next(count int) []byte {
	buf := b.Peek(count)
	b.Skip(count)
	return buf
}

func (b *ByteReader) Peek(count int) []byte {
	_ = b.Buf[b.Ptr + count - 1]
	return b.Buf[b.Ptr:b.Ptr + count]
}

func (b *ByteReader) Skip(count int) {
	b.Ptr += count
}

func (b *ByteReader) Uint8() uint8 {
	i := b.Buf[b.Ptr] // byte = uint8
	b.Ptr++
	return i
}

func (b *ByteReader) Uint16() uint16 {
	return bin.Uint16(b.Next(2))
}

func (b *ByteReader) Uint32() uint32 {
	return bin.Uint32(b.Next(4))
}

func (b *ByteReader) Uint64() uint64 {
	return bin.Uint64(b.Next(8))
}

func (b *ByteReader) Int8() int8 {
	return int8(b.Uint8())
}

func (b *ByteReader) Int16() int16 {
	return int16(b.Uint16())
}

func (b *ByteReader) Int32() int32 {
	return int32(b.Uint32())
}

func (b *ByteReader) Int64() int64 {
	return int64(b.Uint64())
}
