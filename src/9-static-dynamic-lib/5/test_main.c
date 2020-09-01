#include <stdio.h>

void goPrintln(char*);
int number_add_mod(int a, int b, int mode);

int main() {
  // go build -buildmode=c-shared -o number.so
  // gcc test_number.c number.so
  int a=10, b=4, mod=5;

  printf("(%d+%d)%%%d = %d\n", a, b, mod, number_add_mod(a,b,mod));

  goPrintln("hi, test main");

  return 0;
}