#ifndef __ANALYZER_CONDITIONS_H__
#define __ANALYZER_CONDITIONS_H__

#include <ast.h>
#include <analyzer/type.h>

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
    struct analyzable_type result_type;
    struct ast *expression;
    struct ast *body;
};

#endif
