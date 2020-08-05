package main

// static void noreturn() {}
import "C"
import "fmt"

func main() {
	v, e := C.noreturn()
	fmt.Println(v, e)      // [] <nil>
	fmt.Printf("%#v\n", v) // main._Ctype_void{}
}
