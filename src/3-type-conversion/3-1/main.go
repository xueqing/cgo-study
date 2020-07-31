package main

/*
#include <stdio.h>

void printTypeSize() {
	// c   sizeof char:  1 1 1
	printf("c   sizeof char:  %lu %lu %lu\n", sizeof(char), sizeof(signed char), sizeof(unsigned char));
	// c   sizeof int:   2 2 4 4
	printf("c   sizeof int:   %lu %lu %lu %lu\n", sizeof(short), sizeof(unsigned short), sizeof(int), sizeof(unsigned int));
	// c   sizeof long:  8 8 8 8
	printf("c   sizeof long:  %lu %lu %lu %lu\n", sizeof(long), sizeof(unsigned long), sizeof(long long int), sizeof(unsigned long long int));
	// c   sizeof float: 4 8 8
	printf("c   sizeof float: %lu %lu %lu\n", sizeof(float), sizeof(double), sizeof(size_t));
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	C.printTypeSize()

	// cgo sizeof char:  1 1 1
	fmt.Printf("cgo sizeof char:  %v %v %v\n", unsafe.Sizeof(*(new(C.char))), unsafe.Sizeof(*(new(C.schar))), unsafe.Sizeof(*(new(C.uchar))))
	// cgo sizeof int:   2 2 4 4
	fmt.Printf("cgo sizeof int:   %v %v %v %v\n", unsafe.Sizeof(*(new(C.short))), unsafe.Sizeof(*(new(C.ushort))), unsafe.Sizeof(*(new(C.int))), unsafe.Sizeof(*(new(C.uint))))
	// cgo sizeof long:  8 8 8 8
	fmt.Printf("cgo sizeof long:  %v %v %v %v\n", unsafe.Sizeof(*(new(C.long))), unsafe.Sizeof(*(new(C.ulong))), unsafe.Sizeof(*(new(C.longlong))), unsafe.Sizeof(*(new(C.ulonglong))))
	// cgo sizeof float: 4 8 8
	fmt.Printf("cgo sizeof float: %v %v %v\n", unsafe.Sizeof(*(new(C.float))), unsafe.Sizeof(*(new(C.double))), unsafe.Sizeof(*(new(C.size_t))))

	// go  sizeof char:  1 1 1
	fmt.Printf("go  sizeof char:  %v %v %v\n", unsafe.Sizeof(*(new(byte))), unsafe.Sizeof(*(new(int8))), unsafe.Sizeof(*(new(uint8))))
	// go  sizeof int:   2 2 4 4
	fmt.Printf("go  sizeof int:   %v %v %v %v\n", unsafe.Sizeof(*(new(int16))), unsafe.Sizeof(*(new(uint16))), unsafe.Sizeof(*(new(int32))), unsafe.Sizeof(*(new(uint32))))
	// go  sizeof long:  4 4 8 8
	fmt.Printf("go  sizeof long:  %v %v %v %v\n", unsafe.Sizeof(*(new(int32))), unsafe.Sizeof(*(new(uint32))), unsafe.Sizeof(*(new(int64))), unsafe.Sizeof(*(new(uint64))))
	// go  sizeof float: 4 8 8
	fmt.Printf("go  sizeof float: %v %v %v\n", unsafe.Sizeof(*(new(float32))), unsafe.Sizeof(*(new(float64))), unsafe.Sizeof(*(new(uint))))
}
