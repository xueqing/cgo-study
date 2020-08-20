package main

import (
	"fmt"
	"unsafe"

	"github.com/xueqing/cgo-study/src/6-encapsulate-qsort/6-3/qsort"
)

func main() {
	arr := []int32{2, 4, 9, 7, 1}
	qsort.Sort((unsafe.Pointer)(&arr[0]), len(arr), int(unsafe.Sizeof(arr[0])),
		func(a, b unsafe.Pointer) int {
			pa, pb := (*int32)(a), (*int32)(b)
			return int(*pa - *pb)
		},
	)

	fmt.Println(arr)
}
