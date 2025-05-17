%{
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include <ast.h>
#include <str.h>

int yylex(void);
int yyerror(const char *s);
static struct ast *root;
%}


%union {
    struct ast *node;
    struct str str;
    uint64_t num;
    double dec;
}

%token <num> INT_TOK
%token <dec> FLOAT_TOK
%token <str> IDENTIFIER_TOK
%type <node> program statements statement expr expr0 vars vars0 type pointer_type ident ident_chain fns fn_arg_list

%left PLUS_TOK MINUS_TOK
%left STAR_TOK SLASH_TOK

%%

program:
    statements                      { root = program_node($1); }

statements:
    statement ';' statements        { $$ = statement_node($3, $1); }
    | expr0 ';' statements          { $$ = statement_node($3, $1); }
    |                               { $$ = NULL; }
    ;

ident:
    IDENTIFIER_TOK                  { $$ = identifier_node($1); }

ident_chain:
    ident                           { $$ = $1;                            }
    | ident ',' ident_chain         { $$ = identifier_chain_node($3, $1); }
    ;

statement:
    vars                            { $$ = $1; }
    | fns                           { $$ = $1; }
    ;

fns:
    type ident '(' fn_arg_list ')'  { $$ = fn_decl_node($1, $2, $4); }
    ;

fn_arg_list:
    vars0                           { $$ = fn_arg_list_node(NULL, $1); }
    | vars0 ',' fn_arg_list         { $$ = fn_arg_list_node($3, $1);   }
    |                               { $$ = NULL;                       }
    ;

vars0:
    type ident                      { $$ = var_decl_node(type_node($1), $2);    }
    | type ident '=' expr           { $$ = var_def_node(type_node($1), $2, $4); }
    ;

vars:
    type ident_chain                { $$ = var_decl_node(type_node($1), $2);    }
    | type ident '=' expr           { $$ = var_def_node(type_node($1), $2, $4); }
    ;

type:
    ident                           { $$ = pointer_node(NULL, $1); }
    | pointer_type                  { $$ = $1; }
    ;

pointer_type:
    type STAR_TOK                   { $$ = pointer_node($1, NULL); }
    ;

expr0:
    INT_TOK                         { $$ = int_node($1);                }
    | FLOAT_TOK                     { $$ = float_node($1);              }
    | ident                         { $$ = $1;                          }
    | expr0 PLUS_TOK expr0          { $$ = operation_node('+', $1, $3); }
    | expr0 MINUS_TOK expr0         { $$ = operation_node('-', $1, $3); }
    | expr0 SLASH_TOK expr0         { $$ = operation_node('/', $1, $3); }
    | '(' expr ')'                  { $$ = bracket_node($2);            }
    ;

expr:
    expr0                           { $$ = $1;                          }
    | expr STAR_TOK expr            { $$ = operation_node('*', $1, $3); }
    ;
%%


struct ast *parse() {
    yyparse();
    return root;
}

int yyerror(const char *s) {
    fprintf(stderr, "Error: %s\n", s);
    return 0;
}
