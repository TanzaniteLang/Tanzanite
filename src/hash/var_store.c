#include <stdbool.h>
#include <stdlib.h>
#include <string.h>

#include <hash.h>
#include <stack.h>
#include <stddef.h>

#include <analyzer/variable.h>
#include <hash/var_store.h>

HASH_IMPL(var_store_hash, struct analyzable_variable);
STACK_IMPL(var_store, struct var_store_hash);

void var_store_push_frame(struct var_store *store)
{
    struct var_store_hash empty = {0};
    uint32_t it = var_store_push(store);
    stack_value(store, it) = empty;
}

void var_store_pop_frame(struct var_store *store)
{
    var_store_pop(store);
}

struct var_store_res var_store_find(struct var_store *store, const char *key)
{
    struct var_store_res res = {0};
    res.found = false;

    uint32_t hash_iter = stack_top(store);

    for (;; hash_iter--) {
        struct var_store_hash *hash = &stack_value(store, hash_iter);
        uint32_t it = var_store_hash_find(hash, key);
        if (hash_exists(hash, it)) {
            res.found = true;
            res.payload = hash_value(hash, it);
            break;
        }

        if (hash_iter == stack_bottom(store))
            break;
    }

    return res;
}

struct analyzable_variable *var_store_insert(struct var_store *store, const char *key)
{
    struct var_store_hash *hash = &stack_value(store, stack_top(store));
    uint32_t it = var_store_hash_insert(hash, key);

    return &hash_value(hash, it);
}
