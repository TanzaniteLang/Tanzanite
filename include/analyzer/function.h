#ifndef __ANALYZER_FUNCTION_H__
#define __ANALYZER_FUNCTION_H__

#include <stdbool.h>
#include <ast.h>
#include <str.h>
#include <stdint.h>

#include <analyzer/type.h>

struct analyzable_function {
    struct analyzable_type return_type;
    struct str name;
    struct ast *body;

    struct analyzable_fn_arg *args;
    size_t args_count;
    bool immutable;
    /* if true and fn call has more args than args_count, that's fine */
    bool variadic; 
    bool checked;
};

struct analyzable_fn_arg {
    struct analyzable_type type;
    struct str identifier;
    struct ast *default_value;
};

struct analyzable_call_arg {
    struct ast *value;
};

struct analyzable_call {
    struct str identifier;
    struct analyzable_call_arg *args;
    size_t args_count;
};

#endif
