#include <ast.h>
#include <str.h>

#include <stdio.h>
#include <stdlib.h>

struct ast *program_node(struct ast *statement)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = PROGRAM;
    node->u.program = statement;

    return node;
}

struct ast *statement_node(struct ast *list, struct ast *statement)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = STATEMENT;
    node->u.statement.current = statement;
    node->u.statement.next = list;

    return node;
}

struct ast *int_node(uint64_t val)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = INT;
    node->u.number = val;

    return node;
}

struct ast *float_node(double val)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = FLOAT;
    node->u.decimal = val;

    return node;
}

struct ast *identifier_node(struct str ident)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = IDENTIFIER;
    node->u.identifier = ident;

    return node;
}

struct ast *identifier_chain_node(struct ast *list, struct ast *ident)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = IDENTIFIER_CHAIN;
    node->u.identifier_chain.current = ident;
    node->u.identifier_chain.next = list;

    return node;
}

struct ast *operation_node(char op, struct ast *left, struct ast *right)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = OPERATION;
    node->u.operation.op = op;
    node->u.operation.left = left;
    node->u.operation.right = right;

    return node;
}

struct ast *bracket_node(struct ast *expr)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = BRACKETS;
    node->u.bracket = expr;

    return node;
}

struct ast *var_decl_node(struct ast *type, struct ast *ident)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = VAR_DECL;
    node->u.variable_declaration.type = type;
    node->u.variable_declaration.identifier = ident;

    return node;
}

struct ast *var_def_node(struct ast *type, struct ast *ident, struct ast *val)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = VAR_DEF;
    node->u.variable_definition.type = type;
    node->u.variable_definition.identifier = ident;
    node->u.variable_definition.value = val;

    return node;
}

struct ast *fn_decl_node(struct ast *type, struct ast *ident, struct ast *args)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = FN_DECL;
    node->u.function_declaration.return_type = type;
    node->u.function_declaration.ident = ident;
    node->u.function_declaration.arg_list = args;

    return node;
}
struct ast *fn_def_node(struct ast *type, struct ast *ident, struct ast *args, struct ast *body)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = FN_DEF;
    node->u.function_definition.return_type = type;
    node->u.function_definition.ident = ident;
    node->u.function_definition.arg_list = args;
    node->u.function_definition.body = body;

    return node;
}

struct ast *fn_call_node(struct ast *ident, struct ast *first_arg)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = FN_CALL;
    node->u.function_call.ident = ident;
    node->u.function_call.first_arg = first_arg;

    return node;
}

struct ast *fn_arg_list_node(struct ast *list, struct ast *arg)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = FN_ARG;
    node->u.function_argument.next = list;
    node->u.function_argument.current = arg;

    return node;
}

struct ast *type_node(struct ast *type)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = TYPE_NODE;
    node->u.type = type;

    return node;
}

struct ast *pointer_node(struct ast *list, struct ast *type)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = POINTER;
    node->u.pointer.next = list;
    node->u.pointer.current = type;

    return node;
}



static void offset_text(int count)
{
    for (int i = 0; i < count; i++)
        putchar(' ');
}

static void _describe(struct ast *node, int spacing)
{
    if (node == NULL) {
        offset_text(spacing);
        printf("nil\n");
        return;
    }

    switch (node->type) {
    case PROGRAM:
        offset_text(spacing);
        printf("Program {\n");
        _describe(node->u.program, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case STATEMENT:
        offset_text(spacing);
        printf("\e[33mStatement\e[0m {\n");
        _describe(node->u.statement.current, spacing + 2);
        offset_text(spacing);
        if (node->u.statement.next != NULL)
            printf("},\n");
        else
            printf("}\n");
        _describe(node->u.statement.next, spacing);
        break;
    case BRACKETS:
        offset_text(spacing);
        printf("\e[32mBrackets\e[0m (\n");
        _describe(node->u.bracket, spacing + 2);
        offset_text(spacing);
        printf(")\n");
        break;
    case INT:
        offset_text(spacing);
        printf("\e[36mInt\e[0m: %ld\n", node->u.number);
        break;
    case FLOAT:
        offset_text(spacing);
        printf("\e[36mFloat\e[0m: %f\n", node->u.decimal);
        break;
    case IDENTIFIER:
        offset_text(spacing);
        printf("\e[36mIdent\e[0m: %s\n", node->u.identifier.str);
        break;
    case OPERATION:
        offset_text(spacing);
        printf("\e[35mOperation\e[0m: %c\e[0m {\n", node->u.operation.op);
        _describe(node->u.operation.left, spacing + 2);
        _describe(node->u.operation.right, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case VAR_DECL:
        offset_text(spacing);
        printf("\e[34mVar Decl\e[0m {\n");
        _describe(node->u.variable_declaration.type, spacing + 2);
        _describe(node->u.variable_declaration.identifier, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case IDENTIFIER_CHAIN:
        offset_text(spacing);
        printf("\e[33mIdent Chain\e[0m {\n");
        _describe(node->u.identifier_chain.current, spacing + 2);
        offset_text(spacing);
        if (node->u.identifier_chain.next != NULL)
            printf("},\n");
        else
            printf("}\n");
        _describe(node->u.identifier_chain.next, spacing);
        break;
    case TYPE_NODE:
        offset_text(spacing);
        printf("\e[31mType\e[0m {\n");
        _describe(node->u.type, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case POINTER:
        if (node->u.pointer.current != NULL) {
            _describe(node->u.pointer.current, spacing);
        } else {
            offset_text(spacing);
            printf("\e[35mPointer\e[0m {\n");
            _describe(node->u.pointer.next, spacing + 2);
            offset_text(spacing);
            printf("}\n");
        }
        break;
    case VAR_DEF:
        offset_text(spacing);
        printf("\e[34mVar Def\e[0m {\n");
        _describe(node->u.variable_definition.type, spacing + 2);
        _describe(node->u.variable_definition.identifier, spacing + 2);
        _describe(node->u.variable_definition.value, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case FN_DECL:
        offset_text(spacing);
        printf("\e[34mFn Decl\e[0m {\n");
        _describe(node->u.function_declaration.return_type, spacing + 2);
        _describe(node->u.variable_declaration.identifier, spacing + 2);
        spacing += 2;
        offset_text(spacing);
        printf("\e[32mArguments\e[0m (\n");
        _describe(node->u.function_declaration.arg_list, spacing + 2);
        offset_text(spacing);
        printf(")\n");
        spacing -= 2;
        offset_text(spacing);
        printf("}\n");
        break;
    case FN_ARG:
        offset_text(spacing);
        printf("\e[33mFn Args\e[0m {\n");
        _describe(node->u.function_argument.current, spacing + 2);
        offset_text(spacing);
        if (node->u.function_argument.next != NULL)
            printf("},\n");
        else
            printf("}\n");
        _describe(node->u.function_argument.next, spacing);
        break;
    case FN_DEF:
        offset_text(spacing);
        printf("\e[34mFn Def\e[0m {\n");
        _describe(node->u.function_definition.return_type, spacing + 2);
        _describe(node->u.variable_definition.identifier, spacing + 2);
        spacing += 2;
        offset_text(spacing);
        printf("\e[32mArguments\e[0m (\n");
        _describe(node->u.function_definition.arg_list, spacing + 2);
        offset_text(spacing);
        printf(")\n");
        offset_text(spacing);
        printf("Body\e[0m {\n");
        spacing += 2;
        _describe(node->u.function_definition.body, spacing);
        spacing -= 2;
        offset_text(spacing);
        printf("}\n");
        spacing -= 2;
        offset_text(spacing);
        printf("}\n");
        break;
    case FN_CALL:
        offset_text(spacing);
        printf("\e[34mFn Call\e[0m {\n");
        _describe(node->u.function_call.ident, spacing + 2);
        spacing += 2;
        offset_text(spacing);
        printf("\e[32mArguments\e[0m (\n");
        _describe(node->u.function_call.first_arg, spacing + 2);
        offset_text(spacing);
        printf(")\n");
        spacing -= 2;
        offset_text(spacing);
        printf("}\n");
        break;
    }
}

void describe(struct ast *node)
{
    _describe(node, 0);
}
