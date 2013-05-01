package stresc

import (
	"testing"
	"bytes"
)

var tests = 
[]struct{
	in string
	out string
	err bool
}{
	{
		`hej`,
		"hej",
		false,
	},	
	{
		`a\a`,
		"a\a",
		false,
	},	
	{
		`b\b`,
		"b\b",
		false,
	},		
	{
		`f\f`,
		"f\f",
		false,
	},		
	{
		`n\n`,
		"n\n",
		false,
	},		
	{
		`r\r`,
		"r\r",
		false,
	},		
	{
		`t\t`,
		"t\t",
		false,
	},		
	{
		`v\v`,
		"v\v",
		false,
	},			
	{
		`v\'`,
		"v'",
		false,
	},		
	{
		`v\"`,
		"v\"",
		false,
	},	
	{
		`\\\`,
		"",
		true,
	},	
	{
		`\\`,
		"\\",
		false,
	},	
}

func TestAll(t *testing.T) {
	for _, v := range tests {
		bin := []byte(v.in)
		bout:= []byte(v.out)
		b, e := Escape(bin)
		if v.err {
			if e == nil {
				t.Error("Expected error but didn't get one", v.in)
				t.FailNow()
			} else {
				continue
			}
		}
		if e != nil {
			t.Error("Got unexpected error: ", e)
			t.FailNow()
		}
		if !bytes.Equal(b, bout) {
			t.Error("Unexpected output '", string(b), "' for '", v.in, "'")
			t.FailNow()
		}	
	}
}