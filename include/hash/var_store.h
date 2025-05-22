#ifndef __HASH_VAR_STORE_H__
#define __HASH_VAR_STORE_H__

#include <hash.h>
#include <stddef.h>

#include <analyzer/variable.h>

HASH_DECL(var_store, struct analyzable_variable);

#endif
