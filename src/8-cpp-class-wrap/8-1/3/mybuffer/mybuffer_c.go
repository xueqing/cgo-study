package mybuffer

/*
#cgo CXXFLAGS: -std=c++11

#include "mybuffer_c.h"
*/
import "C"

type cgoMyBufferT C.MyBuffer_T

func cgoNewMyBuffer(size int) *cgoMyBufferT {
	p := C.NewMyBuffer(C.int(size))
	return (*cgoMyBufferT)(p)
}

func cgoDeleteMyBuffer(p *cgoMyBufferT) {
	C.DeleteMyBuffer((*C.MyBuffer_T)(p))
}

func cgoMyBufferSize(p *cgoMyBufferT) C.int {
	return C.MyBuffer_Size((*C.MyBuffer_T)(p))
}

func cgoMyBufferData(p *cgoMyBufferT) *C.char {
	return C.MyBuffer_Data((*C.MyBuffer_T)(p))
}
