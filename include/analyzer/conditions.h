#ifndef __ANALYZER_CONDITIONS_H__
#define __ANALYZER_CONDITIONS_H__

#include <analyzer/type.h>

#include <ast.h>

struct analyzable_if {
    struct analyzable_type result_type;
    struct ast *expression;
    struct ast *body;
    bool unless;

    /* can be NULL */
    struct analyzable_elsif *elsifs;
    size_t elsifs_count;

    /* can be NULL */
    struct ast *else_op;
};

struct analyzable_elsif {
    struct ast *expression;
    struct ast *body;
};

#endif
