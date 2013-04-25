package adr

import (
	"bufio"
	"bytes"
	"testing"
)

var tests = []struct {
	In   string
	Toks []int
}{
	{
		In:   "123",
		Toks: []int{NUMBER},
	},
	{
		In:   "(123),;#2 2$.'+.3d",
		Toks: []int{pBSTART, NUMBER, pBEND, COMMA, SEMI, HASH, NUMBER, NUMBER, DOLLAR, DOT, ADRMARK, PLUS, DOT, NUMBER},
	},
}

func TestAll(t *testing.T) {
	for _, v := range tests {
		// Lex the entire thing and check that the expected tokens are extracted
		b := bytes.NewBufferString(v.In)
		l := newLexer(bufio.NewReader(b))
		va := []int{}
		for {
			vv := l.Lex(new(yySymType))
			if vv == 0 {
				break
			}
			va = append(va, vv)
		}
		// Check so that they are equal
		if len(v.Toks) == len(va) {
			for i := range v.Toks {
				if v.Toks[i] != va[i] {
					goto Fail
				}
			}
			continue
		}
	Fail:
		t.Log("\n", v.Toks, "\n", va)
		t.Fail()
	}
}
