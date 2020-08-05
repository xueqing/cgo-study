package main

import "C"

//export addByGo
func addByGo(a, b C.int) C.int {
	return a + b
}
