package encodegen

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type ObjBuffer struct {
	buf Buffer
	lr  io.LimitedReader
	err error // sticky
}

func (b *ObjBuffer) grow(n int)        { b.buf.Grow(n) }
func (b *ObjBuffer) Bytes() []byte     { return b.buf.Bytes() }
func (b *ObjBuffer) next(n int) []byte { return b.buf.Next(n) }

func (b *ObjBuffer) Write(p []byte) {
	b.buf.Write(p)
}

func (b *ObjBuffer) Read(p []byte) {
	if b.err != nil {
		return
	}
	_, b.err = io.ReadFull(&b.buf, p)
}

func (b *ObjBuffer) ReadString() string {
	p := make([]byte, b.ReadPrefix(1))
	b.Read(p)
	return BytesToString(p)
}

func (b *ObjBuffer) WriteString(s string) {
	b.WritePrefix(len(s))
	b.buf.WriteString(s)
}

func (b *ObjBuffer) WriteBool(p bool) {
	if p {
		b.buf.WriteByte(1)
	} else {
		b.buf.WriteByte(0)
	}
}

func (b *ObjBuffer) ReadBool() bool {
	if b.err != nil {
		return false
	}
	c, err := b.buf.ReadByte()
	if err != nil {
		b.err = err
		return false
	}
	if c != 0 && c != 1 {
		b.err = errors.New("invalid boolean")
		return false
	}
	return c == 1
}

func (b *ObjBuffer) WriteByte(c byte) {
	b.buf.WriteByte(c)
}

func (b *ObjBuffer) ReadByte() byte {
	c, err := b.buf.ReadByte()
	if err != nil {
		b.err = err
		return 0
	}
	return c
}

func (b *ObjBuffer) WriteUint64(u uint64) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, u)
	b.buf.Write(buf)
}

func (b *ObjBuffer) ReadUint64() uint64 {
	if b.err != nil {
		return 0
	}
	buf := b.buf.Next(8)
	if len(buf) < 8 {
		b.err = io.EOF
		return 0
	}
	return binary.LittleEndian.Uint64(buf)
}

func (b *ObjBuffer) WritePrefix(i int) {
	b.WriteUint64(uint64(i))
}

func (b *ObjBuffer) ReadPrefix(elemSize int) int {
	n := b.ReadUint64()
	if n > uint64(b.buf.Len()/elemSize) {
		b.err = fmt.Errorf("marshalled object contains invalid length prefix (%v elems x %v bytes/elem > %v bytes left in message)", n, elemSize, b.buf.Len())
		return 0
	}
	return int(n)
}

func (b *ObjBuffer) WritePrefixedBytes(p []byte) {
	b.WritePrefix(len(p))
	b.Write(p)
}

func (b *ObjBuffer) ReadPrefixedBytes() []byte {
	p := make([]byte, b.ReadPrefix(1))
	b.Read(p)
	return p
}

func (b *ObjBuffer) Rewind() {
	b.buf.Seek(-b.buf.Offset())
}

func (b *ObjBuffer) Reset() {
	b.buf.Reset()
	b.err = nil
}

func (b *ObjBuffer) Err() error {
	return b.err
}
