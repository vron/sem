package interpreter

import (
	"math/rand"
	"sams/memfile/gap"
	"sams/parser"
	"testing"
	"unicode"
	"unicode/utf8"
)

func insertDelete(f *File) {
	pos := rand.Intn(f.f.Length())
	f.dot = Adr{pos, pos}
	f.run(parser.Command{Type: parser.C_i, Text: "Hello"}, 0)
	f.dot = Adr{pos, pos + 5}
	f.run(parser.Command{Type: parser.C_d}, 0)
}

func changeSame(f *File) {
	pos := rand.Intn(f.f.Length())
	f.dot = Adr{pos, pos + 1}
	f.run(parser.Command{Type: parser.C_c, Text: "a"}, 0)
}

func xtractSubstitute(f *File) {
	f.dot = Adr{0, f.f.Length()}
	f.run(parser.Command{Type: parser.C_x, Text: ".", Cmds: []parser.Command{{Type: parser.C_s, Text: ".", Sub: "F"}}}, 0)
}

func xtractChangeAll(f *File) {
	f.dot = Adr{0, f.f.Length()}
	f.run(parser.Command{Type: parser.C_x, Text: ".", Cmds: []parser.Command{{Type: parser.C_c, Text: "F"}}}, 0)
}

func bench(b *testing.B, f func(*File), size int, bytes int) {
	fi := createFile(1024 * size)
	if bytes > 0 {
		b.SetBytes(int64(1024 * size))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(fi)
	}
}

func BenchmarkXtractChangeAll_1k(b *testing.B)   { bench(b, xtractChangeAll, 1, 1) }
func BenchmarkXtractChangeAll_10k(b *testing.B)  { bench(b, xtractChangeAll, 10, 1) }
func BenchmarkXtractChangeAll_100k(b *testing.B) { bench(b, xtractChangeAll, 100, 1) }
func BenchmarkXtractChangeAll_1M(b *testing.B)   { bench(b, xtractChangeAll, 1024, 1) }

func BenchmarkXtractSubstitute_1k(b *testing.B)   { bench(b, xtractSubstitute, 1, 1) }
func BenchmarkXtractSubstitute_10k(b *testing.B)  { bench(b, xtractSubstitute, 10, 1) }
func BenchmarkXtractSubstitute_100k(b *testing.B) { bench(b, xtractSubstitute, 100, 1) }
func BenchmarkXtractSubstitute_1M(b *testing.B)   { bench(b, xtractSubstitute, 1024, 1) }

func BenchmarkChangeSame_1k(b *testing.B)   { bench(b, changeSame, 1, 0) }
func BenchmarkChangeSame_10k(b *testing.B)  { bench(b, changeSame, 10, 0) }
func BenchmarkChangeSame_100k(b *testing.B) { bench(b, changeSame, 100, 0) }

func BenchmarkInsertDelete_1k(b *testing.B)   { bench(b, insertDelete, 1, 0) }
func BenchmarkInsertDelete_10k(b *testing.B)  { bench(b, insertDelete, 10, 0) }
func BenchmarkInsertDelete_100k(b *testing.B) { bench(b, insertDelete, 100, 0) }
func BenchmarkInsertDelete_1M(b *testing.B)   { bench(b, insertDelete, 1024, 0) }
func BenchmarkInsertDelete_10M(b *testing.B)  { bench(b, insertDelete, 10*1024, 0) }
func BenchmarkInsertDelete_100M(b *testing.B) { bench(b, insertDelete, 100*1024, 0) }

func createFile(size int) *File {
	// Keep filling with random runes, ensuring that we have somewhat reasonable line lengths
	buf := make([]byte, size)
	var added int
	for added = 0; added < size-2; {
		var r rune

		// Get a random from letters
		lets := unicode.Categories["L"].R16
		cid := rand.Intn(len(lets) - 1)
		re := lets[cid]
		// And get a random one from this stride!
		num := rand.Intn(int((re.Hi - re.Lo) / re.Stride))
		r = rune(re.Lo + uint16(num)*re.Stride)
		// 
		// So now encode the rune into our slice!
		f := rand.Float32()
		if f < 1.0/40.0 {
			r = '\n'
		}
		le := utf8.EncodeRune(buf[added:], r)
		added += le
	}
	// Fill with something we know for sure is 1 byte long
	for i := added; i < len(buf); i++ {
		buf[i] = byte('a')
	}
	f, _ := New(gap.New(buf))
	return f
}
