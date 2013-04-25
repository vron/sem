package test

import (
	"io/ioutil"
	"github.com/vron/sem/memfile"
	//"sams/memfile/blocks"
	"github.com/vron/sem/memfile/gap"
//	"github.com/vron/sem/memfile/slice"
	"testing"
)

var imp = []memfile.File{
	gap.New(nil),
//	slice.New(nil),
}

// TODO: This testing should be made much more rigorous!! (It is an important building block..)

// Test all implementations of memfiles!
func TestAll(t *testing.T) {
	for _, f := range imp {
		// Do some operations and check that we get what we expect!
		f.Change(0, 0, []byte("12345"))
		f.Change(0, 0, []byte("0"))
		f.Change(6, 6, []byte("6"))
		f.Change(1, 5, nil)
		expect(t, f, "056")
		// Add some rows so we may test also the other funcitons
		f.Change(2, 2, []byte("\n"))
		f.Change(1, 1, []byte("\n"))
		f.Change(5, 5, []byte("\nrow"))
		expect(t, f, "0\n5\n6\nrow")
		if o, _ := f.OffsetLine(1, 1); o != 2 {
			t.Error("Wrong line start", o)
		}
		if o, _ := f.OffsetLine(0, 1); o != 0 {
			t.Error("Wrong line start", o)
		}
		if o, _ := f.OffsetLine(0, 8); o != 6 {
			t.Error("Wrong line start", o)
		}
	}
}

func expect(t *testing.T, f memfile.File, s string) {
	f.Seek(0,0)
	b, _ := ioutil.ReadAll(f)
	if string(b) != s {
		t.Errorf("'%v'!='%v'", string(b), s)
		t.FailNow()
	}
}
