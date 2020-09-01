package main

import "C"

import (
	"fmt"

	_ "github.com/xueqing/cgo-study/src/9-static-dynamic-lib/5/number"
)

func main() {
	fmt.Println("Done")
}

//export goPrintln
func goPrintln(s *C.char) {
	fmt.Println("goPrintln:", C.GoString(s))
}
