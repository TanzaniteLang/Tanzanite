#ifndef __AST_H__
#define __AST_H__

#include <stdint.h>
#include <str.h>

enum node_type {
    PROGRAM,
    STATEMENT,
    BRACKETS,

    INT,
    FLOAT,
    IDENTIFIER,
    STRING,
    IDENTIFIER_CHAIN,
    UNARY,

    OPERATION,

    VAR_DECL,
    VAR_DEF,
    TYPE_NODE,
    POINTER,

    FN_DECL,
    FN_DEF,
    FN_ARG,
    FN_CALL,

    IF_COND,
    ELSIF_COND,
    ELSE_COND,
};

struct ast {
    enum node_type type;
    union {
        struct ast *program;
        struct {
            struct ast *current;
            struct ast *next;
        } statement;
        struct ast *bracket;
        uint64_t number;
        double decimal;
        struct str identifier;
        struct str string;
        struct {
            struct ast *current;
            struct ast *next;
        } identifier_chain;
        struct {
            char op;
            struct ast *value;
        } unary;
        struct {
            char op;
            struct ast *left;
            struct ast *right;
        } operation;
        struct {
            struct ast *type;
            struct ast *identifier;
            struct ast *value;
        } variable_definition;
        struct {
            struct ast *type;
            struct ast *identifier;
        } variable_declaration;
        struct ast *type;
        struct {
            struct ast *current;
            struct ast *next;
        } pointer;
        struct {
            struct ast *return_type;
            struct ast *ident;
            struct ast *arg_list;
        } function_declaration;
        struct {
            struct ast *return_type;
            struct ast *ident;
            struct ast *arg_list;
            struct ast *body;
        } function_definition;
        struct {
            struct ast *current;
            struct ast *next;
        } function_argument;
        struct {
            struct ast *ident;
            struct ast *first_arg;
        } function_call;
        struct {
            struct ast *expr;
            struct ast *body;
            struct ast *next;
        } if_statement;
        struct {
            struct ast *expr;
            struct ast *body;
            struct ast *next;
        } elsif_statement;
        struct ast *else_statement;
    } u;
};

struct ast *parse();

struct ast *program_node(struct ast *statement);
struct ast *statement_node(struct ast *list, struct ast *statement);
struct ast *int_node(uint64_t val);
struct ast *float_node(double val);
struct ast *identifier_node(struct str ident);
struct ast *string_node(struct str string);
struct ast *identifier_chain_node(struct ast *list, struct ast *ident);
struct ast *operation_node(char op, struct ast *left, struct ast *right);
struct ast *bracket_node(struct ast *expr);
struct ast *var_decl_node(struct ast *type, struct ast *ident);
struct ast *var_def_node(struct ast *type, struct ast *ident, struct ast *val);
struct ast *fn_decl_node(struct ast *type, struct ast *ident, struct ast *args);
struct ast *fn_def_node(struct ast *type, struct ast *ident, struct ast *args, struct ast *body);
struct ast *fn_call_node(struct ast *ident, struct ast *first_arg);
struct ast *fn_arg_list_node(struct ast *list, struct ast *arg);
struct ast *type_node(struct ast *type);
struct ast *pointer_node(struct ast *list, struct ast *type);
struct ast *if_node(struct ast *expr, struct ast *body, struct ast *next);
struct ast *elsif_node(struct ast *expr, struct ast *body, struct ast *next);
struct ast *else_node(struct ast *body);
struct ast *unary_node(char op, struct ast *val);

void describe(struct ast *node);

#endif
