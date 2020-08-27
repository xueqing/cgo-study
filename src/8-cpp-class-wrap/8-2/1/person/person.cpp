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