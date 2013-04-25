package regexp

import (
	"io"
	"github.com/vron/sem/regexp/syntax"
	"unicode/utf8"
)

func CompileReverse(expr string) (*Regexp, error) {
	return compileReverse(expr, syntax.Perl, false)
}

func compileReverse(expr string, mode syntax.Flags, longest bool) (*Regexp, error) {
	re, err := syntax.Parse(expr, mode)
	if err != nil {
		return nil, err
	}
	maxCap := re.MaxCap()
	capNames := re.CapNames()
	// Reverse it, using our own reversal function that we hope work as it shoudl, but
	// this is just based on a frivolous comment found on the internet.. I don't reallt know
	// the theory of regexps
	re.Reverse()

	re = re.Simplify()
	prog, err := syntax.Compile(re)
	if err != nil {
		return nil, err
	}
	regexp := &Regexp{
		expr:        expr,
		prog:        prog,
		numSubexp:   maxCap,
		subexpNames: capNames,
		cond:        prog.StartCond(),
		longest:     longest,
	}
	regexp.prefix, regexp.prefixComplete = prog.Prefix()
	if regexp.prefix != "" {
		// TODO(rsc): Remove this allocation by adding
		// IndexString to package bytes.
		regexp.PrefixBytes = []byte(regexp.prefix)
		regexp.prefixRune, _ = utf8.DecodeRuneInString(regexp.prefix)
	}
	return regexp, nil
}

// inputReadSeeker scans a ReadSeeker
type inputReadSeeker struct {
	r     io.ReadSeeker
	buf   []byte
	start int
}

// TODO: Is this one actually ever used to read anything but the next
// rune, and if so is the assumption below of setting the position t0
// after that read rune instead of next correct
func (i *inputReadSeeker) Step(pos int) (rune, int) {
	pos = pos + i.start
	ol, _ := i.r.Seek(0, 1)
	// Get the length of the thing
	_, e := i.r.Seek(int64(pos), 0)
	if e != nil {
		i.r.Seek(ol, 0)
		return endOfText, 0
	}

	// So we are now at the correct place, try to decode a rune
	if len(i.buf) < 4 {
		if cap(i.buf) < 4 {
			i.buf = make([]byte, 4)
		} else {
			// TODO: Is this how to extend a slice to its maximum capacity?
			i.buf = i.buf[:4]
		}
	}
	n, _ := i.r.Read(i.buf)
	if n > 0 {
		if i.buf[0] < utf8.RuneSelf {
			i.r.Seek(ol+1, 0)
			return rune(i.buf[0]), 1
		}
		r, n := utf8.DecodeRune(i.buf)
		i.r.Seek(ol+int64(n), 0)
		return r, n
	}
	i.r.Seek(ol, 0)
	return endOfText, 0
}

func (i *inputReadSeeker) CanCheckPrefix() bool {
	return true
}

// Store where we are, look forward for match and seek back..
func (i *inputReadSeeker) HasPrefix(re *Regexp) bool {
	//log.Println("Has prefix:", string(re.PrefixBytes))
	ol, _ := i.r.Seek(0, 1)
	defer i.r.Seek(ol, 0)

	if len(i.buf) < 512 {
		if cap(i.buf) < 512 {
			i.buf = make([]byte, 512)
		} else {
			// TODO: Is this how to extend a slice to its maximum capacity?
			i.buf = i.buf[:cap(i.buf)]
		}
	}

	n, e := io.ReadAtLeast(i.r, i.buf, len(re.PrefixBytes))
	if e != nil || n < len(re.PrefixBytes) {
		return false
	}
	for j := range i.buf {
		if i.buf[j] != re.PrefixBytes[j] {
			return false
		}
	}
	return true
}

// Find first index of re.PrefixBytes in it from pos
func (i *inputReadSeeker) Index(re *Regexp, pos int) int {
	pos = pos + i.start
	//log.Println("Index:", string(re.PrefixBytes), pos)
	ol, _ := i.r.Seek(0, 1)
	defer i.r.Seek(ol, 0)

	_, e1 := i.r.Seek(int64(pos), 0)
	if e1 != nil {
		return -1
	}

	if len(i.buf) < 512 {
		if cap(i.buf) < 512 {
			i.buf = make([]byte, 512)
		} else {
			// TODO: Is this how to extend a slice to its maximum capacity?
			i.buf = i.buf[:cap(i.buf)]
		}
	}

	ipb := 0 // Index in PrefixBytes
	ib := 0  // Index in buffer
	lb := 0  // Number of bytes in buffer
	ii := 0  // Number of bytes read before current buffer
	var e error = nil
	for {
		if ib >= lb {
			if e != nil { // We have handled all input, handle the error
				return -1
			}
			lb, e = i.r.Read(i.buf) // Read in more
			ii += ib
			ib = 0
			continue
		}
		if re.PrefixBytes[ipb] != i.buf[ib] {
			ipb = 0
			ib++
			continue
		}
		// So they were equal, advance and check if it is a matches
		ib++
		ipb++
		if ipb >= len(re.PrefixBytes) {
			return ii + ib - len(re.PrefixBytes)
		}
	}
	return -1
}

// As I understand it this thing should check the rune before and the one
// after the given location? IS THIS THING EVER USED?!
func (i *inputReadSeeker) Context(pos int) syntax.EmptyOp {
	pos = pos + i.start
	//log.Println("Context: ", pos)
	//log.Println(pos)
	ol, _ := i.r.Seek(0, 1)
	defer i.r.Seek(ol, 0)

	r1, r2 := endOfText, endOfText

	if len(i.buf) != 8 {
		if cap(i.buf) < 8 {
			i.buf = make([]byte, 8)
		} else {
			// TODO: Is this how to extend a slice to its maximum capacity?
			i.buf = i.buf[:8]
		}
	}

	// Try to read out the surrounding 8 bytes but take care so that
	// if we are to early we read fewer before
	posa := pos - 4
	if pos < 4 {
		posa = 0
	}
	_, e := i.r.Seek(int64(posa), 0)
	if e != nil {
		return syntax.EmptyOpContext(r1, r2)
	}

	n, e := i.r.Read(i.buf)

	if pos > 0 && n > 0 {
		r1, _ = utf8.DecodeLastRune(i.buf[:pos-posa])
	}
	if n > pos-posa {
		r2, _ = utf8.DecodeRune(i.buf[pos-posa:])
	}
	return syntax.EmptyOpContext(r1, r2)
}




/* BELOW IS EXTRA PUBLIC FUNCTIONS THAT CAN BE CALLED ON ANY Input */
func MatchInput(pattern string, r Input) (matched bool, error error) {
	re, err := Compile(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchInput(r), nil
}

func (re *Regexp) FindInputSubmatchIndex(r Input) []int {
	return re.pad(re.doExecuteInput(r, 0, re.prog.NumCap))
}

func (re *Regexp) FindInputIndex(r Input) (loc []int) {
	a := re.doExecuteInput(r, 0, 2)
	if a == nil {
		return nil
	}
	return a[0:2]
}
func (re *Regexp) MatchInput(r Input) bool {
	return re.doExecuteInput(r, 0, 0) != nil
}

// This should be merged into ordinary doExecute
func (re *Regexp) doExecuteInput(i Input, pos int, ncap int) []int {
	m := re.get()
	m.init(ncap)
	if !m.match(i, pos) {
		re.put(m)
		return nil
	}
	if ncap == 0 {
		re.put(m)
		return empty // empty but not nil
	}
	cap := make([]int, ncap)
	copy(cap, m.matchcap)
	re.put(m)
	return cap
}



/* Below for ReaderSeeker, to be removed! */



// MatchReader checks whether a textual regular expression matches the text
// read by the RuneReader.  More complicated queries need to use Compile and
// the full Regexp interface.
func MatchReadSeeker(pattern string, r io.ReadSeeker) (matched bool, error error) {
	re, err := Compile(pattern)
	if err != nil {
		return false, err
	}
	return re.MatchReadSeeker(r), nil
}

func (re *Regexp) FindReadSeekerSubmatchIndex(r io.ReadSeeker) []int {
	return re.pad(re.doExecuteSeeker(r, 0, re.prog.NumCap))
}

func (re *Regexp) FindReadSeekerIndex(r io.ReadSeeker) (loc []int) {
	a := re.doExecuteSeeker(r, 0, 2)
	if a == nil {
		return nil
	}
	return a[0:2]
}

// MatchReader returns whether the Regexp matches the text read by the
// RuneReader.  The return value is a boolean: true for match, false for no
// match.
func (re *Regexp) MatchReadSeeker(r io.ReadSeeker) bool {
	return re.doExecuteSeeker(r, 0, 0) != nil
}

// This should be merged into ordinary doExecute
func (re *Regexp) doExecuteSeeker(r io.ReadSeeker, pos int, ncap int) []int {
	m := re.get()
	var i Input
	if r != nil {
		ita, _ := r.Seek(0, 1)
		i = &inputReadSeeker{r: r, start: int(ita)} // TODO: This one should be made cached like the others
	}
	m.init(ncap)
	if !m.match(i, pos) {
		re.put(m)
		return nil
	}
	if ncap == 0 {
		re.put(m)
		return empty // empty but not nil
	}
	cap := make([]int, ncap)
	copy(cap, m.matchcap)
	re.put(m)
	return cap
}
