#ifndef __ANALYZER_VALUE_H__
#define __ANALYZER_VALUE_H__

#include <analyzer/type.h>

struct ast;

struct analyzable_value {
    struct analyzable_type result;
    struct ast *value;
};

#endif
