package main

// extern int go_qsort_compare(void *a, void *b);
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/xueqing/cgo-study/src/6-encapsulate-qsort/6-2/qsort"
)

//export go_qsort_compare
func go_qsort_compare(a, b unsafe.Pointer) C.int {
	pa, pb := (*C.int)(a), (*C.int)(b)
	return C.int(*pa - *pb)
}

func main() {
	arr := []int32{2, 4, 9, 7, 1}
	qsort.Sort(unsafe.Pointer(&arr[0]),
		len(arr), int(unsafe.Sizeof(arr[0])),
		qsort.CompareFunc(C.go_qsort_compare))
	fmt.Println(arr)
}
