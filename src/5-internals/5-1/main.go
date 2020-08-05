package main

/*
#include "add.h"

static int addByC(int a, int b) {
	return a+b;
}
*/
import "C"

import "fmt"

func main() {
	// run `go tool cgo main.go add.go` to generate intermediate files
	fmt.Println(C.addByC(1, 1))
	fmt.Println(C.addByGo(1, 1))
	noCgo()
}
