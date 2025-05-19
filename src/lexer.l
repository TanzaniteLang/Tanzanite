%{
#include <ast.h>
#include <str.h>
#include <parser.h>
#include <stdbool.h>

/* https://stackoverflow.com/a/26857402 */
int yycolumn = 1;
int start_column = 0;

#define YY_USER_ACTION                                                   \
  start_line = prev_yylineno; start_column = yycolumn;                   \
  if (yylineno == prev_yylineno) yycolumn += yyleng;                     \
  else {                                                                 \
    for (yycolumn = 1; yytext[yyleng - yycolumn] != '\n'; ++yycolumn) {} \
    prev_yylineno = yylineno;                                            \
  }
%}
%option yylineno

%%
 int start_line = 0;
 int prev_yylineno = yylineno;

 /* Statements */
"="                 return '='; 
"if"                return IF_TOK;
"unless"            return UNLESS_TOK;
"else"              return ELSE_TOK;
"elsif"             return ELSIF_TOK;
"for"               return FOR_TOK;
"while"             return WHILE_TOK;
"until"             return UNTIL_TOK;
"break"             return BREAK_TOK;
"next"              return NEXT_TOK;
"case"              return CASE_TOK;
"when"              return WHEN_TOK;
"def"               return DEF_TOK;
"fun"               return FUN_TOK;
"sizeof"            return SIZEOF_TOK;
"begin"             return BEGIN_TOK;
"return"            return RETURN_TOK;

 /* Constants */
"true"              { yylval.boolean = true; return BOOL_TOK; } 
"false"             { yylval.boolean = false; return BOOL_TOK; }

 /* Delimiters */
"?"                 return '?'; 
"."                 return '.';
"..."               return SPLAT_TOK;
","                 return ',';
":"                 return ':';
";"                 return ';';

 /* Brackets */
"("                 return '('; 
")"                 return ')';
"["                 return '[';
"]"                 return ']';
"{"                 return '{';
"}"                 return '}';

 /* Reserved keywords */
"loop"              return LOOP_TOK; 
"rescue"            return RESCUE_TOK;
"then"              return THEN_TOK;
"do"                return DO_TOK;
"end"               return END_TOK;
"with"              return WITH_TOK;
"auto"              return AUTO_TOK;
"as"                return AS_TOK;

 /* Operators */
"=="                return EQL_TOK;
"+"                 return '+'; 
"++"                return INCREMENT_TOK;
"+="                return ADD_ASSIGN_TOK;
"-"                 return '-';
"--"                return DECREMENT_TOK;
"-="                return SUB_ASSIGN_TOK;
"*"                 return '*';
"*="                return MUL_ASSIGN_TOK;
"/"                 return '/';
"/="                return DIV_ASSIGN_TOK;
"//"                return FLOOR_DIV_TOK;
"//="               return FLOOR_DIV_ASSIGN_TOK;
"%"                 return '%';
"%="                return MOD_ASSIGN_TOK;
"!"                 return '!';
"!="                return NOT_EQL_TOK;
"~"                 return '~';
"~="                return BIT_NOT_ASSIGN_TOK;
"&"                 return '&';
"&="                return BIT_AND_ASSIGN_TOK;
"&&"                return AND_TOK;
"|"                 return '|';
"|="                return BIT_OR_ASSIGN_TOK;
"||"                return OR_TOK;
"|>"                return PIPE_FORWARD_TOK;
"^"                 return '^';
"^="                return XOR_ASSIGN_TOK;
"<"                 return '<';
"<="                return LESS_THAN_EQL_TOK;
"<<"                return LEFT_SHIFT_TOK;
"<<="               return LEFT_SHIFT_ASSIGN_TOK;
">"                 return '>';
">="                return GREATER_THAN_EQL_TOK;
">>"                return RIGHT_SHIFT_TOK;
">>="               return RIGHT_SHIFT_ASSIGN_TOK;

 /* Values */
[A-z]+[A-z0-9]*     { yylval.str = str_init(yytext, yyleng); return IDENTIFIER_TOK;     } 
\"(?:[^"\\]|\\.)*\" { yylval.str = str_init(yytext + 1, yyleng - 2); return STRING_TOK; }
'.+'                { yylval.ch = yytext[1]; return CHAR_TOK;                           }
[0-9]+\.[0-9]+      { yylval.dec = strtod(yytext, NULL); return FLOAT_TOK;              }
[0-9]+              { yylval.num = atol(yytext); return INT_TOK;                        }
[ \t\n]             ;

 /* Dead end */
.                   { fprintf(stderr, "???: %s\n", yytext); abort(); } 

%%

int yywrap(void) {
    return 1;
}
