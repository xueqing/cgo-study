#include <stdio.h>

#include "number.h"

int main() {
  // go build -buildmode=c-shared -o number.so
  // gcc test_number.c number.so
  int a=10, b=4, mod=5;

  printf("(%d+%d)%%%d = %d\n", a, b, mod, number_add_mod(a,b,mod));

  return 0;
}