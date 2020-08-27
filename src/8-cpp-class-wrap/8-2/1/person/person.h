extern "C" {
  #include "./person_c.h"
}

struct person {
  person_handle_t goobj_;

  person(const char *name, int age);
  ~person();

  void set(char *name, int age);
  char* get_name(char *buf, int size);
  int get_age();
};
