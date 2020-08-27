package main

/*
#cgo CXXFLAGS: -std=c++11

#include "person/person_test_c.h"
*/
import "C"

func main() {
	/*
		$ go run main.go
		# command-line-arguments
		/tmp/go-build237897354/b001/_x002.o: In function `_cgo_b7e481d83ce7_Cfunc_test_new_person':
		main.cgo2.c:(.text+0x50): undefined reference to `test_new_person'
		collect2: error: ld returned 1 exit status
	*/
	C.test_new_person()
}
