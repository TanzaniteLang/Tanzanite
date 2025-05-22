#ifndef __ANALYZER_VARIABLE_H__
#define __ANALYZER_VARIABLE_H__

#include <stdbool.h>
#include <ast.h>
#include <str.h>
#include <stdint.h>

#include <analyzer/type.h>

struct analyzable_variable {
    struct analyzable_type type;
    struct str identifier;
    struct ast *value;

    bool is_declaration;
};

#endif
