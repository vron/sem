package perfreg

import (
	"bytes"
	"regexp"
	"testing"
)

var reg *regexp.Regexp

func init() {
	var e error
	reg, e = regexp.Compile("(Su)+")
	if e != nil {
		panic(e.Error())
	}
}

func BenchmarkReader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buf := bytes.NewReader(data)
		_ = buf.Len()
		b.StartTimer()
		for {
			loc := reg.FindReaderIndex(buf)
			if loc == nil {
				break
			}
		}
	}
}
func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buf := bytes.NewReader(data)
		_ = buf.Len()
		b.StartTimer()
		start := 0
		for {
			loc := reg.FindIndex(data[start:])
			if loc == nil {
				break
			}
			start += loc[1]
		}
	}
}

var data = []byte(`
THIS IS A TEXT, TAKEN FROM PKG TESTING DOC

The Go Programming Language

Package testing

import "testing"
Overview
Index
Subdirectories
Overview ▾

Package testing provides support for automated testing of Go packages. It is intended to be used in concert with the “go test” command, which automates execution of any function of the form

func TestXxx(*testing.T)
where Xxx can be any alphanumeric string (but the first letter must not be in [a-z]) and serves to identify the test routine. These TestXxx routines should be declared within the package they are testing.

Functions of the form

func BenchmarkXxx(*testing.B)
are considered benchmarks, and are executed by the "go test" command when the -test.bench flag is provided.

A sample benchmark function looks like this:

func BenchmarkHello(b *testing.B) {
    for i := 0; i < b.N; i++ {
        fmt.Sprintf("hello")
    }
}
The benchmark package will vary b.N until the benchmark function lasts long enough to be timed reliably. The output

testing.BenchmarkHello    10000000    282 ns/op
means that the loop ran 10000000 times at a speed of 282 ns per loop.

If a benchmark needs some expensive setup before running, the timer may be stopped:

func BenchmarkBigLen(b *testing.B) {
    b.StopTimer()
    big := NewBig()
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        big.Len()
    }
}
The package also runs and verifies example code. Example functions may include a concluding comment that begins with "Output:" and is compared with the standard output of the function when the tests are run, as in these examples of an example:

func ExampleHello() {
        fmt.Println("hello")
        // Output: hello
}

func ExampleSalutations() {
        fmt.Println("hello, and")
        fmt.Println("goodbye")
        // Output:
        // hello, and
        // goodbye
}
Example functions without output comments are compiled but not executed.

The naming convention to declare examples for a function F, a type T and method M on type T are:

func ExampleF() { ... }
func ExampleT() { ... }
func ExampleT_M() { ... }
Multiple example functions for a type/function/method may be provided by appending a distinct suffix to the name. The suffix must start with a lower-case letter.

func ExampleF_suffix() { ... }
func ExampleT_suffix() { ... }
func ExampleT_M_suffix() { ... }
The entire test file is presented as the example when it contains a single example function, at least one other function, type, variable, or constant declaration, and no test or benchmark functions.

Index ▾

func Main(matchString func(pat, str string) (bool, error), tests []InternalTest, benchmarks []InternalBenchmark, examples []InternalExample)
func RunBenchmarks(matchString func(pat, str string) (bool, error), benchmarks []InternalBenchmark)
func RunExamples(matchString func(pat, str string) (bool, error), examples []InternalExample) (ok bool)
func RunTests(matchString func(pat, str string) (bool, error), tests []InternalTest) (ok bool)
func Short() bool
type B
    func (c *B) Error(args ...interface{})
    func (c *B) Errorf(format string, args ...interface{})
    func (c *B) Fail()
    func (c *B) FailNow()
    func (c *B) Failed() bool
    func (c *B) Fatal(args ...interface{})
    func (c *B) Fatalf(format string, args ...interface{})
    func (c *B) Log(args ...interface{})
    func (c *B) Logf(format string, args ...interface{})
    func (b *B) ResetTimer()
    func (b *B) SetBytes(n int64)
    func (b *B) StartTimer()
    func (b *B) StopTimer()
type BenchmarkResult
    func Benchmark(f func(b *B)) BenchmarkResult
    func (r BenchmarkResult) NsPerOp() int64
    func (r BenchmarkResult) String() string
type InternalBenchmark
type InternalExample
type InternalTest
type T
    func (c *T) Error(args ...interface{})
    func (c *T) Errorf(format string, args ...interface{})
    func (c *T) Fail()
    func (c *T) FailNow()
    func (c *T) Failed() bool
    func (c *T) Fatal(args ...interface{})
    func (c *T) Fatalf(format string, args ...interface{})
    func (c *T) Log(args ...interface{})
    func (c *T) Logf(format string, args ...interface{})
    func (t *T) Parallel()
Package files

benchmark.go example.go testing.go

func Main

func Main(matchString func(pat, str string) (bool, error), tests []InternalTest, benchmarks []InternalBenchmark, examples []InternalExample)
An internal function but exported because it is cross-package; part of the implementation of the "go test" command.

func RunBenchmarks

func RunBenchmarks(matchString func(pat, str string) (bool, error), benchmarks []InternalBenchmark)
An internal function but exported because it is cross-package; part of the implementation of the "go test" command.

func RunExamples

func RunExamples(matchString func(pat, str string) (bool, error), examples []InternalExample) (ok bool)
func RunTests

func RunTests(matchString func(pat, str string) (bool, error), tests []InternalTest) (ok bool)
func Short

func Short() bool
Short reports whether the -test.short flag is set.

type B

type B struct {
    N int
    // contains filtered or unexported fields
}
B is a type passed to Benchmark functions to manage benchmark timing and to specify the number of iterations to run.

func (*B) Error

func (c *B) Error(args ...interface{})
Error is equivalent to Log() followed by Fail().

func (*B) Errorf

func (c *B) Errorf(format string, args ...interface{})
Errorf is equivalent to Logf() followed by Fail().

func (*B) Fail

func (c *B) Fail()
Fail marks the function as having failed but continues execution.

func (*B) FailNow

func (c *B) FailNow()
FailNow marks the function as having failed and stops its execution. Execution will continue at the next test or benchmark.

func (*B) Failed

func (c *B) Failed() bool
Failed returns whether the function has failed.

func (*B) Fatal

func (c *B) Fatal(args ...interface{})
Fatal is equivalent to Log() followed by FailNow().

func (*B) Fatalf

func (c *B) Fatalf(format string, args ...interface{})
Fatalf is equivalent to Logf() followed by FailNow().

func (*B) Log

func (c *B) Log(args ...interface{})
Log formats its arguments using default formatting, analogous to Println(), and records the text in the error log.

func (*B) Logf

func (c *B) Logf(format string, args ...interface{})
Logf formats its arguments according to the format, analogous to Printf(), and records the text in the error log.

func (*B) ResetTimer

func (b *B) ResetTimer()
ResetTimer sets the elapsed benchmark time to zero. It does not affect whether the timer is running.

func (*B) SetBytes

func (b *B) SetBytes(n int64)
SetBytes records the number of bytes processed in a single operation. If this is called, the benchmark will report ns/op and MB/s.

func (*B) StartTimer

func (b *B) StartTimer()
StartTimer starts timing a test. This function is called automatically before a benchmark starts, but it can also used to resume timing after a call to StopTimer.

func (*B) StopTimer

func (b *B) StopTimer()
StopTimer stops timing a test. This can be used to pause the timer while performing complex initialization that you don't want to measure.

type BenchmarkResult

type BenchmarkResult struct {
    N     int           // The number of iterations.
    T     time.Duration // The total time taken.
    Bytes int64         // Bytes processed in one iteration.
}
The results of a benchmark run.

func Benchmark

func Benchmark(f func(b *B)) BenchmarkResult
Benchmark benchmarks a single function. Useful for creating custom benchmarks that do not use the "go test" command.

func (BenchmarkResult) NsPerOp

func (r BenchmarkResult) NsPerOp() int64
func (BenchmarkResult) String

func (r BenchmarkResult) String() string
type InternalBenchmark

type InternalBenchmark struct {
    Name string
    F    func(b *B)
}
An internal type but exported because it is cross-package; part of the implementation of the "go test" command.

type InternalExample

type InternalExample struct {
    Name   string
    F      func()
    Output string
}
type InternalTest

type InternalTest struct {
    Name string
    F    func(*T)
}
An internal type but exported because it is cross-package; part of the implementation of the "go test" command.

type T

type T struct {
    // contains filtered or unexported fields
}
T is a type passed to Test functions to manage test state and support formatted test logs. Logs are accumulated during execution and dumped to standard error when done.

func (*T) Error

func (c *T) Error(args ...interface{})
Error is equivalent to Log() followed by Fail().

func (*T) Errorf

func (c *T) Errorf(format string, args ...interface{})
Errorf is equivalent to Logf() followed by Fail().

func (*T) Fail

func (c *T) Fail()
Fail marks the function as having failed but continues execution.

func (*T) FailNow

func (c *T) FailNow()
FailNow marks the function as having failed and stops its execution. Execution will continue at the next test or benchmark.

func (*T) Failed

func (c *T) Failed() bool
Failed returns whether the function has failed.

func (*T) Fatal

func (c *T) Fatal(args ...interface{})
Fatal is equivalent to Log() followed by FailNow().

func (*T) Fatalf

func (c *T) Fatalf(format string, args ...interface{})
Fatalf is equivalent to Logf() followed by FailNow().

func (*T) Log

func (c *T) Log(args ...interface{})
Log formats its arguments using default formatting, analogous to Println(), and records the text in the error log.

func (*T) Logf

func (c *T) Logf(format string, args ...interface{})
Logf formats its arguments according to the format, analogous to Printf(), and records the text in the error log.

func (*T) Parallel

func (t *T) Parallel()
Parallel signals that this test is to be run in parallel with (and only with) other parallel tests in this CPU group.

Subdirectories

Name	    	Synopsis
..
iotest	    	Package iotest implements Readers and Writers useful mainly for testing.
quick	    	Package quick implements utility functions to help with black box testing.
Build version go1.0.3.
Except as noted, the content of this page is licensed under the Creative Commons Attribution 3.0 License, and code is licensed under a BSD license.
Terms of Service | Privacy Policy
`)
