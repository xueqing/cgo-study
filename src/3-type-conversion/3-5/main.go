package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var (
		i int
		u uint64
	)

	fmt.Println(i)
	fmt.Println(u)

	i = 10
	u = *(*uint64)(unsafe.Pointer(&i))
	fmt.Println(i)
	fmt.Println(u)

	u = 20
	i = *(*int)(unsafe.Pointer(&u))
	fmt.Println(i)
	fmt.Println(u)
}
