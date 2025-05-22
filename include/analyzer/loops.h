#ifndef __ANALYZER_LOOPS_H__
#define __ANALYZER_LOOPS_H__

#include <analyzer/type.h>
#include <stdbool.h>
#include <str.h>

struct analyzable_for {
    struct analyzable_type payload_type;
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
    struct analyzable_type result_type;
    struct ast *expr;
    struct ast *body;

    bool do_while;
    bool until;
};

#endif
