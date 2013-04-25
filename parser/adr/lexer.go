package adr

import (
	"bufio"
	"errors"
	"strconv"
	"unicode"
)

type lexer struct {
	in *bufio.Reader
}

func newLexer(b *bufio.Reader) *lexer {
	l := &lexer{in: b}

	return l
}

func (l *lexer) Lex(lval *yySymType) int {
	// Try to peak and read until we dont see adr chars
	l.consume(unicode.IsSpace)

	pe, _, e := l.in.ReadRune()
	if e != nil {
		pe = 'x'
	} else {
		l.in.UnreadRune()
	}

	switch {
	case pe == '#':
		l.in.ReadRune()
		return HASH
	case pe == '$':
		l.in.ReadRune()
		return DOLLAR
	case pe == '.':
		l.in.ReadRune()
		return DOT
	case pe == ',':
		l.in.ReadRune()
		return COMMA
	case pe == ';':
		l.in.ReadRune()
		return SEMI
	case pe == '+':
		l.in.ReadRune()
		return PLUS
	case pe == '-':
		l.in.ReadRune()
		return MINUS
	case pe == '(':
		l.in.ReadRune()
		return pBSTART
	case pe == ')':
		l.in.ReadRune()
		return pBEND
	case pe == '\'':
		l.in.ReadRune()
		return ADRMARK
	case unicode.IsDigit(pe):
		r, e := l.consume(unicode.IsDigit)
		if e != nil {
			l.Error(e.Error())
			return 0
		}
		i, _ := strconv.Atoi(string(r))
		lval.val = i
		return NUMBER
	case pe == '/':
		l.in.ReadRune()
		rr, e := l.consume(func(r rune) bool { return r != '/' })
		if e != nil {
			l.Error(e.Error())
			return 0
		}
		l.in.ReadRune()
		lval.reg = string(rr)
		return REG
	default:
		// So end of the adr, return the end token!
		return 0
	}
	return 0
}

func (l *lexer) Error(e string) {
	err = errors.New(e)
}

// Keep consuming runes as long as f returns true
func (l *lexer) consume(f func(rune) bool) ([]rune, error) {
	rr := []rune{}
	for f(l.peek()) {
		r, _, e := l.in.ReadRune()
		if e != nil {
			return rr, e
		}
		rr = append(rr, r)
	}
	return rr, nil
}

func (l *lexer) peek() rune {
	r, _, e := l.in.ReadRune()
	if e != nil {
		return 0
	}
	l.in.UnreadRune()
	return r
}
