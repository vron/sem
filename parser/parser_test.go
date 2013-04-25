package parser

import (
	"testing"
)

var Tests = []struct {
	I, O, Desc string
	C          []string
}{
	{
		Desc: "Some stuff with adr",
		I:    iName,
		C: []string{
			"#3,#4 a/dirlik/",
		},
		O: "Dirlik dirlikdirlik",
	},
	{
		Desc: "Simple append",
		I:    iName,
		C: []string{
			"a/dirlik/",
		},
		O: "Dirlik dirlikdirlik",
	},
	{
		Desc: "Simple insert",
		I:    iName,
		C: []string{
			"i/dirlik/",
		},
		O: "dirlikDirlik dirlik",
	},
	{
		Desc: "Simple replace",
		I:    iName,
		C: []string{
			"c/dirlik/",
		},
		O: "dirlik",
	},
	{
		Desc: "Simple delete",
		I:    iName,
		C: []string{
			"d",
		},
		O: "",
	},
	{
		Desc: "Simple substitute",
		I:    iName,
		C: []string{
			"s/irlik/an/",
		},
		O: "Dan dirlik",
	},
	{
		Desc: "Simple move",
		I:    iName,
		C: []string{
			"m 1",
		},
		O: "Dirlik dirlik",
	},
	{
		Desc: "Simple copy",
		I:    iName,
		C: []string{
			"t 1",
		},
		O: "Dirlik dirlikDirlik dirlik",
	},
	{
		Desc: "Simple extract",
		I:    iName,
		C: []string{
			"x/irlik/d",
		},
		O: "D d",
	},
	{
		Desc: "Addres char",
		I:    iName,
		C: []string{
			"#3",
		},
		O: "Dirlik dirlik",
	},
	{
		Desc: "Addres line",
		I:    iName,
		C: []string{
			"33",
		},
		O: "Dirlik dirlik",
	},
	{
		Desc: "Addres reg",
		I:    iName,
		C: []string{
			"/ir/",
		},
		O: "Dirlik dirlik",
	},
	{
		Desc: "Addres +",
		I:    iName,
		C: []string{
			"1+3",
		},
		O: "Dirlik dirlik",
	},
	{
		Desc: "Implicit +",
		I:    iName,
		C: []string{
			"3,5#100",
		},
		O: "Dirlik dirlik",
	},
	{
		Desc: "+ Missing",
		I:    iName,
		C: []string{
			"+",
		},
		O: "Dirlik dirlik",
	},
}

func TestAll(t *testing.T) {
	for _, v := range Tests {
		for _, c := range v.C {
			ParseString(c)
			//t.Logf("%v: %v\n", v.Desc, c)

		}
	}
}

var (
	iName = `Dirlik dirlik`
)
