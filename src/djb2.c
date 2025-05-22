#include <djb2.h>

unsigned int djb2(const char *bytes, size_t len)
{
    unsigned int hash = 5381;
    for (size_t i = 0; i < len; ++i)
        hash = hash * 33 + bytes[i];
    return hash;
}
