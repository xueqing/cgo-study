package main

import "C"

// export helloString
func helloString(s string) {}

// export helloSlice
func helloSlice(s []byte) {}

func main() {
	println("Hello, world")
}