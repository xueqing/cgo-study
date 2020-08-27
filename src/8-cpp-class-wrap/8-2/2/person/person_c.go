package person

/*
#cgo CXXFLAGS: -std=c++11

#include "./person_c.h"
*/
import "C"
import (
	"unsafe"

	"github.com/xueqing/cgo-study/src/8-cpp-class-wrap/8-2/2/object"
)

//export new_person
func new_person(name *C.char, age C.int) C.person_handle_t {
	id := object.NewID(NewPerson(C.GoString(name), int(age)))
	return C.person_handle_t(id)
}

//export delete_person
func delete_person(h C.person_handle_t) {
	id := object.ID(h)
	id.Free()
}

//export person_set
func person_set(h C.person_handle_t, name *C.char, age C.int) {
	p := object.ID(h).Get().(*Person)
	p.Set(C.GoString(name), int(age))
}

//export person_get_name
func person_get_name(h C.person_handle_t, buf *C.char, size C.int) *C.char {
	p := object.ID(h).Get().(*Person)
	name, _ := p.Get()

	n := int(size) - 1
	bufSlice := ((*[1 << 31]byte)(unsafe.Pointer(buf)))[0:n:n]
	copy(bufSlice, []byte(name))
	bufSlice[n] = 0

	return buf
}

//export person_get_age
func person_get_age(h C.person_handle_t) C.int {
	p := object.ID(h).Get().(*Person)
	_, age := p.Get()

	return C.int(age)
}
