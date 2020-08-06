package main

//int sum(int a, int b);
import "C"
import "fmt"

//export sum
func sum(a, b C.int) C.int {
	return a + b
}

func main() {
	// run `go build -buildmode=c-archive -o sum.a main.go`
	// to generate header(sum.h) and library(sum.a) files
	// run `go tool cgo main.go` to generate intermediate files
	fmt.Println(C.sum(1, 1))
}
