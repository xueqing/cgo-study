package main

/*
#include <stdio.h>

void printint(int v) {
	printf("printint: %d\n", v);
}

// add static will compile error: "undefined reference to `cs'"
const char* cs="hello";
*/
import "C"

import _ "github.com/xueqing/cgo-study/2-cgo-basic/2-1/cgo_helper"

func main() {
	v := 42
	C.printint(C.int(v))

	// ./main.go:21:26: cannot use *_Cvar_cs (type *_Ctype_char) as type *cgo_helper._Ctype_char in argument to cgo_helper.PrintCString
	// cgo_helper.PrintCString(C.cs)
}
