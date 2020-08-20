package main

import (
	"fmt"

	"github.com/xueqing/cgo-study/src/6-encapsulate-qsort/6-4/qsort"
)

func main() {
	arr := []int32{2, 4, 9, 7, 1}
	qsort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})

	fmt.Println(arr)
}
