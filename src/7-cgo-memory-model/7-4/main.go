package main

/*
extern int* getGoPtr();

static void Main() {
	int* p = getGoPtr();
	*p = 42;
}
*/
import "C"

func main() {
	C.Main() //panic: runtime error: cgo result has Go pointer
	// use `GODEBUG=cgocheck=0 go run .` will not throw exception
}

//export getGoPtr
func getGoPtr() *C.int {
	return new(C.int)
}
