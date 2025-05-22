#include <ast.h>
#include <analyzer/context.h>
#include <analyzer.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <float.h>
#include <hash/type_store.h>

struct builtin_types {
    char *name;
    size_t size;
};

static const struct builtin_types types[] = {
    { "bool", 1                 },
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
    { "isize", 8                },
    { "usize", 8                },
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
static struct ast *_prepare_expr(struct analyzer_context *ctx, struct ast *expr);
static struct analyzable_type _get_type(struct analyzer_context *ctx, struct ast *type);
static struct analyzable_type _just_cast(struct analyzable_type current, struct analyzable_type target);
static struct analyzable_type _attempt_cast(struct analyzable_type current, struct analyzable_type target);

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
    struct ast v = {0};
    v.type = ANALYZE_VAR;

    if (var->type == VAR_DECL) {
        v.u.a_var.identifier = var->u.variable_declaration.identifier->u.identifier;
        v.u.a_var.type = _get_type(ctx, var->u.variable_declaration.type);
        v.u.a_var.is_declaration = true;
    } else if (var->type == VAR_DEF) {
        v.u.a_var.identifier = var->u.variable_definition.identifier->u.identifier;
        if (var->u.variable_definition.type == NULL) {
            struct ast *prepared = _prepare_expr(ctx, var->u.variable_definition.value);
            v.u.a_var.type = _get_type(ctx, prepared);
            v.u.a_var.value = prepared;
        } else {
            v.u.a_var.type = _get_type(ctx, var->u.variable_definition.type);
            struct ast *prepared = _prepare_expr(ctx, var->u.variable_definition.value);
            v.u.a_var.type = _attempt_cast(_get_type(ctx, prepared), v.u.a_var.type);
            v.u.a_var.value = prepared;
        }
        v.u.a_var.is_declaration = false;
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

        it = var_store_insert(&ctx->variables, var->u.assignment.left->u.identifier.str);

        v.u.a_var.identifier = var->u.assignment.left->u.identifier;
        v.u.a_var.value = _prepare_expr(ctx, var->u.assignment.right);
        v.u.a_var.type = _get_type(ctx, var->u.assignment.right);
        v.u.a_var.is_declaration = false;

    }

    uint32_t it = var_store_insert(&ctx->variables, v.u.a_var.identifier.str);
    hash_value(&ctx->variables, it) = v.u.a_var;

    return v;
}

static struct analyzable_type _get_type(struct analyzer_context *ctx, struct ast *type)
{
    struct analyzable_type t = {0};

start:

    switch (type->type) {
    case TYPE_NODE:
        type = type->u.type;
        goto start;
        break;
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
    case ANALYZE_VALUE:
        return type->u.a_value.result;
    case ANALYZE_OPERATION:
        return type->u.a_operation.result_type;
    case ANALYZE_TYPE_CAST:
        return type->u.a_cast.target;
    case BRACKETS:
        return _get_type(ctx, type->u.bracket);
    default:
        fprintf(stderr, "expected type nodes, got %d!\n", type->type);
        abort();
    }

    return t;
}

static struct analyzable_type _just_cast(struct analyzable_type current, struct analyzable_type target)
{
    if (current.size > target.size)
        return current;
    return target;
}

static struct analyzable_type _attempt_cast(struct analyzable_type current, struct analyzable_type target)
{
    if (current.size > target.size)
        fprintf(stderr, "warning: target type is smaller than current size, value will be truncated!\n");

    return target;
}

static struct ast *_prepare_expr(struct analyzer_context *ctx, struct ast *expr)
{
    switch (expr->type) {
    case BRACKETS:
        expr->u.bracket = _prepare_expr(ctx, expr->u.bracket);
        break;
    case INT: {
        struct analyzable_value v = {0};
        const char *type = NULL;

        v.value = dup_node(expr);
        int64_t val = expr->u.number;
        if (val > 0) {
            if (val <= INT8_MAX)
                type = "i8";
            else if (val <= INT16_MAX)
                type = "i16";
            else if (val <= INT32_MAX)
                type = "i32";
            else
                type = "i64";
        } else {
            if (val >= INT8_MIN)
                type = "i8";
            else if (val >= INT16_MIN)
                type = "i16";
            else if (val >= INT32_MIN)
                type = "i32";
            else
                type = "i64";
        }
        uint32_t it = type_store_find(&ctx->types, type);

        v.result = hash_value(&ctx->types, it);
        expr->type = ANALYZE_VALUE;
        expr->u.a_value = v;
        }
        break;
    case FLOAT: {
        struct analyzable_value v = {0};
        const char *type = NULL;

        v.value = dup_node(expr);
        double val = expr->u.decimal;
        if (val > 0) {
            if (val <= FLT_MAX)
                type = "f32";
            else
                type = "f64";
        } else {
            if (val >= FLT_MIN)
                type = "f32";
            else
                type = "f64";
        }
        uint32_t it = type_store_find(&ctx->types, type);

        v.result = hash_value(&ctx->types, it);
        expr->type = ANALYZE_VALUE;
        expr->u.a_value = v;
        }
        break;
    case IDENTIFIER: {
        struct analyzable_value v = {0};

        v.value = dup_node(expr);
        uint32_t it = var_store_find(&ctx->variables, expr->u.identifier.str);
        if (it == ctx->variables.cap) {
            fprintf(stderr, "variable %s could not be found!\n", expr->u.identifier.str);
            abort();
        }

        v.result = hash_value(&ctx->variables, it).type;
        expr->type = ANALYZE_VALUE;
        expr->u.a_value = v;
        }
        break;
    case CHAR: {
        struct analyzable_value v = {0};

        v.value = dup_node(expr);
        uint32_t it = type_store_find(&ctx->types, "u8");

        v.result = hash_value(&ctx->types, it);
        expr->type = ANALYZE_VALUE;
        expr->u.a_value = v;
        }
        break;
    case BOOL: {
        struct analyzable_value v = {0};

        v.value = dup_node(expr);
        uint32_t it = type_store_find(&ctx->types, "bool");

        v.result = hash_value(&ctx->types, it);
        expr->type = ANALYZE_VALUE;
        expr->u.a_value = v;
        }
        break;
    case STRING: {
        struct analyzable_value v = {0};

        v.value = dup_node(expr);
        uint32_t it = type_store_find(&ctx->types, "u8");

        v.result = hash_value(&ctx->types, it);
        v.result.pointer_depth++;
        expr->type = ANALYZE_VALUE;
        expr->u.a_value = v;
        }
        break;
    case UNARY: {
        struct analyzable_value v = {0};

        expr->u.unary.value = _prepare_expr(ctx, expr->u.unary.value);
        struct analyzable_type t = _get_type(ctx, expr->u.unary.value);
        if (strcmp(expr->u.unary.op, "&") == 0)
            t.pointer_depth++;
        else if (strcmp(expr->u.unary.op, "sizeof") == 0) {
            uint32_t it = type_store_find(&ctx->types, "usize");
            t = hash_value(&ctx->types, it);
        }
        v.value = dup_node(expr);
        v.result = t;

        expr->type = ANALYZE_VALUE;
        expr->u.a_value = v;
        }
        break;
    case POINTER_DEREF: {
        struct analyzable_value v = {0};
        expr->u.to_deref = _prepare_expr(ctx, expr->u.to_deref);
        struct analyzable_type t = _get_type(ctx, expr->u.to_deref);

        if (t.pointer_depth < 1) {
            fprintf(stderr, "expected pointer type, got %s!\n", t.identifier.str);
            abort();
        }

        t.pointer_depth--;
        v.value = dup_node(expr);
        v.result = t;

        expr->type = ANALYZE_VALUE;
        expr->u.a_value = v;
        }
        break;
    case OPERATION: {
        struct analyzable_operation o = {0};
        expr->u.operation.left = _prepare_expr(ctx, expr->u.operation.left);
        expr->u.operation.right = _prepare_expr(ctx, expr->u.operation.right);

        o.result_type = _just_cast(_get_type(ctx, expr->u.operation.left),
                                      _get_type(ctx, expr->u.operation.right));
        /* TODO: THIS */
        if (strcmp(expr->u.operation.op, "//") == 0) {
            fprintf(stderr, "// is not supported yet!\n");
            abort();
        } else if (strcmp(expr->u.operation.op, "|>") == 0) {
            fprintf(stderr, "|> is not supported yet!\n");
            abort();
        }
        o.left = expr->u.operation.left;
        o.right = expr->u.operation.right;
        o.operation = expr->u.operation.op;

        expr->type = ANALYZE_OPERATION;
        expr->u.a_operation = o;
        }
        break;
    case TYPE_CAST: {
        struct analyzable_cast c = {0};
        expr->u.type_cast.expr = _prepare_expr(ctx, expr->u.type_cast.expr);
        c.target = _attempt_cast(_get_type(ctx, expr->u.type_cast.expr), _get_type(ctx, expr->u.type_cast.type));

        c.value = expr->u.type_cast.expr;

        expr->type = ANALYZE_TYPE_CAST;
        expr->u.a_cast = c;
        }
        break;
    case FIELD_ACCESS:
    default:
        fprintf(stderr, "did not expect %d in expression!\n", expr->type);
        abort();
    }

    return expr;
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
