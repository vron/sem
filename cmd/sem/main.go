package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/vron/sem/memfile/gap"
	"github.com/vron/sem/parser"
	"github.com/vron/sem/interpreter"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	fW bool
	fA string
	fD string
	fH bool
)

func init() {
	flag.BoolVar(&fW, "w", false, "write output to files instead of to stdout")
	flag.StringVar(&fA, "a", ".mod", "file ending to append, if empty replace original files")
	flag.StringVar(&fD, "d", "", "directory to write files to, if empty same as in file")
	flag.BoolVar(&fH, "h", false, "display this information")
}

/* Operational modes:
sem path/to/script/file.sem infile...
	read the file, execute commands and if no error write to stdout
	flags:
		-w	Write output to files
		-a=".mod" 	add the extra ending to the file
*/
func main() {
	flag.Parse()
	if fH {
		fmt.Println("usage: sem [flags] path/to/script [path ...]")
		flag.PrintDefaults()
		return
	}

	// Check so that we can get the sem file and parse it
	if flag.NArg() < 1 {
		fmt.Println("error: need path to a sem file to progress")
		os.Exit(1)
	}
	sf, e := ioutil.ReadFile(flag.Arg(0))
	if e != nil {
		fmt.Println("error: could not read sem file:", e)
		os.Exit(1)
	}
	if bytes.HasPrefix(sf, []byte("#!")) {
		// Ignore the shebang line
		i := bytes.Index(sf, []byte("\n"))
		sf = sf[i:]
	}
	cmds, e := parser.Parse(sf)
	if e != nil {
		fmt.Println("error: could not parse sem file")
		fmt.Println(e)
		os.Exit(1)
	}

	if flag.NArg() == 1 {
		if fW || fD != "" || fA != ".mod" {
			fmt.Println("when using stdin only stdout is supported as output")
			os.Exit(1)
		}
		do(cmds, os.Stdin, os.Stdout)
	}

	// Ok, so we are working on the specified files, loop each one of them and
	// run them through the command
	for i := 1; i < flag.NArg(); i++ {
		dofile(cmds, flag.Arg(i))
	}
	_ = cmds

}

func dofile(cmds []parser.Command, p string) {
	f, e := os.Open(p)
	if e != nil {
		fmt.Printf("error: could not open '%s': %v\n", p, e)
		return
	}
	defer f.Close()
	out := os.Stdout
	if fW {
		if fD == "" {
			fo, e := os.Create(p+fA)
			if e != nil {
				fmt.Printf("error: could not create '%s': %v\n", p+fA, e)
				return
			}
			defer fo.Close()
			out = fo
		} else {
			// Take the file name of the original file and write to this dir
			info, e := os.Stat(fD)
			if e != nil {
				fmt.Printf("error: no directory '%s'; %v\n", fD, e)
				os.Exit(1)
			}
			if !info.IsDir() {
				fmt.Printf("error: not a directory '%s'\n", fD)
				os.Exit(1)
			}
			_, fn := filepath.Split(p)
			fp := filepath.Join(fD, fn, fA)
			fo, e := os.Create(fp)
			if e != nil {
				fmt.Printf("error: could not create '%s': %v\n", fp, e)
				return
			}
			defer fo.Close()
			out = fo
		}
	}
	do(cmds, f, out)
}

// Run all the commands on the data in in and if no error happens
// then do the output to out
func do(cmds []parser.Command, in io.Reader, out io.Writer) {
	f := gap.NewReader(in)
	i, e := interpreter.New(f)
	if e != nil {
		fmt.Println("error: could not create interpreter:", e)
		return
	}
	e = i.Run(cmds)
	if e != nil {
		fmt.Println("error while executing commands:", e)
		return
	}
	_, e = i.File().Seek(0,0)
	if e != nil {
		fmt.Println("error: could not setupt output:", e)
		return
	}
	_, e = io.Copy(out, i.File())
	if e != nil {
		fmt.Println("error: could not write all data to output:", e)
		return
	}
}
