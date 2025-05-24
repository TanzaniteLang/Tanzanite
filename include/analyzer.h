#ifndef __ANALYZER_H__
#define __ANALYZER_H__

#include <ast.h>

#include <analyzer/context.h>

struct ast *prepare(struct analyzer_context *ctx, struct ast *to_process);
#endif
