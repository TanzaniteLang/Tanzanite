#ifndef __AST_H__
#define __AST_H__

#include <stdint.h>
#include <str.h>
#include <stdbool.h>

enum node_type {
    /* Parser nodes */
    PROGRAM,
    STATEMENT,
    BRACKETS,
    INT,
    FLOAT,
    IDENTIFIER,
    CHAR,
    BOOL,
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
    IF_EXPR,
    EXPR_IF,
    ELSIF_COND,
    ELSE_COND,
    FOR,
    WHILE,
    FIELD_ACCESS,
    POINTER_DEREF,
    ASSIGNMENT,
    TYPE_CAST,
    NEXT,
    BREAK,
    VARIADIC,

    /* Analysis special nodes */

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
        char ch;
        bool boolean;
        struct {
            struct ast *current;
            struct ast *next;
        } identifier_chain;
        struct {
            char *op;
            struct ast *value;
        } unary;
        struct {
            char *op;
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
            bool immutable;
        } function_declaration;
        struct {
            struct ast *return_type;
            struct ast *ident;
            struct ast *arg_list;
            struct ast *body;
            bool immutable;
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
            bool unless;
        } if_statement;
        struct {
            struct ast *expr;
            struct ast *val;
            struct ast *else_val;
            bool unless;
        } if_expression;
        struct {
            struct ast *expr;
            struct ast *condition;
            bool unless;
        } expression_if;
        struct {
            struct ast *expr;
            struct ast *body;
            struct ast *next;
        } elsif_statement;
        struct ast *else_statement;
        struct {
            struct ast *expr;
            struct ast *capture;
            struct ast *body;
        } for_statement;
        struct {
            struct ast *expr;
            struct ast *body;
            bool do_while;
            bool until;
        } while_statement;
        struct {
            struct ast *left;
            struct ast *right;
        } field_access;
        struct ast *to_deref;
        struct {
            char *op;
            struct ast *left;
            struct ast *right;
        } assignment;
        struct {
            struct ast *expr;
            struct ast *type;
        } type_cast;
    } u;
};

struct ast *parse();

struct ast *program_node(struct ast *statement);
struct ast *statement_node(struct ast *list, struct ast *statement);
struct ast *int_node(uint64_t val);
struct ast *float_node(double val);
struct ast *identifier_node(struct str ident);
struct ast *string_node(struct str string);
struct ast *char_node(char ch);
struct ast *bool_node(short boolean);
struct ast *identifier_chain_node(struct ast *list, struct ast *ident);
struct ast *operation_node(char *op, struct ast *left, struct ast *right);
struct ast *bracket_node(struct ast *expr);
struct ast *var_decl_node(struct ast *type, struct ast *ident);
struct ast *var_def_node(struct ast *type, struct ast *ident, struct ast *val);
struct ast *fn_decl_node(struct ast *type, struct ast *ident, struct ast *args, bool immutable);
struct ast *fn_def_node(struct ast *type, struct ast *ident, struct ast *args, struct ast *body, bool immutable);
struct ast *fn_call_node(struct ast *ident, struct ast *first_arg);
struct ast *fn_arg_list_node(struct ast *list, struct ast *arg);
struct ast *type_node(struct ast *type);
struct ast *pointer_node(struct ast *list, struct ast *type);
struct ast *if_node(struct ast *expr, struct ast *body, struct ast *next, bool unless);
struct ast *expr_if_node(struct ast *expr, struct ast *condition, bool unless);
struct ast *if_expr_node(struct ast *expr, struct ast *value, struct ast *else_value, bool unless);
struct ast *elsif_node(struct ast *expr, struct ast *body, struct ast *next);
struct ast *else_node(struct ast *body);
struct ast *unary_node(char *op, struct ast *val);
struct ast *for_node(struct ast *expr, struct ast *capture, struct ast *body);
struct ast *while_node(struct ast *expr, struct ast *body, bool do_while, bool until);
struct ast *field_access_node(struct ast *left, struct ast *right);
struct ast *pointer_deref_node(struct ast *expr);
struct ast *assign_node(char *op, struct ast *left, struct ast *right);
struct ast *type_cast_node(struct ast *expr, struct ast *type);
struct ast *break_node();
struct ast *next_node();
struct ast *variadic_node();

void describe(struct ast *node);

#endif
