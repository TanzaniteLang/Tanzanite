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
static void _emit_pointer(struct str_builder *b, struct ast *ptr);
static void _emit_fn_decl(struct str_builder *b, struct ast *decl);
static void _emit_fn_def(struct str_builder *b, struct ast *def);
static void _emit_fn_call(struct str_builder *b, struct ast *call);
static void _emit_if(struct str_builder *b, struct ast *cond);
static void _emit_elsif(struct str_builder *b, struct ast *cond);

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
    case OPERATION:
        if (strcmp(a->u.operation.op, "//") == 0) {
            fprintf(stderr, "// not yet supported!");
            abort();
        } else if (strcmp(a->u.operation.op, "|>") == 0) {
            fprintf(stderr, "|> not yet supported!");
            abort();
        } else {
            _emit_c(b, a->u.operation.left);
            str_builder_printf(b, " %s ", a->u.operation.op);
            _emit_c(b, a->u.operation.right);
        }
        break;
    case TYPE_NODE:
        _emit_c(b, a->u.type);
        break;
    case POINTER:
        _emit_pointer(b, a);
        break;
    case VAR_DECL:
        _emit_c(b, a->u.variable_declaration.type);
        _emit_c(b, a->u.variable_declaration.identifier);
        break;
    case VAR_DEF:
        _emit_c(b, a->u.variable_definition.type);
        _emit_c(b, a->u.variable_definition.identifier);
        str_builder_append_cstr(b, " = ");
        _emit_c(b, a->u.variable_definition.value);
        break;
    case FN_DECL:
        _emit_fn_decl(b, a);
        return false;
        break;
    case VARIADIC:
        str_builder_append_cstr(b, "...");
        break;
    case FN_DEF:
        _emit_fn_def(b, a);
        return false;
        break;
    case FN_CALL:
        _emit_fn_call(b, a);
        break;
    case ASSIGNMENT:
        _emit_c(b, a->u.assignment.left);
        str_builder_printf(b, " %s ", a->u.assignment.op);
        _emit_c(b, a->u.assignment.right);
        break;
    case IF_COND:
        _emit_if(b, a);
        if (a->u.if_statement.next != NULL)
            _emit_c(b, a->u.if_statement.next);
        return false;
        break;
    case ELSIF_COND:
        _emit_elsif(b, a);
        if (a->u.if_statement.next != NULL)
            _emit_c(b, a->u.if_statement.next);
        return false;
        break;
    case ELSE_COND:
        str_builder_append_cstr(b, "else {\n");
        _emit_body(b, a->u.else_statement);
        str_builder_append_cstr(b, "}\n");
        return false;
        break;
    case IDENTIFIER_CHAIN:
    case FN_ARG:
    case IF_EXPR:
    case EXPR_IF:
    case FOR:
    case WHILE:
    case FIELD_ACCESS:
    case POINTER_DEREF:
    case TYPE_CAST:
    case NEXT:
    case BREAK:
    default:
        fprintf(stderr, "Unhandled node type %d!\n", a->type);
        abort();
    }

    return true;
}

static void _emit_elsif(struct str_builder *b, struct ast *cond)
{
    str_builder_append_cstr(b, "else if (");

    _emit_c(b, cond->u.if_statement.expr);

    str_builder_append_cstr(b, ") {\n");
    _emit_body(b, cond->u.elsif_statement.body);
    str_builder_append_char(b, '}');
    if (cond->u.elsif_statement.next == NULL)
        str_builder_append_char(b, '\n');
}

static void _emit_if(struct str_builder *b, struct ast *cond)
{
    str_builder_append_cstr(b, "if (");
    if (cond->u.if_statement.unless)
        str_builder_append_cstr(b, "!(");

    _emit_c(b, cond->u.if_statement.expr);

    if (cond->u.if_statement.unless)
        str_builder_append_char(b, ')');
    str_builder_append_cstr(b, ") {\n");
    _emit_body(b, cond->u.if_statement.body);
    str_builder_append_char(b, '}');
    if (cond->u.if_statement.next == NULL)
        str_builder_append_char(b, '\n');
}

static void _emit_fn_call(struct str_builder *b, struct ast *call)
{
    _emit_c(b, call->u.function_call.ident);
    str_builder_append_char(b, '(');
    struct ast *arg_iter = call->u.function_call.first_arg;
    while (arg_iter != NULL) {
        _emit_c(b, arg_iter->u.function_argument.current);
        arg_iter = arg_iter->u.function_argument.next;
        if (arg_iter != NULL)
            str_builder_append_cstr(b, ", ");
    }
    str_builder_append_char(b, ')');
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

static void _emit_fn_decl(struct str_builder *b, struct ast *decl)
{
    _emit_c(b, decl->u.function_definition.return_type);
    _emit_c(b, decl->u.function_definition.ident);
    str_builder_append_char(b, '(');
    struct ast *arg_iter = decl->u.function_definition.arg_list;
    while (arg_iter != NULL) {
        _emit_c(b, arg_iter->u.function_argument.current);
        arg_iter = arg_iter->u.function_argument.next;
        if (arg_iter != NULL)
            str_builder_append_cstr(b, ", ");
    }
    str_builder_append_cstr(b, ");\n\n");
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

    while (iter != NULL && iter->type == STATEMENT) {
        bool res = _emit_c(b, iter->u.statement.current);
        if (res)
            str_builder_append_cstr(b, ";\n");

        iter = iter->u.statement.next;
    }
}
