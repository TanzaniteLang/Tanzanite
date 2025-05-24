#ifndef __STR_BUILDER_H__
#define __STR_BUILDER_H__

#include <str.h>
#include <stddef.h>

struct str_builder {
    struct str buffer;
    size_t allocated;
};

void str_builder_deinit(struct str_builder *b);
void str_builder_shrink(struct str_builder *b);
struct str str_builder_str(struct str_builder *b);
void str_builder_append_char(struct str_builder *b, char c);
void str_builder_append_cstr(struct str_builder *b, char *str);
void str_builder_append_str(struct str_builder *b, struct str s);
void str_builder_printf(struct str_builder *b, const char *fmt, ...);

#endif
