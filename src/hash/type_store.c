#include <stdbool.h>
#include <stdlib.h>
#include <string.h>

#include <hash.h>

#include <analyzer/type.h>
#include <hash/type_store.h>

HASH_IMPL(type_store, struct analyzable_type);
