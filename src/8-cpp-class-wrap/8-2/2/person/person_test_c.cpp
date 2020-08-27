#include <stdio.h>

#include "./person.h"

extern "C" {
  #include "./person_test_c.h"
}

void test_new_person() {
  auto p = person::New("kiki", 28);
  
  char buf[64];
  char* name = p->get_name(buf, sizeof(buf)-1);
  int age = p->get_age();
  printf("%s %d\n", name, age);

  p->set((char*)"jimmy", 26);
  name = p->get_name(buf, sizeof(buf)-1);
  age = p->get_age();
  printf("%s %d\n", name, age);

  p->Delete();
}