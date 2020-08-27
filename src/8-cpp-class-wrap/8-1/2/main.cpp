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