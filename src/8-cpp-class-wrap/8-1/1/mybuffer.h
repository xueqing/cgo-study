#include <string>

struct MyBuffer {
  std::string* s_;

  MyBuffer(int size);
  ~MyBuffer();
  int Size() const;
  char* Data();
};