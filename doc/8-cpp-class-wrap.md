# C++ 类包装

- [C++ 类包装](#c-类包装)
  - [cgo 不支持 c++ 语法](#cgo-不支持-c-语法)
  - [C++ 类到 go 语言对象](#c-类到-go-语言对象)
    - [构造 一个 C++ 类](#构造-一个-c-类)
    - [用纯 C 函数接口封装 C++ 类](#用纯-c-函数接口封装-c-类)
    - [将纯 C 接口函数转为 go 函数](#将纯-c-接口函数转为-go-函数)
    - [包装为 go 对象](#包装为-go-对象)
  - [go 语言对象到 C++ 类](#go-语言对象到-c-类)
    - [构造一个 go 对象](#构造一个-go-对象)
    - [导出 C 接口](#导出-c-接口)
    - [封装 C++ 类](#封装-c-类)
    - [改进 C++ 类的封装](#改进-c-类的封装)
    - [彻底移除 C++ 的 this 指针](#彻底移除-c-的-this-指针)

## cgo 不支持 c++ 语法

cgo 是 C 语言和 go 语言直接的桥梁，原则上无法直接支持 c++ 类。cgo 不支持 C++ 语法的根本原因是 C++ 没有一个二进制接口规范(ABI)。一个 C++ 类的构造函数在编译生成目标文件时生成链接符号名称、方法在不同平台甚至是 C++ 的不同版本之间都不同。

但是 C++ 兼容 C 语言，所以可以通过增加一组 C 语言函数接口作为 C++ 类和 cgo 之间的桥梁，就可以间接地实现 C++ 和 go 之间的关联。

cgo 只支持 C 语言中值类型的数据类型，所以无法直接使用 C++ 的引用参数等特性。

## C++ 类到 go 语言对象

实现 C++ 类到 go 语言对象的包装步骤：

- 用 纯 C 函数接口包装 C++ 类
- 通过 cgo 将纯 C 函数接口映射到 go 函数
- 实现 go 包装对象，将 C++ 类的方法用 go 对象的方法实现

### 构造 一个 C++ 类

```cpp
// mybuffer.h
#include <string>

struct MyBuffer {
  std::string* s_;

  MyBuffer(int size);
  ~MyBuffer();
  int Size() const;
  char* Data();
};

// mybuffer.cpp
#include "mybuffer.h"

MyBuffer::MyBuffer(int size) {
  this->s_ = new std::string(size, char('\0'));
}

MyBuffer::~MyBuffer() {
  delete this->s_;
}

int MyBuffer::Size() const{
  return this->s_->size();
}

char* MyBuffer::Data() {
  return (char*)this->s_->c_str();
}

// main.cpp
#include <stdio.h>

#include "mybuffer.h"

int main() {
  // g++ -std=c++11 main.cpp mybuffer.cpp
  auto pBuf = new MyBuffer(1024);

  auto data = pBuf->Data();
  auto size = pBuf->Size();
  printf("%d %s\n", size, data);

  delete pBuf;

  return 0;
}
```

### 用纯 C 函数接口封装 C++ 类

在 C 语言中期望使用 C++ 类的方式：

```cpp
// main.cpp
#include <stdio.h>

extern "C" {
  #include "./mybuffer_c.h"
}

int main() {
  // g++ -std=c++11 main.cpp mybuffer_c.cpp mybuffer.cpp
  MyBuffer_T* pBuf = NewMyBuffer(1024);

  char* data = MyBuffer_Data(pBuf);
  int size = MyBuffer_Size(pBuf);
  printf("%d %s\n", size, data);

  DeleteMyBuffer(pBuf);

  return 0;
}
```

定义 C 语言接口规范。这个 C 语言的头文件是 cgo 使用，必须采用 C 语言规范的名字修饰规则：

```c
// mybuffer_c.h
typedef struct MyBuffer_T MyBuffer_T;

MyBuffer_T* NewMyBuffer(int size);
void DeleteMyBuffer(MyBuffer_T *p);

int MyBuffer_Size(MyBuffer_T *p);
char* MyBuffer_Data(MyBuffer_T *p);
```

使用 C 语言封装实现 C++ 类的成员函数。`extern "C"` 指示编译器按 C 语言而不是 C++ 进行编译：

```cpp
// mybuffer_c.cpp
#include "./mybuffer.h"

extern "C" {
  #include "./mybuffer_c.h"
}

struct MyBuffer_T : MyBuffer {
  MyBuffer_T(int size) : MyBuffer(size) {}
  ~MyBuffer_T() {}
};

MyBuffer_T* NewMyBuffer(int size) {
  auto p = new MyBuffer_T(size);
  return p;
}

void DeleteMyBuffer(MyBuffer_T* p) {
  delete p;
}

int MyBuffer_Size(MyBuffer_T* p){
  return p->Size();
}

char* MyBuffer_Data(MyBuffer_T* p) {
  return p->Data();
}
```

### 将纯 C 接口函数转为 go 函数

因为包中包含 C++11 的语法，因此需要通过 `#cgo CXXFLAGS: -std=c++11` 打开 C++11 的选项。

```go
// mybuffer_c.go
package mybuffer

/*
#cgo CXXFLAGS: -std=c++11

#include "mybuffer_c.h"
*/
import "C"

type cgoMyBufferT C.MyBuffer_T

func cgoNewMyBuffer(size int) *cgoMyBufferT {
  p := C.NewMyBuffer(C.int(size))
  return (*cgoMyBufferT)(p)
}

func cgoDeleteMyBuffer(p *cgoMyBufferT) {
  C.DeleteMyBuffer((*C.MyBuffer_T)(p))
}

func cgoMyBufferSize(p *cgoMyBufferT) C.int {
  return C.MyBuffer_Size((*C.MyBuffer_T)(p))
}

func cgoMyBufferData(p *cgoMyBufferT) *C.char {
  return C.MyBuffer_Data((*C.MyBuffer_T)(p))
}
```

### 包装为 go 对象

因为 go 语言的切片本身含有长度信息，可以将 `Data` 和 `Size` 函数合并，返回一个对应底层 C 语言缓存空间的切片。

```go
// mybuffer.go
package mybuffer

import "unsafe"

// MyBuffer ...
type MyBuffer struct {
  cptr *cgoMyBufferT
}

// NewMyBuffer ...
func NewMyBuffer(size int) *MyBuffer {
  return &MyBuffer{
    cptr: cgoNewMyBuffer(size),
  }
}

// Delete ...
func (p *MyBuffer) Delete() {
  cgoDeleteMyBuffer(p.cptr)
}

// Data ...
func (p *MyBuffer) Data() []byte {
  data := cgoMyBufferData(p.cptr)
  size := cgoMyBufferSize(p.cptr)
  return ((*[1 << 31]byte)(unsafe.Pointer(data)))[0:int(size):int(size)]
}
```

之后就可以在 go 语言使用包装后的缓存对象了(底层是基于 C++ 的 `std::string` 实现)。

```go
// main.go
package main

//#include <stdio.h>
import "C"
import (
  "unsafe"

  "github.com/xueqing/cgo-study/src/8-cpp-class-wrap/8-1/3/mybuffer"
)

func main() {
  buf := mybuffer.NewMyBuffer(1024)
  defer buf.Delete()

  copy(buf.Data(), []byte("hello\x00"))
  C.puts((*C.char)(unsafe.Pointer(&(buf.Data()[0]))))
}
```

## go 语言对象到 C++ 类

实现 go 语言对象到 C++ 类的包装步骤：

- 将 go 对象映射到一个 id
- 基于 id 导出对应的 C 接口函数
- 基于 C 接口函数包装为 C++ 类

### 构造一个 go 对象

```go
// person.go
package person

// Person ...
type Person struct {
  name string
  age  int
}

// NewPerson ...
func NewPerson(name string, age int) *Person {
  return &Person{
    name: name,
    age:  age,
  }
}

// Set ...
func (p *Person) Set(name string, age int) {
  p.name = name
  p.age = age
}

// Get ...
func (p *Person) Get() (name string, age int) {
  return p.name, p.age
}

```

要在 C/C++ 中访问 go 对象，需要通过 cgo 导出 C 接口。

### 导出 C 接口

类比签名 C++ 对象到 C 接口的过程，抽象一组 C 接口描述 go 对象。

注意：cgo 导出 C 函数时，输入参数和返回值类型都不支持 const 修饰，也不支持可变参数的函数类型。

```h
// person_c.h
#include <stdint.h>

typedef uintptr_t person_handle_t;

person_handle_t new_person(char *name, int age);
void delete_person(person_handle_t p);

void person_set(person_handle_t p, char *name, int age);
char *person_get_name(person_handle_t p, char *buf, int size);
int person_get_age(person_handle_t p);
```

注意：因为无法在 C/C++ 中直接长期访问 go 内存对象，所以需要使用前面的技术将 go 对象映射为一个整数 id。

```go
// object.go
package object

import "sync"

// ID ...
type ID int32

var refs struct {
  sync.Mutex
  objs map[ID]interface{}
  next ID
}

func init() {
  refs.Lock()
  defer refs.Unlock()

  refs.objs = make(map[ID]interface{})
  refs.next = 1
}

// NewID ...
func NewID(obj interface{}) ID {
  refs.Lock()
  defer refs.Unlock()

  id := refs.next
  refs.next++

  refs.objs[id] = obj
  return id
}

// IsNil ...
func (id ID) IsNil() bool {
  return id == 0
}

// Get ...
func (id ID) Get() interface{} {
  refs.Lock()
  defer refs.Unlock()

  return refs.objs[id]
}

// Free ...
func (id *ID) Free() interface{} {
  refs.Lock()
  defer refs.Unlock()

  obj := refs.objs[*id]
  delete(refs.objs, *id)
  *id = 0

  return obj
}
```

在创建 go 对象之后，`NewObjectId` 将其映射为 id，然后将 id 强制转换成 `person_handle_t` 类型返回。其他接口是根据 `person_handle_t` 所表示的 id，根据 id 解析出对应的 go 对象。

```go
// person_c.go
package person

/*
#cgo CXXFLAGS: -std=c++11

#include "./person_c.h"
*/
import "C"
import (
  "unsafe"

  "github.com/xueqing/cgo-study/src/8-cpp-class-wrap/8-2/1/object"
)

//export new_person
func new_person(name *C.char, age C.int) C.person_handle_t {
  id := object.NewID(NewPerson(C.GoString(name), int(age)))
  return C.person_handle_t(id)
}

//export delete_person
func delete_person(h C.person_handle_t) {
  id := object.ID(h)
  id.Free()
}

//export person_set
func person_set(h C.person_handle_t, name *C.char, age C.int) {
  p := object.ID(h).Get().(*Person)
  p.Set(C.GoString(name), int(age))
}

//export person_get_name
func person_get_name(h C.person_handle_t, buf *C.char, size C.int) *C.char {
  p := object.ID(h).Get().(*Person)
  name, _ := p.Get()

  n := int(size) - 1
  bufSlice := ((*[1 << 31]byte)(unsafe.Pointer(buf)))[0:n:n]
  copy(bufSlice, []byte(name))
  bufSlice[n] = 0

  return buf
}

//export person_get_age
func person_get_age(h C.person_handle_t) C.int {
  p := object.ID(h).Get().(*Person)
  _, age := p.Get()

  return C.int(age)
}
```

### 封装 C++ 类

根据封装好的 C 接口封装C++ 类。

```h
// person.h
extern "C" {
  #include "./person_c.h"
}

struct person {
  person_handle_t goobj_;

  person(const char *name, int age);
  ~person();

  void set(char *name, int age);
  char* get_name(char *buf, int size);
  int get_age();
};
```

```cpp
// person.cpp
#include "person.h"

person::person(const char *name, int age) {
  this->goobj_ = new_person((char*)name, age);
}

person::~person() {
  delete_person(this->goobj_);
}

void person::set(char *name, int age) {
  person_set(this->goobj_, name, age);
}

char* person::get_name(char *buf, int size) {
  return person_get_name(this->goobj_, buf, size);
}

int person::get_age() {
  return person_get_age(this->goobj_);
}
```

### 改进 C++ 类的封装

上述封装实现中，每次通过 `new` 创建一个 `Person` 类实例需要两次内存分配：一次是针对 C++ 版本的，一次是 go 语言版本。其中，C++ 类内部职员一个 `person_handle_t` 类型的成员，用于映射 go 对象。可以优化这个封装，即将 `person_handle_t` 直接当做 C++ 类对象使用。

```h
// person.h
extern "C" {
  #include "./person_c.h"
}

struct person {
  static person* New(const char *name, int age);
  void Delete();

  void set(char *name, int age);
  char* get_name(char *buf, int size);
  int get_age();
};
```

```cpp
// person.cpp
#include "person.h"

person* person::New(const char *name, int age) {
  return (person*)new_person((char*)name, age);
}

void person::Delete() {
  delete_person(person_handle_t(this));
}

void person::set(char *name, int age) {
  person_set(person_handle_t(this), name, age);
}

char* person::get_name(char *buf, int size) {
  return person_get_name(person_handle_t(this), buf, size);
}

int person::get_age() {
  return person_get_age(person_handle_t(this));
}
```

### 彻底移除 C++ 的 this 指针

go 语言的方法绑定到类型。可以基于基本类型 `int` 定义一个新类型，并在不改变原有数据底层内存结构时，自由切换基本类型和重命名类型。

```go
package main

import "fmt"

// MyInt ...
type MyInt int

// Twice ...
func (i MyInt) Twice() int {
  return (int)(i) << 1
}

func main() {
  var x = MyInt(3)
  fmt.Println(int(x))
  fmt.Println(x.Twice())
}
```

C++ 一般会通过 class 包含基本类型作为成员变量，并绑定相关方法。但是类中的 `this` 就是 class 的指针类型，导致不能和基本类型互相转化。

可以从 C 语言角度考虑，即把 `this` 当做一个普通的 `void*` 指针，随意转换为其他类型。因此，C++ 的方法也可以用于普通的非 class 类型，即可以将普通成员函数绑定到基本类型。但是，纯虚方法是绑定到对象的，即接口。

```cpp
#include <stdio.h>

class MyInt {
  int v_;

public:
  MyInt(int v){ this->v_ = v; }
  int Twice() const{ return this->v_<<1; }
};

struct MyInt2 {
  int Twice() {
    const int *p = (int*)this;
    return (*p) << 1;
  }
};

int main() {
  MyInt x(3);
  // warning: format ‘%d’ expects argument of type ‘int’, but argument 2 has type ‘MyInt’ [-Wformat=]
  printf("%d %d\n", x, x.Twice());// 3 6

  int y = 4;
  printf("%d %d\n", y, ((MyInt*)(&y))->Twice());// 4 8

  int y2 = 5;
  printf("%d %d\n", y2, ((MyInt2*)(&y2))->Twice());// 5 10

  return 0;
}
```
