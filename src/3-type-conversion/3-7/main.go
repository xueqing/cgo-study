package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type x struct {
	i, j int
}

type y struct {
	i int
}

func main() {
	p := make([]x, 5, 5)
	q := make([]y, 20, 20)

	for idx := 1; idx <= 15; idx++ {
		q[idx-1] = y{i: idx}
	}

	fmt.Println(len(p), cap(p), p)
	fmt.Println(len(q), cap(q), q)

	pHdr := (*reflect.SliceHeader)(unsafe.Pointer(&p))
	qHdr := (*reflect.SliceHeader)(unsafe.Pointer(&q))
	pHdr.Data = qHdr.Data
	pHdr.Len = qHdr.Len * int(unsafe.Sizeof(q[0])) / int(unsafe.Sizeof(p[0]))
	pHdr.Cap = qHdr.Cap * int(unsafe.Sizeof(q[0])) / int(unsafe.Sizeof(p[0]))

	fmt.Println(len(p), cap(p), p)
	fmt.Println(len(q), cap(q), q)
}
