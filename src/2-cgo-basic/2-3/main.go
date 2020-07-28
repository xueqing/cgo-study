// +build debug

package main

func main() {
	// use `go build` error: "build .: cannot find module for path ."
	// use `go build -tags="debug"`
	// or use `go build -tags="linux debug"`
	println("debug mode")
}
