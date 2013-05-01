package parser

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"github.com/vron/sem/parser/adr"
)

// Constants used to describe what type of a command something is
const (
	pC_ERROR = iota
	C_ADR
	C_a
	C_c
	C_i
	C_d
	C_s
	C_m
	C_t
	C_pipeIn
	C_pipeOut
	C_pipe
	C_bang
	C_x
	C_y
	C_g
	C_v
	C_k
)

// Define a map so that we may go the other way
var cmdString = map[uint]string{
	pC_ERROR:  "err",
	C_ADR:     "addr",
	C_a:       "a",
	C_c:       "c",
	C_i:       "i",
	C_d:       "d",
	C_s:       "s",
	C_m:       "m",
	C_t:       "t",
	C_pipeIn:  "<",
	C_pipeOut: ">",
	C_pipe:    "|",
	C_bang:    "!",
	C_x:       "x",
	C_y:       "y",
	C_g:       "g",
	C_v:       "v",
	C_k:       "k",
}

type cmdDesc struct {
	text        string
	takesText   bool
	takesRegSub bool
	takesAdr    bool
	takesUnix   bool
	takesCmd    bool
	cmdType     uint
	desc        string
}

var commands = []cmdDesc{
	{"a", true, false, false, false, false, C_a, "Append text after dot"},
	{"c", true, false, false, false, false, C_c, "Change text in dot"},
	{"i", true, false, false, false, false, C_i, "Insert text before dot"},

	{"d", false, false, false, false, false, C_d, "Delete text in dot"},

	{"s", false, true, false, false, false, C_s, "Substitute text in dot"},

	{"m", false, false, true, false, false, C_m, "Move text in dot after address"},
	{"t", false, false, true, false, false, C_t, "Copy text in dot after address"},

	{"<", false, false, false, true, false, C_pipeIn, "Replace dot by command"},
	{">", false, false, false, true, false, C_pipeOut, "Send dot to command"},
	{"|", false, false, false, true, false, C_pipe, "Send to and replace dot by command"},
	{"!", false, false, false, true, false, C_bang, "Run the command"},

	{"x", true, false, false, false, true, C_x, "For each math, set dot, run command"},
	{"y", true, false, false, false, true, C_y, "Between matches, set dot, run command"},
	{"g", true, false, false, false, true, C_g, "If it matches, run command"},
	{"v", true, false, false, false, true, C_v, "If it doesnt match, run command"},

	{"k", false, false, false, false, false, C_k, "Store address in dot"},
}

type Command struct {
	err  error
	Type uint   // The type of command, given by the constants
	Text string // Text field of the command
	Sub  string // Second text filed, or substitute field of command
	// Num         int		
	Cmds []Command // List of subcommands
	//Left, Right *Command
	Adr adr.Node // The adress structure if this is a address command or takes an address
}

// Pretty print this command as a string
func (c *Command) String() string {
	return "S"
}

// Tries to parse the given slice of bytes as commands, if an error is encountered
// it will return a slice of all fully parsed commands and an error. Any partly
// parsed command is discarded
func Parse(b []byte) ([]Command, error) {
	return parse(b)
}

// As Parse but takes a string
func ParseString(s string) ([]Command, error) {
	return parse([]byte(s))
}

// As Parse but for a reader
func ParseReader(r io.Reader) ([]Command, error) {
	// TODO: Do not require the entire thing to be read into memory!
	b, e := ioutil.ReadAll(r)
	if e != nil && e != io.EOF {
		return nil, e
	}
	return parse(b)
}

// As Parse but uses it to propagate put one command at a time at the returned
// channels (blocking). If an error is encountered an error is propagated on 
// the error channel and this function returns. However, all commands on cmd
// should be handled before error
func ParseLive(r *bufio.Reader) (chan Command, chan error) {
	cm, er := make(chan Command), make(chan error)
	go func() {
		l := lexer{r, make(chan *Command)}
		go l.run()
		for {
			v := <-l.items
			if v == nil {
				return
			}
			// Check so the command is not an error
			if v.Type == pC_ERROR {
				println("ERR")
				if v.err != io.EOF {
					er <- v.err
					return
				}
			} else {
				cm <- *v
			}
		}
	}()
	return cm, er
}

func parse(b []byte) ([]Command, error) {
	ba := bytes.NewBuffer(b)
	buf := bufio.NewReader(ba)
	l := lexer{buf, make(chan *Command)}
	go l.run()
	// Get the output until EOF
	cmds := []Command{}
	var e error
	for {
		v := <-l.items
		if v == nil {
			break
		}
		// Check so the command is not an error
		if v.Type == pC_ERROR {
			if v.err != io.EOF {
				e = v.err
			}
			break
		} else {
			cmds = append(cmds, *v)
		}
	}
	// TODO: Handle the errors
	return cmds, e
}
