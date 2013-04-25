package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Variables and initilization for flag vars
var (
	fInp  string
	fOut  string
	fCmd  string
	fEdit bool
	fHelp bool
	fVer  bool
)

func init() {
	flag.StringVar(&fInp, "i", "", "File to use as input stream, if empty Stdin")
	flag.StringVar(&fOut, "o", "", "File to use as output stream, if empty Stdout")
	flag.StringVar(&fCmd, "c", "", "File to use as command input")
	flag.BoolVar(&fEdit, "e", false, "Enter interactive (edit) mode")
	flag.BoolVar(&fHelp, "h", false, "Print this information")
	flag.BoolVar(&fVer, "v", false, "Print version information")
}

func main() {
	inp_, cmd_, out_ := setupSession()
	defer inp_.Close()
	defer cmd_.Close()
	defer out_.Close()
	// Set up the io to be buffered
	inp, cmd, out := bufio.NewReader(inp_),
		bufio.NewReader(cmd_),
		bufio.NewWriter(out_)
	defer func() {
		e := out.Flush()
		if e != nil {
			fmt.Fprintln(os.Stderr, "Could not flush output, may be incomplete: ", e)
		}
	}()

	// Create a new lexer and read out everything from it!
	l := NewLexer(cmd)
	b := make([]Token, 1)
	for {
		n, e := l.Read(b)
		if n > 0 {
			fmt.Println(n, b[1])
		}
		if e != nil {
			fmt.Println(e)
			break
		}
	}

	fmt.Println("Done! ", &inp, &cmd, &out)
}

func setupSession() (input, commands io.ReadCloser, output io.WriteCloser) {
	flag.Parse()

	if fHelp || fVer {
		if fVer {
			fmt.Fprintln(os.Stderr, "This is sams, for streaming processing of non-line oriented files!")
		}
		if fHelp {
			flag.PrintDefaults()
		}
		os.Exit(0)
	}

	// Now, we need to check for sanity of the input arguments

	if fEdit {
		// If in interactive mode input can't be stdin, and there can not be a 
		// command file specified
		if fInp == "" {
			fmt.Fprintln(os.Stderr, "Cannot read both commands and input from Stdin")
			os.Exit(1)
		}
		if fCmd != "" {
			fmt.Fprintln(os.Stderr, "Cannot commands from both Stdin and a file")
			os.Exit(1)
		}
		commands = os.Stdin
	}

	if fCmd != "" {
		//Try to open it for reading
		var e error
		commands, e = os.Open(fCmd)
		if e != nil {
			fmt.Fprintln(os.Stderr, "Could not read command file: ", e)
			os.Exit(1)
		}
	}

	// If commands was not given as interactive or as a file, then 
	// it should be taken from the command line, wrapp that in a ReadCloser
	// so that it may be used in the same way as the others
	if commands == nil {
		cs := flag.Args()
		if len(cs) < 1 {
			fmt.Fprintln(os.Stderr, "No commands given, exiting...")
			os.Exit(1)
		}
		str := ""
		for _, s := range cs {
			str += s
		}
		commands = BufferCloser{strings.NewReader(str)}
	}

	// Set up the input and output
	// Be careful to close the earlier files if we quit
	input = os.Stdin
	if fInp != "" {
		var e error
		input, e = os.Open(fInp)
		if e != nil {
			fmt.Fprintln(os.Stderr, "Could not read input file: ", e)
			commands.Close()
			os.Exit(1)
		}
	}
	output = os.Stdout
	if fOut != "" {
		var e error
		output, e = os.Create(fOut)
		if e != nil {
			fmt.Fprintln(os.Stderr, "Could not open output file: ", e)
			commands.Close()
			input.Close()
			os.Exit(1)
		}
	}

	return
}

type BufferCloser struct {
	io.Reader
}

func (_ BufferCloser) Close() error {
	return nil
}
