package parser

import (
	"bufio"
	"fmt"
	"github.com/vron/sem/parser/adr"
	"github.com/vron/sem/stresc"
	"strconv"
	"strings"
	"unicode"
)

type item struct {
	r rune
}

type lexer struct {
	input *bufio.Reader
	items chan *Command
}
type stateFn func(*lexer) stateFn

func (l *lexer) run() {
	for state := lexStart; state != nil; {
		state = state(l)
	}
	close(l.items)
}

// TODO: How do we differentiate between a EOF inside and after a command?
func lexStart(l *lexer) stateFn {
	// Try to parse a command
	c, e := l.parseCommand()
	if e != nil {
		l.items <- &Command{Type: c_error, err: e}
		return nil
	}
	l.items <- c
	return lexStart
}

func (l *lexer) parseCommand() (*Command, error) {
	// Consume spaces
	l.consume(unicode.IsSpace)

	// Read an additional rune and use it to switch 
	// the state to try to understand what we should do
	r := l.peek()

	// Try to parse an address
	addr, e := l.parseAddr()
	if e != nil {
		return nil, e
	}
	if addr != nil {
		// We extracted an address, restart
		return addr, nil
	}

	// So no address, then we work on the current addr,
	// Try to find one of the commands we expect!
	for _, v := range commands {
		if strings.HasPrefix(v.text, string(r)) {
			// TODO: Implement here if we want to support commands longer than 1 letter
			l.ReadRune()

			c := Command{Type: v.cmdType}
			if v.takesText || v.takesReg {
				// Next char is terminator
				term, _, e := l.ReadRune()
				if e != nil {
					return nil, e
				}
				// Read until terminator
				str, e := l.consumeDelim(term)
				if e != nil {
					return nil, e
				}
				// Also take the terminator out
				l.ReadRune()
				if v.takesText {
					// Try to escape the text and set it
					strb, e := stresc.Escape([]byte(string(str)))
					if e != nil {
						return nil, e
					}
					c.Text = string(strb)
				} else {
					c.Text = string(str)
				}
			}
			if v.takesRegSub {
				// Next char is terminator
				term, _, e := l.ReadRune()
				if e != nil {
					return nil, e
				}
				// Read until terminator (twice)
				text, e := l.consumeDelim(term)
				if e != nil {
					return nil, e
				}
				l.ReadRune()
				sub, e := l.consumeDelim(term)
				if e != nil {
					return nil, e
				}
				l.ReadRune()
				// Try to escape the text and set it
				subb, e := stresc.Escape([]byte(string(sub)))
				if e != nil {
					return nil, e
				}
				c.Text = string(text)
				c.Sub = string(subb)
			}
			if v.takesAdr {
				// Try to parse a adr
				adr, e := l.parseAddr()
				if e != nil {
					return nil, e
				}
				if adr == nil {
					return nil, fmt.Errorf("Must be able to read Adr")
				}
				c.Adr = adr.Adr
			}
			if v.takesUnix {
				return nil, fmt.Errorf("Not implemented unix yet")
			}
			if v.takesCmd {
				c.Cmds = make([]Command, 0)
				// Now, a addr by itself is not a complete command, so
				// we keep reading commands until it is not a address command!

				for {
					co, e := l.parseCommand()
					if e != nil {
						return nil, e
					}
					c.Cmds = append(c.Cmds, *co)
					// TODO: Do we need to break here?
					if co.Type != C_adr {
						break
					}
				}
			}

			return &c, nil
		}
	}
	// If we reach here output an error
	r, _, e = l.ReadRune()
	if e != nil {
		return nil, e
	}
	return nil, fmt.Errorf("Unexpected char '%v'", string(r))
}

// nil, nil if there is no address to start to parse
// nil, error if started parsing addr and couldn't finish or
// an other error was encountered
func (l *lexer) parseAddr() (*Command, error) {
	node, err := adr.Parse(l.input)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, nil
	}
	// Create a new command of type adr
	return &Command{Type: C_adr, Adr: *node}, nil
}

func (l *lexer) ReadRune() (rune, int, error) {
	return l.input.ReadRune()
}
func (l *lexer) peek() rune {
	r, _, e := l.input.ReadRune()
	if e != nil {
		return 0
	}
	l.input.UnreadRune()
	return r
}

// Read out an integer greedily
func (l *lexer) getNum() (int, error) {
	l.consume(unicode.IsSpace)
	r, e := l.consume(unicode.IsDigit)
	if e != nil {
		return 0, e
	}
	i, er := strconv.Atoi(string(r))
	if er != nil {
		return 0, er
	}
	return i, nil
}

// Consume a string of text until the delimiter, but where the delimiter
// may be escaped by being prepended by a '\'
func (l *lexer) consumeDelim(del rune) ([]rune, error) {
	rr := []rune{}
	flag := false
	for rp := l.peek(); ; rp = l.peek() {
		if flag {
			if rp == del {
				r, _, e := l.ReadRune()
				if e != nil {
					return rr, e
				}
				rr = append(rr, r)
				flag = false
				continue
			} else {
				flag = false
				rr = append(rr, '\\')
				r, _, e := l.ReadRune()
				if e != nil {
					return rr, e
				}
				rr = append(rr, r)
			}
		} else {
			if rp == '\\' {
				flag = true
				_, _, e := l.ReadRune()
				if e != nil {
					return rr, e
				}
				continue
			} else {
				if rp == del {
					break
				}
				r, _, e := l.ReadRune()
				if e != nil {
					return rr, e
				}
				rr = append(rr, r)
			}
		}
	}
	return rr, nil
}

// Keep consuming runes as long as f returns true
func (l *lexer) consume(f func(rune) bool) ([]rune, error) {
	rr := []rune{}
	for f(l.peek()) {
		r, _, e := l.ReadRune()
		if e != nil {
			return rr, e
		}
		rr = append(rr, r)
	}
	return rr, nil
}
