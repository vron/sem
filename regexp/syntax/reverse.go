package syntax

import (
	"fmt"
	"spew"
)

var p = fmt.Println

// Takes a regexp and creates it opposit, i.e the regexp that
// should be used if we want to find the last match by reading the
// input backwards (ie. last in original order, first in new)
func (x *Regexp) Reverse() {
	spew.Config.DisableMethods = true
	//spew.Dump(x)
	// Recursively swap all my children:
	for _, v := range x.Sub {
		if v != nil {
			v.Reverse()
		}
	}
	for _, v := range x.Sub0 {
		if v != nil {
			v.Reverse()
		}
	}
	// Swap meaning of beg and end:
	switch x.Op {
	case OpBeginLine:
		x.Op = OpEndLine
	case OpEndLine:
		x.Op = OpBeginLine
	case OpBeginText:
		x.Op = OpEndText
	case OpEndText:
		x.Op = OpBeginText
	case OpLiteral:
		// Reverse the rune run
		var rt rune
		for i := 0; i < len(x.Rune)/2; i++ {
			rt = x.Rune[i]
			x.Rune[i] = x.Rune[len(x.Rune)-1-i]
			x.Rune[len(x.Rune)-1-i] = rt
		}
	}

	// Also reverse the order of the sub-expressions:
	var rt *Regexp
	for i := 0; i < len(x.Sub)/2; i++ {
		rt = x.Sub[i]
		x.Sub[i] = x.Sub[len(x.Sub)-1-i]
		x.Sub[len(x.Sub)-1-i] = rt
	}
}
