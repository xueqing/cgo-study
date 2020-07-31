package main

/*
#include <stdint.h>

struct A {
	int i;
	float f;

	int type;// key word in Go, access by 'A._type'
};

struct B {
	int type;
	float _type;// with this, cannot access 'int type'
};

struct D {
	int size:10;
	float arr[];
};

union A1 {
	int i;
	float f;
};

union B1 {
	int8_t i8;
	int64_t i64;
};

enum D1 {
	ONE,
	TWO,
};
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	var a C.struct_A
	fmt.Printf("type: %T\n", a) // main._Ctype_struct_A
	println(a.i)
	println(a.f)
	println(a._type) // A.type
	println("========A end")

	var b C.struct_B
	println(b._type) // B._type
	println("========B end")

	// var d C.struct_D
	// println(d.size) // d.size undefined (type _Ctype_struct_D has no field or method size)
	// println(d.arr)  // d.arr undefined (type _Ctype_struct_D has no field or method arr)
	// println("========D end")

	var a1 C.union_A1
	fmt.Printf("type: %T\n", a1) // [4]uint8
	// println(a1.q)                //a1.q undefined (type [4]byte has no field or method q)
	// println(a1.f)                //a1.f undefined (type [4]byte has no field or method f)
	fmt.Println("========A1 end")

	var b1 C.union_B1
	fmt.Printf("type: %T\n", b1)                // [8]uint8
	fmt.Println(*(*C.int)(unsafe.Pointer(&b1))) // use unsafe to access C union
	fmt.Println(*(*C.float)(unsafe.Pointer(&b1)))
	fmt.Println("========B1 end")

	var d1 C.enum_D1 = C.TWO
	fmt.Printf("type: %T\n", d1) // uint32
	fmt.Println(d1)              // 1
	fmt.Println(C.ONE)           // 0
	fmt.Println(C.TWO)           // 1
	fmt.Println("========D1 end")
}
