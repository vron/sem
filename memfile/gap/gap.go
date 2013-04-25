package gap

import (
	"io"
	"io/ioutil"
	"errors"
	"bytes"
	"sams/memfile"
	"unicode/utf8"
	"sams/regexp"
	"sams/regexp/syntax"
)

const (
	DEFAULT_SIZE = 1024
	GROW_DENOM  = 5
)

type File struct {
	b []byte
	gapStart, gapEnd int // Start and end position of the gap
	pos int				 // Current position (in stream coord, not gap coord)
}


func NewReader(in io.Reader) *File {
	b, e := ioutil.ReadAll(in)
	if e != nil {
		panic(e)
	}
	return New(b)
}


// When we create a new file we start with the gap at
// the beginning of the file since most operations will
// start here!
func New(b []byte) *File {
	if b == nil {
		return &File{
			b: make([]byte, DEFAULT_SIZE, DEFAULT_SIZE),
			gapStart: 0,
			gapEnd: DEFAULT_SIZE,
			pos: 0,
		}
	}
	// So we start with no gap (all other fields zero)
	return &File{b: b}
}

func (f *File) BackwardsReader(start int) io.RuneReader {
	return &br{start, f}
}

type br struct {
	ind int
	f *File
}

// TODO: (skip the gap)
func (b *br) ReadRune() (rune, int, error) {
	// Try to decode the last rune in this part of the string
	if b.ind == 0 {
		return 0, 0, io.EOF
	}
	r, n :=  utf8.DecodeLastRune(b.f.b[:b.ind])
	b.ind -= n
	if r == utf8.RuneError && n == 1 {
		return r, n, errors.New("Invalid")
	}
	return r, n, nil
}

func (f *File) Read(b []byte) (l int, e error) {
	// copy out as many bytes as possible, always
	// breaking around the gap so we can use the original data
	
	if f.pos >= f.gapStart {
		// We can simply copy out as many as possible after the gap
		l = len(b)
		if f.Length()-f.pos < l {
			l = f.Length()-f.pos
			e = io.EOF
		}
		copy(b,f.b[f.pos+(f.gapEnd-f.gapStart):])
		f.pos += l
		return l, e
	}
		// So we are before the gap, check if requested length will go into
		// the buffer if so cut of
		l = len(b)
		if f.gapStart-f.pos < l {
			l = f.gapStart-f.pos
		}
		copy(b, f.b[f.pos:f.pos+l])
		f.pos += l
		return l, e
}

func (f *File) Length() int {
	return len(f.b)-(f.gapEnd-f.gapStart)
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	off := int(offset)
	if whence == 1 {
		off = f.pos+off
	} else if whence == 2 {
		off = len(f.b)+off
	} else if whence != 0 {
		return 0, errors.New("Unknown whence")
	}
	if off < 0 {
		f.pos = 0
		return 0, memfile.OutOfBounds
	}
	if off > len(f.b) {
		return  int64(len(f.b)), errors.New("Position after end")
	}
	f.pos = off
	return int64(off), nil
}
func (f *File) Change(start, end int, data []byte) error {
	if start < 0 || start > len(f.b) || end < start || end > len(f.b) {
		return memfile.OutOfBounds
	}
	
	// Move the gap so start-1 position is at start of gap
	if start > f.gapStart {
		copy(f.b[f.gapStart:], f.b[f.gapEnd:start+(f.gapEnd-f.gapStart)])
		f.gapEnd += start-f.gapStart
		f.gapStart = start  
	} else if start < f.gapStart {
		copy(f.b[f.gapEnd-(f.gapStart-start):f.gapEnd], f.b[start:f.gapStart])
		f.gapEnd += start-f.gapStart
		f.gapStart = start
	}
	
	// So to say remove the replaced part
	f.gapEnd += end-start

	// Check if there is enough space to fill or if extension is needed:
	if (f.gapEnd-f.gapStart) < len(data) {
		reqSize := f.Length()-(end-start)+len(data)
		newSize := reqSize + 1 + reqSize/GROW_DENOM
		tb := make([]byte, newSize, newSize)
		copy(tb, f.b[:f.gapStart])
		copy(tb[f.gapEnd+(newSize-len(f.b)):], f.b[f.gapEnd:])
		f.gapEnd += newSize-len(f.b)
		f.b = tb
	}

	// Copy in the data and update the gap
	copy(f.b[start: start+len(data)], data)
	f.gapStart += len(data)

	return nil
}

func (f *File) Get(start, end int) ([]byte, error) {
	if start < 0 || start > len(f.b) || end < start || end > len(f.b) {
		return nil, memfile.OutOfBounds
	}
	
	// If the requested range is across the gap we need to create a copy
	if end < f.gapStart {
		return f.b[start:end], nil
	} else if start > f.gapStart {
		return f.b[start+(f.gapEnd-f.gapStart):end+(f.gapEnd-f.gapStart)], nil
	}
	b := make([]byte, end-start, end-start)
	//println(start, f.gapStart, len(f.b))
	copy(b, f.b[start:f.gapStart])
	copy(b[f.gapStart-start:], f.b[f.gapEnd:])
	return b, nil
}

// TODO:
func (f *File) OffsetLine(ln, start int) (offset int, e error) {
	if start < 0 || start > len(f.b) {
		return 0, memfile.OutOfBounds
	}
	if ln == 0 {
		i := bytes.LastIndex(f.b[:start], []byte("\n"))
		return i+1, nil
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
		})+1, nil
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
		return va+start+1, nil
	}
	return len(f.b), nil
}

// TODO
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

// We assume that all operations are always performed on full runes, meaning
// that a rune cannot be split across the gap!

// Implementation of the regexp interface:
var endOfText rune = -1

func (f *File) Step(pos int) (r rune, width int) {
	pos += f.pos
	if pos < f.Length() {
		if pos >= f.gapStart {
			pos += f.gapEnd-f.gapStart
		}
		c := f.b[pos]
		if c < utf8.RuneSelf {
			return rune(c), 1
		}
		return utf8.DecodeRune(f.b[pos:])
	}
	return endOfText , 0
}
func (f *File) CanCheckPrefix() bool{
	return true
}
func (f *File) HasPrefix(re *regexp.Regexp) bool{
	if f.pos + len(re.PrefixBytes) <= f.gapStart {
		return bytes.HasPrefix(f.b[f.pos:], re.PrefixBytes)
	} else if f.pos > f.gapStart {
		return bytes.HasPrefix(f.b[f.pos+(f.gapEnd-f.gapStart):], re.PrefixBytes)
	}
	// Otherwise we create a copy
	// TODO: Avoid this garbage!
	b, e := f.Get(f.pos, f.pos+len(re.PrefixBytes))
	if e != nil {
		return false
	}
	return bytes.HasPrefix(b, re.PrefixBytes)
}
func (f *File) Index(re *regexp.Regexp, pos int) int {
	pos += f.pos
	if pos >= f.gapStart {
		return bytes.Index(f.b[pos+(f.gapEnd-f.gapStart):], re.PrefixBytes)
	}
	// So first check before the gap and then after the gap
	ind := bytes.Index(f.b[pos:], re.PrefixBytes)
	if ind >= 0 {
		return ind
	}
	aft := bytes.Index(f.b[f.gapEnd:], re.PrefixBytes)
	if aft < 0 {
		return aft
	}
	return aft + f.gapStart-pos
}
func (f *File) Context(pos int) syntax.EmptyOp{
	pos += f.pos
	r1, r2 := endOfText, endOfText
	if pos > 0 && pos <= f.Length() {
		if pos < f.gapStart {
			r1, _ = utf8.DecodeLastRune(f.b[:pos])
		} else {
			r1, _ = utf8.DecodeLastRune(f.b[:pos+(f.gapEnd-f.gapStart)])
		}
	}
	if pos < f.Length() {
		if pos < f.gapStart {
			r2, _ = utf8.DecodeRune(f.b[pos:])
		} else {
			r2, _ = utf8.DecodeRune(f.b[pos+(f.gapEnd-f.gapStart):])
		}
	}
	return syntax.EmptyOpContext(r1, r2)
}