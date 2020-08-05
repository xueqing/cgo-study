package main

// extern int add(int a, int b);
import "C"
import "fmt"

//export add
func add(a, b C.int) C.int {
	return a + b
}

func main() {
	fmt.Println(C.add(1, 1))
}
