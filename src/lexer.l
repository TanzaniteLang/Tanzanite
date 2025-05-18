%{
#include <ast.h>
#include <str.h>
#include <parser.h>
#include <stdbool.h>
%}

%%
"if"                return IF_TOK;
"elsif"             return ELSIF_TOK;
"else"              return ELSE_TOK;
"then"              return THEN_TOK;
"def"               return DEF_TOK;
"for"               return FOR_TOK;
"do"                return DO_TOK;
"end"               return END_TOK;
"while"             return WHILE_TOK;
"loop"              return LOOP_TOK;
"true"              { yylval.boolean = true; return BOOL_TOK; }
"false"             { yylval.boolean = false; return BOOL_TOK; }
[A-z]+[A-z0-9]*     { yylval.str = str_init(yytext, yyleng); return IDENTIFIER_TOK; }
\"(?:[^"\\]|\\.)*\" { yylval.str = str_init(yytext + 1, yyleng - 2); return STRING_TOK;     }
'.+'                { yylval.ch = yytext[1]; return CHAR_TOK; }
[0-9]+\.[0-9]+      { yylval.dec = strtod(yytext, NULL); return FLOAT_TOK; }
[0-9]+              { yylval.num = atol(yytext); return INT_TOK; }
[ \t\n]             ;
","                 return ',';
"'"                 return '\'';
";"                 return ';';
"="                 return '=';
":"                 return ':';
"{"                 return '{';
"}"                 return '}';
"("                 return '(';
")"                 return ')';
"|"                 return '|';
"||"                return OR_TOK;
"&"                 return AMP_TOK;
"+"                 return PLUS_TOK;
"-"                 return MINUS_TOK;
"*"                 return STAR_TOK;
"/"                 return SLASH_TOK;

.                   { fprintf(stderr, "???: %s\n", yytext); abort(); }

%%

int yywrap(void) {
    return 1;
}
