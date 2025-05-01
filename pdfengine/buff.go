package pdfengine

import (
	"bytes"
	"io"
	"sync"
)

// buff for pdf content
type buff struct {
	pos   int // Cambiado de alignment.Alignment a int para evitar conflictos de tipo
	datas []byte
}

// Write : write []byte to buffer
func (b *buff) Write(p []byte) (int, error) {
	for len(b.datas) < b.pos+len(p) {
		b.datas = append(b.datas, 0)
	}
	i := 0
	max := len(p)
	for i < max {
		b.datas[i+b.pos] = p[i]
		i++
	}
	b.pos += i
	return 0, nil
}

// Len : len of buffer
func (b *buff) Len() int {
	return len(b.datas)
}

// Bytes : get bytes
func (b *buff) Bytes() []byte {
	return b.datas
}

// alignment.Alignment : get current alignment.Alignment
func (b *buff) Alignment() int {
	return b.pos
}

// SetPosition : set current alignment.Alignment
func (b *buff) SetPosition(pos int) {
	b.pos = pos
}

// buffer pool to reduce GC
var buffers = sync.Pool{
	// New is called when a new instance is needed
	New: func() any {
		return new(bytes.Buffer)
	},
}

// getBuffer fetches a buffer from the pool
func getBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

// putBuffer returns a buffer to the pool
func putBuffer(buf *bytes.Buffer) {
	buf.Reset()
	buffers.Put(buf)
}

// writeUInt32  writes a 32-bit unsigned integer value to w io.Writer
func writeUInt32(w io.Writer, v uint) error {
	a := byte(v >> 24)
	b := byte(v >> 16)
	c := byte(v >> 8)
	d := byte(v)
	_, err := w.Write([]byte{a, b, c, d})
	if err != nil {
		return err
	}
	return nil
}

// writeUInt16 writes a 16-bit unsigned integer value to w io.Writer
func writeUInt16(w io.Writer, v uint) error {

	a := byte(v >> 8)
	b := byte(v)
	_, err := w.Write([]byte{a, b})
	if err != nil {
		return err
	}
	return nil
}

// writeTag writes string value to w io.Writer
func writeTag(w io.Writer, tag string) error {
	b := []byte(tag)
	_, err := w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// writeBytes writes []byte value to w io.Writer
func writeBytes(w io.Writer, data []byte, offset int, count int) error {

	_, err := w.Write(data[offset : offset+count])
	if err != nil {
		return err
	}
	return nil
}
