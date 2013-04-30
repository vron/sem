package interpreter

import (
	"fmt"
	"github.com/vron/sem/parser"
	//"spew"
	"github.com/vron/sem/memfile/gap"
	"testing"
)

var a = fmt.Println

var Tests = []struct {
	// Description, Input, Command, Output
	desc, i,c,o string
}{
	{
		"Simple append",
		`Dirlik dirlik`,
		"a/dirlik/",
		"Dirlik dirlikdirlik",
	},
	{
		"Simple insert",
		`Dirlik dirlik`,
		"i/dirlik/",
		"dirlikDirlik dirlik",
	},
	{
		"Simple replace",
		`Dirlik dirlik`,
		"c/dirlik/",
		"dirlik",
	},
	{
		"Simple delete",
		`Dirlik dirlik`,
		"d",
		"",
	},
	{
		"Simple substitute",
		`Dirlik dirlik`,
		"s/irlik/an/",
		"Dan dirlik",
	},
	{
		"Simple copy",
		`Dirlik dirlik`,
		"t 1",
		"Dirlik dirlikDirlik dirlik",
	},
	{
		"Simple extract",
		`Dirlik dirlik`,
		"x/irlik/d",
		"D d",
	},
	{
		"Addres char",
		`Dirlik dirlik`,
		"#2,#3d",
		"Dilik dirlik",
	},
	{
		"Addres line",
		"Dirlik dirlik\nZhao zhao",
		"2d",
		"Dirlik dirlik\n",
	},
	{
		"Addres reg",
		`Dirlik dirlik`,
		"/ir/d",
		"Dlik dirlik",
	},
	{
		"Addres +",
		`Dirlik dirlik`,
		"#1+#1,#3d",
		"Dilik dirlik",
	},
	{
		"+ Missing",
		`Dirlik dirlik`,
		"+",
		"Dirlik dirlik",
	},
	{
		"Extract change all chars",
		`Dirlik dirlik`,
		"x/./c/F/",
		"FFFFFFFFFFFFF",
	},
	{
		`Testing \n for newline`,
		`Dirlik dirlik`,
		`s/k/\n/`,
		"Dirli\n dirlik",
	},
	{
		"Line find 0",
		"a\nb\nc\nd\n",
		"0 a/X/",
		"Xa\nb\nc\nd\n",
	},
	{
		"Line find 3",
		"a\nb\nc\nd\n",
		"3 c/X/",
		"a\nb\nXd\n",
	},
	{
		"Line find Max",
		"a\nb\nc\nd\n",
		"30 c/X/",
		"a\nb\nc\nd\nX",
	},
	{
		"Line find find 1",
		"a\nb\nc\nd\n",
		"1+2 c/X/",
		"a\nb\nXd\n",
	},
	{
		"Backwards rune count 0",
		`Dirlik dirlik`,
		"#6-#0 a/X/",
		"DirlikX dirlik",
	},
	{
		"Backwards rune count 3",
		`Dirlik dirlik`,
		"#6-#3 a/X/",
		"DirXlik dirlik",
	},
	{
		"Backwards rune count max",
		`Dirlik dirlik`,
		"#6-#300 a/X/",
		"XDirlik dirlik",
	},
	{
		"Backwards line count 0",
		"a\nb\nc\nd\n",
		"3-0 a/X/",
		"a\nb\nXc\nd\n",
	},
	{
		"Backwards line count 1",
		"a\nb\nc\nd\n",
		"3-1 c/X/",
		"a\nXc\nd\n",
	},
	{
		"Backwards line count 2",
		"a\nb\nc\nd\n",
		"4-2 c/X/",
		"a\nXc\nd\n",
	},
	{
		"Backwards regexp",
		`Dirlik dirlik`,
		"$-/ir/d",
		"Dirlik dlik",
	},
	{
		"Line idiom",
		"a\nb\nc\nabc",
		"x/.*/ d",
		"\n\n\n",
	},
	{
		"Remove dubble lines",
		"a\n\n\n\n\nb\n",
		`x/\n+/c/\n/`,
		"a\nb\n",
	},
/* Below here are reported errors from user files */
	{
		"Bug 1",
		"Dirlik dirlik\n\t\t\tZah zah\n",
		"x/irlik/d",
		"D d\n\t\t\tZah zah\n",
	},
}

func TestAll(t *testing.T) {
	for _, v := range Tests {
		cc, e := parser.Parse([]byte(v.c))
		if e != nil {
			t.Log(e)
			t.Fail()
		}
		fb := gap.New([]byte(v.i))
		f, _ := New(fb)
		e = f.Run(cc)
		if e != nil {
			t.Log(e)
			t.Fail()
		}
		// Check so that the output is correct
		bs, _ := f.f.Get(0, f.f.Length())
		so := string(bs)
		if len(so) == len(v.o) {
			for i := range so {
				if so[i] != v.o[i] {
					goto Fail
				}
			}
			// goto Fail
			continue
		}
	Fail:
		t.Log(v.desc)
		t.Log("Not equal output")
		t.Log("Expected: `", v.o, "`")
		t.Log("Got:      `", so, "`")
		t.FailNow()
	}
}

