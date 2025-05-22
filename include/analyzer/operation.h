#ifndef __ANALYZER_OPERATION_H__
#define __ANALYZER_OPERATION_H__

#include <analyzer/type.h>

struct ast;

struct analyzable_operation {
    struct analyzable_type result_type;
    const char *operation;

    struct ast *left;
    struct ast *right;
};

#endif
