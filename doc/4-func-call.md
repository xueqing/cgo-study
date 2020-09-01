# 函数调用

- [函数调用](#函数调用)
  - [Go 调用 C 函数](#go-调用-c-函数)
  - [C 函数的返回值](#c-函数的返回值)
  - [void 函数的返回值](#void-函数的返回值)
  - [C 调用 Go 导出函数](#c-调用-go-导出函数)

## Go 调用 C 函数

对于一个启用 cgo 特性的程序，cgo 会构造一个虚拟的 C 包。通过这个虚拟的 C 包可以调用 C 语言函数。

## C 函数的返回值

对于有返回值的 C 函数，我们可以正常获取返回值。

C 语言不支持返回多个结果，因此 `<errno.h>` 标准库提供了一个 `errno` 宏用于返回错误状态。可以近似地将 `errno` 看成一个线程安全的全局变量，可以用于记录最近一次错误的状态码。

cgo 针对 `<errno.h>` 标准库的 `errno` 宏做的特殊支持：在 cgo 调用 C 函数时如果有两个返回值，那么第二个返回值将对应 `errno` 错误状态。

cgo 调用 C 函数的第二个返回值是可忽略的 `error` 接口类型，底层对应 `syscall.Errno` 错误类型。

```go
package main

/*
#include <errno.h>

static int div(int a, int b) {
  if (b == 0) {
    errno = EINVAL;
    return 0;
  }
  return a/b;
}
*/
import "C"
import "fmt"

func main() {
  v1, e1 := C.div(4, 2)
  fmt.Println(v1, e1) // 2 <nil>

  v2, e2 := C.div(4, 0)
  fmt.Println(v2, e2) // 0 invalid argument
}
```

## void 函数的返回值

一般情况下，无法获取 `void` 类型函数的返回值，因为没有返回值可以获取。前面的例子中提到，cgo 对 `errno` 做了特殊处理，可以通过第二个返回值来获取 C 语言的错误状态。对于 `void` 类型函数，这个特性依然有效。

```go
package main

// static void noreturn() {}
import "C"
import "fmt"

func main() {
  v, e := C.noreturn()
  fmt.Println(v, e)      // [] <nil>
  fmt.Printf("%#v\n", v) // main._Ctype_void{}
}
```

在 cgo 生成的代码中，`_Ctype_void` 类型对应一个 0 长的数组类型 `[0]byte`，因此 `fmt.Println` 输出的是一个表示空数值的方括弧。

## C 调用 Go 导出函数

cgo 可将 Go 函数导出为 C 语言函数。这样的话，我们可以定义 C 语言接口，然后通过 Go 语言实现。

```go
package main

// extern int add(int a, int b);
import "C"
import "fmt"

//export add
func add(a, b C.int) C.int {
  return a + b
}

func main() {
  fmt.Println(C.add(1, 1))
}
```

`add` 函数名以小写字母开头，对于 Go 语言来说是包内的私有函数。但是从 C 语言角度来看，导出的 `add` 函数是一个可全局访问的 C 语言函数。如果在两个不同的 Go 语言包内，都存在一个同名的要导出为 C 语言函数的 `add` 函数，那么在最终的链接阶段将会出现符号重名的问题。

使用 `//export` (**注意：双斜线和 export 之间没有空格**) 表示需要导出的 Go 函数。

当导出 C 语言接口时，需要保证函数的参数和返回值类型都是 C 语言友好的类型，同时返回值不得直接或间接包含 Go 语言内存空间的指针。
