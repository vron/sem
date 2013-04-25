/*
	Samstream is used to run sam commands on texts for easy use in "scripty"
	applications written in go
*/
package samx

import (
	"io"
	//	"spew"
	"github.com/vron/sem/interpreter"
	"github.com/vron/sem/memfile"
	"github.com/vron/sem/memfile/gap"
	"github.com/vron/sem/parser"
)

type Stream struct {
	file memfile.File
}

// Create a new stream by reading all ontents of the reader
func New(in io.Reader) *Stream {
	s := &Stream{}
	// Read in to the file
	s.file = gap.NewReader(in)
	return s
}

// Runs the commands given by the command string on the file and
// returns the error if one is found
func (s *Stream) Run(cmd string) error {
	// First parse the entire command string
	cm, er := parser.ParseString(cmd)
	if er != nil {
		return er
	}
	//spew.Dump(cm)
	// So now we run these commands on the file
	fi, er := interpreter.New(s.file)
	if er != nil {
		return er
	}
	return fi.Run(cm)
}

// Return a reader that will read contents of the stream (only usefull before doinga ny modifications again)
func (s *Stream) Reader() io.Reader {
	// TODO: This might modify file, return a new reader instead
	s.file.Seek(0, 0)
	return s.file
}
