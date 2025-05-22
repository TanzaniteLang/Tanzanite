#include <stdbool.h>
#include <stdlib.h>
#include <string.h>

#include <hash.h>
#include <stddef.h>

#include <analyzer/variable.h>
#include <hash/var_store.h>

HASH_IMPL(var_store, struct analyzable_variable);
