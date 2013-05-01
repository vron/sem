/*
Package for translating to and from string literals,
that is translating for example "\t" to a string of just
the tab characted. See definition of string literal
at golang.org to understand syntax and what is supporter..
*/
package stresc

import (
	"errors"
	"unicode/utf8"
	"unicode"
)

// Take a string to escape and return a string
func StrByte(s string) ([]byte, error) {
	return Escape([]byte(s))
}

// Take a byte and escape everything, note that a copy is created
func Escape(b []byte)  ([]byte, error) {
	out := make([]byte, 0, len(b))
	// Loop through each unicode character
	offset := 0
	for offset < len(b) {
		r, w, e := next(b, offset)
		offset += w
		if e != nil {
			return nil, e
		}
		
		// So now we have the next rune and its length
		switch {
		case r=='\\':
			// Go into escaping mode until we are done
			r, w, e = next(b, offset)
			offset += w
			if e != nil {
				return nil, e
			}
			var i, base, max uint32
			switch r {
			case 'a':
				out = append(out, 7)
			case 'b':
				out = append(out, 8)
			case 'f':
				out = append(out, 12)
			case 'n':
				out = append(out, 10)
			case 'r':
				out = append(out, 13)
			case 't':
				out = append(out, 9)
			case 'v':
				out = append(out, 11)
			case '\\':
				out = append(out, 92)
			case '"':
				out = append(out, 34)
			case '\'':
				out = append(out, 39)
 			case '0', '1', '2', '3', '4', '5', '6', '7':
 				i, base, max = 3, 8, 255
 			case 'x':
   				i, base, max = 2, 16, 255
   			case 'u':
   				i, base, max = 4, 16, unicode.MaxRune
   			case 'U':
   				i, base, max = 8, 16, unicode.MaxRune
			default:
				return nil, errors.New("unexpected escape sequence start")
			}
			// Keep reading exactly as many as required and accumulate into the unicode point
			var x uint32
			for ;i > 0; i-- {
				r, w, e = next(b, offset)
				offset += w
				if e != nil {
					return nil, e
				}
				d := uint32(digitVal(r))
				if d > base {
					return nil, errors.New("illegal character in escape sequence")
				}
				x = x*base+d
			}
			if x > max || 0xd800 <= x && x < 0xe000 {
				return nil, errors.New("escape sequence is invalid Unicode point")
			}
			
			if base > 0 { // Don't add if we allready have
				bt := make([]byte,4)
				n := utf8.EncodeRune(bt, rune(x))
				out = append(out,bt[:n]...)
			}
		default:
			out = append(out, b[offset-w:offset]...)
		}
	}
	return out, nil
}

func next(b []byte, offset int) (rune, int, error) {
	if offset > len(b)-1 {
		return 0,0,errors.New("need more characters in string literal to progress")
	}
	r, w := rune(b[offset]),1
		switch {
		case r == 0:
			return 0, 0, errors.New("illegal character NUL")
		case r >= 0x80:
			// Multi byte rune, decode it
			r, w = utf8.DecodeRune(b[offset:])
			if r == utf8.RuneError && w == 1 {
				return 0, 0, errors.New("illegal multibyte utf8 encoding")
			}
		}
	return r, w, nil
}

func digitVal(ch rune) int {
	switch {
    case '0' <= ch && ch <= '9':
   			return int(ch - '0')
   		case 'a' <= ch && ch <= 'f':
   			return int(ch - 'a' + 10)
   		case 'A' <= ch && ch <= 'F':
   			return int(ch - 'A' + 10)
   		}
   		return 16 // larger than any legal digit val
 }