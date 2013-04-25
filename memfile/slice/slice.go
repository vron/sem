package slice

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"github.com/vron/sem/memfile"
	"unicode/utf8"
)

type File struct {
	b   []byte
	pos int
}

func NewReader(in io.Reader) *File {
	b, e := ioutil.ReadAll(in)
	if e != nil {
		panic(e)
	}
	return &File{b: b}
}

func New(b []byte) *File {
	if b == nil {
		return &File{b: []byte{}}
	}
	return &File{b: b}
}

func (f *File) BackwardsReader(start int) io.RuneReader {
	return &br{start, f}
}

type br struct {
	ind int
	f   *File
}

func (b *br) ReadRune() (rune, int, error) {
	// Try to decode the last rune in this part of the string
	if b.ind == 0 {
		return 0, 0, io.EOF
	}
	r, n := utf8.DecodeLastRune(b.f.b[:b.ind])
	b.ind -= n
	if r == utf8.RuneError && n == 1 {
		return r, n, errors.New("Invalid")
	}
	return r, n, nil
}

func (f *File) Read(b []byte) (l int, e error) {
	// copy out as many bytes as possible
	l = len(b)
	if len(f.b)-f.pos < l {
		l = len(f.b) - f.pos
		e = io.EOF
	}
	copy(b, f.b[f.pos:])
	f.pos += l
	return l, e
}

func (f *File) Length() int {
	return len(f.b)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	off := int(offset)
	if whence == 1 {
		off = f.pos + off
	} else if whence == 2 {
		off = len(f.b) + off
	} else if whence != 0 {
		return 0, errors.New("Unknown whence")
	}
	if off < 0 {
		f.pos = 0
		return 0, memfile.OutOfBounds
	}
	if off > len(f.b) {
		return int64(len(f.b)), errors.New("Position after end")
	}
	f.pos = off
	return int64(off), nil
}
func (f *File) Change(start, end int, data []byte) error {
	if start < 0 || start > len(f.b) || end < start || end > len(f.b) {
		return memfile.OutOfBounds
	}
	// Fastpath if they are the same size
	if len(data) == end-start {
		copy(f.b[start:end], data)
		return nil
	}
	temp := append(data, f.b[end:]...)
	f.b = append(f.b[:start], temp...)
	return nil
}

func (f *File) Get(start, end int) ([]byte, error) {
	if start < 0 || start > len(f.b) || end < start || end > len(f.b) {
		return nil, memfile.OutOfBounds
	}
	return f.b[start:end], nil
}

func (f *File) OffsetLine(ln, start int) (offset int, e error) {
	if start < 0 || start > len(f.b) {
		return 0, memfile.OutOfBounds
	}
	if ln == 0 {
		i := bytes.LastIndex(f.b[:start], []byte("\n"))
		return i + 1, nil
	}
	if ln < 0 {
		i := 0
		return bytes.LastIndexFunc(f.b[:start], func(r rune) bool {
			if r == '\n' {
				if i == ln {
					return true
				}
				i--
			}
			return false
		}) + 1, nil
	}
	i := 0
	va := bytes.IndexFunc(f.b[start:], func(r rune) bool {
		if r == '\n' {
			i++
			if i == ln {
				return true
			}
		}
		return false
	})
	if va != -1 {
		return va + start + 1, nil
	}
	return len(f.b), nil
}

// Got to hate utf8 for making it complicated... Guess there is not much to do..
func (f *File) OffsetRune(cn, start int) (offset int, e error) {
	if start < 0 || start > len(f.b) {
		return 0, memfile.OutOfBounds
	}
	if cn <= 0 {
		ind := start
		for ; cn < 0; cn++ {
			if ind < 1 {
				return 0, nil
			}
			_, s := utf8.DecodeLastRune(f.b[:ind])
			ind -= s
		}
		return ind, nil
	}
	ind := start
	for ; cn > 0; cn-- {
		if ind >= len(f.b)-1 {
			return len(f.b), nil
		}
		_, s := utf8.DecodeRune(f.b[ind:])
		ind += s
	}
	return ind, nil
}
