package main

import "fmt"
import "unsafe"

type b struct {
	a int
	b [10][2]int
}

func main() {
	a := b{}
	fmt.Println(unsafe.Sizeof(a))
}
