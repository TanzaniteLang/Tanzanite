#ifndef __ANALYZER_CONTEXT_H__
#define __ANALYZER_CONTEXT_H__

#include <hash/type_store.h>
#include <hash/var_store.h>
#include <hash/function_store.h>

struct analyzer_context {
    struct type_store types;
    struct var_store variables;
    struct function_store functions;
};

#endif
