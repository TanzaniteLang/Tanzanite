#include <ast.h>
#include <analyzer/context.h>
#include <analyzer.h>
#include <stdio.h>
#include <stdlib.h>

#include <stdint.h>
#include <hash/type_store.h>

struct builtin_types {
    char *name;
    size_t size;
};

static const struct builtin_types types[] = {
    { "i8", 1                   },
    { "u8", 1                   },
    { "i16", 2                  },
    { "u16", 2                  },
    { "i32", 4                  },
    { "u32", 4                  },
    { "i64", 8                  },
    { "u64", 8                  },
    { "f32", 4                  },
    { "f64", 8                  },
    { "void", 0                 },
    { "char", sizeof(char)      },
    { "short", sizeof(short)    },
    { "int", sizeof(int)        },
    { "long", sizeof(long)      },
    { "size_t", sizeof(size_t)  },
    { "float", sizeof(float)    },
    { "double", sizeof(double)  },
    { NULL, 0 },
};

static void _prepare_global_statement(struct analyzer_context *ctx, struct ast *stmt);
static struct ast _prepare_vars(struct analyzer_context *ctx, struct ast *var);
static struct analyzable_type _get_type(struct analyzer_context *ctx, struct ast *type);

struct ast *prepare(struct analyzer_context *ctx, struct ast *to_process)
{
    const struct builtin_types *type_iter = types;
    while (type_iter->name != NULL) {
        uint32_t it = type_store_insert(&ctx->types, type_iter->name);
        hash_value(&ctx->types, it).identifier.str = type_iter->name;
        hash_value(&ctx->types, it).size = type_iter->size;
        hash_value(&ctx->types, it).pointer_depth = 0;

        type_iter++;
    }

    if (to_process->type != PROGRAM) {
        fprintf(stderr, "expected PROGRAM, got %d!\n", to_process->type);
        abort();
    }

    struct ast *iter = to_process->u.program;
    while (iter != NULL) {
        if (iter->type != STATEMENT) {
            fprintf(stderr, "expected STATEMENT, got %d!\n", iter->type);
            abort();
        }

        _prepare_global_statement(ctx, iter->u.statement.current);

        iter = iter->u.statement.next;
    }

    return to_process;
}

static struct ast _prepare_vars(struct analyzer_context *ctx, struct ast *var)
{
    (void) ctx;
    struct ast v = {0};
    v.type = ANALYZE_VAR;
    if (var->type == VAR_DECL) {
        v.u.a_var.identifier = var->u.variable_declaration.identifier->u.identifier;
        v.u.a_var.type = _get_type(ctx, var->u.variable_declaration.type);
        v.u.a_var.is_declaration = true;
    } else if (var->type == VAR_DEF) {
    } else if (var->type == ASSIGNMENT) {
        if (var->u.assignment.left->type != IDENTIFIER) {
            fprintf(stderr, "expected IDENT on left side of assignment!\n");
            abort();
        }

        uint32_t it = var_store_find(&ctx->variables, var->u.assignment.left->u.identifier.str);
        if (hash_exists(&ctx->variables, it)) {
            fprintf(stderr, "variable %s already exists!\n", var->u.assignment.left->u.identifier.str);
            abort();
        }

        v.u.a_var.identifier = var->u.assignment.left->u.identifier;
        v.u.a_var.value = var->u.assignment.right;
        v.u.a_var.type = _get_type(ctx, v.u.a_var.value);
        v.u.a_var.is_declaration = false;

        hash_value(&ctx->variables, it) = v.u.a_var;
    }

    return v;
}

static struct analyzable_type _get_type(struct analyzer_context *ctx, struct ast *type)
{
    struct analyzable_type t = {0};

    type = type->u.type;

    switch (type->type) {
    case POINTER: {
            struct ast *ptr_iter = type;
            uint8_t depth = 0;

            while (ptr_iter->u.pointer.next != NULL) {
                depth++;
                ptr_iter = ptr_iter->u.pointer.next;
            }

            type = ptr_iter->u.pointer.current;
            uint32_t it = type_store_find(&ctx->types, type->u.identifier.str);
            if (!hash_exists(&ctx->types, it)) {
                fprintf(stderr, "unable to resolve type: %s!\n", type->u.identifier.str);
                abort();
            }

            t = hash_value(&ctx->types, it);
            t.pointer_depth = depth;
        }
        break;
    default:
        fprintf(stderr, "expected type nodes, got %d!\n", type->type);
        abort();
    }

    return t;
}

static void _prepare_global_statement(struct analyzer_context *ctx, struct ast *stmt)
{
    switch (stmt->type) {
    case ASSIGNMENT:
    case VAR_DECL:
    case VAR_DEF:
        *stmt = _prepare_vars(ctx, stmt);
        break;
    case FN_DECL:
    case FN_DEF:
    default:
        fprintf(stderr, "did not expect %d in global scope!\n", stmt->type);
        abort();
    }
}
