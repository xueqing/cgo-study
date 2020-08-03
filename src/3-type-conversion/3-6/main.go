package main

import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	var (
		i, j int32
		c    *C.char
	)

	fmt.Println(i, j)
	fmt.Println(c)

	i = 10
	c = (*C.char)(unsafe.Pointer(uintptr(i)))
	fmt.Println(i, j)
	fmt.Println(c)

	j = (int32)(uintptr(unsafe.Pointer(c)))
	fmt.Println(i, j)
	fmt.Println(c)
}
