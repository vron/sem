package parser

import (
	"bufio"
	"bytes"
	"github.com/vron/sem/parser/adr"
	"io"
	"io/ioutil"
	"strings"
)

// Constants used to describe what type of a command something is
const (
	c_error = iota
	C_adr
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

// Build up a map so that we can go from command to it's string representation
// etc, all so that we can pretty print
var cmdInfo map[uint]int

func init() {
	cmdInfo = make(map[uint]int, len(commands))
	for i, v := range commands {
		cmdInfo[v.cmdType] = i
	}
}

// List of all the commands that we support
var commands = []struct {
	text        string
	takesText   bool
	takesRegSub bool
	takesAdr    bool
	takesUnix   bool
	takesCmd    bool
	cmdType     uint
	desc        string
}{
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
	Type uint      // The type of command, given by the constants
	err  error     // The error that occured if this is a error command
	Text string    // Text field of the command
	Sub  string    // Second text filed, or substitute field of command
	Cmds []Command // List of subcommands
	Adr  adr.Node  // The adress structure if this is a address command or takes an address
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
			if v.Type == c_error {
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
		if v.Type == c_error {
			if v.err != io.EOF {
				e = v.err
			}
			break
		} else {
			cmds = append(cmds, *v)
		}
	}
	return cmds, e
}

func (c *Command) String() string {
	// This thing will recursively print the command so we better create a buffer to 
	// write to..
	buf := bytes.NewBuffer(nil)
	c.recString(buf, true)
	return buf.String()
}

func (c *Command) recString(w io.Writer, terminate bool) {
	// Get the info of what I am
	id, ok := cmdInfo[c.Type]
	if !ok {
		if c.Type != C_adr {
			panic("Unimplemented")
		}
		// So this is an adr, format it nicely by calling into the adr package
		io.WriteString(w, c.Adr.String())
		io.WriteString(w, " ")
		return
	}
	desc := commands[id]
	io.WriteString(w, desc.text)
	if desc.takesText {
		sep, str, _ := getSep(c.Text, "")
		io.WriteString(w, sep)
		io.WriteString(w, str)
		io.WriteString(w, sep)
	}
	if desc.takesRegSub {
		sep, a, b := getSep(c.Text, c.Sub)
		io.WriteString(w, sep)
		io.WriteString(w, a)
		io.WriteString(w, sep)
		io.WriteString(w, b)
		io.WriteString(w, sep)
	}
	if desc.takesAdr {
		io.WriteString(w, " ")
		io.WriteString(w, c.Adr.String())
	}
	if desc.takesUnix {
		panic("TODO")
	}
	if desc.takesCmd {
		io.WriteString(w, " ")
		c.Cmds[0].recString(w, false)
	}

	if terminate {
		io.WriteString(w, "\n")
	}
}

// Priority order for how separator is choosen, try first first etc.
// if none of them can be used without escaping the first is choosen
// and all instances are escaped
var seps = []string{"/", "|", "-", "'", " "}

// Returns the best choice of separator for the given strings, as well
// as the input string escaped so that it is pretty! (sep, escaped a, escaped b)
func getSep(a, b string) (string, string, string) {
	for _, v := range seps {
		if !strings.Contains(a, v) && !strings.Contains(b, v) {
			return v, a, b
		}
	}
	// No match found, then escape each occurance
	return seps[0], strings.Replace(a, seps[0], `\`+seps[0], -1),
		strings.Replace(b, seps[0], `\`+seps[0], -1)
}
