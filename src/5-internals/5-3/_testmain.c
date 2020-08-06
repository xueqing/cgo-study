#include <stdio.h>
#include "sum.h"

int main() {
  // run `gcc _testmain.c -o _testmain ./sum.a -lpthread`
  extern int sum(int a, int b);
  printf("%d\n", sum(1, 2));
  return 0;
}