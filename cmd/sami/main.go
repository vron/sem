package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sams/interpreter"
	"sams/memfile/slice"
	"sams/parser"
)

var p = fmt.Println

func main() {
	log.SetFlags(0)
	flag.Parse()
	// We expect exactly one argument which is the file name, this
	// is the file that we will work on by reading it in, and each time
	// a command is run saving the chagnes to, to allow for some sort of
	// gui based operation by using a editor with capability to reload on
	// change

	if len(flag.Args()) != 1 {
		log.Println("Expect exactly on file name to edit")
		os.Exit(1)
	}
	// Try to open the given file
	f, e := os.Open(flag.Args()[0])
	if e != nil {
		log.Println("Could not open file ", flag.Args()[0])
		os.Exit(1)
	}
	f.Close()

	// Start into the command reading loop:
	doCmd(flag.Args()[0])
}

func doCmd(p string) bool {
	// Start reading on the std input!
	inp := bufio.NewReader(os.Stdin)
	cm, er := parser.ParseLive(inp)
	// So start reading command by command!
	// Make sure we keep the dot between loads
	var a *interpreter.Adr
	for {
		select {
		case c := <-cm:
			// So we have parsed a command, try to execure this command by reading, doing, writing
			b, e := ioutil.ReadFile(p)
			if e != nil {
				log.Println("--> ", e)
				break
			}
			f := slice.New(b)
			fi, e := interpreter.New(f)
			if e != nil {
				log.Println("--> ", e)
				break
			}
			if a != nil {
				fi.SetDot(*a)
			}
			e = fi.Run([]parser.Command{c})
			if e != nil {
				log.Println("--> ", e)
				break
			}
			ap := fi.Dot()
			a = &ap
			// Write the new output
			bb, _ := f.Get(0, f.Length())
			e = ioutil.WriteFile(p, bb, 0666)
			if e != nil {
				log.Println("--> ", e)
				break
			}
			log.Println("--> File loaded, changed and written")
		case e := <-er:
			log.Println(e)
			return false
		}
	}

	return false
}
