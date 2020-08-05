package main

/*
#include <errno.h>

static int div(int a, int b) {
	if (b == 0) {
		errno = EINVAL;
		return 0;
	}
	return a/b;
}
*/
import "C"
import "fmt"

func main() {
	v1, e1 := C.div(4, 2)
	fmt.Println(v1, e1) // 2 <nil>

	v2, e2 := C.div(4, 0)
	fmt.Println(v2, e2) // 0 invalid argument
}
