/* Format sem sources */
package main

import (
	"flag"
	"fmt"
	"github.com/vron/sem/parser"
	"io"
	"io/ioutil"
	"os"
)

var (
	fTabs     bool
	fTabwidth int
	fW        bool
	fH        bool
)

func init() {
	flag.BoolVar(&fTabs, "tabs", true, "indent with tabs")
	flag.IntVar(&fTabwidth, "tabwidth", 4, "tab width")
	flag.BoolVar(&fW, "w", false, "write result to (source) file instead of stdout")
	flag.BoolVar(&fH, "h", false, "show this help")
}

var tab []byte = []byte("\t")

func main() {
	flag.Parse()
	if fH {
		flag.PrintDefaults()
	}
	if !fTabs {
		if fTabwidth < 1 {
			fmt.Fprintln(os.Stderr, "Tabwidth must be larger than 1")
			os.Exit(1)
		}
		tab = make([]byte, fTabwidth)
		for i := range tab {
			tab[i] = byte(' ')
		}
	}
	// Loop through each of the files we are asked to fmt
	for i := 0; i < flag.NArg(); i++ {
		dofile(flag.Arg(i))
	}
}

// Take a file and try to parse it and write it back
func dofile(f string) {
	// Try to read the file, 
	b, e := ioutil.ReadFile(f)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Could no open file '%s': %v", f, e)
		return
	}
	fixfile(b, f)
}

// Try to parse the file and output the formatted file back,
// if file name is empty write to stout instead
func fixfile(b []byte, f string) {
	cmds, e := parser.Parse(b)
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error occured whil parsing file '%s':\n", f)
		fmt.Fprintln(os.Stderr, e)
		fmt.Fprintln(os.Stderr, "Ignoring this file")
		return
	}
	// If we should output it to the file we try to do that
	w := os.Stdout
	if fW && f != "" {
		fi, e := os.Create(f)
		if e != nil {
			fmt.Fprintln(os.Stderr, "Could not open '%s' for writing", f)
			return
		}
		defer fi.Close()
		w = fi
	}
	// Pretty print the command to the writer
	for _, c := range cmds {
		io.WriteString(w,c.String())
		io.WriteString(w,"\n")
	}
}
