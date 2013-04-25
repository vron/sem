package samx

import (
	"io"
	"os"
	"testing"
)

func BenchmarkChangeTiny(b *testing.B)   { bench(b, "inTiny.txt", "x/\\-/c/Z/") }
func BenchmarkChangeSmall(b *testing.B)  { bench(b, "inSmall.txt", "x/\\-/c/Z/") }
func BenchmarkChangeNormal(b *testing.B) { bench(b, "inNormal.txt", "x/\\-/c/Z/") }

func BenchmarkChangeLineTiny(b *testing.B)   { bench(b, "inTiny.txt", "x/.*\\n/a/Z/") }
func BenchmarkChangeLineSmall(b *testing.B)  { bench(b, "inSmall.txt", "x/.*\\n/a/Z/") }
func BenchmarkChangeLineNormal(b *testing.B) { bench(b, "inNormal.txt", `x/A*/c/num/`) }

func BenchmarkAppendTiny(b *testing.B)   { bench(b, "inTiny.txt", "x/\\-/a/Z/") }
func BenchmarkAppendSmall(b *testing.B)  { bench(b, "inSmall.txt", "x/\\-/a/Z/") }
func BenchmarkAppendNormal(b *testing.B) { bench(b, "inNormal.txt", "x/\\-/a/Z/") }

func bench(b *testing.B, file, cmd string) {
	fi, e := os.Stat(file)
	if e != nil {
		b.Fatal(e)
	}
	b.SetBytes(fi.Size())
	var s *Stream
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Read the file
		f, e := os.Open(file)
		if e != nil {
			b.Fatal(e)
		}
		s = New(f)
		f.Close()
		// And run the loop cmd on this input
		b.StartTimer()
		s.Run(cmd)
		b.StopTimer()
	}
	if true {
		b.StopTimer()
		fo, _ := os.Create("out.txt")
		defer fo.Close()
		io.Copy(fo, s.Reader())
	}
}
