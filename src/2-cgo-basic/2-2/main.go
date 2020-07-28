package main

/*
#cgo windows CFLAGS: -DCGO_WINDOWS=1
#cgo darwin CFLAGS: -DCGO_DARWIN=1
#cgo linux CFLAGS: -DCGO_LINUX=1

#if defined(CGO_WINDOWS)
	const char *os="windows";
#elif defined(CGO_DARWIN)
	const char *os="darwin";
#elif defined(CGO_LINUX)
	const char *os="linux";
#else
	const char *os="unknown";
#endif
*/
import "C"

func main() {
	println(C.GoString(C.os))
}
