#ifndef __HASH_VAR_STORE_H__
#define __HASH_VAR_STORE_H__

#include <hash.h>
#include <stack.h>
#include <stddef.h>

#include <analyzer/variable.h>

struct var_store_res {
    bool found;
    struct analyzable_variable payload;
};

HASH_DECL(var_store_hash, struct analyzable_variable);
STACK_DECL(var_store, struct var_store_hash);

void var_store_push_frame(struct var_store *store);
void var_store_pop_frame(struct var_store *store);
struct var_store_res var_store_find(struct var_store *store, const char *key);
struct analyzable_variable *var_store_insert(struct var_store *store, const char *key);

#endif
