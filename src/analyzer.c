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
static struct ast _prepare_vars(struct analyzer_context *ctx, struct ast *var, bool fn_arg);
static struct ast _prepare_fns(struct analyzer_context *ctx, struct ast *fun);
static struct ast *_prepare_expr(struct analyzer_context *ctx, struct ast *expr);
static struct analyzable_type _get_type(struct analyzer_context *ctx, struct ast *type);
static struct analyzable_type _just_cast(struct analyzable_type current, struct analyzable_type target);
static struct analyzable_type _attempt_cast(struct analyzable_type current, struct analyzable_type target);
static bool _expect_type(struct analyzable_type current, const char *name);
static void _assign_args_to_fn(struct analyzer_context *ctx, struct ast *fn, struct ast *first_arg);

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

    var_store_push_frame(&ctx->variables);

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

static struct ast _prepare_vars(struct analyzer_context *ctx, struct ast *var, bool fn_arg)
{
    struct ast variable = {0};
    variable.type = ANALYZE_VAR;

    if (var->type == VAR_DECL) {
        if (fn_arg)
            goto skip1;

        struct var_store_res it = var_store_find(&ctx->variables, var->u.variable_declaration.identifier->u.identifier.str);
        if (it.found) {
            fprintf(stderr, "variable %s already exists!\n", var->u.assignment.left->u.identifier.str);
            abort();
        }
skip1:
        variable.u.a_var.identifier = var->u.variable_declaration.identifier->u.identifier;
        variable.u.a_var.type = _get_type(ctx, var->u.variable_declaration.type);
        variable.u.a_var.is_declaration = true;
    } else if (var->type == VAR_DEF) {
        if (fn_arg)
            goto skip2;

        struct var_store_res it = var_store_find(&ctx->variables, var->u.function_definition.ident->u.identifier.str);
        if (it.found) {
            fprintf(stderr, "variable %s already exists!\n", var->u.assignment.left->u.identifier.str);
            abort();
        }
skip2:
        variable.u.a_var.identifier = var->u.variable_definition.identifier->u.identifier;
        if (var->u.variable_definition.type == NULL) {
            struct ast *prepared = _prepare_expr(ctx, var->u.variable_definition.value);
            variable.u.a_var.type = _get_type(ctx, prepared);
            variable.u.a_var.value = prepared;
        } else {
            variable.u.a_var.type = _get_type(ctx, var->u.variable_definition.type);
            struct ast *prepared = _prepare_expr(ctx, var->u.variable_definition.value);
            variable.u.a_var.type = _attempt_cast(_get_type(ctx, prepared), variable.u.a_var.type);
            variable.u.a_var.value = prepared;
        }
        variable.u.a_var.is_declaration = false;
    } else if (var->type == ASSIGNMENT) {
        if (var->u.assignment.left->type != IDENTIFIER) {
            fprintf(stderr, "expected IDENT on left side of assignment!\n");
            abort();
        }

        if (fn_arg)
            goto skip3;

        struct var_store_res it = var_store_find(&ctx->variables, var->u.assignment.left->u.identifier.str);
        if (it.found) {
            return *var;
        }
skip3:

        variable.u.a_var.identifier = var->u.assignment.left->u.identifier;
        variable.u.a_var.value = _prepare_expr(ctx, var->u.assignment.right);
        variable.u.a_var.type = _get_type(ctx, var->u.assignment.right);
        variable.u.a_var.is_declaration = false;

    }

    if (!fn_arg) {
        struct analyzable_variable *it = var_store_insert(&ctx->variables, variable.u.a_var.identifier.str);
        *it = variable.u.a_var;
    }

    return variable;
}

static struct ast _prepare_fns(struct analyzer_context *ctx, struct ast *fun)
{
    struct ast fn = {0};
    fn.type = ANALYZE_FN;

    if (fun->type == FN_DECL) {
        uint32_t it = function_store_find(&ctx->functions, fun->u.function_declaration.ident->u.identifier.str);
        if (hash_exists(&ctx->functions, it)) {
            struct analyzable_function f = hash_value(&ctx->functions, it);
            fprintf(stderr, "function %s has already been %s!\n", f.name.str, f.declaration ? "declared" : "defined");
            abort();
        }

        fn.u.a_fn.return_type = _get_type(ctx, fun->u.function_declaration.return_type);
        fn.u.a_fn.name = fun->u.function_declaration.ident->u.identifier;
        fn.u.a_fn.body = NULL;

        _assign_args_to_fn(ctx, &fn, fun->u.function_declaration.arg_list);

        fn.u.a_fn.immutable = fun->u.function_declaration.immutable;
        fn.u.a_fn.declaration = true;
        fn.u.a_fn.checked = false;

        it = function_store_insert(&ctx->functions, fn.u.a_fn.name.str);
        hash_value(&ctx->functions, it) = fn.u.a_fn;
    } else if (fun->type == FN_DEF) {
        uint32_t it = function_store_find(&ctx->functions, fun->u.function_definition.ident->u.identifier.str);
        if (hash_exists(&ctx->functions, it)) {
            struct analyzable_function f = hash_value(&ctx->functions, it);
            if (f.declaration == false) {
                fprintf(stderr, "function %s has already been %s!\n", f.name.str, f.declaration ? "declared" : "defined");
                abort();
            }
        }

        fn.u.a_fn.return_type = _get_type(ctx, fun->u.function_definition.return_type);
        fn.u.a_fn.name = fun->u.function_definition.ident->u.identifier;
        fn.u.a_fn.body = fun->u.function_definition.body;

        _assign_args_to_fn(ctx, &fn, fun->u.function_definition.arg_list);

        fn.u.a_fn.immutable = fun->u.function_definition.immutable;
        fn.u.a_fn.declaration = false;
        fn.u.a_fn.checked = false;

        if (hash_exists(&ctx->functions, it)) {
            struct analyzable_function f = hash_value(&ctx->functions, it);
            if (strcmp(f.return_type.identifier.str, fn.u.a_fn.return_type.identifier.str) != 0) {
                fprintf(stderr, "return type missmatch! expected %s got %s!\n", f.return_type.identifier.str,
                    fn.u.a_fn.return_type.identifier.str);
                abort();
            }

            if (f.immutable != fn.u.a_fn.immutable) {
                fprintf(stderr, "function type missmatch! expected %s got %s!\n", f.immutable ? "C" : "Tanzanite",
                    fn.u.a_fn.immutable ? "C" : "Tanzanite");
                abort();
            }

            if (f.variadic != fn.u.a_fn.variadic) {
                fprintf(stderr, "function variadic missmatch! expected %s got %s!\n", f.variadic ? "yes" : "no",
                    fn.u.a_fn.variadic ? "yes" : "no");
                abort();
            }

            if (f.args_count != fn.u.a_fn.args_count) {
                fprintf(stderr, "argument count missmatch! expected %ld got %ld!\n", f.args_count, fn.u.a_fn.args_count);
                abort();
            }

            for (size_t i = 0; i < f.args_count; i++) {
                struct analyzable_fn_arg *decl = f.args + i;
                struct analyzable_fn_arg *def = fn.u.a_fn.args + i;

                if (strcmp(decl->identifier.str, def->identifier.str) != 0) {
                    fprintf(stderr, "%ld. arg name missmatch! expected %s got %s!\n", i + 1, decl->identifier.str,
                        def->identifier.str);
                    abort();
                }

                if (strcmp(decl->type.identifier.str, def->type.identifier.str) != 0) {
                    fprintf(stderr, "%ld. arg type missmatch! expected %s got %s!\n", i + 1, decl->type.identifier.str,
                        def->type.identifier.str);
                    abort();
                }

                if (decl->type.pointer_depth != def->type.pointer_depth) {
                    fprintf(stderr, "%ld. arg pointer depth missmath! expected %ld got %ld!\n", i + 1,
                        decl->type.pointer_depth, def->type.pointer_depth);
                    abort();
                }
            }

            free(fn.u.a_fn.args);
            fn.u.a_fn.args = f.args;
        }

        if (!hash_exists(&ctx->functions, it))
            it = function_store_insert(&ctx->functions, fn.u.a_fn.name.str);

        hash_value(&ctx->functions, it) = fn.u.a_fn;
    }

    return fn;
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
    case ANALYZE_IF:
        return type->u.a_if.result_type;
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
    if (expr == NULL)
        return NULL;

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
        struct var_store_res it = var_store_find(&ctx->variables, expr->u.identifier.str);
        if (!it.found) {
            fprintf(stderr, "variable %s could not be found!\n", expr->u.identifier.str);
            abort();
        }

        v.result = it.payload.type;
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
        } else if (strcmp(expr->u.operation.op, "==") == 0)
            goto assign_bool;
        else if (strcmp(expr->u.operation.op, "!=") == 0)
            goto assign_bool;
        else if (strcmp(expr->u.operation.op, "&&") == 0)
            goto assign_bool;
        else if (strcmp(expr->u.operation.op, "||") == 0) {
assign_bool:
            uint32_t it = type_store_find(&ctx->types, "bool");
            o.result_type = hash_value(&ctx->types, it);
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
    case IF_EXPR: {
        struct analyzable_if cond = {0};
        cond.expression = _prepare_expr(ctx, expr->u.if_expression.expr);
        if (!_expect_type(_get_type(ctx, cond.expression), "bool")) {
            fprintf(stderr, "if/unless expects a bool operation!\n");
            abort();
        }
        cond.body = _prepare_expr(ctx, expr->u.if_expression.val);
        cond.unless = expr->u.if_expression.unless;
        cond.else_op = _prepare_expr(ctx, expr->u.if_expression.else_val);
        cond.result_type = _just_cast(_get_type(ctx, cond.body), _get_type(ctx, cond.else_op));

        expr->type = ANALYZE_IF;
        expr->u.a_if = cond;
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
        *stmt = _prepare_vars(ctx, stmt, false);
        break;
    case FN_DECL:
    case FN_DEF:
        *stmt = _prepare_fns(ctx, stmt);
        break;
    default:
        fprintf(stderr, "did not expect %d in global scope!\n", stmt->type);
        abort();
    }
}

static bool _expect_type(struct analyzable_type current, const char *name)
{
    if (strcmp(current.identifier.str, name) == 0)
        return true;

    return false;
}

static void _assign_args_to_fn(struct analyzer_context *ctx, struct ast *fn, struct ast *first_arg)
{
    bool needs_def_val = false;
    size_t arg_count = 0;

    struct ast *iter = first_arg;
    while (iter != NULL && iter->type == FN_ARG) {
        if (iter->u.function_argument.current->type == VARIADIC) {
            fn->u.a_fn.variadic = true;
            break;
        }
        arg_count++;
        iter = iter->u.function_argument.next;
    }

    fn->u.a_fn.args_count = arg_count;
    fn->u.a_fn.args = calloc(arg_count, sizeof(struct analyzable_fn_arg));

    iter = first_arg;
    for (size_t i = 0; i < arg_count; i++) {
        struct analyzable_fn_arg *ptr = fn->u.a_fn.args + i;
        struct ast prepared = _prepare_vars(ctx, iter->u.function_argument.current, true);
        ptr->type = prepared.u.a_var.type;
        ptr->identifier = prepared.u.a_var.identifier;
        ptr->default_value = prepared.u.a_var.value;

        if (needs_def_val && ptr->default_value == NULL) {
            fprintf(stderr, "%ld. arg %s is expected to have default value!\n", i + 1, ptr->identifier.str);
            abort();
        }

        if (ptr->default_value != NULL)
            needs_def_val = true;

        iter = iter->u.function_argument.next;
    }
}
