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