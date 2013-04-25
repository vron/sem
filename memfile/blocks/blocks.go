package blocks

import (
	"io"
)

const (
	B_SIZE     = 8
	DEFAULT_FR = 0.8
	MIN_FR     = 0.2
)

type File struct {
	first *block
	last  *block
	len   int // Total length of the file
	ind   *index
}

type block struct {
	prev  *block
	next  *block
	len   int
	bytes [B_SIZE]byte
}

func New(b []byte) *File {
	f := File{len: len(b)}
	// Create and fill new parts to the default fill degree
	i := 0
	lastb := nil
	for p := 0; p <= len(b)/int(B_SIZE*DEFAULT_FR); p++ {
		b := block{}
		j := 0
		for ; j < int(B_SIZE*DEFAULT_FR) && i < len(b); j++ {
			b.bytes[j] = b[i]
		}
		b.prev = lastb
		b.len = j
		if lastb != nil {
			lastb.next = b
		} else {
			f.first = b
		}
		lastb = b
	}
	f.last = lastb
	return f
}

func (f *File) Reader(start int) io.RuneReader {
}

func (f *File) BackwardsReader(start int) io.RuneReader {
}

func (f *File) Change(start, end int, data []byte) error {
}

func (f *File) Get(start, end int) ([]byte, error) {

}

func (f *File) OffsetLine(ln, start int) (offset int, e error) {

}

func (f *File) OffsetRune(cn, start int) (offset int, e error) {

}

func (f *File) Length() int {
	return f.len
}
