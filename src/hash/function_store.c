#include <stdbool.h>
#include <stdlib.h>
#include <string.h>

#include <hash.h>
#include <stddef.h>

#include <analyzer/function.h>
#include <hash/function_store.h>

HASH_IMPL(function_store, struct analyzable_function);
