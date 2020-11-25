package gzseek

import (
	"bufio"
	"io"
)

// tellReader is a bufio.Reader that also tracks its offset
// (number of bytes read) within the underlying data.
type tellReader struct {
	*bufio.Reader
	offset int64
}

func newTellReader(r io.Reader) *tellReader {
	return &tellReader{
		Reader: bufio.NewReader(r),
	}
}

func (r *tellReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.offset += int64(n)
	return
}

func (r *tellReader) Offset() int64 {
	return r.offset
}
