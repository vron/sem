package adr

import (
	"bufio"
	"unicode"
	"strconv"
	"strings"
)

type Node struct {
	Type        int    // Type of this node
	Left, Right *Node  // Subnodes if this is a composit node
	Val         int    // Numeric value of this node if appropriate
	Reg         string // String of the regexp if such a node
}

// TODO: We should here simplify to idioms etc and not print more brackets than needed
func (n *Node) String() string {
	return n.recString(0)
}

// Recursively travers left depth first and print this thing!
func (n* Node) recString(depth int) string {
	switch n.Type {
	case HASH:
		return "#" + strconv.Itoa(n.Val)
	case NUMBER:
		return strconv.Itoa(n.Val)
	case REG:
		return "/" + strings.Replace(n.Reg, "/", "\\/", -1) + "/"
	case DOLLAR:
		return "$"
	case DOT:
		return "."
	case ADRMARK:
		return "'"
	case PLUS:
		if depth != 0 {
			return "(" + n.Left.recString(depth+1) + "+" + n.Right.recString(depth+1) + ")"
		} else {
			return n.Left.recString(depth+1) + "+" + n.Right.recString(depth+1)
		}
	case MINUS:
		if depth != 0 {
			return "(" + n.Left.recString(depth+1) + "-" + n.Right.recString(depth+1) + ")"
		} else {
			return n.Left.recString(depth+1) + "-" + n.Right.recString(depth+1)
		}
	case COMMA:
		if depth != 0 {
			return "(" + n.Left.recString(depth+1) + "," + n.Right.recString(depth+1) + ")"
		} else {
			return n.Left.recString(depth+1) + "," + n.Right.recString(depth+1)
		}
	case SEMI:
		if depth != 0 {
			return "(" + n.Left.recString(depth+1) + ";" + n.Right.recString(depth+1) + ")"
		} else {
			return n.Left.recString(depth+1) + ";" + n.Right.recString(depth+1)
		}
	default:
		panic("Unimplemented adr item")
	}
	return ""
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
