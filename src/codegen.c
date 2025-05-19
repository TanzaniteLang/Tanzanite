#include <ast.h>
#include <codegen.h>
#include <str.h>
#include <str_builder.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static void _emit_c(struct str_builder *b, struct ast *a);
static void _emit_pointer(struct str_builder *b, struct ast *ptr);

struct str emit_c(struct ast *ast)
{
    struct str_builder b = {0};

    if (ast->type != PROGRAM) {
        fprintf(stderr, "Expected node type PROGRAM!\n");
        abort();
    }

    struct ast *iter = ast->u.program;

    while (iter && iter->type == STATEMENT) {
        _emit_c(&b, iter->u.statement.current);
        str_builder_append_cstr(&b, ";\n");

        iter = iter->u.statement.next;
    }

    struct str s = str_builder_str(&b);
    return s;
}


static void _emit_c(struct str_builder *b, struct ast *a)
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

    case IDENTIFIER_CHAIN:
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
    case ASSIGNMENT:
    case TYPE_CAST:
    case NEXT:
    case BREAK:
    case VARIADIC:
    default:
        fprintf(stderr, "Unhandled node type %d!\n", a->type);
        abort();
    }
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
