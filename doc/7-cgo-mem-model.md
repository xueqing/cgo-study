# cgo 内存模型

- [cgo 内存模型](#cgo-内存模型)
  - [差异](#差异)
  - [go 访问 C 内存](#go-访问-c-内存)
  - [C 临时访问传入的 go 内存](#c-临时访问传入的-go-内存)
  - [C 长期持有 go 指针对象](#c-长期持有-go-指针对象)
  - [导出 C 函数不能返回 go 内存](#导出-c-函数不能返回-go-内存)

## 差异

C 语言的内存在分配之后是稳定的，但是 go 语言因为函数栈的动态伸缩可能导致栈中内存地址的移动。如果 C 语言持有的是移动之前的 go 指针，那么再次访问 go 对象时会导致程序崩溃。

## go 访问 C 内存

C 语言空间的内存是稳定的，只要不是人为提前释放，go 语言就可以一直使用。

```go
package main

/*
#include <stdlib.h>

void* makeslice(size_t memsize) {
  return malloc(memsize);
}
*/
import "C"
import (
  "fmt"
  "unsafe"
)

func makeByteSlice(n int) []byte {
  p := C.makeslice(C.size_t(n))
  return ((*[1 << 31]byte)(p))[0:n:n] //panic: runtime error: slice bounds out of range [::4294967297] with length 2147483648
}

func freeByteSlice(p []byte) {
  C.free(unsafe.Pointer(&p[0]))
}

func main() {
  s := makeByteSlice(1<<32 + 1)
  s[len(s)-1] = 255
  fmt.Println(s[len(s)-1])
  freeByteSlice(s)
}
```

## C 临时访问传入的 go 内存

方法 1：在 C 语言空间先分配相同大小的内存，然后将 go 的内存填充到 C 的内存空间；返回的内存也是如此。

问题：效率低(要多次分配内存并逐个拷贝元素)，代码繁琐

```go
package main

/*
#include <stdlib.h>
#include <stdio.h>

void printString(const char *s) {
  printf("%s\n", s);
}
*/
import "C"
import "unsafe"

func printString(s string) {
  cs := C.CString(s) // copy go string to C string(memory allocated by C)
  defer C.free(unsafe.Pointer(cs))

  C.printString(cs)
}

func main() {
  printString("hello")
}
```

为了简化并高效处理这个问题，cgo 针对此场景定义了专门的规则：在 cgo 调用的 C 语言函数返回前，cgo 保证传入的 go 语言内存在此期间不会发生移动。

```go
package main

/*
#include <stdio.h>

void printString(const char *s, int n) {
  int i;
  for(i=0; i<n; i++) {
    putchar(s[i]);
  }
  putchar('\n');
}
*/
import "C"
import (
  "reflect"
  "unsafe"
)

func printString(s string) {
  p := (*reflect.StringHeader)(unsafe.Pointer(&s))
  C.printString((*C.char)(unsafe.Pointer(p.Data)), C.int(len(s)))
}

func main() {
  printString("hello")
}
```

存在隐患：如果调用的 C 语言函数运行时间较长，导致引用的 go 语言内存不能移动，间接导致这个 go 内存堆栈对应的 goroutine 不能动态伸缩栈内存，可能导致这个 goroutine 被阻塞。因此，在需要长时间运行的 C 语言函数(尤其是需要等待 CPU 运算资源之外的资源而不确定时间时)，需要谨慎使用 go 语言内存。

注意：在取得 go 内存后要马上传入 C 语言函数，不能保存到临时变量再间接传给 C 语言函数。因为 cgo 只能保证在 C 函数调用之后被传入的 go 内存不会移动，不能保证在传入 C 函数之前内存不发生变化。

```go
// wrong: tmp is not a pointer, cannnot update when moving underlying memory
tmp := uintptr(unsafe.Pointer(&x))
pb := (*int16)(unsafe.Pointer(tmp))
*pb = 42
```

## C 长期持有 go 指针对象

C 调用 go 函数时，C 函数是调用方，go 语言函数的 go 对象的内存生命周期超出了 go 语言的运行时管理。因此，不能再 C 函数直接使用 go 语言对象的内存。

如果需要在 C 语言中访问 go 语言内存对象，可以将 go 语言内存对象在 go 语言空间映射为一个 int 类型的 id，然后通过此 id 来间接访问和控制 go 语言对象。

```go
// objectid.go
package main

import "sync"

type ObjectId int32

var refs struct {
  sync.Mutex
  objs map[ObjectId]interface{}
  next ObjectId
}

func init() {
  refs.Lock()
  defer refs.Unlock()

  refs.objs = make(map[ObjectId]interface{})
  refs.next = 1000
}

func NewObjectId(obj interface{}) ObjectId {
  refs.Lock()
  defer refs.Unlock()

  id := refs.next
  refs.next++

  refs.objs[id] = obj
  return id
}

func (id ObjectId) IsNil() bool {
  return id == 0
}

func (id ObjectId) Get() interface{} {
  refs.Lock()
  defer refs.Unlock()

  return refs.objs[id]
}

func (id *ObjectId) Free() interface{} {
  refs.Lock()
  defer refs.Unlock()

  obj := refs.objs[*id]
  delete(refs.objs, *id)
  *id = 0

  return obj
}
```

```go
// main.go
package main

/*
#include <stdlib.h>

extern char* NewGoString(char* );
extern void FreeGoString(char* );
extern void PrintGoString(char* );

static void printString(const char* s) {
    char* gs = NewGoString(s);
    PrintGoString(gs);
    FreeGoString(gs);
}
*/
import "C"
import (
  "unsafe"
)

//export NewGoString
func NewGoString(s *C.char) *C.char {
  gs := C.GoString(s)
  id := NewObjectId(gs)
  return (*C.char)(unsafe.Pointer(uintptr(id)))
}

//export FreeGoString
func FreeGoString(p *C.char) {
  id := ObjectId(uintptr(unsafe.Pointer(p)))
  id.Free()
}

//export PrintGoString
func PrintGoString(s *C.char) {
  id := ObjectId(uintptr(unsafe.Pointer(s)))
  gs := id.Get().(string)
  print(gs)
}

func printString(s string) {
  cs := C.CString(s) // copy go string to C string(memory allocated by C)
  defer C.free(unsafe.Pointer(cs))

  C.printString(cs)
}

func main() {
  printString("hello")
}
```

## 导出 C 函数不能返回 go 内存

go 语言是从一个固定的虚拟地址空间分配内存。而 C 语言分配的内存则不能使用 go 语言保留的虚拟内存空间。在 cgo 环境，go 语言运行时默认会检查导出函数返回的内存是否是 go 分配的，如果是则抛出异常。

```go
package main

/*
extern int* getGoPtr();

static void Main() {
  int* p = getGoPtr();
  *p = 42;
}
*/
import "C"

func main() {
  C.Main() //panic: runtime error: cgo result has Go pointer
  // use `GODEBUG=cgocheck=0 go run .` will not throw exception
}

//export getGoPtr
func getGoPtr() *C.int {
  return new(C.int) //memory is allocated by go
}
```
