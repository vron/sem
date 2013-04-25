package adr

import (
	"bufio"
	"bytes"
	"spew"
	"testing"
)

var tst = []struct {
	Desc   string
	In     string
	IsAddr bool
	IsErr  bool
}{
	{
		Desc:   "Some stuff ",
		In:     "123,34+#2+1",
		IsAddr: true,
		IsErr:  false,
	},
	{
		Desc:   "Some stuff fail",
		In:     ",##",
		IsAddr: true,
		IsErr:  true,
	},
	{
		Desc:   "Some stuff fail",
		In:     "s,123,34+#2+1",
		IsAddr: false,
		IsErr:  true,
	},
	{
		Desc:   "Fill it up baby",
		In:     ",",
		IsAddr: true,
		IsErr:  false,
	},
}

func TestParse(t *testing.T) {
	for i, v := range tst {
		a, e := Parse(bufio.NewReader(bytes.NewBufferString(v.In)))

		if !v.IsAddr {
			if a != nil || e != nil {
				t.Log("Expect correct output: ", v.Desc)
				t.Fail()
			}
			continue
		}

		if !v.IsErr && e != nil {
			t.Fail()
			t.Log(v.Desc, ":", e)
			continue
		}

		if v.IsErr && e == nil || i < 0 {
			t.Fail()
			t.Log("Expected error: ", v.Desc)
			spew.Dump(a)
			continue
		}
	}
}
