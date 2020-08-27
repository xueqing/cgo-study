package main

//#include <stdio.h>
import "C"
import (
	"unsafe"

	"github.com/xueqing/cgo-study/src/8-cpp-class-wrap/8-1/3/mybuffer"
)

func main() {
	buf := mybuffer.NewMyBuffer(1024)
	defer buf.Delete()

	copy(buf.Data(), []byte("hello\x00"))
	C.puts((*C.char)(unsafe.Pointer(&(buf.Data()[0]))))
}
