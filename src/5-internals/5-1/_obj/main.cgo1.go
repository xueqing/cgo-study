// Code generated by cmd/cgo; DO NOT EDIT.

//line main.go:1:1
package main

/*
#include "add.h"

static int addByC(int a, int b) {
	return a+b;
}
*/
import _ "unsafe"

import "fmt"

func main() {
	// run `go tool cgo main.go add.go` to generate intermediate files
	fmt.Println(( /*line :16:14*/_Cfunc_addByC /*line :16:21*/)(1, 1))
	fmt.Println(( /*line :17:14*/_Cfunc_addByGo /*line :17:22*/)(1, 1))
	noCgo()
}