package main

//#include "hello.h"
import "C"

func main() {
	C.SayHello(C.CString("hello, world\n"))
}
