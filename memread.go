package memread

import (
	"errors"
	"io"
	"math"
	"unsafe"
)

type Reader struct {
	off uint64
	i   int64 // Current index within memory
}

func NewReader(at uint64) *Reader {
	return &Reader{at, 0}
}

func (r *Reader) eof() bool {
	return r.off+uint64(r.i) >= math.MaxUint64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	c := cap(p)
	if c > 0 {
		for n = 0; n < c; n++ {
			if r.eof() {
				err = io.EOF
				break
			}
			p[n] = *(*byte)(unsafe.Pointer(uintptr(r.off + uint64(r.i))))
			r.i++
		}
	}
	return
}

func (r *Reader) ReadAt(p []byte, off int64) (n int, err error) {
	if off < 0 {
		return 0, errors.New("Memread.Reader.ReadAt: negative offset")
	}
	c := cap(p)
	if c > 0 {
		for n = 0; n < c; n++ {
			if r.off+uint64(n)+uint64(off) >= math.MaxUint64 {
				err = io.EOF
				break
			}
			p[n] = *(*byte)(unsafe.Pointer(uintptr(r.off + uint64(off) + uint64(n))))
		}
	}
	return
}
