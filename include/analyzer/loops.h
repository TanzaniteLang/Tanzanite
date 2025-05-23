#ifndef __ANALYZER_LOOPS_H__
#define __ANALYZER_LOOPS_H__

#include <analyzer/type.h>
#include <stdbool.h>
#include <str.h>

struct analyzable_for {
    struct ast *expr;

    struct analyzable_payload *payloads;
    size_t payload_count;

    struct ast *body;
};

struct analyzable_payload {
    struct analyzable_type type;
    struct str identifier;
};

struct analyzable_while {
    struct ast *expr;
    struct ast *body;

    bool infinite;
    bool until;
};

#endif
