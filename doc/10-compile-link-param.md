# 编译和链接参数

- [编译和链接参数](#编译和链接参数)
  - [编译参数 CFLAGS/CPPFLAGS/CXXFLAGS](#编译参数-cflagscppflagscxxflags)
  - [链接参数 LDFALGS](#链接参数-ldfalgs)
  - [pkg-config](#pkg-config)
  - [go get 链](#go-get-链)
  - [多个非 main 包导出 C 函数](#多个非-main-包导出-c-函数)

## 编译参数 CFLAGS/CPPFLAGS/CXXFLAGS

编译参数主要是头文件的检索路径、预定义的宏等参数。因为 C++ 对 C 语言做了很多兼容，因此 C 和 C++ 有很多共享编译参数。

- CFLAGS 对应 C 语言(.c 后缀)编译参数
- CPPFLAGS 对应 C/C++ 语言(.c/cc/cpp/cxx 后缀)编译参数
- CXXFLAGS 对应纯 C++ 语言(.cc/cpp/cxx 后缀)编译参数

## 链接参数 LDFALGS

链接参数主要包括链接库的检索目录、链接库的名字。由于历史遗留问题，链接库不支持相对路径，所以必须为链接库指定绝对路径。cgo 中的 `${SRCDIR}` 表示当前目录的绝对路径。

经过编译后的 C 和 C++ 目标文件格式相同，因此 LDFLAGS 对应 C/C++ 共同的链接参数。

## pkg-config

可以通过 `#cgo pkg-config xxx` 目录来生成 xxx 库需要的编译和链接参数，其底层通过调用 `pkg-config xxx --cflags` 生成编译参数，通过 `pkg-config xxx --libs` 命令生成链接参数。**注意**：`pkg-config` 工具生成的编译和链接参数是公用的。

很多非标准的 C/C++ 库并没有实现对 `pkg-config` 的支持。可以手动为 `pkg-config` 工具创建对应库的编译和链接参数实现支持。

先看 ffmpeg 库编译生成的文件 `libavutil.pc`，位于 ffmpeg 库安装目录下的 `lib/pkgconfig` 目录。

```sh
kiki@ubuntu:~/Documents/ffmpeg/ffmpeg-4.1/lib/pkgconfig$ cat libavutil.pc
prefix=/home/kiki/Documents/ffmpeg/ffmpeg-4.1/
exec_prefix=${prefix}
libdir=/home/kiki/Documents/ffmpeg/ffmpeg-4.1//lib
includedir=/home/kiki/Documents/ffmpeg/ffmpeg-4.1//include

Name: libavutil
Description: FFmpeg utility library
Version: 56.19.100
Requires:
Requires.private:
Conflicts:
Libs: -L${libdir}  -lavutil
Libs.private: -pthread -lva -lva-drm -lva -lva-x11 -lX11 -lvdpau -lX11 -lm -lva
Cflags: -I${includedir}
```

假如有一个 xxx 的 C/C++ 库位于 `/usr/local/lib`，对应头文件在 `/usr/local/include`，可以手动创建 `/usr/local/lib/pkgconfig/xxx.pc` 文件：

```txt
Name: xxx
Cflags:-I/usr/local/include
Libs:-L/usr/local/lib –lxxx
```

其中，`Name` 是库名，`Cflags` 对应使用库需要的编译参数，`Libs` 对应使用库需要的链接参数。如果 pc 文件位于其他目录，可通过 `PKG_CONFIG_PATH` 环境变量指定 `pkg-config` 工具的检索目录。

对于 cgo，可以通过 `PKG_CONFIG_PATH` 环境变量指定自定义的 `pkg-config` 程序。 如果是自己实现 cgo 专用的 pkg-config 程序，只要处理 `--cflags` 和 `--libs` 两个参数。

## go get 链

在使用 `go get` 获取 go 语言包的同时会获取包依赖的包。比如 A 包依赖 B 包，B 包依赖 C 包，C 包依赖 D 包。`go get` 获取 A 包之后会依次获取 B/C/D 包。如果获取 B 包之后构建失败，将导致链条断裂，从而导致 A 包构建失败。

链条断裂常见的原因：

- 不支持某些系统，编译失败
- 依赖 cgo，用户未安装 gcc
- 依赖 cgo，用户未安装依赖库
- 依赖 pkg-config，windows 未安装
- 依赖 pkg-config，未找到对应的 pc 文件
- 依赖自定义的 pkg-config，需要额外配置
- 依赖 swig，用户未安装 swig，或版本不对

## 多个非 main 包导出 C 函数

非 `main` 包导出的 go 函数也是有效的。导出后的 go 函数可以当做 C 函数使用。但是不同包导出的 go 函数将在同一个全局的名字空间，需要注意函数重名的问题。

如果是从不同的包导出 go 函数到 C 语言空间，那么 cgo 自动生成的头文件将无法包含所有的函数声明。必须通过手动写头文件的方式导出全部函数。
