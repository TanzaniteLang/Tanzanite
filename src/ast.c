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

struct ast *string_node(struct str string)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = STRING;
    node->u.string = string;

    return node;
}

struct ast *char_node(char ch)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = CHAR;
    node->u.ch = ch;

    return node;
}

struct ast *bool_node(short boolean)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = BOOL;
    node->u.boolean = boolean;

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

struct ast *operation_node(char *op, struct ast *left, struct ast *right)
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

struct ast *fn_decl_node(struct ast *type, struct ast *ident, struct ast *args, bool immutable)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = FN_DECL;
    node->u.function_declaration.return_type = type;
    node->u.function_declaration.ident = ident;
    node->u.function_declaration.arg_list = args;
    node->u.function_declaration.immutable = immutable;

    return node;
}
struct ast *fn_def_node(struct ast *type, struct ast *ident, struct ast *args, struct ast *body, bool immutable)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = FN_DEF;
    node->u.function_definition.return_type = type;
    node->u.function_definition.ident = ident;
    node->u.function_definition.arg_list = args;
    node->u.function_definition.body = body;
    node->u.function_definition.immutable = immutable;

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

struct ast *if_node(struct ast *expr, struct ast *body, struct ast *next, bool unless)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = IF_COND;
    node->u.if_statement.expr = expr;
    node->u.if_statement.body = body;
    node->u.if_statement.next = next;
    node->u.if_statement.unless = unless;

    return node;
}
struct ast *expr_if_node(struct ast *expr, struct ast *condition, bool unless)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = EXPR_IF;
    node->u.expression_if.expr = expr;
    node->u.expression_if.condition = condition;
    node->u.expression_if.unless = unless;

    return node;
}

struct ast *if_expr_node(struct ast *expr, struct ast *value, struct ast *else_value, bool unless)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = IF_EXPR;
    node->u.if_expression.expr = expr;
    node->u.if_expression.val = value;
    node->u.if_expression.else_val = else_value;
    node->u.if_expression.unless = unless;

    return node;
}

struct ast *elsif_node(struct ast *expr, struct ast *body, struct ast *next)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = ELSIF_COND;
    node->u.elsif_statement.expr = expr;
    node->u.elsif_statement.body = body;
    node->u.elsif_statement.next = next;

    return node;
}

struct ast *else_node(struct ast *body)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = ELSE_COND;
    node->u.else_statement = body;

    return node;
}

struct ast *unary_node(char *op, struct ast *val)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = UNARY;
    node->u.unary.op = op;
    node->u.unary.value = val;

    return node;
}

struct ast *for_node(struct ast *expr, struct ast *capture, struct ast *body)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = FOR;
    node->u.for_statement.expr = expr;
    node->u.for_statement.capture = capture;
    node->u.for_statement.body = body;

    return node;
}

struct ast *while_node(struct ast *expr, struct ast *body, bool do_while, bool until)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = WHILE;
    node->u.while_statement.expr = expr;
    node->u.while_statement.body = body;
    node->u.while_statement.do_while = do_while;
    node->u.while_statement.until = until;

    return node;
}

struct ast *field_access_node(struct ast *left, struct ast *right)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = FIELD_ACCESS;
    node->u.field_access.left = left;
    node->u.field_access.right = right;

    return node;
}
struct ast *pointer_deref_node(struct ast *ptr)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = POINTER_DEREF;
    node->u.to_deref = ptr;

    return node;
}

struct ast *assign_node(char *op, struct ast *left, struct ast *right)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = ASSIGNMENT;
    node->u.assignment.op = op;
    node->u.assignment.left = left;
    node->u.assignment.right = right;

    return node;
}

struct ast *type_cast_node(struct ast *expr, struct ast *type)
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = TYPE_CAST;
    node->u.type_cast.expr = expr;
    node->u.type_cast.type = type;

    return node;
}

struct ast *break_node()
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = BREAK;

    return node;
}

struct ast *next_node()
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = NEXT;

    return node;
}

struct ast *variadic_node()
{
    struct ast *node = calloc(1, sizeof(*node));
    node->type = VARIADIC;

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
    case STRING:
        offset_text(spacing);
        printf("\e[36mStr\e[0m: %s\n", node->u.string.str);
        break;
    case CHAR:
        offset_text(spacing);
        printf("\e[36mChar\e[0m: %c\n", node->u.ch);
        break;
    case BOOL:
        offset_text(spacing);
        printf("\e[36mBool\e[0m: %s\n", node->u.boolean ? "true" : "false");
        break;
    case OPERATION:
        offset_text(spacing);
        printf("\e[35mOperation\e[0m: %s {\n", node->u.operation.op);
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
        if (node->u.function_declaration.immutable) {
            offset_text(spacing);
            printf("C Function: Yes\n");
        } else {
            offset_text(spacing);
            printf("C Function: No\n");
        }
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
        _describe(node->u.function_definition.ident, spacing + 2);
        spacing += 2;
        if (node->u.function_definition.immutable) {
            offset_text(spacing);
            printf("C Function: Yes\n");
        } else {
            offset_text(spacing);
            printf("C Function: No\n");
        }
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
    case IF_COND:
        offset_text(spacing);
        printf("\e[33m%s\e[0m {\n", node->u.if_statement.unless ? "Unless" : "If");
        _describe(node->u.if_statement.expr, spacing + 2);
        _describe(node->u.if_statement.body, spacing + 2);
        _describe(node->u.if_statement.next, spacing);
        break;
    case ELSIF_COND:
        offset_text(spacing);
        printf("\e[33mElsif\e[0m {\n");
        _describe(node->u.if_statement.expr, spacing + 2);
        _describe(node->u.if_statement.body, spacing + 2);
        _describe(node->u.if_statement.next, spacing);
        break;
    case ELSE_COND:
        offset_text(spacing);
        printf("\e[33mElse\e[0m {\n");
        _describe(node->u.else_statement, spacing + 2);
        break;
    case UNARY:
        offset_text(spacing);
        printf("\e[35mUnary\e[0m: %s\e[0m {\n", node->u.unary.op);
        _describe(node->u.unary.value, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case FOR:
        offset_text(spacing);
        printf("\e[34mFor\e[0m {\n");
        _describe(node->u.for_statement.expr, spacing + 2);
        spacing += 2;
        offset_text(spacing);
        printf("\e[32mCapture\e[0m |\n");
        _describe(node->u.for_statement.capture, spacing + 2);
        offset_text(spacing);
        printf("|\n");
        offset_text(spacing);
        printf("Body\e[0m {\n");
        _describe(node->u.for_statement.body, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        spacing -= 2;
        offset_text(spacing);
        printf("}\n");
        break;
    case IF_EXPR:
        offset_text(spacing);
        printf("\e[34m%s Expr\e[0m {\n", node->u.if_expression.unless ? "Unless" : "If");
        _describe(node->u.if_expression.expr, spacing + 2);
        _describe(node->u.if_expression.val, spacing + 2);
        _describe(node->u.if_expression.else_val, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case EXPR_IF:
        offset_text(spacing);
        printf("\e[34mExpr %s\e[0m {\n", node->u.expression_if.unless ? "Unless" : "If");
        _describe(node->u.expression_if.expr, spacing + 2);
        _describe(node->u.expression_if.condition, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case WHILE:
        offset_text(spacing);
        printf("\e[34m%s\e[0m {\n", node->u.while_statement.until ? "Until" : "While");
        _describe(node->u.while_statement.expr, spacing + 2);
        spacing += 2;
        offset_text(spacing);
        printf("Body\e[0m {\n");
        _describe(node->u.while_statement.body, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        offset_text(spacing);
        printf("\e[36mIs do-while\e[0m: %s\n", node->u.while_statement.do_while ? "true" : "false");
        spacing -= 2;
        offset_text(spacing);
        printf("}\n");
        break;
    case FIELD_ACCESS:
        offset_text(spacing);
        printf("\e[35mAccess\e[0m {\n");
        _describe(node->u.field_access.left, spacing + 2);
        _describe(node->u.field_access.right, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case POINTER_DEREF:
        offset_text(spacing);
        printf("\e[35mDeref\e[0m {\n");
        _describe(node->u.program, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case ASSIGNMENT:
        offset_text(spacing);
        printf("\e[34mAssignment\e[0m: %s {\n", node->u.assignment.op);
        _describe(node->u.assignment.left, spacing + 2);
        _describe(node->u.assignment.right, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case TYPE_CAST:
        offset_text(spacing);
        printf("\e[31mType Cast\e[0m {\n");
        _describe(node->u.type_cast.expr, spacing + 2);
        _describe(node->u.type_cast.type, spacing + 2);
        offset_text(spacing);
        printf("}\n");
        break;
    case NEXT:
        offset_text(spacing);
        printf("\e[36mNext\e[0m\n");
        break;
    case BREAK:
        offset_text(spacing);
        printf("\e[36mBreak\e[0m\n");
        break;
    case VARIADIC:
        offset_text(spacing);
        printf("\e[36mVariadic\e[0m\n");
        break;
    }
}

void describe(struct ast *node)
{
    _describe(node, 0);
}
