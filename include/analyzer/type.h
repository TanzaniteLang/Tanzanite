#ifndef __ANALYZER_TYPE_H__
#define __ANALYZER_TYPE_H__

#include <str.h>
#include <stddef.h>

struct ast;

struct analyzable_type {
    struct str identifier;
    size_t pointer_depth;
    uint32_t size;
};

struct analyzable_cast {
    struct analyzable_type target;
    struct ast *value;
};

#endif
