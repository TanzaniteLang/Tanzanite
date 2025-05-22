#ifndef __HASH_H__
#define __HASH_H__

#include <stdint.h>
#include <djb2.h>

#define HASH_MIN_CAP 16

/* Inspired from https://harry.pm/blog/lets_write_a_hashmap/ */

enum hash_state {
    HASH_EMPTY,
    HASH_VALID,
    HASH_FREE,
};

#define hash_begin(hm) ((uint32_t)(0))
#define hash_end(hm) (((hm)->cap))
#define hash_states(hm, it) ((hm)->buckets[(it)].state)
#define hash_key(hm, it) ((hm)->buckets[(it)].key)
#define hash_value(hm, it) ((hm)->buckets[(it)].value)
#define hash_exists(hm, it) ((it) < (hm)->cap && hash_states((hm), (it)) == HASH_VALID)

#define HASH_DECL(name, type)\
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
void name##_remove(struct name *hm, uint32_t it);\
uint32_t name##_find(struct name *hm, const char *key);\
bool name##_resize(struct name *hm);

#define HASH_IMPL(name, type)\
void name##_free(struct name *hm) {\
    if (hm == NULL)\
        return;\
    if (hm->cap > 0)\
        free(hm->buckets);\
    memset(hm, 0, sizeof(*hm));\
}\
uint32_t name##_insert(struct name *hm, const char *key) {\
    if (!name##_resize(hm))\
        return hm->cap;\
    uint32_t it = djb2(key, strlen(key)) % hm->cap;\
    while (hm->buckets[it].state == HASH_VALID && strcmp(key, hm->buckets[it].key))\
        it = (it + 1) % hm->cap;\
    if (hm->buckets[it].state != HASH_VALID)\
        hm->len++;\
    hm->buckets[it].state = HASH_VALID;\
    hm->buckets[it].key = key;\
    return it;\
}\
void name##_remove(struct name *hm, uint32_t it) {\
    if (hash_exists(hm, it)) {\
        hm->buckets[it].state = HASH_FREE;\
        hm->len--;\
    }\
    name##_resize(hm);\
}\
uint32_t name##_find(struct name *hm, const char *key) {\
    if (hm->cap == 0)\
        return hm->cap;\
    uint32_t it = djb2(key, strlen(key)) % hm->cap;\
    while (hm->buckets[it].state == HASH_FREE || (hm->buckets[it].state == HASH_VALID && strcmp(key, hm->buckets[it].key)))\
        it = (it + 1) % hm->cap;\
    if (hm->buckets[it].state != HASH_VALID)\
        return hm->cap;\
    return it;\
}\
bool name##_resize(struct name *hm) {\
    uint32_t old_cap = hm->cap;\
    uint32_t new_cap;\
    if (!hm->cap || hm->len * 4 > hm->cap * 3)\
        new_cap = old_cap > 0 ? old_cap * 2 : HASH_MIN_CAP;\
    else if (hm->cap > HASH_MIN_CAP && hm->len * 4 < hm->cap)\
        new_cap = old_cap / 2;\
    else\
        return true;\
    struct name##_bucket *new_buckets = calloc(new_cap, sizeof(*hm->buckets));\
    if (new_buckets == NULL)\
        return false;\
    for (uint32_t i = 0; i < old_cap; ++i) {\
        if (hm->buckets[i].state != HASH_VALID)\
            continue;\
        const char *key = hm->buckets[i].key;\
        uint32_t it = djb2(key, strlen(key)) % hm->cap;\
        while (new_buckets[it].state == HASH_VALID)\
            it = (it + 1) % new_cap;\
        new_buckets[it] = hm->buckets[i];\
    }\
    free(hm->buckets);\
    hm->buckets = new_buckets;\
    hm->cap = new_cap;\
    return true;\
}

#endif
