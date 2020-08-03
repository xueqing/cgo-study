package main

/*
#include <string.h>
char arr[10]={1, 2, 3};
char *s="Hello";
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	// use reflect.SliceHeader
	var arr0 []byte
	var arr0Hdr = (*reflect.SliceHeader)(unsafe.Pointer(&arr0))
	arr0Hdr.Data = uintptr(unsafe.Pointer(&C.arr[0]))
	arr0Hdr.Len = 10
	arr0Hdr.Cap = 10
	fmt.Println(arr0, len(arr0), cap(arr0), reflect.TypeOf(arr0))

	// use slice directly
	arr1 := (*[31]byte)(unsafe.Pointer(&C.arr[0]))[:10:10]
	fmt.Println(arr1, len(arr1), cap(arr1), reflect.TypeOf(arr1))
	arr2 := (*[31]byte)(unsafe.Pointer(&C.arr[0]))[:15:20]
	fmt.Println(arr2, len(arr2), cap(arr2), reflect.TypeOf(arr2))
	arr3 := (*[31]byte)(unsafe.Pointer(&C.arr[0]))[:5:20]
	fmt.Println(arr3, len(arr3), cap(arr3), reflect.TypeOf(arr3))

	// use reflect.StringHeader
	var s0 string
	var s0Hdr = (*reflect.StringHeader)(unsafe.Pointer(&s0))
	s0Hdr.Data = uintptr(unsafe.Pointer(C.s))
	s0Hdr.Len = int(C.strlen(C.s))
	fmt.Println(s0, len(s0), reflect.TypeOf(s0))

	// use string directly
	sLen := int(C.strlen(C.s))
	s1 := string((*[31]byte)(unsafe.Pointer(C.s))[:sLen:sLen])
	fmt.Println(s1, len(s1), reflect.TypeOf(s1))
}
