#ifndef __STR_H__
#define __STR_H__

#include <stdint.h>

struct str {
    char *str;
    uint64_t size;
};

struct str str_init(const char *string, uint64_t len);
void str_free(struct str *str);

#endif
