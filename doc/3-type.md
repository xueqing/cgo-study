# 类型转换

- [类型转换](#类型转换)
  - [数值类型](#数值类型)
  - [Go 字符串和切片](#go-字符串和切片)
  - [结构体、联合、枚举类型](#结构体联合枚举类型)
    - [结构体](#结构体)
    - [联合体](#联合体)
    - [枚举](#枚举)
  - [数组、字符串和切片](#数组字符串和切片)
  - [指针间的转换](#指针间的转换)
  - [数值和指针的转换](#数值和指针的转换)
  - [切片间的转换](#切片间的转换)

## 数值类型

Go 语言通过虚拟的 “C” 包访问 C 语言的符号，比如 `C.int` 对应 C 语言的 `int` 类型。通过虚拟 “C” 包访问 C 语言类型时名称部分不能有空格字符，cgo 为 C 语言的基础数值类型提供了相应转换规则，比如 `C.uint` 对应 C 语言的 `unsigned int`。

Go 语言和 C 语言数据类型对比：

| C 语言类型 | cgo 类型 | Go 语言类型 |
| --- | --- | --- |
| char | C.char | byte |
| singed char | C.schar | int8 |
| unsigned char | C.uchar | uint8 |
| short | C.short | int16 |
| unsigned short | C.ushort | uint16 |
| int | C.int | int32 |
| unsigned int | C.uint | uint32 |
| long | C.long | int32 |
| unsigned long | C.ulong | uint32 |
| long long int | C.longlong | int64 |
| unsigned long long int | C.ulonglong | uint64 |
| float | C.float | float32 |
| double | C.double | float64 |
| size_t | C.size_t | uint |

Go 语言类型 `<stdint.h>` 头文件类型对比：

| C 语言类型 | cgo 类型 | Go 语言类型 |
| --- | --- | --- |
| int8_t | C.int8_t | int8 |
| uint8_t | C.uint8_t | uint8 |
| int16_t | C.int16_t | int16 |
| uint16_t | C.uint16_t | uint16 |
| int32_t | C.int32_t | int32 |
| uint32_t | C.uint32_t | uint32 |
| int64_t | C.int64_t | int64 |
| uint64_t | C.uint64_t | uint64 |

在 `_cgo_export.h` 头文件中，每个基本的 Go 数值类型都定义了对应的 C 语言类型，它们一般都是以单词 `Go` 为前缀。

## Go 字符串和切片

在 cgo 生成的 `_cgo_export.h` 头文件中还会为 Go 语言的字符串、切片、字典、接口和管道等特有的数据类型生成对应的 C 语言类型：

```go
typedef struct { const char *p; GoInt n; } GoString;
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;
```

但是只有字符串和切片再 cgo 中有一定的使用价值：

- cgo 为字符串和切片的某些 Go 语言版本的操作函数生成了 C 语言版本，因此二者可以在 Go 调用 C 语言函数时马上使用;而 cgo 并未针对其他的类型提供相关的辅助函数
- Go 语言特有的内存模型导致我们无法保持这些由 Go 语言管理的内存指针，所以 C 语言环境并无使用价值

在导出的 C 语言函数中可直接使用 Go 字符串和切片。

## 结构体、联合、枚举类型

C 语言的结构体、联合、枚举类型不能作为匿名成员嵌入到 Go 语言的结构体。Go 语言通过 `C.struct_xxx` 访问 C 语言定义的 `struct xxx` 类型。结构体内存布局按照 C 语言的通用对齐规则。cgo 无法访问指定了特殊对齐规则的结构体。

### 结构体

如果 C 语言的结构体成员命名是 Go 的关键字，可通过在成员名开头添加下划线进行访问；如果 C 语言的结构体刚好包含另外一个以下划线开头的同名成员，那么以 Go 关键字命名的成员将被屏蔽，在 Go 中无法访问。

C 语言无法访问 Go 中定义的结构体类型。

### 联合体

可通过 `C.union_XXX` 访问 C 语言定义的 `union XXX` 类型。但是 Go 语言不支持 C 语言联合类型，会将其转为对应大小的字节数组。

操作 C 语言联合类型变量有三种方法：

- 在 C 语言定义辅助函数
- 通过 Go 语言的 `encoding/binary` 手工解码成员(注意大端小端问题)
- 使用 `unsafe` 包强制转换为对应类型(性能最好)

### 枚举

可通过 `C.enum_XXX` 访问 C 语言定义的 `enum XXX` 类型。

C 语言枚举类型底层对应 `int` 类型，支持负数类型的值。可通过 `C.ONE`/`C.TWO` 直接访问定义的枚举值。

## 数组、字符串和切片

Go 语言和 C 语言的数组、字符串和切片之间的相互转换可以监护为 Go 的 切片和 C 中指向一定长度内存的指针之间的转换。

CGO的C虚拟包提供了以下一组函数，用于Go语言和C语言之间数组和字符串的双向转换：

```go
// Go string to C string
// The C string is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be
// freed, such as by calling C.free (be sure to include stdlib.h
// if C.free is needed).
func C.CString(string) *C.char

// Go []byte slice to C array
// The C array is allocated in the C heap using malloc.
// It is the caller's responsibility to arrange for it to be
// freed, such as by calling C.free (be sure to include stdlib.h
// if C.free is needed).
func C.CBytes([]byte) unsafe.Pointer

// C string to Go string
func C.GoString(*C.char) string

// C data with explicit length to Go string
func C.GoStringN(*C.char, C.int) string

// C data with explicit length to Go []byte
func C.GoBytes(unsafe.Pointer, C.int) []byte
```

`C.CString` 从输入的 Go 字符串克隆一个 C 语言格式的字符串，返回的字符串由 `C.malloc` 函数分配，不使用时需要通过 `C.free` 函数释放。`C.CBytes` 函数的功能和 `C.CString` 类似，用于从输入的 Go 语言字节切片克隆一个 C 语言版本的字节数组，同样返回的数组需要在合适的时候释放。`C.GoString` 用于从 NULL 结尾的 C 语言字符串克隆一个 Go 语言字符串。`C.GoStringN` 是另一个字符数组克隆函数。`C.GoBytes` 用于从 C 语言数组克隆一个 Go 语言字节切片。

这些辅助函数都是以克隆的方式运行。当 Go 语言字符串和切片向 C 语言转换时，克隆的内存由 `C.malloc` 函数分配，最终可以通过 `C.free` 函数释放。当 C 语言字符串或数组向 Go 语言转换时，克隆的内存由 Go 语言分配管理。通过该组转换函数，转换前和转换后的内存依然在各自的语言环境中，它们并没有跨越 Go 语言和 C 语言。克隆方式实现转换的优点是接口和内存管理都很简单，缺点是克隆需要分配新的内存和复制操作都会导致额外的开销。

在 `reflect` 包中有字符串和切片的定义：

```go
type StringHeader struct {
  Data uintptr
  Len  int
}

type SliceHeader struct {
  Data uintptr
  Len  int
  Cap  int
}
```

如果不希望单独分配内存，可以在 Go 语言中直接访问 C 语言的内存空间：

```go
package main

/*
#include <string.h>
char arr[10]={1, 2, 3};
char *s="Hello";
*/
import "C"
import (
  "fmt"
  "reflect"
  "unsafe"
)

func main() {
  // use reflect.SliceHeader
  var arr0 []byte
  var arr0Hdr = (*reflect.SliceHeader)(unsafe.Pointer(&arr0))
  arr0Hdr.Data = uintptr(unsafe.Pointer(&C.arr[0]))
  arr0Hdr.Len = 10
  arr0Hdr.Cap = 10
  fmt.Println(arr0, len(arr0), cap(arr0), reflect.TypeOf(arr0))

  // use slice directly
  arr1 := (*[31]byte)(unsafe.Pointer(&C.arr[0]))[:10:10]
  fmt.Println(arr1, len(arr1), cap(arr1), reflect.TypeOf(arr1))
  arr2 := (*[31]byte)(unsafe.Pointer(&C.arr[0]))[:15:20]
  fmt.Println(arr2, len(arr2), cap(arr2), reflect.TypeOf(arr2))
  arr3 := (*[31]byte)(unsafe.Pointer(&C.arr[0]))[:5:20]
  fmt.Println(arr3, len(arr3), cap(arr3), reflect.TypeOf(arr3))

  // use reflect.StringHeader
  var s0 string
  var s0Hdr = (*reflect.StringHeader)(unsafe.Pointer(&s0))
  s0Hdr.Data = uintptr(unsafe.Pointer(C.s))
  s0Hdr.Len = int(C.strlen(C.s))
  fmt.Println(s0, len(s0), reflect.TypeOf(s0))

  // use string directly
  sLen := int(C.strlen(C.s))
  s1 := string((*[31]byte)(unsafe.Pointer(C.s))[:sLen:sLen])
  fmt.Println(s1, len(s1), reflect.TypeOf(s1))
}
```

因为 Go 语言的字符串是只读的，用户需要自己保证 Go 字符串在使用期间，底层对应的 C 字符串内容不会发生变化、内存不会被提前释放掉。

在 cgo 中，会为字符串和切片生成和上面结构对应的 C 语言版本的结构体：

```go
typedef struct { const char *p; GoInt n; } GoString;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;
```

在 C 语言中可以通过 `GoString` 和 `GoSlice` 来访问 Go 语言的字符串和切片。如果是 Go 语言中数组类型，可以将数组转为切片后再行转换。如果字符串或切片对应的底层内存空间由 Go 语言运行时管理，那么在 C 语言中不能长时间保存 Go 内存对象。

## 指针间的转换

以下代码演示了 Go 中如何将 X 类型的指针转化为 Y 类型的指针：

```go
var p *X
var q *Y

q = (*Y)(unsafe.Pointer(p)) // *X => *Y
p = (*X)(unsafe.Pointer(q)) // *Y => *X
```

为了实现 X 类型指针到 Y 类型指针的转换，需要借助 `unsafe.Pointer` 作为中间桥接类型实现不同类型指针之间的转换。`unsafe.Pointer` 指针类型类似 C 语言中的 `void*` 类型的指针。

任何类型的指针都可以通过强制转换为 `unsafe.Pointer` 指针类型去掉原有的类型信息，然后再重新赋予新的指针类型而达到指针间的转换的目的。

## 数值和指针的转换

Go 语言禁止将数值类型直接转为指针类型。可以 `uintptr` 为中介，实现数值类型到 `unsafe.Pointr` 指针类型的转换。再结合前面提到的方法，就可以实现数值和指针的转换了。

```go
package main

import "C"

import (
  "fmt"
  "unsafe"
)

func main() {
  var (
    i, j int32
    c    *C.char
  )

  fmt.Println(i, j)
  fmt.Println(c)

  i = 10
  c = (*C.char)(unsafe.Pointer(uintptr(i)))
  fmt.Println(i, j)
  fmt.Println(c)

  j = (int32)(uintptr(unsafe.Pointer(c)))
  fmt.Println(i, j)
  fmt.Println(c)
}
```

## 切片间的转换

Go 语言的 `reflect` 包提供了切片类型的底层结构，再结合前面讨论到不同类型之间的指针转换技术就可以实现 `[]X` 和 `[]Y` 类型的切片转换：

```go
package main

import (
  "fmt"
  "reflect"
  "unsafe"
)

type x struct {
  i, j int
}

type y struct {
  i int
}

func main() {
  p := make([]x, 5, 5)
  q := make([]y, 20, 20)

  for idx := 1; idx <= 15; idx++ {
    q[idx-1] = y{i: idx}
  }

  fmt.Println(len(p), cap(p), p)
  fmt.Println(len(q), cap(q), q)

  pHdr := (*reflect.SliceHeader)(unsafe.Pointer(&p))
  qHdr := (*reflect.SliceHeader)(unsafe.Pointer(&q))
  pHdr.Data = qHdr.Data
  pHdr.Len = qHdr.Len * int(unsafe.Sizeof(q[0])) / int(unsafe.Sizeof(p[0]))
  pHdr.Cap = qHdr.Cap * int(unsafe.Sizeof(q[0])) / int(unsafe.Sizeof(p[0]))

  fmt.Println(len(p), cap(p), p)
  fmt.Println(len(q), cap(q), q)
}
```

不同切片类型之间转换的思路：先构造一个空的目标切片，然后用原有的切片底层数据填充目标切片。如果 `X` 和 `Y` 类型的大小不同，需要重新设置 `Len` 和 `Cap` 属性。需要注意的是，如果 `X` 或 `Y` 是空类型，上述代码中可能导致除 0 错误，实际代码需要根据情况酌情处理。
