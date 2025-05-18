%{
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include <ast.h>
#include <str.h>


extern int yylineno; // Line number from the lexer
extern char *yytext; // Current token text

int yylex(void);
int yyerror(const char *s);
static struct ast *root;

// FIXME: .* as deref
%}


%union {
    struct ast *node;
    struct str str;
    short boolean;
    char ch;
    uint64_t num;
    double dec;
}

%token <num> INT_TOK
%token <dec> FLOAT_TOK
%token <str> IDENTIFIER_TOK STRING_TOK
%token <ch> CHAR_TOK
%token <boolean> BOOL_TOK
%type <node> program statements statement expr ident vars type pointer_type fns fn_args body call_args value unary 
%type <node> if_cond elsif_branch else_branch fors ident_chain whiles expr1

%token DEF_TOK END_TOK IF_TOK THEN_TOK ELSIF_TOK ELSE_TOK FOR_TOK DO_TOK WHILE_TOK LOOP_TOK

%right '='
%left PLUS_TOK MINUS_TOK
%left STAR_TOK SLASH_TOK
%left AMP_TOK OR_TOK

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
    | if_cond                       { $$ = $1; }
    | fors                          { $$ = $1; }
    | whiles                        { $$ = $1; }
    ;

whiles:
    WHILE_TOK expr DO_TOK body END_TOK { $$ = while_node($2, $4, 0); }
    | LOOP_TOK DO_TOK body END_TOK     { $$ = while_node(NULL, $3, 1); }
    ;

ident_chain:
    ident                           { $$ = identifier_chain_node(NULL, $1); }
    | ident ',' ident_chain         { $$ = identifier_chain_node($3, $1);   }
    |                               { $$ = NULL; };
    ;

fors:
    FOR_TOK expr DO_TOK body END_TOK { $$ = for_node($2, NULL, $4); }
    | FOR_TOK expr OR_TOK DO_TOK body END_TOK { $$ = for_node($2, NULL, $5); }
    | FOR_TOK expr '|' ident_chain '|' DO_TOK body END_TOK { $$ = for_node($2, $4, $7); }
    ;

if_cond:
    IF_TOK expr1 THEN_TOK body END_TOK           { $$ = if_node($2, $4, NULL); }
    | IF_TOK expr1 THEN_TOK body elsif_branch    { $$ = if_node($2, $4, $5);   }    
    | IF_TOK expr1 THEN_TOK body else_branch     { $$ = if_node($2, $4, $5);   }
    ;

elsif_branch:
    ELSIF_TOK expr THEN_TOK body END_TOK        { $$ = elsif_node($2, $4, NULL); }
    | ELSIF_TOK expr THEN_TOK body elsif_branch { $$ = elsif_node($2, $4, $5); }
    | ELSIF_TOK expr THEN_TOK body else_branch  { $$ = elsif_node($2, $4, $5); }
    ;

else_branch:
    ELSE_TOK body END_TOK           { $$ = else_node($2); }
    ;

fns:
    DEF_TOK ident '(' fn_args ')' body END_TOK            { $$ = fn_def_node(type_node(NULL), $2, $4, $6); }
    | DEF_TOK ident '(' fn_args ')' ':' type body END_TOK { $$ = fn_def_node(type_node($7), $2, $4, $8); }
    ;

call_args:
    expr                            { $$ = fn_arg_list_node(NULL, $1); }
    | expr ',' call_args            { $$ = fn_arg_list_node($3, $1);   }
    |                               { $$ = NULL; }
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

value:
    INT_TOK                         { $$ = int_node($1);    }
    | FLOAT_TOK                     { $$ = float_node($1);  }
    | STRING_TOK                    { $$ = string_node($1); }
    | CHAR_TOK                      { $$ = char_node($1);   }
    | BOOL_TOK                      { $$ = bool_node($1);   }
    | ident                         { $$ = $1;              }
    ;

unary:
    value                            { $$ = $1; }
    | PLUS_TOK expr1                 { $$ = unary_node('+', $2); }
    | MINUS_TOK expr1                { $$ = unary_node('-', $2); }
    | AMP_TOK expr1                  { $$ = unary_node('&', $2); }

expr:
    expr1                                        { $$ = $1; }
    | expr1 IF_TOK expr1                         { $$ = expr_if_node($1, $3);     }
    | IF_TOK expr1 THEN_TOK expr1 ELSE_TOK expr1 { $$ = if_expr_node($2, $4, $6); }
    ;

expr1:
    unary                             { $$ = $1; }
    | expr1 PLUS_TOK expr1            { $$ = operation_node('+', $1, $3); }
    | expr1 MINUS_TOK expr1           { $$ = operation_node('-', $1, $3); }
    | expr1 STAR_TOK expr1            { $$ = operation_node('*', $1, $3); }
    | expr1 SLASH_TOK expr1           { $$ = operation_node('/', $1, $3); }
    | ident '(' call_args ')'         { $$ = fn_call_node($1, $3);        }
    | '(' expr ')' '(' call_args ')'  { $$ = fn_call_node($2, $5);        }
    | '(' expr ')'                    { $$ = bracket_node($2);            }
    ;
%%


struct ast *parse() {
    yyparse();
    return root;
}

int yyerror(const char *s) {
    fprintf(stderr, "Error at line %d: %s: '%s'\n", yylineno, s, yytext);
    return 0;
}
