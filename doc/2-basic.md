# 基础

- [基础](#基础)
  - [import "C"](#import-c)
  - [cgo 语句](#cgo-语句)
  - [build tag 条件编译](#build-tag-条件编译)

## import "C"

`import "C"` 表示使用 cgo 特性，紧跟在这行语句签名的注释包含的是 C 语言代码。启用 cgo 时，还可在当前目录包含 C/C++ 源文件。

注意：

- 通过虚拟 C 包导入的 C 语言符号无需是大写字母开头，不受 Go 语言的导出规则约束。
- 不同的 Go 包通过 cgo 导入到虚拟 C 包的类型不通用，是包私有的。

## cgo 语句

在 `import "C"` 前的注释中科通过 `#cgo` 语句设置编译和链接阶段的相关参数：

- 编译阶段参数主要用于定义相关宏和指定头文件检索路径
- 链接阶段参数主要用于指定库文件检索路径和要链接的库文件
  - 可使用变量 `${SRCDIR}` 表示当前包目录的绝对路径

```go
// #cgo CFLAGS: -DPNG_DEBUG=1 -I./include
// #cgo LDFLAGS: -L/usr/local/lib -lpng
// #include <png.h>
import "C"
```

`#cgo` 支持条件选择，当满足某个操作系统或某个 CPU 架构时后面的编译或链接选项生效：

```go
// #cgo windows CFLAGS: -DX86=1
// #cgo !windows LDFLAGS: -lm
```

如果在不同系统 cgo 对应不同的 C 代码，可先使用 `#cgo` 指令定义不同的 C 语言宏，然后通过宏区分不同的代码。

```go
/*
#cgo windows CFLAGS: -DCGO_OS_WINDOWS=1
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo linux CFLAGS: -DCGO_OS_LINUX=1

#if defined(CGO_OS_WINDOWS)
    const char* os = "windows";
#elif defined(CGO_OS_DARWIN)
    static const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
    static const char* os = "linux";
#else
#    error(unknown os)
#endif
*/
```

## build tag 条件编译

build tag 是在 Go 或 cgo 环境下的 C/C++ 文件开头的一种特殊注释。下面的源文件只有在设置 `debug` 构建标志时才会被构建：

```go
// +build debug

package main
```

可以用以下命令构建：

```sh
go build -tags="debug"
go build -tags="windows debug"
```

我们可通过 `-tags` 命令行参数同时指定多个 build 标志，它们之间用空格分隔。当有多个 build tag 时，我们将多个标志通过逻辑操作的规则来组合使用。比如以下的构建标志表示只有在 “linux/386” 或 “darwin 平台下非 cgo 环境” 才进行构建。

```go
// +build linux,386 darwin,!cgo
```

其中 `linux,386` 中 linux 和 386 用逗号连接表示 AND；而 `linux,386` 和 `darwin,!cgo` 之间通过空白分割来表示 OR。
