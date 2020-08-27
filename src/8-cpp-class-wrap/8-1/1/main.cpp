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