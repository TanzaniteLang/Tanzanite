%{
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include <ast.h>
#include <str.h>


extern int yylineno;
extern int start_column;
extern int yycolumn;
extern char *yytext;

int yylex(void);
int yyerror(const char *s);
static struct ast *root;

/* TODO: ++ and -- */
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
%type <node> if_cond elsif_branch else_branch fors ident_chain whiles expr1 field_access assignment

%token IF_TOK UNLESS_TOK ELSE_TOK ELSIF_TOK FOR_TOK WHILE_TOK UNTIL_TOK BREAK_TOK NEXT_TOK CASE_TOK WHEN_TOK DEF_TOK
%token FUN_TOK SIZEOF_TOK BEGIN_TOK RETURN_TOK LOOP_TOK RESCUE_TOK THEN_TOK DO_TOK END_TOK WITH_TOK AUTO_TOK

%left '+' INCREMENT_TOK '-' DECREMENT_TOK '*' '/' FLOOR_DIV_TOK '%' '~' '&' '|' '^' PIPE_FORWARD_TOK
%left LEFT_SHIFT_TOK RIGHT_SHIFT_TOK
%left EQL_TOK '!' NOT_EQL_TOK AND_TOK OR_TOK '<' LESS_THAN_EQL_TOK '>' GREATER_THAN_EQL_TOK
%left '.'
%right '=' ADD_ASSIGN_TOK SUB_ASSIGN_TOK MUL_ASSIGN_TOK DIV_ASSIGN_TOK FLOOR_DIV_ASSIGN_TOK MOD_ASSIGN_TOK
%right BIT_NOT_ASSIGN_TOK BIT_AND_ASSIGN_TOK BIT_OR_ASSIGN_TOK XOR_ASSIGN_TOK
%right LEFT_SHIFT_ASSIGN_TOK RIGHT_SHIFT_ASSIGN_TOK

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
    | FOR_TOK expr WITH_TOK OR_TOK DO_TOK body END_TOK { $$ = for_node($2, NULL, $6); }
    | FOR_TOK expr WITH_TOK '|' ident_chain '|' DO_TOK body END_TOK { $$ = for_node($2, $5, $8); }
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
    ident ':' AUTO_TOK '=' expr     { $$ = var_def_node(NULL, $1, $5); }
    | ident ':' type                { $$ = var_decl_node(type_node($3), $1); }
    | ident ':' type '=' expr       { $$ = var_def_node(type_node($3), $1, $5); }
    ;

type:
    ident                           { $$ = pointer_node(NULL, $1); }
    | pointer_type                  { $$ = $1; }
    ;

pointer_type:
    type '*'                   { $$ = pointer_node($1, NULL); }
    ;

ident:
    IDENTIFIER_TOK                  { $$ = identifier_node($1); }
    ;

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
    | '+' value                      { $$ = unary_node("+", $2);      }
    | '+' '(' expr ')'               { $$ = unary_node("+", $3);      }
    | INCREMENT_TOK value            { $$ = unary_node("++", $2);     }
    | INCREMENT_TOK '(' expr ')'     { $$ = unary_node("++", $3);     }
    | '-' value                      { $$ = unary_node("-", $2);      }
    | '-' '(' expr ')'               { $$ = unary_node("-", $3);      }
    | DECREMENT_TOK value            { $$ = unary_node("--", $2);     }
    | DECREMENT_TOK '(' expr ')'     { $$ = unary_node("--", $3);     }
    | '!' value                      { $$ = unary_node("!", $2);      }
    | '!' '(' expr ')'               { $$ = unary_node("!", $3);      }
    | '~' value                      { $$ = unary_node("~", $2);      }
    | '~' '(' expr ')'               { $$ = unary_node("~", $3);      }
    | '&' value                      { $$ = unary_node("&", $2);      }
    | '&' '(' expr ')'               { $$ = unary_node("&", $3);      }
    | SIZEOF_TOK value               { $$ = unary_node("sizeof", $2); }
    | SIZEOF_TOK '(' expr ')'        { $$ = unary_node("sizeof", $3); }
    ;

expr:
    expr1                                        { $$ = $1; }
    | expr1 IF_TOK expr1                         { $$ = expr_if_node($1, $3);     }
    | IF_TOK expr1 THEN_TOK expr1 ELSE_TOK expr1 { $$ = if_expr_node($2, $4, $6); }
    ;

field_access:
    expr1 '.' ident                   { $$ = field_access_node($1, $3);    }
    | expr1 '.' '*'                   { $$ = pointer_deref_node($1);       }
    ;

assignment:
    expr1 '=' expr1                         { $$ = assign_node("=", $1, $3);   }
    | expr1 ADD_ASSIGN_TOK expr1            { $$ = assign_node("+=", $1, $3);  }
    | expr1 SUB_ASSIGN_TOK expr1            { $$ = assign_node("-=", $1, $3);  }
    | expr1 MUL_ASSIGN_TOK expr1            { $$ = assign_node("*=", $1, $3);  }
    | expr1 DIV_ASSIGN_TOK expr1            { $$ = assign_node("/=", $1, $3);  }
    | expr1 FLOOR_DIV_ASSIGN_TOK expr1      { $$ = assign_node("//=", $1, $3); }
    | expr1 MOD_ASSIGN_TOK expr1            { $$ = assign_node("&=", $1, $3);  }
    | expr1 LEFT_SHIFT_ASSIGN_TOK expr1     { $$ = assign_node("<<=", $1, $3); }
    | expr1 RIGHT_SHIFT_ASSIGN_TOK expr1    { $$ = assign_node(">>=", $1, $3); }
    | expr1 BIT_NOT_ASSIGN_TOK expr1        { $$ = assign_node("~=", $1, $3);  }
    | expr1 BIT_AND_ASSIGN_TOK expr1        { $$ = assign_node("&=", $1, $3);  }
    | expr1 BIT_OR_ASSIGN_TOK expr1         { $$ = assign_node("|=", $1, $3);  }
    | expr1 XOR_ASSIGN_TOK expr1            { $$ = assign_node("^=", $1, $3);  }
    ;

expr1:
    field_access                      { $$ = $1; }
    | unary                           { $$ = $1;                           }
    | expr1 '+' expr1                 { $$ = operation_node("+", $1, $3);  }
    | expr1 '-' expr1                 { $$ = operation_node("-", $1, $3);  }
    | expr1 '*' expr1                 { $$ = operation_node("*", $1, $3);  }
    | expr1 '/' expr1                 { $$ = operation_node("/", $1, $3);  }
    | expr1 FLOOR_DIV_TOK expr1       { $$ = operation_node("//", $1, $3); }
    | expr1 '%' expr1                 { $$ = operation_node("%", $1, $3);  }
    | expr1 LEFT_SHIFT_TOK expr1      { $$ = operation_node("<<", $1, $3); }
    | expr1 RIGHT_SHIFT_TOK expr1     { $$ = operation_node(">>", $1, $3); }
    | expr1 EQL_TOK expr1             { $$ = operation_node("==", $1, $3); }
    | expr1 NOT_EQL_TOK expr1         { $$ = operation_node("!=", $1, $3); }
    | expr1 '&' expr1                 { $$ = operation_node("&", $1, $3);  }
    | expr1 '^' expr1                 { $$ = operation_node("^", $1, $3);  }
    | expr1 '|' expr1                 { $$ = operation_node("|", $1, $3);  }
    | expr1 AND_TOK expr1             { $$ = operation_node("&&", $1, $3); }
    | expr1 OR_TOK expr1              { $$ = operation_node("||", $1, $3); }
    | expr1 PIPE_FORWARD_TOK expr1    { $$ = operation_node("|>", $1, $3); }
    | assignment                      { $$ = $1; }
    | ident '(' call_args ')'         { $$ = fn_call_node($1, $3);         }
    | field_access '(' call_args ')'  { $$ = fn_call_node($1, $3);         }
    | '(' expr ')' '(' call_args ')'  { $$ = fn_call_node($2, $5);         }
    | '(' expr ')'                    { $$ = bracket_node($2);             }
    ;
%%


struct ast *parse() {
    yyparse();
    return root;
}

int yyerror(const char *s) {
    fprintf(stderr, "Error at line (%d:%d): %s: '%s'\n", yylineno, start_column, s, yytext);
    return 0;
}
