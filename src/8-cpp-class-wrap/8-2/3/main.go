package main

import "fmt"

// MyInt ...
type MyInt int

// Twice ...
func (i MyInt) Twice() int {
	return (int)(i) << 1
}

func main() {
	var x = MyInt(3)
	fmt.Println(int(x))
	fmt.Println(x.Twice())
}
