package interpreter

import (
	"errors"
	"fmt"
	"github.com/vron/sem/parser"
	padr "github.com/vron/sem/parser/adr"
	"github.com/vron/sem/regexp"
	"github.com/davecgh/go-spew/spew"
	"time"
)

// Store the global value of the adr mark
// TODO: Put this into file
var adrmark Adr

// fix is a int that refere to a fixed position in the file, it is returned updated
// to keep track for e.g. extract when it does modification
func (f *File) run(cmd parser.Command, fix int) (fa int, e error) {
	var (
		adr Adr
		reg *regexp.Regexp
	)

	f.dot = f.sanitiseAdr(f.dot)

	// TODO: Update fix locations
	switch cmd.Type {
	case parser.C_ADR: // Set dot command
		f.dot, e = f.EvalAdr(cmd.Adr)
	case parser.C_a: // Append
		f.f.Change(f.dot.End, f.dot.End, []byte(cmd.Text))
		if fix >= f.dot.End {
			fix += len(cmd.Text)
		}
	case parser.C_c: // Change
		f.f.Change(f.dot.Start, f.dot.End, []byte(cmd.Text))
		if fix >= f.dot.End {
			fix += len(cmd.Text) - (f.dot.End - f.dot.Start)
		} else if fix > f.dot.Start && fix < f.dot.End {
			fix = -1 // Fix point not well defined with this change
		}
	case parser.C_i: // Insert
		f.f.Change(f.dot.Start, f.dot.Start, []byte(cmd.Text))
		if fix > f.dot.Start {
			fix += len(cmd.Text)
		}
	case parser.C_d: // Delete
		f.f.Change(f.dot.Start, f.dot.End, nil)
		if fix >= f.dot.End {
			fix -= (f.dot.End - f.dot.Start)
		} else if fix > f.dot.Start && fix < f.dot.End {
			fix = -1 // Fix point not well defined with this change
		}
	case parser.C_m: // Move
		adr, e = f.EvalAdr(cmd.Adr)
		if e != nil {
			break
		}
		te, _ := f.f.Get(f.dot.Start, f.dot.End)
		f.f.Change(adr.End, adr.End, te)
		f.f.Change(f.dot.Start, f.dot.End, nil)
	case parser.C_t: // Copy
		adr, e = f.EvalAdr(cmd.Adr)
		if e != nil {
			break
		}
		te, _ := f.f.Get(f.dot.Start, f.dot.End)
		f.f.Change(adr.End, adr.End, te)
	case parser.C_s: // Substitute
		reg, e = regexp.Compile(cmd.Text) // |Text|Sub|
		if e != nil {
			break
		}
		f.f.Seek(int64(f.dot.Start), 0)
		loc := reg.FindInputSubmatchIndex(f.f)
		if loc == nil {
			break
		}
		for i := range loc {
			loc[i] += f.dot.Start
		}
		t, _ := f.f.Get(loc[0], loc[1])
		re := reg.Expand([]byte{}, []byte(cmd.Sub), t, loc)
		// Now we have it so now substitute
		f.f.Change(loc[0], loc[1], re)
	case parser.C_x: // Extract
		// TODO: Think about how this should work for after change overlapping etc. ... (need usage to know I think)
		// Note that it may create eternal loop, if the user is not carefull... We consider this ok...
		reg, e = regexp.Compile(cmd.Text)
		if e != nil {
			break
		}
		var fp = f.dot.Start
		f.f.Seek(int64(fp), 0)
		var lastloc int  = -1
		for loc := reg.FindInputIndex(f.f); loc != nil; loc = reg.FindInputIndex(f.f) {
			println(loc[0], loc[1], fp, f.f.Length())
			baj, _ := f.f.Get(0,f.f.Length());
			println(string(baj))
			time.Sleep(1*time.Second)
			if len(loc) > 1 && loc[0]==loc[1] {
				// If we have null match we need to check so we don't keep stamping at the
				// same location
				if lastloc < 0 {
					lastloc	= loc[0]
				} else {
					lastloc = -1;
					fp++ // We should advance one rune note one byte!
					continue
				}
			}
			// Create a new file with dot set and run the provided cmd on that
			fn := File{f: f.f, dot: Adr{Start: loc[0] + fp, End: loc[1] + fp}}
			fp = fp + loc[1]
			println("d", fp, fn.dot.Start, fn.dot.End)
			fp, e = fn.run(cmd.Cmds[0], fp) // TODO: What if this is nill, loop loop? i.e. Run all subcommands not just one!
			if e != nil {
				break
			} else if fp < 0 {
				e = errors.New("Undefined inner extract change")
			}
			// TODO: Maybe we could return negative fix point to indicate that we had an overlapping edit?
			if fp >= f.f.Length() {
				break
			}
			f.f.Seek(int64(fp), 0)
		}
	case parser.C_y: // Extract between matches
		e = errors.New("y not implemented yet")
	case parser.C_g: // Guard
		reg, e = regexp.Compile(cmd.Text)
		if e != nil {
			break
		}
		f.f.Seek(int64(f.dot.Start), 0)
		loc := reg.FindInputIndex(f.f)
		if loc != nil && loc[1]+f.dot.Start <= f.dot.End {
			return f.run(cmd.Cmds[0], fix)
		}
	case parser.C_v: // Guard not
		reg, e = regexp.Compile(cmd.Text)
		if e != nil {
			break
		}
		f.f.Seek(int64(f.dot.Start), 0)
		loc := reg.FindInputIndex(f.f)
		if loc == nil || loc[1]+f.dot.Start > f.dot.End {
			return f.run(cmd.Cmds[0], fix)
		}
	case parser.C_k: // Store dot in adr mark
		adrmark = f.dot
	default:
		return fix, fmt.Errorf("Unkown command %d", cmd.Type)
	}
	return fix, e
}

// Return the location of the evaluated addres in this file the returned
// addres is guarenteed to be sanitised when it is created, an error is generated
// only if there is an error (not for ivalid adr, that is just sanitised)
func (f *File) EvalAdr(pa padr.Node) (Adr, error) {
	// TODO: Make sure that the adr is sanitised!
	a, e := f.evalAdr(pa, false, 0)
	return f.sanitiseAdr(a), e
}

func (f *File) evalAdr(pa padr.Node, reverse bool, start int) (Adr, error) {
	if reverse {
		pa.Val = -pa.Val
	}
	switch pa.Type {
	case padr.DOLLAR:
		return Adr{Start: f.f.Length(), End: f.f.Length()}, nil
	case padr.DOT:
		return f.dot, nil
	case padr.ADRMARK:
		return adrmark, nil
	case padr.HASH:
		no, _ := f.f.OffsetRune(pa.Val, start)
		return Adr{Start: no, End: no}, nil
	case padr.NUMBER:
		if pa.Val == 0 {
			ind, _ := f.f.OffsetLine(0, start)
			return Adr{Start: ind, End: ind}, nil
		} else if pa.Val > 0 {
			pa.Val--
			st, _ := f.f.OffsetLine(pa.Val, start)
			en, _ := f.f.OffsetLine(1, st)
			return Adr{Start: st, End: en}, nil
		} else if pa.Val < 0 {
			pa.Val++
			en, _ := f.f.OffsetLine(pa.Val, start)
			st, _ := f.f.OffsetLine(-1, en)
			return Adr{Start: st, End: en}, nil
		}
	case padr.REG:
		if reverse {
			reg, e := regexp.CompileReverse(pa.Reg)
			if e != nil {
				return Adr{}, e
			}
			rd := f.f.BackwardsReader(start)
			ind := reg.FindReaderIndex(rd)
			if ind == nil { // No match
				return Adr{Start: f.f.Length(), End: f.f.Length()}, nil
			}
			a := Adr{Start: -ind[1] + start, End: -ind[0] + start}
			return a, nil
		}
		reg, e := regexp.Compile(pa.Reg)
		if e != nil {
			return Adr{}, e
		}
		f.f.Seek(int64(start), 0)
		ind := reg.FindInputIndex(f.f)
		if ind == nil { // No match
			return Adr{Start: f.f.Length(), End: f.f.Length()}, nil
		}
		return Adr{Start: ind[0] + start, End: ind[1] + start}, nil
	case padr.PLUS:
		la, le := f.evalAdr(*pa.Left, false, start)
		if le != nil {
			return Adr{}, le
		}
		la = f.sanitiseAdr(la)
		return f.evalAdr(*pa.Right, false, la.End)
	case padr.MINUS:
		la, le := f.evalAdr(*pa.Left, false, start)
		if le != nil {
			return Adr{}, le
		}
		la = f.sanitiseAdr(la)
		return f.evalAdr(*pa.Right, true, la.Start)
	case padr.COMMA:
		la, le := f.evalAdr(*pa.Left, false, start)
		if le != nil {
			return Adr{}, le
		}
		ra, re := f.evalAdr(*pa.Right, false, start)
		if re != nil {
			return Adr{}, re
		}
		return Adr{Start: la.Start, End: ra.End}, nil
	case padr.SEMI:
		la, le := f.evalAdr(*pa.Left, false, start)
		if le != nil {
			return Adr{}, le
		}
		f.dot = f.sanitiseAdr(la)
		ra, re := f.evalAdr(*pa.Right, false, start)
		if re != nil {
			return Adr{}, re
		}
		return Adr{Start: la.Start, End: ra.End}, nil
	default:
		panic("Unkown adr operator")
	}

	spew.Dump(pa)
	return Adr{}, nil
}

// This clamps it to the bounds of the file
func (f *File) sanitiseAdr(a Adr) Adr {
	if a.Start > f.f.Length() {
		return Adr{f.f.Length(), f.f.Length()}
	}
	if a.Start < 0 {
		a.Start = 0
	}
	if a.End < a.Start {
		a.End = a.Start
	}
	if a.End > f.f.Length() {
		a.End = f.f.Length()
	}
	return a
}
