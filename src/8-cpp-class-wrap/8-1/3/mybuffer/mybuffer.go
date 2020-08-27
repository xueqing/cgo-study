package mybuffer

import "unsafe"

// MyBuffer ...
type MyBuffer struct {
	cptr *cgoMyBufferT
}

// NewMyBuffer ...
func NewMyBuffer(size int) *MyBuffer {
	return &MyBuffer{
		cptr: cgoNewMyBuffer(size),
	}
}

// Delete ...
func (p *MyBuffer) Delete() {
	cgoDeleteMyBuffer(p.cptr)
}

// Data ...
func (p *MyBuffer) Data() []byte {
	data := cgoMyBufferData(p.cptr)
	size := cgoMyBufferSize(p.cptr)
	return ((*[1 << 31]byte)(unsafe.Pointer(data)))[0:int(size):int(size)]
}
