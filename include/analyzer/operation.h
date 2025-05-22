#ifndef __ANALYZER_OPERATION_H__
#define __ANALYZER_OPERATION_H__

#include <ast.h>
#include <analyzer/type.h>

struct analyzable_operation {
    struct analyzable_type result_type;
    const char *operation;

    struct ast *left;
    struct ast *right;
};

#endif
