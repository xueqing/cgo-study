extern "C" {
  #include "./person_c.h"
}

struct person {
  static person* New(const char *name, int age);
  void Delete();

  void set(char *name, int age);
  char* get_name(char *buf, int size);
  int get_age();
};
