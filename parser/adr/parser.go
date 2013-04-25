package adr

import (
	"bufio"
	"unicode"
)

type Node struct {
	Type        int    // Type of this node
	Left, Right *Node  // Subnodes if this is a composit node
	Val         int    // Numeric value of this node if appropriate
	Reg         string // String of the regexp if such a node
}

// Parse a address greedily from the buffer, if the first non-witespace
// token found is not a valid start of an adr nil, nil is returned
// If a valid adr is parsed, !nil, nil is returned
// If there is an error whilst parsing the addr nil, !nil is returned
// Note that it will consume all spaces in the buffer to check the next char
func Parse(b *bufio.Reader) (*Node, error) {
	// TODO: Change from bytes.Bufer to something nicer... Best is an interface..
	//Check for start, is it the start of an adr
	tl := newLexer(b)
	tl.consume(unicode.IsSpace)
	r, _, e := b.ReadRune()
	if e != nil {
		return nil, nil
	}
	b.UnreadRune()
	if !(unicode.IsDigit(r) ||
		r == '/' ||
		r == '+' ||
		r == '-' ||
		r == '.' ||
		r == '$' ||
		r == '\'' ||
		r == '#' ||
		r == ',' ||
		r == ';') {
		return nil, nil
	}

	fullNode = nil
	l := newLexer(b)
	ok := yyParse(l)
	if ok != 0 {
		// There was an error whils parsing:
		return nil, err
	}
	return fullNode, nil
}
