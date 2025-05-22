#ifndef __HASH_H__
#define __HASH_H__

#include <stdint.h>

/* Inspired from https://harry.pm/blog/lets_write_a_hashmap/ */

enum hash_state {
    HASH_EMPTY,
    HASH_VALID,
    HASH_FREE,
};

#define hash_begin(hm) ((size_t)(0))
#define hash_end(hm) (((hm)->cap))
#define hash_states(hm, it) ((hm)->buckets[(it)].state)
#define hash_key(hm, it) ((hm)->buckets[(it)].key)
#define hash_value(hm, it) ((hm)->buckets[(it)].value)
#define hash_exists(hm, it) ((it) < (hm)->cap && hash_states((hm), (it)) == HASH_VALID)

#define HASH_DECL(name, type) \
struct name##_bucket {\
    enum hash_state state;\
    const char *key;\
    type value;\
};\
\
struct name {\
    uint32_t len;\
    uint32_t cap;\
    struct name##_bucket *buckets;\
};\
void name##_free(struct name *hm);\
uint32_t name##_insert(struct name *hm, const char *key);\
void name##_remove(struct name *hm, const char *key);\
uint32_t name##_find(struct name *hm, const char *key);

#endif
