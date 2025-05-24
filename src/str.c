#include <str.h>

#include <stdlib.h>
#include <stdlib.h>
#include <string.h>

static uint64_t getlen(const char *str, uint64_t len);
static char *clonestr(const char *str, uint64_t len);

struct str str_init(const char *string, uint64_t len)
{
    struct str s = {0};
    s.str = clonestr(string, len);
    s.size = getlen(string, len);

    return s;
}

void str_free(struct str *str)
{
    free(str->str);
    str->str = NULL;
    str->size = 0;
}


static uint64_t getlen(const char *str, uint64_t len)
{
    uint64_t l = 0;
    const char *iter = str;

    for (; *iter != '\0' && len > 0; iter++, len--, l++);

    return l;
}

static char *clonestr(const char *str, uint64_t len)
{
    uint64_t l = getlen(str, len);
    char *dup = calloc(l + 1, sizeof(char));

    memcpy(dup, str, l);
    return dup;
}
