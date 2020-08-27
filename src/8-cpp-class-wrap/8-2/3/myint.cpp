#include <stdio.h>

class MyInt {
  int v_;

public:
  MyInt(int v){ this->v_ = v; }
  int Twice() const{ return this->v_<<1; }
};

struct MyInt2 {
  int Twice() {
    const int *p = (int*)this;
    return (*p) << 1;
  }
};

int main() {
  MyInt x(3);
  // warning: format ‘%d’ expects argument of type ‘int’, but argument 2 has type ‘MyInt’ [-Wformat=]
  printf("%d %d\n", x, x.Twice());// 3 6

  int y = 4;
  printf("%d %d\n", y, ((MyInt*)(&y))->Twice());// 4 8

  int y2 = 5;
  printf("%d %d\n", y2, ((MyInt2*)(&y2))->Twice());// 5 10

  return 0;
}