# 快速入门

- [快速入门](#快速入门)

通过 `import "C"` 启用 cgo 特性。

```go
// +build go1.10

package main

//void SayHello(_GoString_ s);
import "C"

import "fmt"

func main() {
  C.SayHello("Hello, world\n")
}

//export SayHello
func SayHello(s string) {
  fmt.Print(s)
}
```

思考题:上面 `main` 函数和 `SayHello` 函数是否在同一个 `goroutine` 里执行？

- 参考 [C code and goroutine scheduling](https://stackoverflow.com/questions/28354141/c-code-and-goroutine-scheduling)
- 参考 [go/src/runtime/cgocall.go](https://github.com/golang/go/blob/master/src/runtime/cgocall.go)
