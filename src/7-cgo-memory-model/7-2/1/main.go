package main

/*
#include <stdlib.h>
#include <stdio.h>

void printString(const char *s) {
	printf("%s\n", s);
}
*/
import "C"
import "unsafe"

func printString(s string) {
	cs := C.CString(s) // copy go string to C string(memory allocated by C)
	defer C.free(unsafe.Pointer(cs))

	C.printString(cs)
}

func main() {
	printString("hello")
}
