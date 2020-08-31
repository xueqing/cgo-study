package main

/*
#cgo CFLAGS: -I./number
#cgo LDFLAGS: -L${SRCDIR}/number -lnumber

#include "number.h"
*/
import "C"
import "fmt"

func main() {
	fmt.Println(int(C.number_add_mod(5, 10, 4)))
}
