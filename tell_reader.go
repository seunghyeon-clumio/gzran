package gzran

import (
	"bufio"
	"io"
)

// tellReader is a bufio.Reader that also tracks its offset
// (number of bytes read) within the underlying data.
// tellReader implements flate.Reader.
type tellReader struct {
	r      *bufio.Reader
	offset int64
}

func newTellReader(r io.Reader) *tellReader {
	return &tellReader{
		r: bufio.NewReader(r),
	}
}

func (r *tellReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	r.offset += int64(n)
	return
}

func (r *tellReader) ReadByte() (byte, error) {
	b, err := r.r.ReadByte()
	if err == nil {
		r.offset++
	}
	return b, err
}

func (r *tellReader) Offset() int64 {
	return r.offset
}
