#include <stdio.h>

#include "number.h"

int main() {
  // go build -buildmode=c-archive -o number.a
  // gcc test_number.c number.a -lpthread
  int a=10, b=4, mod=5;

  printf("(%d+%d)%%%d = %d\n", a, b, mod, number_add_mod(a,b,mod));

  return 0;
}