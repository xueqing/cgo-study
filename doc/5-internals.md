# 内部机制

- [内部机制](#内部机制)
  - [cgo 生成的中间文件](#cgo-生成的中间文件)
  - [Go 调用 C 函数](#go-调用-c-函数)
  - [C 调用 Go 函数](#c-调用-go-函数)

cgo 特性主要是通过一个叫 cgo 的命令行工具来辅助输出 Go 和 C 之间的桥接代码。

## cgo 生成的中间文件

在一个包有 4 个文件

```h
//add.h
int addByGo(int a, int b);
```

```go
//add.go
package main

import "C"

//export addByGo
func addByGo(a, b C.int) C.int {
  return a + b
}
```

```go
//nocgo.go
package main

import "fmt"

func noCgo() {
  fmt.Println("no cgo here")
}

```

```go
//main.go
package main

/*
#include "add.h"

static int addByC(int a, int b) {
  return a+b;
}
*/
import "C"

import "fmt"

func main() {
  // run `go tool cgo main.go add.go` to generate intermediate files
  fmt.Println(C.addByC(1, 1))
  fmt.Println(C.addByGo(1, 1))
  noCgo()
}
```

运行 `go tool cgo main.go add.go` 在当前目录的 `_obj` 文件夹下保存生成的中间文件:

- cgo 命令会为每个包含了 cgo 代码的 Go 文件创建 2 个中间文件，比如 `main.go` 会分别创建 `main.cgo1.go` 和 `main.cgo2.c` 两个中间文件。
- 然后会为整个包创建一个 `_cgo_gotypes.go` Go 文件，其中包含 Go 语言部分辅助代码。
- 还会创建一个 `_cgo_export.h` 和 `_cgo_export.c` 文件，对应 Go 语言导出到 C 语言的类型和函数。

## Go 调用 C 函数

```go
package main

// int sum(int a, int b) { return a+b; }
import "C"

func main() {
  println(C.sum(1, 1))
}
```

运行 `go tool cgo main.go` 在当前目录的 `_obj` 文件夹下保存生成的中间文件:

```sh
$ ls _obj | awk '{print $NF}'
_cgo_export.c
_cgo_export.h
_cgo_flags
_cgo_gotypes.go
_cgo_main.c
_cgo_.o
main.cgo1.go
main.cgo2.c
```

- 其中 `_cgo_.o`、`_cgo_flags` 和 `_cgo_main.c` 文件和代码没有直接的逻辑关联，可以暂时忽略。
- `main.cgo1.go` 是 `main.go` 文件展开虚拟 C 包相关函数和变量后的 Go 代码：
  - `C.xxx` 形式的函数会被替换为 `_Cfunc_xxx` 格式的纯 Go 函数，其中前缀 `_Cfunc_` 表示这是一个 C 函数，对应一个私有的 Go 桥接函数。
  - `_Cfunc_sum` 函数在 `_cgo_gotypes.go` 中定义。`_Ctype_int` 类型对应 `C.int` 类型，命名规则和 `_Cfunc_xxx` 类似，不同的前缀用于区分函数和类型。
- 被传入 C 语言函数 `_cgo_d473794c4ea1_Cfunc_sum` 也是 cgo 生成的中间函数。函数在 `main.cgo2.c` 定义
  - 函数参数只有一个 `void` 范型的指针，函数没有返回值。真实的 `sum 函数的函数参数和返回值均通过唯一的参数指针类实现。
  - `_cgo_topofstack` 函数相关的代码用于C函数调用后恢复调用栈。
  - `_cgo_tsan_acquire` 和 `_cgo_tsan_release` 则是用于扫描 cgo 相关的函数，是对 cgo 相关函数的指针做相关检查。
- 函数调用顺序：
  - `main.go`: `C.sum`
  - `main.cgo1/go`: `_Cfunc_sum`
  - `_cgo_gotypes.go`: `_cgo_runtime_cgocall`
  - `runtime.cgocall`: `_cgo_d473794c4ea1_Cfunc_sum`
  - `main.cgo2.c`: `sum`
  - `_cgo_export.c`:

## C 调用 Go 函数

```go
//main.go
package main

//int sum(int a, int b);
import "C"
import "fmt"

//export sum
func sum(a, b C.int) C.int {
  return a + b
}

func main() {
  fmt.Println(C.sum(1, 1))
}
```

运行 `go tool cgo main.go` 在当前目录的 `_obj` 文件夹下保存生成的中间文件:

```sh
$ ls _obj | awk '{print $NF}'
_cgo_export.c
_cgo_export.h
_cgo_flags
_cgo_gotypes.go
_cgo_main.c
_cgo_.o
main.cgo1.go
main.cgo2.c
```

```c
//_testmain.c
#include <stdio.h>
#include "sum.h"

int main() {
  // run `gcc _testmain.c -o _testmain ./sum.a -lpthread`
  extern int sum(int a, int b);
  printf("%d\n", sum(1, 2));
  return 0;
}
```

函数调用顺序：

- `_testmain.c`: `sum`
- `_cgo_export.c`: `_cgoexp_4930e30f86ac_sum`
- `_cgo_gotypes.go`:  `_cgoexpwrap_4930e30f86ac_sum` ,`_cgo_runtime_cgocallback`
- `runtime.cgocallback`:
- `main.go`: `sum`
