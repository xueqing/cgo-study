#include <stdint.h>

typedef uintptr_t person_handle_t;

person_handle_t new_person(char *name, int age);
void delete_person(person_handle_t p);

void person_set(person_handle_t p, char *name, int age);
char *person_get_name(person_handle_t p, char *buf, int size);
int person_get_age(person_handle_t p);