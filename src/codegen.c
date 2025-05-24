#include <ast.h>
#include <codegen.h>
#include <str.h>
#include <str_builder.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static bool _emit_c(struct str_builder *b, struct ast *a);
static void _emit_body(struct str_builder *b, struct ast *body);
static void _emit_fn(struct str_builder *b, struct analyzable_function *fn);
static void _emit_type(struct str_builder *b, struct analyzable_type *type);
static void _emit_type_cast(struct str_builder *b, struct analyzable_type *type);
static void _emit_fn_call(struct str_builder *b, struct analyzable_call *call);
static void _emit_for(struct str_builder *b, struct analyzable_for *loop);
static void _emit_while(struct str_builder *b, struct analyzable_while *loop);
static void _emit_if(struct str_builder *b, struct analyzable_if *cond);
static void _emit_elsif(struct str_builder *b, struct analyzable_elsif *cond);

static void _emit_fn_decl(struct str_builder *b, struct ast *decl);
static void _emit_fn_def(struct str_builder *b, struct ast *def);

struct str emit_c(struct ast *ast)
{
    struct str_builder b = {0};

    if (ast->type != PROGRAM) {
        fprintf(stderr, "Expected node type PROGRAM!\n");
        abort();
    }

    _emit_body(&b, ast->u.program);

    struct str s = str_builder_str(&b);
    return s;
}


static bool _emit_c(struct str_builder *b, struct ast *a)
{
    switch (a->type) {
    case INT:
        str_builder_printf(b, "%ld", a->u.number);
        break;
    case FLOAT:
        str_builder_printf(b, "%f", a->u.decimal);
        break;
    case IDENTIFIER:
        str_builder_printf(b, "%s", a->u.identifier);
        break;
    case CHAR:
        str_builder_printf(b, "'%c'", a->u.ch);
        break;
    case BOOL:
        str_builder_printf(b, "%s", a->u.boolean ? "true" : "false");
        break;
    case STRING:
        str_builder_printf(b, "\"%s\"", a->u.string.str);
        break;
    case BRACKETS:
        str_builder_append_char(b, '(');
        _emit_c(b, a->u.bracket);
        str_builder_append_char(b, ')');
        break;
    case UNARY:
        str_builder_printf(b, "%s", a->u.unary.op);
        _emit_c(b, a->u.unary.value);
        break;
    case ANALYZE_VALUE:
        _emit_type_cast(b, &a->u.a_value.result);
        _emit_c(b, a->u.a_value.value);
        break;
    case ANALYZE_VAR:
        _emit_type(b, &a->u.a_var.type);
        str_builder_append_char(b, ' ');
        str_builder_append_cstr(b, a->u.a_var.identifier.str);
        if (!a->u.a_var.is_declaration) {
            str_builder_append_cstr(b, " = ");
            _emit_c(b, a->u.a_var.value);
        }
        break;
    case ANALYZE_FN:
        _emit_fn(b, &a->u.a_fn);
        return false;
        break;
    case ANALYZE_FN_CALL:
        _emit_type_cast(b, &a->u.a_fn_call.result_type);
        _emit_fn_call(b, &a->u.a_fn_call);
        break;
    case ANALYZE_OPERATION:
        _emit_type_cast(b, &a->u.a_operation.result_type);
        _emit_c(b, a->u.a_operation.left);
        str_builder_append_cstr(b, a->u.a_operation.operation);
        _emit_c(b, a->u.a_operation.right);
        break;
    case ANALYZE_TYPE_CAST:
        _emit_type_cast(b, &a->u.a_cast.target);
        _emit_c(b, a->u.a_cast.value);
        break;
    case ANALYZE_FOR:
        _emit_for(b, &a->u.a_for);
        return false;
        break;
    case ANALYZE_WHILE:
        _emit_while(b, &a->u.a_while);
        return false;
        break;
    case ASSIGNMENT:
        _emit_c(b, a->u.assignment.left);
        str_builder_printf(b, " %s ", a->u.assignment.op);
        _emit_c(b, a->u.assignment.right);
        break;
    case ANALYZE_IF:
        _emit_if(b, &a->u.a_if);
        return false;
        break;
    case NEXT:
        str_builder_append_cstr(b, "continue");
        break;
    case BREAK:
        str_builder_append_cstr(b, "break");
        break;
    case STATEMENT:
    case IDENTIFIER_CHAIN:
    case OPERATION:
    case VAR_DECL:
    case VAR_DEF:
    case TYPE_NODE:
    case POINTER:
    case FN_DECL:
    case FN_DEF:
    case FN_ARG:
    case FN_CALL:
    case IF_COND:
    case IF_EXPR:
    case EXPR_IF:
    case ELSIF_COND:
    case ELSE_COND:
    case FOR:
    case WHILE:
    case FIELD_ACCESS:
    case POINTER_DEREF:
    case TYPE_CAST:
    case VARIADIC:
    case RANGE:
        fprintf(stderr, "Unhandled node type %d!\n", a->type);
        abort();
    }

    return true;
}

static void _emit_fn(struct str_builder *b, struct analyzable_function *fn)
{
    _emit_type(b, &fn->return_type);
    str_builder_printf(b, " %s(",  fn->name.str);
    for (size_t i = 0; i < fn->args_count; i++) {
        struct analyzable_fn_arg *arg = fn->args + i;
        _emit_type(b, &arg->type);
        str_builder_printf(b, " %s",  arg->identifier.str);
        if (i + 1 < fn->args_count || fn->variadic)
            str_builder_append_cstr(b, ", ");
    }

    if (fn->variadic)
        str_builder_append_cstr(b, "...");

    str_builder_append_char(b, ')');

    if (fn->declaration) {
        str_builder_append_cstr(b, ";\n\n");
        return;
    }
    str_builder_append_cstr(b, "\n{\n");
    _emit_body(b, fn->body);
    str_builder_append_cstr(b, "}\n\n");
}

static void _emit_type(struct str_builder *b, struct analyzable_type *type)
{
    str_builder_append_cstr(b, type->identifier.str);
    for (size_t i = 0; i < type->pointer_depth; i++)
        str_builder_append_char(b, '*');
}

static void _emit_type_cast(struct str_builder *b, struct analyzable_type *type)
{
    str_builder_append_char(b, '(');
    _emit_type(b, type);
    str_builder_append_char(b, ')');
}

static void _emit_fn_call(struct str_builder *b, struct analyzable_call *call)
{
    str_builder_append_cstr(b, call->identifier.str);
    str_builder_append_char(b, '(');
    for (size_t i = 0; i < call->args_count; i++) {
        struct analyzable_call_arg *arg = call->args + i;
        _emit_c(b, arg->value);
        if (i + 1 < call->args_count)
            str_builder_append_cstr(b, ", ");
    }

    str_builder_append_char(b, ')');
}

static void _emit_for(struct str_builder *b, struct analyzable_for *loop)
{
    if (loop->payload_count == 1 && loop->expr->type == RANGE) {
        struct analyzable_type t = loop->payloads[0].type;
        str_builder_append_cstr(b, "for (");
        _emit_type(b, &t);
        str_builder_printf(b, " %s = %ld;", loop->payloads[0].identifier.str, loop->expr->u.range.start);
        str_builder_printf(b, " %s <= %ld;", loop->payloads[0].identifier.str, loop->expr->u.range.end);
        str_builder_printf(b, " %s++) {\n", loop->payloads[0].identifier.str);
        _emit_body(b, loop->body);
        str_builder_append_cstr(b, "}\n");
    } else {
        fprintf(stderr, "XXX: very limited, only to range with payload!\n");
        abort();
    }
}

static void _emit_while(struct str_builder *b, struct analyzable_while *loop)
{
    if (loop->infinite) {
        str_builder_append_cstr(b, "while (true) {\n");
    } else {
        if (loop->until)
            str_builder_append_cstr(b, "until (");
        else
            str_builder_append_cstr(b, "while (");
        _emit_c(b, loop->expr);
        str_builder_append_cstr(b, ") {\n");
    }

    _emit_body(b, loop->body);
    str_builder_append_cstr(b, "}\n");
}

static void _emit_if(struct str_builder *b, struct analyzable_if *cond)
{
    if (cond->unless)
        str_builder_append_cstr(b, "unless (");
    else
        str_builder_append_cstr(b, "if (");

    _emit_c(b, cond->expression);
    str_builder_append_cstr(b, ") {\n");
    _emit_body(b, cond->body);
    str_builder_append_cstr(b, "} ");

    for (size_t i = 0; i < cond->elsifs_count; i++) {
        _emit_elsif(b, cond->elsifs + i);
    }

    str_builder_append_cstr(b, "\n");
}

static void _emit_elsif(struct str_builder *b, struct analyzable_elsif *cond)
{
    str_builder_append_cstr(b, "else if (");

    _emit_c(b, cond->expression);

    str_builder_append_cstr(b, ") {\n");
    _emit_body(b, cond->body);
    str_builder_append_cstr(b, "} ");
}




static void _emit_fn_def(struct str_builder *b, struct ast *def)
{
    _emit_c(b, def->u.function_declaration.return_type);
    _emit_c(b, def->u.function_declaration.ident);
    str_builder_append_char(b, '(');
    struct ast *arg_iter = def->u.function_declaration.arg_list;
    while (arg_iter != NULL) {
        _emit_c(b, arg_iter->u.function_argument.current);
        arg_iter = arg_iter->u.function_argument.next;
        if (arg_iter != NULL)
            str_builder_append_cstr(b, ", ");
    }
    str_builder_append_cstr(b, ")\n{\n");
    _emit_body(b, def->u.function_definition.body);
    str_builder_append_cstr(b, "}\n\n");
}

static void _emit_pointer(struct str_builder *b, struct ast *ptr)
{
    struct ast *iter = ptr;
    while (iter->type == POINTER) {
        if (iter->u.pointer.next == NULL)
            iter = iter->u.pointer.current;
        else
            iter = iter->u.pointer.next;
    }

    _emit_c(b, iter);
    str_builder_append_char(b, ' ');
    iter = ptr->u.pointer.next;
    while (iter != NULL && iter->type == POINTER) {
        str_builder_append_char(b, '*');
        iter = iter->u.pointer.next;
    }
}

static void _emit_body(struct str_builder *b, struct ast *body)
{
    struct ast *iter = body;

    if (iter != NULL && iter->type != STATEMENT) {
        bool res = _emit_c(b, body);
        if (res)
            str_builder_append_cstr(b, ";\n");
    }

    while (iter != NULL && iter->type == STATEMENT) {
        bool res = _emit_c(b, iter->u.statement.current);
        if (res)
            str_builder_append_cstr(b, ";\n");

        iter = iter->u.statement.next;
    }
}
