#ifndef __ANALYZER_FUNCTION_H__
#define __ANALYZER_FUNCTION_H__

#include <stdbool.h>
#include <str.h>
#include <stdint.h>

#include <analyzer/type.h>

struct ast;

struct analyzable_function {
    struct analyzable_type return_type;
    struct str name;
    /* can be NULL */
    struct ast *body;

    struct analyzable_fn_arg *args;
    size_t args_count;
    bool immutable;
    bool declaration;
    /* if true and fn call has more args than args_count, that's fine */
    bool variadic; 
    /* the analyzer has checked the function */
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
