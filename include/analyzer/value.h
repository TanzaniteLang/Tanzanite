#ifndef __ANALYZER_VALUE_H__
#define __ANALYZER_VALUE_H__

#include <ast.h>
#include <analyzer/type.h>

struct analyzable_value {
    struct analyzable_type result;
    struct ast *value;
};

#endif
