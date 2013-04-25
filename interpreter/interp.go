package interpreter

import (
	"errors"
	"sams/memfile"
	"sams/parser"
)

type File struct {
	f   memfile.File // The backing file we operate on
	dot Adr          // Current adr
}

type Adr struct {
	Start, End int
}

// Return a new File that commands may be run on, with the dot defaulting
// to the entire file. Note that if the backing memfile is changed the dot
// will no longer be up to date!
func New(f memfile.File) (*File, error) {
	if f != nil {
		return &File{f: f, dot: Adr{Start: 0, End: f.Length()}}, nil
	}
	return nil, errors.New("Need a backing file")
}

// Run the given commands on the file
func (f *File) Run(cmds []parser.Command) error {
	if f == nil {
		return errors.New("Cannot run on nil file")
	}
	for _, v := range cmds {
		_, e := f.run(v, 0)
		if e != nil {
			return e
		}
	}
	return nil
}

// Get the value of current dot, not well defined whilst Run is running!
func (f *File) Dot() Adr {
	return f.dot
}

// Set dot to the given error, returns an error if the adr is not valid
func (f *File) SetDot(a Adr) error {
	if a.Start < 0 || a.Start > a.End || a.End > f.f.Length() {
		return errors.New("Ivalid addr")
	}
	f.dot = a
	return nil
}

// Get the underlying memfile that this File is operating on
func (f *File) File() memfile.File {
	return f.f
}
