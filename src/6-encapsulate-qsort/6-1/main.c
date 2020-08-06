#include <stdio.h>
#include <stdlib.h>

#define DIM(x) (sizeof(x)/sizeof((x)[0]))

static int cmp(const void *a, const void *b) {
  const int *pa = (int*)a;
  const int *pb = (int*)b;
  return *pa - *pb;
}

int main() {
  int arr[] = { 42, 8, 109, 97, 23, 25 };
  int i;

  qsort(arr, DIM(arr), sizeof(arr[0]), cmp);

  for(i=0; i<DIM(arr); i++) {
    printf("%d ", arr[i]);
  }
  printf("\n");

  return 0;
}