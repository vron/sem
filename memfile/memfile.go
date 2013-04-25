package memfile

/*
Defines a interface that all things that want to be memfiles need to implement.
A memfile is a utf8 encoded text file stored in memory that som operations can be performed on!
*/

import (
	"errors"
	"io"
	"github.com/vron/sem/regexp"
)

// The error that should be returned if this is the issue
var OutOfBounds = errors.New("Range out of bounds")

// TODO: Add a method that returns a ReadSeeker that reads the File backwards!

type File interface {
	// Get access to raw reading of bytes from the file
	Read(p []byte) (n int, e error)
	Seek(offset int64, whence int) (int64, error)
	//io.ReadSeeker

	//Reader(start int) io.RuneReader

	// Returns a RuneReader that can be used to read runes backwards from the given
	// starting point
	BackwardsReader(start int) io.RuneReader

	// Change the bytes at [start:end] to the bytes in data. Note that no sanity checking
	// on utf8 runes is perfromed at the two ends. Returns and error if start || end is out
	// of bounds ot end < start.
	Change(start, end int, data []byte) error

	// Return a byte slice for the given position, note that it MAY be a copy and MAY be
	// a slicer pointing to the original data, so use only for reading!
	Get(start, end int) ([]byte, error)

	// Return the offset (in absolute terms) of the start of the given line number from position start. Line 0 is the
	// start of the line start is on, 1 start of next line (i.e. position after first found line end)
	// ln < 0 is not supported
	OffsetLine(ln, start int) (offset int, e error)

	// Return the offset (in absolute terms) of the start of the given rune number from position start.
	OffsetRune(cn, start int) (offset int, e error)

	// Return the length in bytes
	Length() int
	
	// A Memfile should also implement the regexp interface to be happy!
	regexp.Input
}

type TrackedFile struct {
	f File

	// If we are in a transaction and how many changed applied
	trans bool
	noCh  int
}

// Takes a File and wraps it, still implementing File but now adding functionality
// to track all changes done through the returned structures Change method, to group
// several changes as one transaction use Start, End but dont forget to call End or
// all successive changes will be logged as one!
func TrackChanges(f File) *TrackedFile {
	return &TrackedFile{f: f}
}

// TODO: This type should also export a way to later look at revisions!

// Return an error if a transaction is allready started. A transaction is considered
// started first when a change is applied following a Start call.
func (tf *TrackedFile) Start() error {
	if tf.trans {
		return errors.New("Transaction allready started")
	}
	tf.noCh = 0
	tf.trans = true
	return nil
}

// Return the number of changes applied in this transaction. An error is returned
// if no transaction was started.
func (tf *TrackedFile) End() (int, error) {
	if !tf.trans {
		return 0, errors.New("No transaction to end")
	}
	tf.trans = false
	return tf.noCh, nil
}

func (tf *TrackedFile) Read(b []byte) (int, error) {
	return tf.f.Read(b)
}
func (tf *TrackedFile) Seek(offset int64, whence int) (int64, error) {
	return tf.f.Seek(offset, whence)
}
func (tf *TrackedFile) Change(start, end int, data []byte) error {
	if tf.trans {
		tf.noCh++
	}
	return tf.f.Change(start, end, data)
}
func (tf *TrackedFile) OffsetLine(ln, start int) (offset int, e error) {
	return tf.f.OffsetLine(ln, start)
}
func (tf *TrackedFile) OffsetRune(cn, start int) (offset int, e error) {
	return tf.f.OffsetRune(cn, start)
}
