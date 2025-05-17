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
%type <node> program statements statement expr ident vars type pointer_type fns fn_args body fn_call call_args

%token DEF_TOK END_TOK

%left PLUS_TOK MINUS_TOK
%left STAR_TOK SLASH_TOK

%%

program:
    statements                      { root = program_node($1); }
    ;

statements:
    statement statements            { $$ = statement_node($2, $1); }
    | expr ';' statements           { $$ = statement_node($3, $1); }
    |                               { $$ = NULL; }
    ;

statement:
    vars ';'                        { $$ = $1; }
    | fns                           { $$ = $1; }
    | fn_call ';'                   { $$ = $1; }
    ;

fns:
    DEF_TOK ident '(' fn_args ')' body END_TOK            { $$ = fn_def_node(type_node(NULL), $2, $4, $6); }
    | DEF_TOK ident '(' fn_args ')' ':' type body END_TOK   { $$ = fn_def_node(type_node($7), $2, $4, $8); }
    ;

fn_call:
    ident call_args                 { $$ = fn_call_node($1, $2); }
    | ident '(' call_args ')'       { $$ = fn_call_node($1, $3); }
    ;

call_args:
    expr                            { $$ = fn_arg_list_node(NULL, $1); }
    | expr ',' call_args            { $$ = fn_arg_list_node($3, $1);   }
    ;

body:
    statements                      { $$ = $1; }
    ;

fn_args:
    vars                            { $$ = fn_arg_list_node(NULL, $1); }
    | vars ',' fn_args              { $$ = fn_arg_list_node($3, $1);   }
    |                               { $$ = NULL;                       }
    ;

vars:
    ident '=' expr                  { $$ = var_def_node(NULL, $1, $3); }
    | ident ':' type                { $$ = var_decl_node(type_node($3), $1); }
    | ident ':' type '=' expr       { $$ = var_def_node(type_node($3), $1, $5); }
    ;

type:
    ident                           { $$ = pointer_node(NULL, $1); }
    | pointer_type                  { $$ = $1; }
    ;

pointer_type:
    type STAR_TOK                   { $$ = pointer_node($1, NULL); }
    ;

ident:
    IDENTIFIER_TOK                  { $$ = identifier_node($1); }

expr:
    INT_TOK                         { $$ = int_node($1);                }
    | FLOAT_TOK                     { $$ = float_node($1);              }
    | ident                         { $$ = $1;                          }
    | expr PLUS_TOK expr            { $$ = operation_node('+', $1, $3); }
    | expr MINUS_TOK expr           { $$ = operation_node('-', $1, $3); }
    | expr STAR_TOK expr            { $$ = operation_node('*', $1, $3); }
    | expr SLASH_TOK expr           { $$ = operation_node('/', $1, $3); }
    | ident '(' call_args ')'       { $$ = fn_call_node($1, $3);        }
    | '(' expr ')'                  { $$ = bracket_node($2);            }
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
