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