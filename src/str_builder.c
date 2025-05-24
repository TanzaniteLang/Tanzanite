#include <str.h>
#include <str_builder.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdarg.h>

const int INIT_SIZE = 32;

static void expand_buffer(struct str_builder *b);

void str_builder_deinit(struct str_builder *b)
{
    str_free(&b->buffer);
    b->allocated = 0;
}

void str_builder_shrink(struct str_builder *b)
{
    if (b->buffer.str != NULL) {
        b->allocated = b->buffer.size;
        b->buffer.str = realloc(b->buffer.str, b->allocated);
    }
}

struct str str_builder_str(struct str_builder *b)
{
    struct str s = { 0 };

    if (b->buffer.str != NULL) {
        s = b->buffer;
        b->buffer.str = NULL;
        b->buffer.size = 0;
        b->allocated = 0;
    }

    return s;
}

void str_builder_append_char(struct str_builder *b, char c)
{
    while (b->buffer.size + 1 >= b->allocated)
        expand_buffer(b);

    b->buffer.str[b->buffer.size] = c; 
    b->buffer.size++;
    b->buffer.str[b->buffer.size] = 0;
}

void str_builder_append_cstr(struct str_builder *b, const char *str)
{
    int len = strlen(str);
    while (b->buffer.size + len + 1 >= b->allocated)
        expand_buffer(b);

    memcpy(b->buffer.str + b->buffer.size, str, len);
    b->buffer.size += len;
    b->buffer.str[b->buffer.size] = 0;
}

void str_builder_append_str(struct str_builder *b, struct str s)
{
    while (b->buffer.size + s.size + 1 >= b->allocated)
        expand_buffer(b);

    memcpy(b->buffer.str + b->buffer.size, s.str, s.size);
    b->buffer.size += s.size;
    b->buffer.str[b->buffer.size] = 0;
}

void str_builder_printf(struct str_builder *b, const char *fmt, ...)
{
    va_list args;
    va_start(args, fmt);

    size_t length = vsnprintf(NULL, 0, fmt, args);
    va_end(args);
    va_start(args, fmt);
    char buf[length + 1];
    memset(buf, 0, length + 1);
    vsnprintf(buf, length + 1, fmt, args);
    str_builder_append_cstr(b, buf);
    va_end(args);
}




static void expand_buffer(struct str_builder *b)
{
    if (b->buffer.str != NULL) {
        b->allocated *= 2;
        b->buffer.str = realloc(b->buffer.str, b->allocated * sizeof(char));
    } else {
        b->allocated = INIT_SIZE;
        b->buffer.str = calloc(b->allocated, sizeof(char));
    }
}
